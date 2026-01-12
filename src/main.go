package main

import (
	"time"

	"github.com/robfig/cron/v3"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/server"
	"gd-webhook/src/service"
)

func main() {
	// 1. Initialize configuration
	cfgManager := config.NewManager()
	cfgManager.LoadConfig()

	// 2. Initialize logging
	cfg := cfgManager.GetConfig()
	logger.InitLogging(&cfg)
	
	// Get app name and version
	appName := config.GetAppName()
	appVersion := config.GetAppVersion()
	logger.Info("🚀 %s v%s starting...", appName, appVersion)

	// 3. Initialize services
	driveService := service.NewDriveService(cfgManager)
	fileTree := service.NewFileTree(driveService)
	rcloneService := service.NewRcloneService(cfgManager)
	symediaService := service.NewSymediaService(cfgManager)
	syncService := service.NewSyncService(cfgManager, driveService, fileTree, rcloneService, symediaService)

	// 4. Initialize Cron (log cleanup and rotation)
	cronRunner := cron.New(cron.WithSeconds())
	cronRunner.Start()

	// Setup log cleanup cron
	if cfg.Advanced.LogCleanupCron == "" {
		cfg.Advanced.LogCleanupCron = "0 0 3 * * ?"
	}
	_, err := cronRunner.AddFunc(cfg.Advanced.LogCleanupCron, func() {
		// Re-fetch config to apply latest retention policy
		currentCfg := cfgManager.GetConfig()
		logger.CleanupLogs(&currentCfg)
	})
	if err != nil {
		logger.Error("❌ Cron format error: %v", err)
	} else {
		logger.Info("⏰ Log cleanup cron scheduled: [%s]", cfg.Advanced.LogCleanupCron)
	}

	// 5. Initialize HTTP server (middleware & handler)
	middleware := server.NewMiddleware(cfgManager)
	handler := server.NewHandler(cfgManager, driveService, syncService, rcloneService, symediaService)
	srv := server.NewServer(cfgManager, handler, middleware)

	// 6. Start background tasks
	// Log rotation check (every minute)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			c := cfgManager.GetConfig()
			logger.CheckLogRotation(&c)
		}
	}()

	// Drive service initialization and sync loop startup
	go func() {
		if err := driveService.InitDriveService(); err != nil {
			logger.Warning("⚠️ Drive service not initialized, please login via WebUI")
		} else {
			logger.Verbose(1, "⏳ Preloading file tree...")
			syncService.BuildFileTreeSkeleton(false)
			logger.Info("✅ Preload complete (nodes: %d)", fileTree.CountNodes())

			driveService.EnsureStartPageToken()
			token := driveService.GetStartPageToken()
			driveService.RegisterWatch(token)
		}
		// Start web server
		srv.Start()
	}()

	// Webhook auto-renewal (every 6 days to avoid 7-day expiry)
	go func() {
		ticker := time.NewTicker(6 * 24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if driveService.Srv != nil {
				token := driveService.GetStartPageToken()
				logger.Info("🔄 [Scheduled] Renewing webhook...")
				driveService.RegisterWatch(token)
			}
		}
	}()

	// Start sync processing loop
	go syncService.StartProcessLoop()

	// 7. Block main thread
	select {}
}
