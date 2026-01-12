package service

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"
)

// DriveService wraps Google Drive API operations
type DriveService struct {
	Srv            *drive.Service
	Limiter        *rate.Limiter
	DriveNameCache sync.Map
	DriveCacheLoad sync.Map
	OAuthConfig    *oauth2.Config
	ConfigManager  *config.Manager
}

// NewDriveService creates a new DriveService
func NewDriveService(cm *config.Manager) *DriveService {
	cfg := cm.GetConfig()
	return &DriveService{
		Limiter:       rate.NewLimiter(rate.Limit(cfg.Google.RateLimitQPS), cfg.Google.RateLimitQPS),
		ConfigManager: cm,
	}
}

// InitOAuthConfig loads credentials.json
func (s *DriveService) InitOAuthConfig() error {
	logger.Info("🔐 [InitOAuthConfig] Loading credentials from: %s", model.CredFile)
	
	// Check if file exists
	if _, statErr := os.Stat(model.CredFile); os.IsNotExist(statErr) {
		logger.Error("🔐 [InitOAuthConfig] File does not exist: %s", model.CredFile)
		return statErr
	}
	
	b, err := os.ReadFile(model.CredFile)
	if err != nil {
		logger.Error("🔐 [InitOAuthConfig] Failed to read file: %v", err)
		return err
	}
	
	logger.Info("🔐 [InitOAuthConfig] File read successfully, size: %d bytes", len(b))
	
	// Log partial content for debugging (hide sensitive data)
	if len(b) > 0 {
		logger.Info("🔐 [InitOAuthConfig] File content preview: %s...", string(b[:min(100, len(b))]))
	}
	
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		logger.Error("🔐 [InitOAuthConfig] Failed to parse JSON: %v", err)
		return err
	}
	
	s.OAuthConfig = config
	logger.Info("🔐 [InitOAuthConfig] OAuth config loaded successfully")
	logger.Info("🔐 [InitOAuthConfig] ClientID: %s...", config.ClientID[:min(20, len(config.ClientID))])
	logger.Info("🔐 [InitOAuthConfig] RedirectURL: %s", config.RedirectURL)
	return nil
}

// InitDriveService initializes the Drive client
func (s *DriveService) InitDriveService() error {
	if s.OAuthConfig == nil {
		if err := s.InitOAuthConfig(); err != nil {
			return err
		}
	}

	tok, err := s.TokenFromFile(model.TokenFile)
	if err != nil {
		logger.Warning("⚠️ Token not found or invalid, please login via WebUI")
		return err
	}
	client := s.OAuthConfig.Client(context.Background(), tok)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logger.Error("Failed to create Drive service: %v", err)
		return err
	}
	s.Srv = srv
	logger.Info("✅ Drive service initialized")
	return nil
}

// TokenFromFile loads token from file
func (s *DriveService) TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveToken saves token to file
func (s *DriveService) SaveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		logger.Error("Unable to cache OAuth token: %v", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// WaitRateLimit waits for rate limiter
func (s *DriveService) WaitRateLimit() {
	if s.Limiter != nil {
		_ = s.Limiter.Wait(context.Background())
	}
}

// GetDriveName gets drive name (with cache)
func (s *DriveService) GetDriveName(driveID string) string {
	if driveID == "" {
		cfg := s.ConfigManager.GetConfig()
		name := cfg.Google.PersonalDriveName
		if name == "" {
			return "Cloud Drive"
		}
		return name
	}
	if val, ok := s.DriveNameCache.Load(driveID); ok {
		return val.(string)
	}
	if _, loaded := s.DriveCacheLoad.LoadOrStore(driveID, true); loaded {
		time.Sleep(100 * time.Millisecond)
		if val, ok := s.DriveNameCache.Load(driveID); ok {
			return val.(string)
		}
		return driveID
	}
	defer s.DriveCacheLoad.Delete(driveID)

	s.WaitRateLimit()
	if s.Srv == nil {
		return driveID
	}

	d, err := s.Srv.Drives.Get(driveID).Fields("name").Do()
	if err != nil {
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "insufficient authentication scopes") {
			logger.Warning("⚠️ Insufficient permissions to get shared drive name (ID: %s), showing ID.", driveID)
			s.DriveNameCache.Store(driveID, driveID)
			return driveID
		}
		logger.Error("Failed to get shared drive (ID: %s): %v", driveID, err)
		return driveID
	}
	s.DriveNameCache.Store(driveID, d.Name)
	logger.Verbose(model.LogLevelInfo, "💾 Found new shared drive: [%s] -> %s", driveID, d.Name)
	return d.Name
}

// ListAllDrives lists all shared drives
func (s *DriveService) ListAllDrives() ([]*drive.Drive, error) {
	s.WaitRateLimit()
	if s.Srv == nil {
		return nil, nil // Not initialized
	}

	var allDrives []*drive.Drive
	pageToken := ""

	for {
		q := s.Srv.Drives.List().
			PageSize(100).
			Fields("nextPageToken, drives(id, name)")

		if pageToken != "" {
			q.PageToken(pageToken)
		}

		r, err := q.Do()
		if err != nil {
			return nil, err
		}

		allDrives = append(allDrives, r.Drives...)

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
	return allDrives, nil
}

// RegisterWatch registers webhook listener
func (s *DriveService) RegisterWatch(pageToken string) {
	if s.Srv == nil {
		return
	}
	id := uuid.New().String()
	exp := time.Now().Add(7 * 24 * time.Hour).UnixMilli()

	cfg := s.ConfigManager.GetConfig()
	domain := strings.TrimRight(cfg.Server.PublicURL, "/")
	path := strings.TrimLeft(cfg.Server.WebhookPath, "/")
	fullAddr := domain + "/" + path

	ch := &drive.Channel{Id: id, Type: "web_hook", Address: fullAddr, Expiration: exp}
	s.WaitRateLimit()
	_, err := s.Srv.Changes.Watch(pageToken, ch).IncludeItemsFromAllDrives(true).SupportsAllDrives(true).Do()
	if err != nil {
		logger.Error("Failed to register Watch: %v", err)
		return
	}
	logger.Info("✅ Webhook registered: %s", fullAddr)
}

// GetStartPageToken gets the locally saved PageToken
func (s *DriveService) GetStartPageToken() string {
	f, _ := os.ReadFile(model.StartTokenFile)
	return string(f)
}

// EnsureStartPageToken ensures there's an initial PageToken
func (s *DriveService) EnsureStartPageToken() {
	f, err := os.ReadFile(model.StartTokenFile)
	if err == nil && len(f) > 0 {
		return
	}
	s.WaitRateLimit()
	if s.Srv == nil {
		return
	}
	r, err := s.Srv.Changes.GetStartPageToken().SupportsAllDrives(true).Do()
	if err != nil {
		logger.Error("Failed to get StartPageToken: %v", err)
		return
	}
	s.SaveTokenStr(r.StartPageToken)
}

// SaveTokenStr saves PageToken string
func (s *DriveService) SaveTokenStr(t string) {
	err := os.WriteFile(model.StartTokenFile, []byte(t), 0644)
	if err != nil {
		logger.Error("Failed to save StartToken: %v", err)
	}
}
