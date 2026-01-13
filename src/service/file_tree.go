package service

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"

	"gd-webhook/src/logger"
	"gd-webhook/src/model"
)

// FileTree is the in-memory file tree structure
type FileTree struct {
	sync.RWMutex
	nodes    map[string]*model.FileNode            // ID -> Node
	children map[string]map[string]*model.FileNode // ParentID -> ChildID -> Node
	ds       *DriveService
}

// NewFileTree creates a new file tree
func NewFileTree(ds *DriveService) *FileTree {
	return &FileTree{
		nodes:    make(map[string]*model.FileNode),
		children: make(map[string]map[string]*model.FileNode),
		ds:       ds,
	}
}

// SetTargetDrives updates the target drive list
func (t *FileTree) SetTargetDrives(ids []string) {
	t.Lock()
	defer t.Unlock()
}

// ReplaceWith atomically replaces the current tree data with another tree's data
func (t *FileTree) ReplaceWith(other *FileTree) {
	t.Lock()
	other.RLock()
	defer t.Unlock()
	defer other.RUnlock()

	// Deep copy nodes
	t.nodes = make(map[string]*model.FileNode, len(other.nodes))
	for k, v := range other.nodes {
		t.nodes[k] = v
	}

	// Deep copy children
	t.children = make(map[string]map[string]*model.FileNode, len(other.children))
	for pID, kids := range other.children {
		newKids := make(map[string]*model.FileNode, len(kids))
		for k, v := range kids {
			newKids[k] = v
		}
		t.children[pID] = newKids
	}

	logger.Info("🌳 Tree replaced atomically. Nodes: %d", len(t.nodes))
}

// GetNode returns a copy of the node for the given ID
func (t *FileTree) GetNode(id string) (*model.FileNode, bool) {
	t.RLock()
	defer t.RUnlock()
	if n, ok := t.nodes[id]; ok {
		// Return copy to prevent external mutation
		node := *n
		return &node, true
	}
	return nil, false
}

// Save saves the file tree to disk (NDJSON format)
func (t *FileTree) Save() error {
	t.RLock()
	defer t.RUnlock()
	return t.saveInternal()
}

func (t *FileTree) saveInternal() error {
	dir := filepath.Dir(model.TreeCacheFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	tmpFile := model.TreeCacheFile + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)

	for _, node := range t.nodes {
		if err := enc.Encode(node); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	f.Close()

	return os.Rename(tmpFile, model.TreeCacheFile)
}

// Load loads the file tree from disk (Supports NDJSON and Legacy JSON)
func (t *FileTree) Load() error {
	t.Lock()
	defer t.Unlock()

	if err := t.loadStreaming(); err == nil {
		return nil
	} else {
		logger.Warning("⚠️ Streaming load failed (%v), trying legacy format...", err)
	}

	if err := t.loadLegacy(); err != nil {
		return err
	}

	logger.Info("🔄 Migrating legacy cache to new format...")
	if err := t.saveInternal(); err != nil {
		logger.Error("❌ Failed to migrate cache: %v", err)
	} else {
		logger.Info("✅ Cache migrated successfully!")
	}

	// Force GC to release buffer memory immediately
	runtime.GC()
	debug.FreeOSMemory()

	return nil
}

func (t *FileTree) loadStreaming() error {
	f, err := os.Open(model.TreeCacheFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use json.Decoder for stream parsing (better handling of large tokens vs Scanner)
	decoder := json.NewDecoder(bufio.NewReader(f))

	// Pre-clear nodes
	t.nodes = make(map[string]*model.FileNode)
	count := 0

	for decoder.More() {
		var node model.FileNode
		if err := decoder.Decode(&node); err != nil {
			// If EOF or other error, break. But More() should handle logical EOF.
			// If corruption, we might log and continue or failing.
			// Decoder might stop on error.
			if err == io.EOF {
				break
			}
			logger.Warning("⚠️ Corrupt JSON token in cache: %v", err)
			// Try to recover? Decoder might be stuck.
			// Usually safe to abort streaming load and fallback if too many errors.
			// For now, return error to trigger fallback?
			// Or simple continue? JSON stream decoder state might be invalid.
			return err
		}
		t.nodes[node.ID] = &node
		count++
	}

	if count == 0 {
		return os.ErrNotExist // Treat empty as not found/invalid
	}

	t.rebuildChildren()
	logger.Info("📂 Loaded %d nodes from cache stream", count)
	return nil
}

func (t *FileTree) loadLegacy() error {
	data, err := os.ReadFile(model.TreeCacheFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &t.nodes); err != nil {
		return err
	}

	t.rebuildChildren()
	logger.Info("📂 Loaded %d nodes from legacy cache", len(t.nodes))
	return nil
}

// rebuildChildren rebuilds the children index
func (t *FileTree) rebuildChildren() {
	t.children = make(map[string]map[string]*model.FileNode)
	for _, node := range t.nodes {
		if node.ParentID != "" {
			if t.children[node.ParentID] == nil {
				t.children[node.ParentID] = make(map[string]*model.FileNode)
			}
			t.children[node.ParentID][node.ID] = node
		}
	}
	logger.Info("✅ Rebuilt index from cache, %d nodes", len(t.nodes))
}

// UpdateNode updates or adds a node
func (t *FileTree) UpdateNode(id, name, parentID string, isDir bool, driveID string) {
	t.Lock()
	defer t.Unlock()

	// If node exists and parent changed, remove from old parent's children list
	if oldNode, exists := t.nodes[id]; exists {
		if oldNode.ParentID != parentID {
			if kids, ok := t.children[oldNode.ParentID]; ok {
				delete(kids, id)
				if len(kids) == 0 {
					delete(t.children, oldNode.ParentID)
				}
			}
		}
	}

	node := &model.FileNode{ID: id, Name: name, ParentID: parentID, IsDir: isDir, DriveID: driveID}
	t.nodes[id] = node

	if parentID != "" {
		if t.children[parentID] == nil {
			t.children[parentID] = make(map[string]*model.FileNode)
		}
		t.children[parentID][id] = node
	}
}

// RemoveNode removes a node
func (t *FileTree) RemoveNode(id string) {
	t.Lock()
	defer t.Unlock()

	node, exists := t.nodes[id]
	if !exists {
		return
	}

	if node.ParentID != "" {
		if kids, ok := t.children[node.ParentID]; ok {
			delete(kids, id)
			if len(kids) == 0 {
				delete(t.children, node.ParentID)
			}
		}
	}

	delete(t.children, id)
	delete(t.nodes, id)
}

// GetPath gets the full path of a node
func (t *FileTree) GetPath(id string) (string, bool) {
	t.RLock()
	defer t.RUnlock()
	return t.getPathLocked(id)
}

// getPathLocked internal recursive path resolution
func (t *FileTree) getPathLocked(id string) (string, bool) {
	node, ok := t.nodes[id]
	if !ok {
		return "", false
	}

	if node.ParentID == "" {
		driveName := t.ds.GetDriveName(node.DriveID)

		// [Fix] If this is a shared drive's root node (ID == DriveID)
		// Return /DriveName directly to avoid duplication (e.g. /DriveName/DriveName)
		if node.DriveID != "" && node.ID == node.DriveID {
			return "/" + driveName, true
		}

		// [Fix] If this is My Drive's root node (ID == "root")
		// Return /DriveName directly
		if node.ID == "root" {
			return "/" + driveName, true
		}

		// Other cases (e.g. orphan files in My Drive, or weird structure)
		return "/" + driveName + "/" + node.Name, true
	}

	parentPath, parentOk := t.getPathLocked(node.ParentID)
	if !parentOk {
		// If parent node not found, tree is incomplete (parent folder not synced or loaded)
		// Return false, let upper layer ResolvePathWithFallback fetch it
		return "", false
	}

	return filepath.Join(parentPath, node.Name), true
}

// GetDescendants gets all descendant nodes
func (t *FileTree) GetDescendants(rootID string) []model.DescendantInfo {
	t.RLock()
	defer t.RUnlock()

	var results []model.DescendantInfo
	var scan func(currentID string)

	scan = func(currentID string) {
		if node, ok := t.nodes[currentID]; ok {
			if p, ok := t.getPathLocked(currentID); ok {
				results = append(results, model.DescendantInfo{
					ID:      currentID,
					Path:    p,
					IsDir:   node.IsDir,
					DriveID: node.DriveID,
				})
			}
			if kids, ok := t.children[currentID]; ok {
				for kidID := range kids {
					scan(kidID)
				}
			}
		}
	}
	scan(rootID)
	return results
}

// ResolvePathWithFallback attempts to get path, falls back to API if not found
func (t *FileTree) ResolvePathWithFallback(id string) string {
	p, ok := t.GetPath(id)
	if ok {
		return p
	}

	t.ds.WaitRateLimit()
	if t.ds.Srv == nil {
		return "/UNKNOWN/" + id
	}

	f, err := t.ds.Srv.Files.Get(id).Fields("id,name,parents,mimeType,driveId").SupportsAllDrives(true).Do()
	if err != nil {
		logger.Warning("⚠️ [Fallback] API query failed (ID: %s): %v", id, err)
		return "/UNKNOWN_API_ERROR/" + id
	}

	pid := ""
	if len(f.Parents) > 0 {
		pid = f.Parents[0]
	}

	t.UpdateNode(id, f.Name, pid, f.MimeType == "application/vnd.google-apps.folder", f.DriveId)
	logger.Verbose(model.LogLevelDebug, "🧩 [Fallback] Added node: %s (Parent: %s)", f.Name, pid)

	if pid != "" {
		if _, ok := t.GetPath(pid); !ok {
			logger.Verbose(model.LogLevelDebug, "   ↳ Recursively fetching parent: %s", pid)
			t.ResolvePathWithFallback(pid)
		}
	}

	p, _ = t.GetPath(id)
	if p == "" {
		// [Fix] If path is still empty, parent folder is missing, return error path instead of root
		// Try to get parent folder name (if cached)
		parentName := "UNKNOWN_PARENT"
		if pid != "" {
			t.RLock()
			if pn, ok := t.nodes[pid]; ok {
				parentName = pn.Name
			}
			t.RUnlock()
		}
		logger.Warning("⚠️ [Fallback] Path resolution failed (ID: %s, Parent: %s), returning error path", id, pid)
		return "/UNRESOLVED_PATH/" + parentName + "/" + f.Name
	}
	return p
}

// CountNodes returns total node count
func (t *FileTree) CountNodes() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.nodes)
}
