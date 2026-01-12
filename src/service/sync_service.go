package service

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"

	"google.golang.org/api/drive/v3"
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
	activeTasks           int64     // Active task count
	todayCompletedTasks   int64     // Today's completed task count
	historyCompletedTasks int64     // Historical completed task count (total)
	lastResetDate         string    // Last date for daily counter reset
	isProcessing          bool      // Whether processing a task
}

// isIgnoredParent checks if parent folder ID is in ignore list
func (s *SyncService) isIgnoredParent(parentID string) bool {
	s.ConfigManager.Lock.RLock()
	defer s.ConfigManager.Lock.RUnlock()
	for _, id := range s.ConfigManager.Cfg.Google.IgnoredParents {
		if id == parentID {
			return true
		}
	}
	return false
}

// NewSyncService creates a new sync service
func NewSyncService(
	cm *config.Manager,
	ds *DriveService,
	tree *FileTree,
	rc *RcloneService,
	sy *SymediaService,
) *SyncService {
	// Initialize ignore list
	if cm.Cfg.Google.IgnoredParents != nil {
		tree.SetIgnoredParents(cm.Cfg.Google.IgnoredParents)
	}

	return &SyncService{
		ConfigManager: cm,
		DriveInfo:     ds,
		Tree:          tree,
		Rclone:        rc,
		Symedia:       sy,
		TriggerChan:   make(chan struct{}, 20),
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
			s.todayCompletedTasks = 0
			s.lastResetDate = today
		}
		
		// Increment counters
		s.todayCompletedTasks++
		s.historyCompletedTasks++
		s.mu.Unlock()
	}
}

// BuildFileTreeSkeleton pre-builds the file tree
func (s *SyncService) BuildFileTreeSkeleton(forceRebuild bool) {
	if !forceRebuild {
		// Try loading from cache
		if err := s.Tree.Load(); err == nil && s.Tree.CountNodes() > 0 {
			logger.Info("📂 Loaded cached file tree")
			return
		}
		logger.Info("⚠️ Cache not found or invalid, starting full build...")
	} else {
		logger.Info("♻️ Force rebuilding file tree...")
	}

	s.DriveInfo.WaitRateLimit()

	// Update FileTree ignore list (prevent stale config)
	s.ConfigManager.Lock.RLock()
	ignored := s.ConfigManager.Cfg.Google.IgnoredParents
	s.ConfigManager.Lock.RUnlock()
	s.Tree.SetIgnoredParents(ignored)

	// [Debug] List all drives for user filtering
	if drives, err := s.DriveInfo.ListAllDrives(); err == nil {
		logger.Info("📋 [Debug] Found %d drives:", len(drives))
		for _, d := range drives {
			logger.Info("   - ID: %s | Name: %s", d.Id, d.Name)
		}
	} else {
		logger.Warning("⚠️ Unable to get drives list: %v", err)
	}

	pt := ""
	count := 0
	for {
		q := s.DriveInfo.Srv.Files.List().
			Q("trashed = false").
			Fields("nextPageToken, files(id, name, parents, mimeType, driveId)").
			PageSize(1000).
			SupportsAllDrives(true).
			IncludeItemsFromAllDrives(true)

		if pt != "" {
			q.PageToken(pt)
		}

		r, err := q.Do()
		if err != nil {
			logger.Error("❌ Failed to list files: %v", err)
			return
		}

		for _, f := range r.Files {
			pid := ""
			if len(f.Parents) > 0 {
				pid = f.Parents[0]
			}
			isDir := f.MimeType == "application/vnd.google-apps.folder"
			s.Tree.UpdateNode(f.Id, f.Name, pid, isDir, f.DriveId)
			count++
		}
		logger.Verbose(model.LogLevelInfo, "   ↳ Loaded %d nodes...", count)

		if r.NextPageToken == "" {
			break
		}
		pt = r.NextPageToken
	}
	logger.Info("✅ Raw nodes loaded: %d, starting to prune ignored list...", s.Tree.CountNodes())

	// Prune ignored folders (recursive delete)
	if len(ignored) > 0 {
		prunedCount := 0
		for _, ignoreID := range ignored {
			// Get all descendants first
			descendants := s.Tree.GetDescendants(ignoreID)
			// Delete descendants
			for _, d := range descendants {
				s.Tree.RemoveNode(d.ID)
				prunedCount++
			}
			// Delete self
			if _, ok := s.Tree.GetPath(ignoreID); ok {
				s.Tree.RemoveNode(ignoreID)
				prunedCount++
			}
		}
		logger.Info("✂️ Pruned %d ignored nodes", prunedCount)
	}

	logger.Info("✅ File tree build complete, final node count: %d", s.Tree.CountNodes())

	// Save to cache
	if err := s.Tree.Save(); err != nil {
		logger.Error("❌ Failed to save file tree cache: %v", err)
	} else {
		logger.Info("💾 File tree cache saved")
	}
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

	// Sync ignore list to Tree for ResolvePathWithFallback
	s.ConfigManager.Lock.RLock()
	ignored := s.ConfigManager.Cfg.Google.IgnoredParents
	s.ConfigManager.Lock.RUnlock()
	s.Tree.SetIgnoredParents(ignored)
	logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "📋 [Diag] ignored_parents list: %v", ignored)

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

	// [Diag] Log change details
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
					notifs = append(notifs, Notif{d.Path, "delete", d.IsDir})
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

		// Check ignore list
		if s.isIgnoredParent(pid) {
			logger.Debug(s.ConfigManager.Cfg.Advanced.LogLevel, "🚫 [Diag] Skipping change: FileID=%s, Name=%s, DriveID=%s, ParentID=%s (ParentID in ignore list)", fileID, f.Name, f.DriveId, pid)
			continue
		}

		isDirBool := f.MimeType == "application/vnd.google-apps.folder"

		s.Tree.UpdateNode(fileID, f.Name, pid, isDirBool, f.DriveId)
		processedIDs[fileID] = true
		newPath := s.Tree.ResolvePathWithFallback(fileID)

		if !foundOld {
			logger.Info("🆕 [Create] %s", newPath)
			logger.WriteHistory(s.ConfigManager.Cfg, "CREATE", newPath)
			rcloneDirs[filepath.Dir(newPath)] = true
			notifs = append(notifs, Notif{newPath, "create", isDirBool})
		} else if oldPath != newPath {
			logger.Info("✏️ [Move] %s -> %s", oldPath, newPath)
			logger.WriteHistory(s.ConfigManager.Cfg, "MOVE", newPath)
			rcloneDirs[filepath.Dir(oldPath)] = true
			rcloneDirs[filepath.Dir(newPath)] = true

			// For moves, send old path delete and new path create
			notifs = append(notifs, Notif{oldPath, "delete", isDirBool})
			notifs = append(notifs, Notif{newPath, "create", isDirBool})

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
					notifs = append(notifs, Notif{oldChildPath, "delete", d.IsDir})
					notifs = append(notifs, Notif{d.Path, "create", d.IsDir})
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
			s.Symedia.SendWebhook(n.Path, n.Action, n.IsDir)
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
	if s.lastResetDate != today {
		todayTasks = 0 // New day, today's count is 0
	}
	
	return TaskStats{
		ActiveTasks:           s.activeTasks,
		TodayCompletedTasks:   todayTasks,
		HistoryCompletedTasks: s.historyCompletedTasks,
	}
}
