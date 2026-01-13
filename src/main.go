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
	cfgManager := config.NewManager()
	cfgManager.LoadConfig()

	cfg := cfgManager.GetConfig()
	logger.InitLogging(&cfg)

	appName := config.GetAppName()
	appVersion := config.GetAppVersion()
	logger.Info("🚀 %s v%s starting...", appName, appVersion)

	driveService := service.NewDriveService(cfgManager)
	fileTree := service.NewFileTree(driveService)
	rcloneService := service.NewRcloneService(cfgManager)
	symediaService := service.NewSymediaService(cfgManager)
	syncService := service.NewSyncService(cfgManager, driveService, fileTree, rcloneService, symediaService)

	cronRunner := cron.New(cron.WithSeconds())
	cronRunner.Start()

	if cfg.Advanced.LogCleanupCron == "" {
		cfg.Advanced.LogCleanupCron = "0 0 3 * * ?"
	}
	_, err := cronRunner.AddFunc(cfg.Advanced.LogCleanupCron, func() {
		currentCfg := cfgManager.GetConfig()
		logger.CleanupLogs(&currentCfg)
	})
	if err != nil {
		logger.Error("❌ Cron format error: %v", err)
	} else {
		logger.Info("⏰ Log cleanup cron scheduled: [%s]", cfg.Advanced.LogCleanupCron)
	}

	middleware := server.NewMiddleware(cfgManager)
	handler := server.NewHandler(cfgManager, driveService, syncService, rcloneService, symediaService)
	srv := server.NewServer(cfgManager, handler, middleware)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			c := cfgManager.GetConfig()
			logger.CheckLogRotation(&c)
		}
	}()

	go func() {
		initErr := driveService.InitDriveService()
		if initErr != nil {
			logger.Warning("⚠️ Drive service not initialized, please login via WebUI")
		}

		if len(cfgManager.GetConfig().Google.TargetDriveIDs) == 0 {
			drives, err := driveService.ListAllDrives()
			if err != nil && driveService.Srv != nil {
				logger.Warning("⚠️ Failed to list drives: %v", err)
			} else if len(drives) > 0 {
				logger.Info("📋 [Drive List] Found %d available drives:", len(drives))
				for _, d := range drives {
					logger.Info("  - Name: %-20s | ID: %s", d.Name, d.Id)
				}
			}
		}

		if initErr == nil {
			logger.Verbose(1, "⏳ Loading file tree cache...")
			if err := fileTree.Load(); err == nil {
				logger.Info("📂 Cache loaded (nodes: %d)", fileTree.CountNodes())
			}

			go syncService.BuildFileTreeSkeleton(false)

			driveService.EnsureStartPageToken()
			token := driveService.GetStartPageToken()
			driveService.RegisterWatch(token)
		}
		srv.Start()
	}()

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

	go syncService.StartProcessLoop()

	select {}
}
