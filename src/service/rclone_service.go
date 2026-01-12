package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gd-webhook/src/config"
	"gd-webhook/src/logger"
	"gd-webhook/src/model"
)

// RcloneService handles Rclone integration
type RcloneService struct {
	ConfigManager *config.Manager
	LimitChan     chan struct{}
}

// NewRcloneService creates a new Rclone service
func NewRcloneService(cm *config.Manager) *RcloneService {
	return &RcloneService{
		ConfigManager: cm,
		LimitChan:     make(chan struct{}, 5),
	}
}

// Refresh triggers Rclone VFS refresh
func (s *RcloneService) Refresh(originPath string) {
	s.ConfigManager.Lock.RLock()
	instances := s.ConfigManager.Cfg.Rclone
	regexRulesMap := s.ConfigManager.RcloneRegexRules // map[int][]*regexp.Regexp
	logLevel := s.ConfigManager.Cfg.Advanced.LogLevel
	s.ConfigManager.Lock.RUnlock()

	if len(instances) == 0 {
		return
	}

	// Refresh all matching instances concurrently
	for idx, instance := range instances {
		go func(i int, inst model.RcloneInstance) {
			regexRules := regexRulesMap[i]
			finalPath := originPath
			matched := false

			// Check regex rules
			for j, rule := range inst.Mapping {
				if j < len(regexRules) && regexRules[j].MatchString(originPath) {
					finalPath = regexRules[j].ReplaceAllString(originPath, rule.Replacement)
					matched = true
					logger.Debug(logLevel, "🔍 [Rclone-%s] Regex matched: %s -> %s", inst.Name, originPath, finalPath)
					break
				}
			}

			if !matched {
				// Try smart root directory matching
				tempPath := originPath
				if !strings.HasSuffix(tempPath, "/") {
					tempPath += "/"
				}
				for j, rule := range inst.Mapping {
					if j < len(regexRules) && regexRules[j].MatchString(tempPath) {
						finalPath = regexRules[j].ReplaceAllString(tempPath, rule.Replacement)
						finalPath = strings.TrimRight(finalPath, "/")
						if finalPath == "" {
							finalPath = "/"
						}
						matched = true
						logger.Debug(logLevel, "🔍 [Rclone-%s] Smart root match: %s -> %s", inst.Name, originPath, finalPath)
						break
					}
				}
			}

			if !matched {
				// If no rule matched, skip this refresh for this instance
				// In multi-instance setup, each manages its own paths
				return
			}

			rcHost := inst.Host
			rcEp := inst.Endpoint
			if rcEp == "" {
				rcEp = "/vfs/refresh"
			}

			// Build URL with _async=true parameter
			fullURL := strings.TrimRight(rcHost, "/") + "/" + strings.TrimLeft(rcEp, "/")
			if strings.Contains(fullURL, "?") {
				fullURL += "&_async=true"
			} else {
				fullURL += "?_async=true"
			}

			payload := map[string]string{"dir": finalPath, "recursive": "true"}
			data, _ := json.Marshal(payload)

			logger.Info("🔄 [Rclone-%s] Refreshing: %s", inst.Name, finalPath)

			// Limit concurrent requests
			s.LimitChan <- struct{}{}
			defer func() { <-s.LimitChan }()

			// Create request
			req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(data))
			if err != nil {
				logger.Error("❌ [Rclone-%s] Failed to create request: %v", inst.Name, err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			// Print detailed request info in debug mode
			if logLevel >= model.LogLevelDebug {
				logger.Debug(logLevel, "   👉 Method: %s", req.Method)
				logger.Debug(logLevel, "   👉 URL: %s", fullURL)
				logger.Debug(logLevel, "   👉 Headers:")
				for key, values := range req.Header {
					for _, value := range values {
						logger.Debug(logLevel, "      %s: %s", key, value)
					}
				}
				logger.Debug(logLevel, "   👉 Body: %s", string(data))
			}

			cl := &http.Client{Timeout: 30 * time.Second}
			resp, err := cl.Do(req)
			if err != nil {
				logger.Error("❌ [Rclone-%s] Refresh failed: %v", inst.Name, err)
				return
			}
			defer resp.Body.Close()

			// Read response body for debug
			respBody, _ := io.ReadAll(resp.Body)

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				logger.Info("✅ [Rclone-%s] Refresh successful [%s]", inst.Name, resp.Status)
				if logLevel >= model.LogLevelDebug {
					logger.Debug(logLevel, "   👈 Response Status: %s", resp.Status)
					logger.Debug(logLevel, "   👈 Response Body: %s", string(respBody))
				}
			} else {
				logger.Error("❌ [Rclone-%s] Refresh error [%s]", inst.Name, resp.Status)
				if logLevel >= model.LogLevelDebug {
					logger.Debug(logLevel, "   👈 Response Status: %s", resp.Status)
					logger.Debug(logLevel, "   👈 Response Body: %s", string(respBody))
				}
			}
		}(idx, instance)
	}
}

// WaitForCooldown waits for Rclone cooldown period
func (s *RcloneService) WaitForCooldown() {
	s.ConfigManager.Lock.RLock()
	wait := time.Duration(s.ConfigManager.Cfg.Advanced.RcloneWaitSeconds) * time.Second
	s.ConfigManager.Lock.RUnlock()

	if wait > 0 {
		logger.Verbose(model.LogLevelInfo, "⏳ Rclone cooldown (%v)...", wait)
		time.Sleep(wait)
	}
}
