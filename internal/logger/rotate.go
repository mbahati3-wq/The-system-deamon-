package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RotationConfig holds log rotation configuration
type RotationConfig struct {
	MaxSize    int64
	MaxBackups int
	MaxAge     int // days
}

// RotateBySize rotates the log file if it exceeds the size limit
func (l *Logger) RotateBySize(maxSize int64) error {
	fi, err := os.Stat(l.logPath)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if fi.Size() < maxSize {
		return nil
	}

	return l.Rotate()
}

// RotateByAge rotates log files older than maxAge days
func (l *Logger) RotateByAge(maxAge int) error {
	logDir := filepath.Dir(l.logPath)
	logName := filepath.Base(l.logPath)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if file matches pattern
		if !isLogFile(entry.Name(), logName) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		age := now.Sub(info.ModTime()).Hours() / 24
		if age > float64(maxAge) {
			filePath := filepath.Join(logDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove old log file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// CleanupOldLogs removes old backup logs, keeping only the maxBackups most recent
func (l *Logger) CleanupOldLogs(maxBackups int) error {
	logDir := filepath.Dir(l.logPath)
	logName := filepath.Base(l.logPath)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var backups []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && isLogFile(entry.Name(), logName) && entry.Name() != logName {
			backups = append(backups, entry)
		}
	}

	// Sort by modification time, newest first
	sort.Slice(backups, func(i, j int) bool {
		iInfo, _ := backups[i].Info()
		jInfo, _ := backups[j].Info()
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	// Remove old backups
	for i := maxBackups; i < len(backups); i++ {
		filePath := filepath.Join(logDir, backups[i].Name())
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove old log backup %s: %w", filePath, err)
		}
	}

	return nil
}

// isLogFile checks if a filename belongs to the log file pattern
func isLogFile(filename, logName string) bool {
	return filename == logName || (len(filename) > len(logName) && filename[:len(logName)] == logName)
}
