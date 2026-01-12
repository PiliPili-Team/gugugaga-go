package config

import (
	"encoding/json"
	"os"
	"regexp"
	"sync"

	"gd-webhook/src/logger"
	"gd-webhook/src/model"
)

// GetAppName returns the application name, reads from APP_NAME env var, defaults to "GD Watcher"
func GetAppName() string {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		return "GD Watcher"
	}
	return appName
}

// GetAppVersion returns the application version, reads from APP_VERSION env var, defaults to "4.0"
func GetAppVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "4.0"
	}
	return version
}

// Manager manages configuration loading and updates
type Manager struct {
	Cfg              *model.Config
	Lock             sync.RWMutex
	SARegexRules     []*regexp.Regexp         // Symedia regex rules cache
	RcloneRegexRules map[int][]*regexp.Regexp // Rclone regex rules cache (Index -> Rules)
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		Cfg:              &model.Config{},
		RcloneRegexRules: make(map[int][]*regexp.Regexp),
	}
}

// LoadConfig loads configuration from disk
func (m *Manager) LoadConfig() {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	// Set default values
	m.Cfg.Advanced.DebounceSeconds = 5
	m.Cfg.Advanced.RcloneWaitSeconds = 15
	m.Cfg.Advanced.LogSaveEnabled = true
	m.Cfg.Advanced.LogMaxSizeMB = 10
	m.Cfg.Advanced.LogRetentionDays = 7
	m.Cfg.Advanced.LogCleanupCron = "0 0 3 * * ?"
	m.Cfg.Server.ListenPort = 8448
	m.Cfg.Server.WebhookPath = "/gd-webhook"
	if m.Cfg.Google.PersonalDriveName == "" {
		m.Cfg.Google.PersonalDriveName = "My Drive"
	}

	f, err := os.ReadFile(model.ConfigFile)
	if err != nil {
		logger.Error("Failed to read config file: %v", err)
		return
	}
	if err := json.Unmarshal(f, m.Cfg); err != nil {
		// Try old format compatibility (if rclone is object instead of array)
		// For simplicity, assume user has updated config format or using new config
		logger.Error("Config file JSON format error (ensure rclone config is an array): %v", err)
		return
	}

	// Compatibility check
	if m.Cfg.Advanced.LogRetentionDays <= 0 {
		m.Cfg.Advanced.LogRetentionDays = 7
	}
	if m.Cfg.Advanced.LogCleanupCron == "" {
		m.Cfg.Advanced.LogCleanupCron = "0 0 3 * * ?"
	}

	if m.Cfg.Server.WebhookPath == "" {
		m.Cfg.Server.WebhookPath = "/gd-webhook"
	}

	// Compile regex rules
	m.SARegexRules = nil
	for _, mapping := range m.Cfg.Mapping {
		r, err := regexp.Compile(mapping.Regex)
		if err == nil {
			m.SARegexRules = append(m.SARegexRules, r)
		}
	}

	m.RcloneRegexRules = make(map[int][]*regexp.Regexp)
	for idx, instance := range m.Cfg.Rclone {
		var rules []*regexp.Regexp
		for _, mapping := range instance.Mapping {
			r, err := regexp.Compile(mapping.Regex)
			if err == nil {
				rules = append(rules, r)
			}
		}
		m.RcloneRegexRules[idx] = rules
	}

	logger.Verbose(model.LogLevelInfo, "📜 Loaded %d SA rules, %d Rclone instances", len(m.SARegexRules), len(m.Cfg.Rclone))
}

// SaveConfig saves configuration to disk
func (m *Manager) SaveConfig() error {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	f, err := os.Create(model.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(m.Cfg)
}

// GetConfig returns a copy of the configuration
func (m *Manager) GetConfig() model.Config {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	return *m.Cfg
}

// UpdateConfig updates configuration and recompiles regex rules
func (m *Manager) UpdateConfig(newCfg model.Config) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	*m.Cfg = newCfg

	// Recompile regex rules
	m.SARegexRules = nil
	for _, mapping := range m.Cfg.Mapping {
		r, err := regexp.Compile(mapping.Regex)
		if err == nil {
			m.SARegexRules = append(m.SARegexRules, r)
		}
	}

	m.RcloneRegexRules = make(map[int][]*regexp.Regexp)
	for idx, instance := range m.Cfg.Rclone {
		var rules []*regexp.Regexp
		for _, mapping := range instance.Mapping {
			r, err := regexp.Compile(mapping.Regex)
			if err == nil {
				rules = append(rules, r)
			}
		}
		m.RcloneRegexRules[idx] = rules
	}
}

// SaveCredentialsFile regenerates credentials.json
func (m *Manager) SaveCredentialsFile(id, secret, redirect string) {
	cred := model.GoogleCredJSON{}
	cred.Web.ClientID = id
	cred.Web.ClientSecret = secret
	cred.Web.RedirectURIs = []string{redirect}
	cred.Web.AuthURI = "https://accounts.google.com/o/oauth2/auth"
	cred.Web.TokenURI = "https://oauth2.googleapis.com/token"
	cred.Web.AuthProviderX509CertURL = "https://www.googleapis.com/oauth2/v1/certs"
	cred.Web.ProjectID = "gd-watcher"

	f, err := os.Create(model.CredFile)
	if err != nil {
		logger.Error("Failed to create credentials file: %v", err)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cred); err != nil {
		logger.Error("Failed to write credentials file: %v", err)
		return
	}
	logger.Info("🔐 credentials.json regenerated")
}
