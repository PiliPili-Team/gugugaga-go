package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/web"
)

type Server struct {
	ConfigManager *config.Manager
	Handler       *Handler
	Middleware    *Middleware
	StaticFS      http.Handler
}

func NewServer(cm *config.Manager, h *Handler, m *Middleware) *Server {
	logger.Info("🔍 Checking embedded assets...")
	entries, _ := web.StaticFS.ReadDir("static")
	for _, e := range entries {
		logger.Info("   - Found in embed: %s (IsDir: %v)", e.Name(), e.IsDir())
	}

	subFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		logger.Error("❌ Failed to create sub-filesystem for static assets: %v", err)
		subFS = web.StaticFS
	} else {
		logger.Info("✅ Static sub-filesystem created successfully")
	}

	// Set Handler's Middleware reference
	h.SetMiddleware(m)

	return &Server{
		ConfigManager: cm,
		Handler:       h,
		Middleware:    m,
		StaticFS:      http.FileServer(http.FS(subFS)),
	}
}

func (s *Server) Start() {
	// Setup router Mux
	mux := http.NewServeMux()

	// Auth API routes (login/logout)
	mux.HandleFunc("/api/auth/login", s.Handler.HandleLogin)
	mux.HandleFunc("/api/auth/logout", s.Handler.HandleLogout)

	// Wallpaper API
	mux.HandleFunc("/api/bing/wallpaper", s.Handler.HandleBingWallpaper)
	mux.HandleFunc("/api/tmdb/wallpaper", s.Handler.HandleTMDBWallpaper)
	
	// System status API
	mux.HandleFunc("/api/status", s.Handler.HandleSystemStatus)

	// API routes
	mux.HandleFunc("/api/logs", s.Handler.HandleLogs)
	mux.HandleFunc("/api/logs/clear_mem", s.Handler.HandleClearMem)
	mux.HandleFunc("/api/logs/clear_files", s.Handler.HandleClearFiles)
	mux.HandleFunc("/api/auth/login_url", s.Handler.HandleOAuthLoginURL)
	mux.HandleFunc("/api/auth/callback", s.Handler.HandleOAuthCallback)
	mux.HandleFunc("/api/config/get", s.Handler.HandleConfigGet)
	mux.HandleFunc("/api/config/update", s.Handler.HandleConfigUpdate)
	mux.HandleFunc("/api/trigger", s.Handler.HandleTrigger)
	mux.HandleFunc("/api/rclone_full", s.Handler.HandleRcloneFull)
	mux.HandleFunc("/api/test_symedia", s.Handler.HandleTestSymedia)
	mux.HandleFunc("/api/tree/refresh", s.Handler.HandleTreeRefresh)

	// Webhook route
	mux.HandleFunc(s.ConfigManager.Cfg.Server.WebhookPath, s.Handler.HandleWebhook)

	// Static assets
	mux.Handle("/", s.StaticFS)

	// Middleware chain: HostCheck -> Auth -> Mux
	// Note: AuthMiddleware internally skips webhook path and login auth
	handler := s.Middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
	handler = s.Middleware.HostCheckMiddleware(handler)

	cfg := s.ConfigManager.GetConfig()
	port := cfg.Server.ListenPort
	ssl := cfg.Server.SSL
	addr := fmt.Sprintf(":%d", port)

	if ssl.Enabled {
		if ssl.CertPath == "" || ssl.KeyPath == "" {
			logger.Error("❌ SSL enabled but missing cert/key path, falling back to HTTP")
			logger.Info("📡 Dashboard (HTTP): http://localhost%s", addr)
			if err := http.ListenAndServe(addr, handler); err != nil {
				logger.Error("Failed to start: %v", err)
			}
		} else {
			logger.Info("🔐 SSL enabled")
			logger.Info("📡 Dashboard (HTTPS): https://localhost%s", addr)
			if err := http.ListenAndServeTLS(addr, ssl.CertPath, ssl.KeyPath, handler); err != nil {
				logger.Error("❌ SSL startup failed: %v", err)
				logger.Warning("⚠️ Attempting fallback to HTTP mode...")
				logger.Info("📡 Dashboard (HTTP): http://localhost%s", addr)
				if err := http.ListenAndServe(addr, handler); err != nil {
					logger.Error("Failed to start: %v", err)
				}
			}
		}
	} else {
		logger.Info("📡 Dashboard (HTTP): http://localhost%s", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			logger.Error("Failed to start: %v", err)
		}
	}
}
