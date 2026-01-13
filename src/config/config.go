package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"

	"gd-webhook/src/model"
)

// generateRandomPassword generates a random 16-character password
func generateRandomPassword() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

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

	// Ensure config directory exists
	_ = os.MkdirAll(model.ConfigDir, 0755)
	_ = os.MkdirAll(model.DataDir, 0755)

	// Set default values
	m.Cfg.Advanced.DebounceSeconds = 5
	m.Cfg.Advanced.RcloneWaitSeconds = 15
	m.Cfg.Advanced.LogSaveEnabled = true
	m.Cfg.Advanced.LogMaxSizeMB = 10
	m.Cfg.Advanced.LogRetentionDays = 7
	m.Cfg.Advanced.LogCleanupCron = "0 0 3 * * ?"
	m.Cfg.Server.ListenPort = 8448
	m.Cfg.Server.WebhookPath = "/gd-webhook"
	m.Cfg.Google.PersonalDriveName = "My Drive"

	// Track if we need to save config (for first-time setup)
	needSave := false
	isFirstRun := false

	f, err := os.ReadFile(model.ConfigFile)
	if err != nil {
		// Config file doesn't exist - first time setup
		isFirstRun = true
		fmt.Println("📝 Config file not found, creating default configuration...")
		needSave = true
	} else {
		if err := json.Unmarshal(f, m.Cfg); err != nil {
			// Try old format compatibility (if rclone is object instead of array)
			// For simplicity, assume user has updated config format or using new config
			fmt.Printf("❌ Config file JSON format error (ensure rclone config is an array): %v\n", err)
			return
		}
	}

	// Generate default credentials if not set
	if m.Cfg.Auth.Username == "" {
		m.Cfg.Auth.Username = "admin"
		needSave = true
	}
	if m.Cfg.Auth.Password == "" {
		m.Cfg.Auth.Password = generateRandomPassword()
		needSave = true
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println("🔐 Generated default login credentials:")
		fmt.Printf("   Username: %s\n", m.Cfg.Auth.Username)
		fmt.Printf("   Password: %s\n", m.Cfg.Auth.Password)
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println("⚠️  Please change your password after first login!")
		fmt.Println("═══════════════════════════════════════════════════════════")
	} else if isFirstRun {
		// First run but password was somehow set (shouldn't happen, but just in case)
		fmt.Printf("🔐 Using configured credentials - Username: %s\n", m.Cfg.Auth.Username)
	}

	// Save config if needed (first run or missing credentials)
	if needSave {
		if err := m.saveConfigWithoutLock(); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
		} else if isFirstRun {
			fmt.Printf("✅ Config saved to: %s\n", model.ConfigFile)
		}
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

	// Set defaults for Symedia timeout (Default 60s, Max 120s)
	if m.Cfg.Symedia.Timeout <= 0 {
		m.Cfg.Symedia.Timeout = 60
	} else if m.Cfg.Symedia.Timeout > 120 {
		m.Cfg.Symedia.Timeout = 120
	}

	// Set defaults for Google ListDelay (Min 1000ms)
	if m.Cfg.Google.ListDelay < 1000 {
		m.Cfg.Google.ListDelay = 1000
	}
	// Strict: BatchSleepInterval (Min 300s)
	if m.Cfg.Google.BatchSleepInterval < 300 {
		m.Cfg.Google.BatchSleepInterval = 300
	}

	// Set defaults for Rclone timeouts (Default 60s, Max 120s)
	for i := range m.Cfg.Rclone {
		if m.Cfg.Rclone[i].Timeout <= 0 {
			m.Cfg.Rclone[i].Timeout = 60
		} else if m.Cfg.Rclone[i].Timeout > 120 {
			m.Cfg.Rclone[i].Timeout = 120
		}
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

	fmt.Printf("📜 Loaded %d SA rules, %d Rclone instances\n", len(m.SARegexRules), len(m.Cfg.Rclone))
}

// saveConfigWithoutLock saves configuration without acquiring lock (for internal use)
func (m *Manager) saveConfigWithoutLock() error {
	f, err := os.Create(model.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(m.Cfg)
}

// SaveConfig saves configuration to disk
func (m *Manager) SaveConfig() error {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	return m.saveConfigWithoutLock()
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
	// Validate/Fix configuration
	if newCfg.Symedia.Timeout > 120 {
		newCfg.Symedia.Timeout = 120
	}
	for i := range newCfg.Rclone {
		if newCfg.Rclone[i].Timeout > 120 {
			newCfg.Rclone[i].Timeout = 120
		}
	}
	if newCfg.Google.ListDelay < 1000 {
		newCfg.Google.ListDelay = 1000
	}

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
		fmt.Printf("❌ Failed to create credentials file: %v\n", err)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cred); err != nil {
		fmt.Printf("❌ Failed to write credentials file: %v\n", err)
		return
	}
	fmt.Println("🔐 credentials.json regenerated")
}
