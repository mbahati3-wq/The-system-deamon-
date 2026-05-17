package pidfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PIDFile struct {
	path string
	pid  int
}

// New creates a new PID file
func New(path string) (*PIDFile, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create pidfile directory: %w", err)
	}

	// Check if process is already running
	if pid, err := Read(path); err == nil {
		if IsProcessRunning(pid) {
			return nil, fmt.Errorf("process already running with PID %d", pid)
		}
	}

	// Write current PID
	currentPID := os.Getpid()
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", currentPID)), 0644); err != nil {
		return nil, fmt.Errorf("failed to write pidfile: %w", err)
	}

	return &PIDFile{
		path: path,
		pid:  currentPID,
	}, nil
}

// Read reads the PID from a PID file
func Read(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read pidfile: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid PID in pidfile: %w", err)
	}

	return pid, nil
}

// Remove removes the PID file
func (pf *PIDFile) Remove() error {
	if err := os.Remove(pf.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove pidfile: %w", err)
	}
	return nil
}

// GetPID returns the PID stored in this PID file
func (pf *PIDFile) GetPID() int {
	return pf.pid
}

// Path returns the path to the PID file
func (pf *PIDFile) Path() string {
	return pf.path
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	// Try to send signal 0 (which doesn't actually send a signal)
	// If it succeeds, the process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(os.Signal(nil))
	return err == nil
}
