package server

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gd-webhook/src/config"
)

type Middleware struct {
	ConfigManager *config.Manager
	Sessions      *SessionStore
}

func NewMiddleware(cm *config.Manager) *Middleware {
	return &Middleware{
		ConfigManager: cm,
		Sessions:      NewSessionStore(),
	}
}

// SessionStore handles session storage
type SessionStore struct {
	sessions map[string]time.Time
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	s := &SessionStore{
		sessions: make(map[string]time.Time),
	}
	// Start goroutine to cleanup expired sessions
	go s.cleanup()
	return s
}

func (s *SessionStore) Create() string {
	b := make([]byte, 32)
	rand.Read(b)
	h := sha256.Sum256(b)
	token := hex.EncodeToString(h[:])

	s.mu.Lock()
	s.sessions[token] = time.Now().Add(24 * time.Hour) // 24 hour validity
	s.mu.Unlock()

	return token
}

func (s *SessionStore) Valid(token string) bool {
	if token == "" {
		return false
	}
	s.mu.RLock()
	expiry, exists := s.sessions[token]
	s.mu.RUnlock()

	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, expiry := range s.sessions {
			if now.After(expiry) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// HostCheckMiddleware checks Host header to restrict IP access
func (m *Middleware) HostCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.ConfigManager.Lock.RLock()
		restrict := m.ConfigManager.Cfg.Server.SSL.RestrictToDomain
		publicURL := m.ConfigManager.Cfg.Server.PublicURL
		m.ConfigManager.Lock.RUnlock()

		if restrict && publicURL != "" {
			u, err := url.Parse(publicURL)
			if err == nil {
				allowedHost := u.Hostname()
				reqHost, _, _ := net.SplitHostPort(r.Host)
				if reqHost == "" {
					reqHost = r.Host
				}

				if reqHost != allowedHost && reqHost != "localhost" && reqHost != "127.0.0.1" {
					http.Error(w, "403 Forbidden: IP access restricted, please use configured domain.", http.StatusForbidden)
					return
				}
			}
		}
		next(w, r)
	}
}

// AuthMiddleware handles authentication (supports Cookie Session and Basic Auth)
func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.ConfigManager.Lock.RLock()
		user := m.ConfigManager.Cfg.Auth.Username
		pass := m.ConfigManager.Cfg.Auth.Password
		webhookPath := m.ConfigManager.Cfg.Server.WebhookPath
		m.ConfigManager.Lock.RUnlock()

		// Skip auth for webhook path and Bing wallpaper
		if r.URL.Path == webhookPath || r.URL.Path == "/api/bing/wallpaper" {
			next(w, r)
			return
		}

		// Skip auth for login endpoint
		if r.URL.Path == "/api/auth/login" {
			next(w, r)
			return
		}

		// Skip auth for static resources (allow frontend pages and assets)
		// If not an API path, treat as static resource
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next(w, r)
			return
		}

		// Skip if no password configured
		if user == "" || pass == "" {
			next(w, r)
			return
		}

		// 1. Check Cookie Session
		if cookie, err := r.Cookie("gd_session"); err == nil {
			if m.Sessions.Valid(cookie.Value) {
				next(w, r)
				return
			}
		}

		// 2. Check Basic Auth (retained for backward compatibility)
		reqUser, reqPass, ok := r.BasicAuth()
		if ok && subtle.ConstantTimeCompare([]byte(reqUser), []byte(user)) == 1 && subtle.ConstantTimeCompare([]byte(reqPass), []byte(pass)) == 1 {
			next(w, r)
			return
		}

		// Unauthorized: return 401 (no longer triggers Basic Auth dialog)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
