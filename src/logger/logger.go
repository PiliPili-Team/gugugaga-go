package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gd-webhook/src/model"
)

var (
	memLogs        []string
	memLogLock     sync.RWMutex
	logFileHandle  *os.File
	currentLogDate string
	historyLock    sync.Mutex
)

// MemLogger is a custom Writer that stores logs in memory
type MemLogger struct{}

func (w *MemLogger) Write(p []byte) (n int, err error) {
	memLogLock.Lock()
	defer memLogLock.Unlock()
	if len(memLogs) >= model.MaxWebLogs {
		memLogs = memLogs[model.MaxWebLogs/10:] // Remove old logs to prevent unlimited memory growth
	}
	memLogs = append(memLogs, string(p))
	return len(p), nil
}

// InitLogging initializes the logging system
func InitLogging(cfg *model.Config) {
	enabled := cfg.Advanced.LogSaveEnabled
	baseDir := cfg.Advanced.LogDir

	if baseDir == "" {
		baseDir = "userdata/logs"
	}

	sysDir := filepath.Join(baseDir, "system")
	histDir := filepath.Join(baseDir, "history")
	_ = os.MkdirAll(sysDir, 0755)
	_ = os.MkdirAll(histDir, 0755)

	if logFileHandle != nil {
		_ = logFileHandle.Close()
		logFileHandle = nil
	}

	currentLogDate = time.Now().Format("2006-01-02")
	filePath := filepath.Join(sysDir, fmt.Sprintf("app.%s.log", currentLogDate))

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if enabled {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logFileHandle = f
			writers = append(writers, f)
		} else {
			fmt.Printf("❌ Failed to open log file: %v\n", err)
		}
	}
	writers = append(writers, &MemLogger{})
	log.SetOutput(io.MultiWriter(writers...))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// CheckLogRotation checks if log rotation is needed (daily)
func CheckLogRotation(cfg *model.Config) {
	newDate := time.Now().Format("2006-01-02")
	if newDate != currentLogDate {
		Info("📅 Rotating log for new day...")
		InitLogging(cfg)
	}
}

// GetMemLogs returns memory logs (returns a copy)
func GetMemLogs(startIdx int) ([]string, int) {
	memLogLock.RLock()
	defer memLogLock.RUnlock()

	if startIdx < 0 || startIdx > len(memLogs) {
		startIdx = 0
	}

	var result []string
	if startIdx < len(memLogs) {
		// Create copy to avoid concurrent read/write after return
		result = make([]string, len(memLogs)-startIdx)
		copy(result, memLogs[startIdx:])
	}
	return result, len(memLogs)
}

// ClearMemLogs clears memory logs
func ClearMemLogs() {
	memLogLock.Lock()
	defer memLogLock.Unlock()
	memLogs = []string{}
}

// GetLogFileHandle returns the current log file handle (if needed)
func GetLogFileHandle() *os.File {
	return logFileHandle
}

// Log helper functions

func Info(f string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf(f, v...))
}

func Verbose(currentLevel int, f string, v ...interface{}) {
	if currentLevel >= model.LogLevelInfo {
		_ = log.Output(2, fmt.Sprintf(f, v...))
	}
}

func Debug(currentLevel int, f string, v ...interface{}) {
	if currentLevel >= model.LogLevelDebug {
		_ = log.Output(2, fmt.Sprintf("🐞 "+f, v...))
	}
}

func Error(f string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("❌ "+f, v...))
}

func Warning(f string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("⚠️ "+f, v...))
}

// WriteHistory writes to history CSV file
func WriteHistory(cfg *model.Config, action, path string) {
	if !cfg.Advanced.LogSaveEnabled {
		return
	}
	baseDir := cfg.Advanced.LogDir
	if baseDir == "" {
		baseDir = "userdata/logs"
	}

	histPath := filepath.Join(baseDir, "history", fmt.Sprintf("history.%s.csv", time.Now().Format("2006-01-02")))

	historyLock.Lock()
	defer historyLock.Unlock()

	f, err := os.OpenFile(histPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer f.Close()
		_, _ = f.WriteString(fmt.Sprintf("%s,%s,%s\n", time.Now().Format(time.RFC3339), action, path))
	} else {
		Error("Failed to write history: %v", err)
	}
}

// CleanupLogs cleans old logs based on retention days
func CleanupLogs(cfg *model.Config) {
	if !cfg.Advanced.LogCleanupEnabled {
		return
	}

	days := cfg.Advanced.LogRetentionDays
	baseDir := cfg.Advanced.LogDir
	if baseDir == "" {
		baseDir = "userdata/logs"
	}
	if days <= 0 {
		days = 7
	}

	Info("🧹 Starting log cleanup (keeping last %d days)...", days)

	expiration := time.Now().AddDate(0, 0, -days)
	deletedCount := 0
	dirs := []string{filepath.Join(baseDir, "system"), filepath.Join(baseDir, "history")}

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(expiration) {
				// Don't delete the currently active log file
				if logFileHandle != nil && filepath.Base(logFileHandle.Name()) == file.Name() {
					continue
				}
				if err := os.Remove(filepath.Join(dir, file.Name())); err == nil {
					deletedCount++
					Debug(cfg.Advanced.LogLevel, "🗑️ Deleted expired log: %s", file.Name())
				}
			}
		}
	}
	if deletedCount > 0 {
		Info("✅ Cleaned up %d files", deletedCount)
	} else {
		Info("✅ Cleanup complete, no expired files")
	}
}
