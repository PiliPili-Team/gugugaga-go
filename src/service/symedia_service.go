package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"
)

// SymediaService handles webhooks to external media systems (Emby)
type SymediaService struct {
	ConfigManager *config.Manager
}

// NewSymediaService creates a new Symedia service
func NewSymediaService(cm *config.Manager) *SymediaService {
	return &SymediaService{
		ConfigManager: cm,
	}
}

// SendWebhook sends a webhook notification
func (s *SymediaService) SendWebhook(originPath, action string, isDir bool) {
	s.ConfigManager.Lock.RLock()
	regexRules := s.ConfigManager.SARegexRules
	cfg := s.ConfigManager.Cfg
	s.ConfigManager.Lock.RUnlock()

	finalPath := originPath
	matched := false

	for i, rule := range cfg.Mapping {
		if regexRules[i].MatchString(originPath) {
			finalPath = regexRules[i].ReplaceAllString(originPath, rule.Replacement)
			matched = true
			logger.Debug(cfg.Advanced.LogLevel, "🔍 [SA] Regex matched: %s -> %s", originPath, finalPath)
			break
		}
	}
	if !matched {
		logger.Warning("⚠️ [SA] Regex not matched: %s", originPath)
	}

	notify := cfg.Symedia.NotifyUnmatched
	syHost := cfg.Symedia.Host
	syEp := cfg.Symedia.Endpoint
	bodyT := cfg.Symedia.BodyTemplate
	heads := cfg.Symedia.Headers

	if !matched && !notify {
		logger.Debug(cfg.Advanced.LogLevel, "🚫 [SA] Skipping: %s", originPath)
		return
	}

	fullURL := strings.TrimRight(syHost, "/") + "/" + strings.TrimLeft(syEp, "/")
	u, _ := url.Parse(fullURL)
	q := u.Query()

	replacements := map[string]interface{}{
		"{{FILE_PATH}}": finalPath,
		"{{ACTION}}":    action,
		"{{IS_DIR}}":    isDir,
	}

	// Process template in query parameters
	for k, v := range bodyT {
		if str, ok := v.(string); ok && !strings.Contains(str, "{{") {
			q.Set(k, str)
		}
	}
	u.RawQuery = q.Encode()

	finalBody := s.processTemplate(bodyT, replacements)

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	if m, ok := finalBody.(map[string]interface{}); ok {
		m["event_time"] = ts
		m["send_time"] = ts
	}

	b, _ := json.Marshal(finalBody)

	logger.Info("📡 Sending notification: %s", finalPath)
	if cfg.Advanced.LogLevel >= model.LogLevelDebug {
		logger.Debug(cfg.Advanced.LogLevel, "   👉 URL: %s", u.String())
		logger.Debug(cfg.Advanced.LogLevel, "   👉 Headers: %v", heads)
		logger.Debug(cfg.Advanced.LogLevel, "   👉 Body: %s", string(b))
	}

	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	// Set user-configured headers
	for k, v := range heads {
		req.Header.Set(k, v)
	}

	// If user didn't configure authorization, add default basic auth
	if _, exists := heads["authorization"]; !exists {
		req.Header.Set("authorization", "basic usernamepassword")
	}

	// If user didn't configure user-agent, add default user-agent
	if _, exists := heads["user-agent"]; !exists {
		req.Header.Set("user-agent", "clouddrive2/0.9.8")
	}

	timeout := cfg.Symedia.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	cl := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		logger.Error("Webhook failed: %v", err)
	} else {
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			logger.Info("✅ [SA] Push successful [%s]", resp.Status)
		} else {
			logger.Error("❌ [SA] Push error [%s]", resp.Status)
		}
	}
}

// processTemplate recursively processes template replacement
func (s *SymediaService) processTemplate(data interface{}, replacements map[string]interface{}) interface{} {
	switch v := data.(type) {
	case string:
		for k, val := range replacements {
			if v == k {
				return val
			}
			if strings.Contains(v, k) {
				strVal := fmt.Sprintf("%v", val)
				v = strings.Replace(v, k, strVal, -1)
			}
		}
		return v
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, val := range v {
			newMap[k] = s.processTemplate(val, replacements)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, val := range v {
			newSlice[i] = s.processTemplate(val, replacements)
		}
		return newSlice
	default:
		return v
	}
}
