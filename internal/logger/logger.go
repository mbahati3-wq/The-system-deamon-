package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	file       *os.File
	logger     *log.Logger
	logPath    string
	maxSize    int64
	maxBackups int
	debug      bool
}

// New creates a new logger instance
func New(logPath string, debug bool) (*Logger, error) {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Setup multi-writer to write to both file and stderr if debug
	var writers []io.Writer
	writers = append(writers, file)
	if debug {
		writers = append(writers, os.Stderr)
	}

	l := &Logger{
		file:       file,
		logger:     log.New(io.MultiWriter(writers...), "[DAEMON] ", log.LstdFlags|log.Lshortfile),
		logPath:    logPath,
		maxSize:    100 * 1024 * 1024, // 100MB default
		maxBackups: 10,
		debug:      debug,
	}

	return l, nil
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Printf("[INFO] %s", msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Printf("[ERROR] %s", msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Printf("[WARN] %s", msg)
}

// Debug logs a debug message (only if debug is enabled)
func (l *Logger) Debug(msg string) {
	if l.debug {
		l.logger.Printf("[DEBUG] %s", msg)
	}
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Rotate performs log rotation
func (l *Logger) Rotate() error {
	// Check if log file needs rotation
	fi, err := os.Stat(l.logPath)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if fi.Size() < l.maxSize {
		return nil // No rotation needed
	}

	// Close current file
	if err := l.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}

	// Rotate files
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := l.logPath + "." + timestamp

	if err := os.Rename(l.logPath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Reopen log file
	file, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen log file: %w", err)
	}

	l.file = file
	return nil
}
