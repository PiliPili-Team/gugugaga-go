package server

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"
	"gd-webhook/src/service"
)

// Server start time
var serverStartTime = time.Now()

// Handler contains all HTTP handler functions
type Handler struct {
	ConfigManager *config.Manager
	DriveInfo     *service.DriveService
	Sync          *service.SyncService
	Rclone        *service.RcloneService
	Symedia       *service.SymediaService
	Middleware    *Middleware
	TotalMemory   uint64
}

// NewHandler creates a Handler instance
func NewHandler(
	cm *config.Manager,
	ds *service.DriveService,
	ss *service.SyncService,
	rc *service.RcloneService,
	sy *service.SymediaService,
) *Handler {
	return &Handler{
		ConfigManager: cm,
		DriveInfo:     ds,
		Sync:          ss,
		Rclone:        rc,
		Symedia:       sy,
		TotalMemory:   getTotalMemory(),
	}
}

func getTotalMemory() uint64 {
	// Mac implementation
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err == nil {
		s := strings.TrimSpace(string(out))
		if v, err := strconv.ParseUint(s, 10, 64); err == nil {
			return v
		}
	}
	return 0
}

// SetMiddleware sets middleware reference (for session access)
func (h *Handler) SetMiddleware(m *Middleware) {
	h.Middleware = m
}

// HandleLogin handles login requests
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	h.ConfigManager.Lock.RLock()
	user := h.ConfigManager.Cfg.Auth.Username
	pass := h.ConfigManager.Cfg.Auth.Password
	h.ConfigManager.Lock.RUnlock()

	// Verify username and password
	if subtle.ConstantTimeCompare([]byte(req.Username), []byte(user)) != 1 ||
		subtle.ConstantTimeCompare([]byte(req.Password), []byte(pass)) != 1 {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session
	if h.Middleware != nil {
		token := h.Middleware.Sessions.Create()
		http.SetCookie(w, &http.Cookie{
			Name:     "gd_session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400, // 24 hours
			SameSite: http.SameSiteLaxMode,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HandleLogout handles logout requests
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("gd_session"); err == nil && h.Middleware != nil {
		h.Middleware.Sessions.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "gd_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HandleBingWallpaper fetches Bing daily wallpaper
func (h *Handler) HandleBingWallpaper(w http.ResponseWriter, r *http.Request) {
	// Bing wallpaper API
	bingAPI := "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(bingAPI)
	if err != nil {
		http.Error(w, "Failed to fetch Bing wallpaper", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var bingResp struct {
		Images []struct {
			URL       string `json:"url"`
			Copyright string `json:"copyright"`
			Title     string `json:"title"`
		} `json:"images"`
	}

	if err := json.Unmarshal(body, &bingResp); err != nil || len(bingResp.Images) == 0 {
		http.Error(w, "Failed to parse Bing response", http.StatusInternalServerError)
		return
	}

	img := bingResp.Images[0]
	// Use HD version (1920x1080)
	hdURL := img.URL
	// Replace resolution for higher quality
	if len(hdURL) > 0 {
		hdURL = "https://www.bing.com" + hdURL + "&w=1920&h=1080&rs=1&c=4"
	}

	result := map[string]string{
		"url":       hdURL,
		"copyright": img.Copyright,
		"title":     img.Title,
		"source":    "bing",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache") // No cache, fetch fresh wallpaper
	json.NewEncoder(w).Encode(result)
}

// HandleTMDBWallpaper fetches TMDB trending movie/TV backdrop
func (h *Handler) HandleTMDBWallpaper(w http.ResponseWriter, r *http.Request) {
	// TMDB API - Get trending
	tmdbAPI := "https://api.themoviedb.org/3/trending/all/day?language=zh-CN"

	// Use free API key (TMDB official example)
	apiKey := "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI0MzE5YTBjYjlhMGJjZGY5ZjZhZjA2ZWY2OWVmMWU5YiIsIm5iZiI6MTczOTU0MjQyMC44MzUsInN1YiI6IjY3YWVjMzc0YjJhNDVlNmJmNzExNTFkYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.Y-_mEy0-5jqwb3D9VJ7FqvuA4G6o_eGHr6fmvYEPcqQ"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", tmdbAPI, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch TMDB data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var tmdbResp struct {
		Results []struct {
			BackdropPath string  `json:"backdrop_path"`
			Title        string  `json:"title"`
			Name         string  `json:"name"`
			Overview     string  `json:"overview"`
			VoteAverage  float64 `json:"vote_average"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &tmdbResp); err != nil || len(tmdbResp.Results) == 0 {
		http.Error(w, "Failed to parse TMDB response", http.StatusInternalServerError)
		return
	}

	// Randomly select a result with backdrop
	var validResults []struct {
		BackdropPath string
		Title        string
		Name         string
		Overview     string
		VoteAverage  float64
	}

	for _, r := range tmdbResp.Results {
		if r.BackdropPath != "" {
			validResults = append(validResults, struct {
				BackdropPath string
				Title        string
				Name         string
				Overview     string
				VoteAverage  float64
			}{r.BackdropPath, r.Title, r.Name, r.Overview, r.VoteAverage})
		}
	}

	if len(validResults) == 0 {
		http.Error(w, "No valid backdrop found", http.StatusNotFound)
		return
	}

	// Random selection
	rand.Seed(time.Now().UnixNano())
	selected := validResults[rand.Intn(len(validResults))]

	// Build HD image URL (original size)
	title := selected.Title
	if title == "" {
		title = selected.Name
	}

	result := map[string]interface{}{
		"url":       "https://image.tmdb.org/t/p/original" + selected.BackdropPath,
		"copyright": title,
		"title":     title,
		"source":    "tmdb",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(result)
}

// HandleSystemStatus returns system status
func (h *Handler) HandleSystemStatus(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(serverStartTime)

	// Format uptime
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	var uptimeStr string
	if days > 0 {
		uptimeStr = fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		uptimeStr = fmt.Sprintf("%dh %dm", hours, minutes)
	} else if minutes > 0 {
		uptimeStr = fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		uptimeStr = fmt.Sprintf("%ds", seconds)
	}

	// Get task statistics
	var taskStats service.TaskStats
	if h.Sync != nil {
		taskStats = h.Sync.GetTaskStats()
	}

	// Get memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate memory usage percentage (Process Sys / Machine Total)
	memUsage := 0.0
	if h.TotalMemory > 0 {
		memUsage = float64(memStats.Sys) / float64(h.TotalMemory) * 100
	} else if memStats.Sys > 0 {
		// Fallback: Alloc / Sys
		memUsage = float64(memStats.Alloc) / float64(memStats.Sys) * 100
	}

	// Get goroutine count as a proxy for CPU activity
	numGoroutines := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// Estimate CPU usage based on goroutine activity (simplified)
	// This is not real CPU usage, just an estimation
	cpuUsage := float64(numGoroutines) / float64(numCPU*100) * 100
	if cpuUsage > 100 {
		cpuUsage = 100
	}

	result := map[string]interface{}{
		"status":                  "online",
		"uptime_seconds":          int64(uptime.Seconds()),
		"uptime_display":          uptimeStr,
		"start_time":              serverStartTime.Format(time.RFC3339),
		"app_name":                config.GetAppName(),
		"app_version":             config.GetAppVersion(),
		"today_completed_tasks":   taskStats.TodayCompletedTasks,
		"history_completed_tasks": taskStats.HistoryCompletedTasks,
		"cpu_usage":               cpuUsage,
		"memory_usage":            memUsage,
		"memory_alloc_mb":         float64(memStats.Alloc) / 1024 / 1024,
		"memory_sys_mb":           float64(memStats.Sys) / 1024 / 1024,
		"goroutines":              numGoroutines,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleLogs returns real-time memory logs
func (h *Handler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	sinceStr := r.URL.Query().Get("since")
	sinceIdx, _ := strconv.Atoi(sinceStr)

	logs, total := logger.GetMemLogs(sinceIdx)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.LogsResponse{Logs: logs, NextIdx: total})
}

// HandleClearMem clears memory logs
func (h *Handler) HandleClearMem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	logger.ClearMemLogs()
	_, _ = w.Write([]byte("ok"))
}

// HandleClearFiles clears disk log files
func (h *Handler) HandleClearFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	h.ConfigManager.Lock.RLock()
	baseDir := h.ConfigManager.Cfg.Advanced.LogDir
	h.ConfigManager.Lock.RUnlock()

	if baseDir == "" {
		baseDir = "userdata/logs"
	}
	sysDir := filepath.Join(baseDir, "system")
	histDir := filepath.Join(baseDir, "history")

	// Simple file deletion, preserve current log
	cleanDir(sysDir, true)
	cleanDir(histDir, false)

	logger.Info("🧹 Server history logs cleared")
	_, _ = w.Write([]byte("ok"))
}

func cleanDir(dir string, protectCurrentLog bool) {
	entries, _ := os.ReadDir(dir)
	currentLog := logger.GetLogFileHandle()
	var currentLogName string
	if currentLog != nil {
		currentLogName = filepath.Base(currentLog.Name())
	}

	for _, e := range entries {
		if !e.IsDir() {
			if protectCurrentLog && e.Name() == currentLogName {
				continue
			}
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}

// HandleConfigGet returns current configuration
func (h *Handler) HandleConfigGet(w http.ResponseWriter, r *http.Request) {
	h.ConfigManager.Lock.RLock()
	cfg := *h.ConfigManager.Cfg
	h.ConfigManager.Lock.RUnlock()

	// Return full config (including password) since protected by auth middleware
	_ = json.NewEncoder(w).Encode(cfg)
}

// HandleConfigUpdate updates configuration
func (h *Handler) HandleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var newCfg model.Config
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, "JSON format error", 400)
		return
	}

	h.ConfigManager.Lock.RLock()
	oldCfg := *h.ConfigManager.Cfg
	h.ConfigManager.Lock.RUnlock()

	if newCfg.Auth.Password == "" {
		newCfg.Auth.Password = oldCfg.Auth.Password
	}

	logChanged := (newCfg.Advanced.LogDir != oldCfg.Advanced.LogDir) ||
		(newCfg.Advanced.LogSaveEnabled != oldCfg.Advanced.LogSaveEnabled)

	addrChanged := (newCfg.Server.PublicURL != oldCfg.Server.PublicURL) ||
		(newCfg.Server.WebhookPath != oldCfg.Server.WebhookPath)

	oauth := newCfg.OAuthConfig
	if oauth.ClientID != "" && oauth.ClientSecret != "" && oauth.RedirectURI != "" {
		h.ConfigManager.SaveCredentialsFile(oauth.ClientID, oauth.ClientSecret, oauth.RedirectURI)
		// Async reload OAuth config
		go func() { _ = h.DriveInfo.InitOAuthConfig() }()
	}

	h.ConfigManager.UpdateConfig(newCfg)
	_ = h.ConfigManager.SaveConfig()

	if logChanged {
		logger.InitLogging(&newCfg)
	}

	if addrChanged {
		logger.Info("🔄 Address change detected, re-registering webhook...")
		go func() {
			time.Sleep(1 * time.Second)
			if h.DriveInfo.Srv != nil {
				token := h.DriveInfo.GetStartPageToken()
				h.DriveInfo.RegisterWatch(token)
				logger.Info("✅ Re-registration complete")
			}
		}()
	}

	logger.Info("💾 Config updated (restart if SSL/port changed)")
	_, _ = w.Write([]byte("ok"))
}

// HandleTrigger manually triggers sync
func (h *Handler) HandleTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		logger.Verbose(model.LogLevelInfo, "👋 Manual sync triggered")
		select {
		case h.Sync.TriggerChan <- struct{}{}:
		default:
		}
		_, _ = w.Write([]byte("ok"))
	}
}

// HandleRcloneFull forces full Rclone refresh
func (h *Handler) HandleRcloneFull(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		logger.Info("👋 Force full Rclone refresh")
		go h.Rclone.Refresh("/")
		_, _ = w.Write([]byte("ok"))
	}
}

// HandleTestSymedia sends test webhook
func (h *Handler) HandleTestSymedia(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var p model.TestSymediaRequest
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			return
		}
		go h.Symedia.SendWebhook(p.Path, "create", false, "")
		_, _ = w.Write([]byte("ok"))
	}
}

// HandleWebhook handles Google Drive webhook callback
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	state := r.Header.Get("X-Goog-Resource-State")

	h.ConfigManager.Lock.RLock()
	logLevel := h.ConfigManager.Cfg.Advanced.LogLevel
	h.ConfigManager.Lock.RUnlock()

	if logLevel >= model.LogLevelDebug {
		logger.Debug(logLevel, "🔔 Received Google Webhook (%s):", r.URL.Path)
		logger.Debug(logLevel, "   - Resource State: %s", state)
		logger.Debug(logLevel, "   - Resource ID: %s", r.Header.Get("X-Goog-Resource-Id"))
		logger.Debug(logLevel, "   - Channel ID: %s", r.Header.Get("X-Goog-Channel-Id"))
	}

	if state != "" {
		select {
		case h.Sync.TriggerChan <- struct{}{}:
		default:
		}
	}
	w.WriteHeader(http.StatusOK)
}

// HandleOAuthLoginURL returns OAuth login URL
func (h *Handler) HandleOAuthLoginURL(w http.ResponseWriter, r *http.Request) {
	logger.Info("🔗 [OAuth] GetLoginURL request received")

	if h.DriveInfo.OAuthConfig == nil {
		logger.Info("🔗 [OAuth] OAuthConfig is nil, attempting to load...")
		err := h.DriveInfo.InitOAuthConfig()
		if err != nil {
			logger.Error("🔗 [OAuth] InitOAuthConfig failed: %v", err)
		}
		if h.DriveInfo.OAuthConfig == nil {
			logger.Error("🔗 [OAuth] OAuthConfig still nil after init attempt")
			http.Error(w, "OAuth config initialization failed. Please check credentials.", http.StatusBadRequest)
			return
		}
	}

	logger.Info("🔗 [OAuth] Generating auth URL...")
	authLoginUrl := h.DriveInfo.OAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	logger.Info("🔗 [OAuth] Auth URL generated successfully (length: %d)", len(authLoginUrl))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.OAuthLoginURLResponse{URL: authLoginUrl})
}

// HandleOAuthCallback handles OAuth callback
func (h *Handler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", 400)
		return
	}
	if h.DriveInfo.OAuthConfig == nil {
		_ = h.DriveInfo.InitOAuthConfig()
	}

	token, err := h.DriveInfo.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Exchange failed: "+err.Error(), 500)
		logger.Error("❌ OAuth Exchange failed: %v", err)
		return
	}

	h.DriveInfo.SaveToken(model.TokenFile, token)
	logger.Info("🎉 Authorization successful! Token saved.")

	// Re-initialize service with new token
	_ = h.DriveInfo.InitDriveService()

	go func() {
		if h.DriveInfo.Srv != nil {
			logger.Verbose(model.LogLevelInfo, "⏳ Initializing file tree...")
			h.Sync.BuildFileTreeSkeleton(true)
			h.DriveInfo.EnsureStartPageToken()
			h.DriveInfo.RegisterWatch(h.DriveInfo.GetStartPageToken())
			logger.Info("✅ System ready")
		}
	}()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// HandleTreeRefresh forces file tree refresh
func (h *Handler) HandleTreeRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.Sync.ForceRebuild()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "Tree refresh started"}`))
}
