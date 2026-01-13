package model

const (
	DataDir        = "userdata/data"
	ConfigDir      = "userdata/config"
	StartTokenFile = "userdata/data/start_token.txt"
	TokenFile      = "userdata/config/token.json"
	CredFile       = "userdata/config/credentials.json"
	ConfigFile     = "userdata/config/config.json"
	TreeCacheFile  = "userdata/data/tree_cache.json"

	// MaxWebLogs is the max log lines displayed in frontend
	MaxWebLogs = 500
)

const (
	LogLevelQuiet = 0 // Core changes only
	LogLevelInfo  = 1 // Flow information
	LogLevelDebug = 2 // Debug information
)

// Config represents all application configuration
type Config struct {
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`

	OAuthConfig struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
	} `json:"oauth_config"`

	Advanced struct {
		LogLevel          int    `json:"log_level"`
		LogSaveEnabled    bool   `json:"log_save_enabled"`
		LogDir            string `json:"log_dir"`
		LogMaxSizeMB      int    `json:"log_max_size_mb"`
		DebounceSeconds   int    `json:"debounce_seconds"`
		RcloneWaitSeconds int    `json:"rclone_wait_seconds"`

		// Log cleanup config
		LogCleanupEnabled bool   `json:"log_cleanup_enabled"` // Enable/disable
		LogRetentionDays  int    `json:"log_retention_days"`  // Retention days
		LogCleanupCron    string `json:"log_cleanup_cron"`    // Cron expression

		// Task statistics persistence
		TaskStats struct {
			TodayCompleted   int64  `json:"today_completed"`
			HistoryCompleted int64  `json:"history_completed"`
			LastResetDate    string `json:"last_reset_date"`
		} `json:"task_stats"`
	} `json:"advanced"`

	Server struct {
		ListenPort  int    `json:"listen_port"`
		PublicURL   string `json:"public_url"`   // Domain only, e.g. https://mydomain.com
		WebhookPath string `json:"webhook_path"` // Path only, e.g. /gd-webhook
		SSL         struct {
			Enabled          bool   `json:"enabled"`
			CertPath         string `json:"cert_path"`
			KeyPath          string `json:"key_path"`
			RestrictToDomain bool   `json:"restrict_to_domain"` // Restrict to domain access
		} `json:"ssl"`
	} `json:"server"`

	Google struct {
		RateLimitQPS       int               `json:"rate_limit_qps"`
		PersonalDriveName  string            `json:"personal_drive_name"`
		TargetDriveIDs     []string          `json:"target_drive_ids"`     // Target Team Drive IDs
		TargetDriveRemarks map[string]string `json:"target_drive_remarks"` // Remarks for target drives
		ListDelay          int               `json:"list_delay"`           // Milliseconds, min 1000
		BatchSleepInterval int               `json:"batch_sleep_interval"` // Sleep seconds every 1000 items
	} `json:"google"`

	Rclone  []RcloneInstance `json:"rclone"`
	Symedia struct {
		Host            string                 `json:"host"`
		Endpoint        string                 `json:"endpoint"`
		NotifyUnmatched bool                   `json:"notify_unmatched"`
		Headers         map[string]string      `json:"headers"`
		BodyTemplate    map[string]interface{} `json:"body_template"`
		Timeout         int                    `json:"timeout"` // Seconds
	} `json:"symedia"`
	Mapping []MappingRule `json:"path_mapping"`
}

// RcloneInstance represents Rclone instance configuration
type RcloneInstance struct {
	Name     string        `json:"name"`
	Host     string        `json:"host"`
	Endpoint string        `json:"endpoint"`
	Timeout  int           `json:"timeout"` // Seconds
	Mapping  []MappingRule `json:"mapping"`
}

// MappingRule represents path mapping rule
type MappingRule struct {
	Regex       string `json:"regex"`       // Search pattern
	Replacement string `json:"replacement"` // Replacement text (empty to delete)
}

// GoogleCredJSON represents Google credentials JSON structure
type GoogleCredJSON struct {
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectURIs            []string `json:"redirect_uris"`
	} `json:"web"`
}

// FileNode represents a node in the file tree
type FileNode struct {
	ID       string
	Name     string
	ParentID string
	IsDir    bool
	DriveID  string
}

// DescendantInfo contains traversal result information
type DescendantInfo struct {
	ID      string
	Path    string
	IsDir   bool
	DriveID string
}

// LogsResponse represents logs API response
type LogsResponse struct {
	Logs    []string `json:"logs"`
	NextIdx int      `json:"next_idx"`
}

// OAuthLoginURLResponse represents login URL response
type OAuthLoginURLResponse struct {
	URL string `json:"url"`
}

// TestSymediaRequest represents test webhook request body
type TestSymediaRequest struct {
	Path string `json:"path"`
}
