package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// SyncService is the core synchronization logic service
type SyncService struct {
	ConfigManager *config.Manager
	DriveInfo     *DriveService
	Tree          *FileTree
	Rclone        *RcloneService
	Symedia       *SymediaService
	TriggerChan   chan struct{}

	// Task statistics
	mu                    sync.RWMutex
	activeTasks           int64  // Active task count
	todayCompletedTasks   int64  // Today's completed task count
	historyCompletedTasks int64  // Historical completed task count (total)
	lastResetDate         string // Last date for daily counter reset
	isProcessing          bool   // Whether processing a task

	buildMu sync.Mutex // Mutex for BuildFileTreeSkeleton
}

// NewSyncService creates a new sync service
func NewSyncService(
	cm *config.Manager,
	ds *DriveService,
	tree *FileTree,
	rc *RcloneService,
	sy *SymediaService,
) *SyncService {
	// Load persisted task stats from config
	cm.Lock.RLock()
	todayCompleted := cm.Cfg.Advanced.TaskStats.TodayCompleted
	historyCompleted := cm.Cfg.Advanced.TaskStats.HistoryCompleted
	lastResetDate := cm.Cfg.Advanced.TaskStats.LastResetDate
	cm.Lock.RUnlock()

	logger.Debug(model.LogLevelDebug, "[Stats] Loaded from config - Today: %d, History: %d, LastReset: %s", todayCompleted, historyCompleted, lastResetDate)

	today := time.Now().Format("2006-01-02")
	if lastResetDate != "" && lastResetDate != today {
		// Add today's completed tasks to history before resetting
		historyCompleted += todayCompleted
		todayCompleted = 0
		lastResetDate = today
	}
	if lastResetDate == "" {
		lastResetDate = today
	}

	return &SyncService{
		ConfigManager:         cm,
		DriveInfo:             ds,
		Tree:                  tree,
		Rclone:                rc,
		Symedia:               sy,
		TriggerChan:           make(chan struct{}, 20),
		todayCompletedTasks:   todayCompleted,
		historyCompletedTasks: historyCompleted,
		lastResetDate:         lastResetDate,
	}
}

// StartProcessLoop starts the main event loop
func (s *SyncService) StartProcessLoop() {
	for range s.TriggerChan {
		s.ConfigManager.Lock.RLock()
		db := time.Duration(s.ConfigManager.Cfg.Advanced.DebounceSeconds) * time.Second
		s.ConfigManager.Lock.RUnlock()

		if db == 0 {
			db = 5 * time.Second
		}

		logger.Verbose(model.LogLevelInfo, "⏰ Change detected, debouncing %v...", db)
		time.Sleep(db)

		// Drain extra events
	drain:
		for {
			select {
			case <-s.TriggerChan:
			default:
				break drain
			}
		}

		// Mark task as started
		s.mu.Lock()
		s.isProcessing = true
		s.activeTasks = 1 // Currently processing 1 task
		s.mu.Unlock()

		// Execute sync
		s.SyncOnce()

		// Mark task as completed
		s.mu.Lock()
		s.isProcessing = false
		s.activeTasks = 0 // Task complete, no active tasks

		// Check if we need to reset today's counter (new day)
		today := time.Now().Format("2006-01-02")
		if s.lastResetDate != today {
			// Add today's completed tasks to history before resetting
			s.historyCompletedTasks += s.todayCompletedTasks
			s.todayCompletedTasks = 0
			s.lastResetDate = today
		}

		// Increment counters
		s.todayCompletedTasks++
		s.historyCompletedTasks++

		// Persist stats to config
		s.ConfigManager.Lock.Lock()
		s.ConfigManager.Cfg.Advanced.TaskStats.TodayCompleted = s.todayCompletedTasks
		s.ConfigManager.Cfg.Advanced.TaskStats.HistoryCompleted = s.historyCompletedTasks
		s.ConfigManager.Cfg.Advanced.TaskStats.LastResetDate = s.lastResetDate
		s.ConfigManager.Lock.Unlock()

		// Save config asynchronously to avoid blocking
		go s.ConfigManager.SaveConfig()

		s.mu.Unlock()
	}
}

// BuildFileTreeSkeleton pre-builds the file tree
func (s *SyncService) BuildFileTreeSkeleton(forceRebuild bool) {
	if !s.buildMu.TryLock() {
		logger.Warning("🔒 BuildFileTreeSkeleton already running, skipping")
		return
	}
	defer s.buildMu.Unlock()

	if !forceRebuild {
		// [Optimization] If tree is already loaded (e.g. by main.go), skip
		if cnt := s.Tree.CountNodes(); cnt > 0 {
			logger.Info("📂 File tree already loaded in memory (nodes: %d), skipping build", cnt)
			return
		}

		// Try loading from cache
		if err := s.Tree.Load(); err == nil && s.Tree.CountNodes() > 0 {
			logger.Info("📂 Loaded cached file tree, nodes: %d", s.Tree.CountNodes())
			return
		} else if err != nil {
			logger.Warning("⚠️ Failed to load cache: %v", err)
		} else {
			logger.Warning("⚠️ Cache is empty")
		}
		logger.Info("⚠️ Cache not found or invalid, starting full build...")
	} else {
		logger.Info("♻️ Force rebuilding file tree...")
	}

	// Disk-Buffered Build Logic
	tmpFile := model.TreeCacheFile + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		logger.Error("❌ Failed to create temp cache file: %v", err)
		return
	}
	defer func() {
		f.Close()
		_ = os.Remove(tmpFile) // Cleanup if failed
	}()

	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)

	s.ConfigManager.Lock.RLock()
	targets := s.ConfigManager.Cfg.Google.TargetDriveIDs
	s.ConfigManager.Lock.RUnlock()

	// Helper to scan a specific scope with buffering
	scanScope := func(driveID string, scopeName string) {
		logger.Info("🔍 Scanning %s...", scopeName)
		count := 0
		lastLogCount := 0
		progressInterval := 1000    // Log every 1000 files
		bufferFlushInterval := 1000 // Flush to disk every 1000 records

		err := s.DriveInfo.ListFiles(context.Background(), "trashed = false", "nextPageToken, incompleteSearch, files(id, name, parents, mimeType, driveId)", driveID, func(f *drive.File) bool {
			pid := ""
			if len(f.Parents) > 0 {
				pid = f.Parents[0]
				if len(f.Parents) > 1 {
					logger.Warning("⚠️ [Multi-Parent] Node %s (%s) has %d parents, using first: %s", f.Name, f.Id, len(f.Parents), pid)
				}
			}
			isDir := f.MimeType == "application/vnd.google-apps.folder"

			// Create node
			node := model.FileNode{
				ID:       f.Id,
				Name:     f.Name,
				ParentID: pid,
				IsDir:    isDir,
				DriveID:  f.DriveId,
			}

			// Stream write to disk
			if err := enc.Encode(node); err != nil {
				logger.Error("❌ Error writing node to buffer: %v", err)
				return true // Try continue?
			}

			count++

			// Flush buffer periodically
			if count%bufferFlushInterval == 0 {
				if err := w.Flush(); err != nil {
					logger.Error("❌ Error flushing buffer: %v", err)
				}

				// [Strict Rate Limit] Pause 5min every 1000 items
				batchSleepSec := s.ConfigManager.Cfg.Google.BatchSleepInterval
				if batchSleepSec < 300 {
					batchSleepSec = 300 // Min 5 minutes
				}
				logger.Warning("⏳ [Risk Control] Scanned %d items. Pausing for %d seconds...", count, batchSleepSec)
				time.Sleep(time.Duration(batchSleepSec) * time.Second)
				logger.Info("▶️ Resuming scan...")
			}

			// Log progress
			if count-lastLogCount >= progressInterval {
				logger.Info("   📈 Progress: %d scanned (Buffered)...", count)
				lastLogCount = count
			}

			return true // Continue
		})

		if err != nil {
			if apiErr, ok := err.(*googleapi.Error); ok {
				logger.Error("❌ Failed to scan %s: [NetCode: %d] %v", scopeName, apiErr.Code, apiErr.Message)
			} else {
				logger.Error("❌ Failed to scan %s: %v", scopeName, err)
			}
		} else {
			logger.Info("✅ Scanned %s: %d nodes buffered", scopeName, count)
		}
	}

	if len(targets) > 0 {
		logger.Info("🎯 Target Mode: Scanning %d specific drives", len(targets))

		// Read remarks
		s.ConfigManager.Lock.RLock()
		remarks := s.ConfigManager.Cfg.Google.TargetDriveRemarks
		s.ConfigManager.Lock.RUnlock()

		for _, tid := range targets {
			name := s.DriveInfo.GetDriveName(tid)
			if remark, ok := remarks[tid]; ok && remark != "" {
				name = fmt.Sprintf("%s (%s)", name, remark)
			}
			scanScope(tid, name)
		}
	} else {
		logger.Warning("⚠️ No target drives configured. Skipping global scan.")
		logger.Info("💡 Please configure Target Drives in the dashboard to start syncing.")
		logger.Info("   (Go to Dashboard -> Targets to add Google Drive/Folder IDs)")
	}

	// Final flush
	if err := w.Flush(); err != nil {
		logger.Error("❌ Final flush failed: %v", err)
		return
	}
	f.Close() // Explicit close to ensure write

	// Atomic Swap on Disk
	if err := os.Rename(tmpFile, model.TreeCacheFile); err != nil {
		logger.Error("❌ Failed to commit cache file: %v", err)
		return
	}

	logger.Info("💾 Cache commited to disk. Loading into memory...")

	// Load the new tree from disk
	newTree := NewFileTree(s.DriveInfo)
	if err := newTree.Load(); err != nil {
		logger.Error("❌ Failed to load new tree: %v", err)
		return
	}

	// Atomic Swap in Memory
	s.Tree.ReplaceWith(newTree)

	logger.Info("✅ File tree build complete, final node count: %d", s.Tree.CountNodes())

	// Save to cache async
	go func() {
		if err := s.Tree.Save(); err != nil {
			logger.Error("❌ Failed to save file tree cache: %v", err)
		} else {
			logger.Info("💾 File tree cache saved")
		}
	}()
}

// ForceRebuild forces a file tree rebuild
func (s *SyncService) ForceRebuild() {
	go s.BuildFileTreeSkeleton(true)
}

// SyncOnce performs a single sync check
func (s *SyncService) SyncOnce() {
	if s.DriveInfo.Srv == nil {
		return
	}
	token := s.DriveInfo.GetStartPageToken()
	if token == "" {
		logger.Warning("⚠️ [Diag] PageToken is empty, skipping sync check")
		return
	}
	logger.Verbose(model.LogLevelInfo, "🔄 Checking changes...")
	logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📋 [Diag] Current PageToken: %s", token)
	s.DriveInfo.WaitRateLimit()

	// [Fix] Paginate to get all changes
	var allChanges []*drive.Change
	pageToken := token
	pageCount := 0
	var newStartPageToken string

	for {
		pageCount++
		call := s.DriveInfo.Srv.Changes.List(pageToken).
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, newStartPageToken, changes(fileId, removed, file(name, parents, mimeType, trashed, driveId))").
			PageSize(500)

		r, err := call.Do()
		if err != nil {
			logger.Error("❌ [Diag] Changes API query failed (possible permission issue): %v", err)
			return
		}

		logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📄 [Diag] Page %d: got %d changes, NextPageToken=%v, NewStartPageToken=%v",
			pageCount, len(r.Changes), r.NextPageToken != "", r.NewStartPageToken != "")

		allChanges = append(allChanges, r.Changes...)

		// Save NewStartPageToken (only on last page)
		if r.NewStartPageToken != "" {
			newStartPageToken = r.NewStartPageToken
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📋 [Diag] Got new PageToken: %s", newStartPageToken)
			break
		}

		// Continue if there's a next page
		if r.NextPageToken != "" {
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📄 [Diag] Fetching next page...")
			pageToken = r.NextPageToken
			s.DriveInfo.WaitRateLimit()
		} else {
			// No NewStartPageToken or NextPageToken - shouldn't happen
			logger.Warning("⚠️ [Diag] Changes API has neither NewStartPageToken nor NextPageToken")
			break
		}
	}

	logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📊 [Diag] Total %d pages, %d changes", pageCount, len(allChanges))

	for i, change := range allChanges {
		if change.File != nil {
			driveID := change.File.DriveId
			if driveID == "" {
				driveID = "My Drive"
			}
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "   [%d] FileID=%s, DriveID=%s, Name=%s, Trashed=%v",
				i, change.FileId, driveID, change.File.Name, change.File.Trashed)
		} else {
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "   [%d] FileID=%s (removed=%v, file=nil)",
				i, change.FileId, change.Removed)
		}
	}

	if len(allChanges) == 0 {
		if newStartPageToken != "" && newStartPageToken != token {
			s.DriveInfo.SaveTokenStr(newStartPageToken)
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "💾 [Diag] Saved new PageToken: %s", newStartPageToken)
		}
		logger.Verbose(model.LogLevelInfo, "💤 No changes")
		return
	}

	rcloneDirs := make(map[string]bool)
	type Notif struct {
		Path, Action string
		IsDir        bool
		DriveID      string
	}
	var notifs []Notif
	processedIDs := make(map[string]bool)

	for _, change := range allChanges {
		fileID := change.FileId
		if processedIDs[fileID] {
			continue
		}
		oldPath, foundOld := s.Tree.GetPath(fileID)
		isDeleted := change.Removed || (change.File != nil && change.File.Trashed)

		if isDeleted {
			if foundOld {
				descendants := s.Tree.GetDescendants(fileID)
				// Sort by path length descending, delete children first
				sort.Slice(descendants, func(i, j int) bool {
					return len(descendants[i].Path) > len(descendants[j].Path)
				})

				for _, d := range descendants {
					processedIDs[d.ID] = true
					logger.Info("🗑️ [Delete] %s", d.Path)
					logger.WriteHistory(s.ConfigManager.Cfg, "DELETE", d.Path)
					notifs = append(notifs, Notif{d.Path, "delete", d.IsDir, d.DriveID})
					rcloneDirs[filepath.Dir(d.Path)] = true
					s.Tree.RemoveNode(d.ID)
				}
			}
			continue
		}

		if change.File == nil {
			continue
		}

		f := change.File
		pid := ""
		if len(f.Parents) > 0 {
			pid = f.Parents[0]
		}

		// [Strict Scope Check]
		// Ensure we ONLY process changes from Target Drives.
		// If a file moves out of a Target Drive, we silent-delete it.
		// If an event is from a non-Target Drive, we ignore it completely.
		s.ConfigManager.Lock.RLock()
		targets := s.ConfigManager.Cfg.Google.TargetDriveIDs
		s.ConfigManager.Lock.RUnlock()

		if len(targets) > 0 {
			// Determine DriveID
			dID := change.DriveId
			if dID == "" && change.File != nil {
				dID = change.File.DriveId
			}

			// If dID matches a target, we process it.
			// "root" or empty usually maps to My Drive, check if that's targeted.
			checkID := dID
			if checkID == "" {
				checkID = "root"
			}

			isTarget := false
			for _, t := range targets {
				if t == checkID {
					isTarget = true
					break
				}
			}

			// For Deletes (File is nil), we must look up the old node to check its origin DriveID.
			if change.File == nil {
				if oldNode, ok := s.Tree.GetNode(change.FileId); ok {
					// Check if the NODE was in a target drive
					oldDID := oldNode.DriveID
					oldCheckID := oldDID
					if oldCheckID == "" {
						oldCheckID = "root"
					}
					isOldTarget := false
					for _, t := range targets {
						if t == oldCheckID {
							isOldTarget = true
							break
						}
					}

					if !isOldTarget {
						// It was tracked but apparently not in target list anymore?
						// Or it was never in target list (legacy)?
						// Either way, if it's not a target drive, we ignore the external notification.
						// But since it IS in our tree, we should probably remove it silently to clean up.
						// Let's assume strict mode: ignore notification.
						// Should we remove from tree? Yes, to be safe.
						s.Tree.RemoveNode(change.FileId)
						logger.Verbose(model.LogLevelInfo, "🗑️ [Scope] Silently removed non-target node %s (DriveID: %s)", change.FileId, oldDID)
						continue
					}
					// If it WAS in target, proceed to Delete logic below.
				} else {
					// Not in tree, and not a file update. Ignore.
					continue
				}
			} else {
				// It's a File Update/Create/Move
				if !isTarget {
					// New state is NON-TARGET.
					// Check if we were tracking it.
					if _, exists := s.Tree.GetPath(change.FileId); exists {
						// It moved OUT of scope.
						// Silent delete.
						s.Tree.RemoveNode(change.FileId)
						logger.Info("📤 [Scope] Node %s moved out of target scope (DriveID: %s). Silently removing.", change.FileId, dID)
					} else {
						// We never knew it, and it's not a target. Fully ignore.
						// verbose log only
						// logger.Verbose(model.LogLevelDebug, "💤 [Scope] Ignored non-target change %s (DriveID: %s)", change.FileId, dID)
					}
					continue
				}
				// It IS a target. Proceed to Update logic.
			}
		}

		isDirBool := f.MimeType == "application/vnd.google-apps.folder"

		s.Tree.UpdateNode(fileID, f.Name, pid, isDirBool, f.DriveId)
		processedIDs[fileID] = true
		newPath := s.Tree.ResolvePathWithFallback(fileID)

		if !foundOld {
			logger.Info("🆕 [Create] %s", newPath)
			logger.WriteHistory(s.ConfigManager.Cfg, "CREATE", newPath)
			rcloneDirs[filepath.Dir(newPath)] = true
			notifs = append(notifs, Notif{newPath, "create", isDirBool, f.DriveId})
		} else if oldPath != newPath {
			logger.Info("✏️ [Move] %s -> %s", oldPath, newPath)
			logger.WriteHistory(s.ConfigManager.Cfg, "MOVE", newPath)
			rcloneDirs[filepath.Dir(oldPath)] = true
			rcloneDirs[filepath.Dir(newPath)] = true

			// For moves, send old path delete and new path create
			notifs = append(notifs, Notif{oldPath, "delete", isDirBool, f.DriveId})
			notifs = append(notifs, Notif{newPath, "create", isDirBool, f.DriveId})

			if isDirBool {
				descendants := s.Tree.GetDescendants(fileID)
				for _, d := range descendants {
					if d.ID == fileID {
						continue
					}
					processedIDs[d.ID] = true
					relPath := strings.TrimPrefix(d.Path, newPath)
					oldChildPath := oldPath + relPath
					logger.Info("   ↳ [ChildMove] %s -> %s", oldChildPath, d.Path)
					notifs = append(notifs, Notif{oldChildPath, "delete", d.IsDir, d.DriveID})
					notifs = append(notifs, Notif{d.Path, "create", d.IsDir, d.DriveID})
				}
			}
		}
	}

	if len(rcloneDirs) > 0 {
		logger.Info("🚀 Refreshing %d Rclone directories...", len(rcloneDirs))
		var wg sync.WaitGroup
		for dir := range rcloneDirs {
			wg.Add(1)
			go func(d string) {
				defer wg.Done()
				s.Rclone.Refresh(d)
			}(dir)
		}
		wg.Wait()
		s.Rclone.WaitForCooldown()
	}

	if len(notifs) > 0 {
		logger.Info("📡 Sending %d notifications...", len(notifs))
		for _, n := range notifs {
			s.Symedia.SendWebhook(n.Path, n.Action, n.IsDir, n.DriveID)
		}
	}

	if newStartPageToken != "" {
		s.DriveInfo.SaveTokenStr(newStartPageToken)
		logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "💾 [Diag] Final save of new PageToken: %s", newStartPageToken)
	}
}

// TaskStats holds task statistics
type TaskStats struct {
	ActiveTasks           int64 // Currently active tasks
	TodayCompletedTasks   int64 // Tasks completed today
	HistoryCompletedTasks int64 // Total tasks completed historically
}

// GetTaskStats returns task statistics
func (s *SyncService) GetTaskStats() TaskStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if we need to reset today's counter (new day check)
	today := time.Now().Format("2006-01-02")
	todayTasks := s.todayCompletedTasks
	historyTasks := s.historyCompletedTasks
	if s.lastResetDate != "" && s.lastResetDate != today {
		// When querying after midnight, today's tasks should show as 0
		// and history should include the previous day's tasks
		historyTasks += todayTasks
		todayTasks = 0
	}

	return TaskStats{
		ActiveTasks:           s.activeTasks,
		TodayCompletedTasks:   todayTasks,
		HistoryCompletedTasks: historyTasks,
	}
}
