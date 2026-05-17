package test

import (
	"testing"
	"time"

	"github.com/mbahati3-wq/the-system-daemon/internal/daemon"
	"github.com/mbahati3-wq/the-system-daemon/internal/monitor"
)

// TestDaemonStartStop tests basic daemon start and stop
func TestDaemonStartStop(t *testing.T) {
	cfg := daemon.DefaultConfig()
	cfg.Debug = true

	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create daemon: %v", err)
	}

	// Shutdown daemon immediately
	err = d.Shutdown()
	if err != nil {
		t.Fatalf("Failed to shutdown daemon: %v", err)
	}
}

// TestMonitorInit tests monitor initialization
func TestMonitorInit(t *testing.T) {
	cfg := &monitor.Config{
		CPUInterval:     1 * time.Second,
		MemoryInterval:  1 * time.Second,
		DiskInterval:    5 * time.Second,
		NetworkInterval: 5 * time.Second,
	}

	mon, err := monitor.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create monitor: %v", err)
	}

	if mon == nil {
		t.Fatal("Monitor is nil")
	}
}

// TestCPUMonitoring tests CPU monitoring functionality
func TestCPUMonitoring(t *testing.T) {
	usage, err := monitor.GetCPUUsage()
	if err != nil {
		t.Fatalf("Failed to get CPU usage: %v", err)
	}

	if usage < 0 || usage > 100 {
		t.Fatalf("Invalid CPU usage percentage: %v", usage)
	}

	t.Logf("CPU Usage: %.2f%%", usage)
}

// TestMemoryMonitoring tests memory monitoring functionality
func TestMemoryMonitoring(t *testing.T) {
	percent, err := monitor.GetMemoryPercent()
	if err != nil {
		t.Fatalf("Failed to get memory percentage: %v", err)
	}

	if percent < 0 || percent > 100 {
		t.Fatalf("Invalid memory percentage: %v", percent)
	}

	t.Logf("Memory Usage: %.2f%%", percent)
}

// TestDiskMonitoring tests disk monitoring functionality
func TestDiskMonitoring(t *testing.T) {
	partitions, err := monitor.GetAllPartitions()
	if err != nil {
		t.Fatalf("Failed to get disk partitions: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatal("No disk partitions found")
	}

	t.Logf("Found %d disk partitions", len(partitions))
}
