package service

import (
	"strings"
	"time"

	"gd-webhook/src/logger"
)

// Retry helper for API calls
func (s *DriveService) retryRequest(operation func() error) error {
	maxRetries := 5
	baseDelay := 1 * time.Second
	var err error

	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		// Check if error is retryable (429, 5xx)
		isRetryable := false
		errMsg := err.Error()
		if strings.Contains(errMsg, "429") ||
			strings.Contains(errMsg, "500") ||
			strings.Contains(errMsg, "502") ||
			strings.Contains(errMsg, "503") ||
			strings.Contains(errMsg, "rateLimitExceeded") ||
			strings.Contains(errMsg, "userRateLimitExceeded") {
			isRetryable = true
		}

		if !isRetryable {
			return err
		}

		sleepTime := baseDelay * time.Duration(1<<i) // Exponential backoff: 1, 2, 4, 8, 16
		logger.Warning("⚠️ API error: %v. Retrying in %v...", err, sleepTime)

		time.Sleep(sleepTime)
	}
	return err
}
