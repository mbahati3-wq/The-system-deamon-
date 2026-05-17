package monitor

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Config struct {
	CPUInterval     time.Duration
	MemoryInterval  time.Duration
	DiskInterval    time.Duration
	NetworkInterval time.Duration
}

type Monitor struct {
	config        *Config
	cpuMonitor    *CPUMonitor
	memMonitor    *MemoryMonitor
	diskMonitor   *DiskMonitor
	netMonitor    *NetworkMonitor
	ticker        *time.Ticker
	stopChan      chan struct{}
	running       bool
}

// New creates a new monitor instance
func New(cfg *Config) (*Monitor, error) {
	if cfg == nil {
		cfg = &Config{
			CPUInterval:     5 * time.Second,
			MemoryInterval:  5 * time.Second,
			DiskInterval:    10 * time.Second,
			NetworkInterval: 10 * time.Second,
		}
	}

	return &Monitor{
		config:     cfg,
		cpuMonitor: &CPUMonitor{},
		memMonitor: &MemoryMonitor{},
		diskMonitor: &DiskMonitor{},
		netMonitor: &NetworkMonitor{},
		stopChan:   make(chan struct{}),
		running:    false,
	}, nil
}

// Start starts the monitor
func (m *Monitor) Start(ctx context.Context) error {
	if m.running {
		return fmt.Errorf("monitor already running")
	}

	m.running = true
	m.ticker = time.NewTicker(time.Second)

	go func() {
		cpuTicker := time.NewTicker(m.config.CPUInterval)
		memTicker := time.NewTicker(m.config.MemoryInterval)
		diskTicker := time.NewTicker(m.config.DiskInterval)
		netTicker := time.NewTicker(m.config.NetworkInterval)

		defer cpuTicker.Stop()
		defer memTicker.Stop()
		defer diskTicker.Stop()
		defer netTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				m.running = false
				return

			case <-cpuTicker.C:
				if err := m.cpuMonitor.Monitor(); err != nil {
					log.Printf("CPU monitoring error: %v", err)
				}

			case <-memTicker.C:
				if err := m.memMonitor.Monitor(); err != nil {
					log.Printf("Memory monitoring error: %v", err)
				}

			case <-diskTicker.C:
				if err := m.diskMonitor.Monitor(); err != nil {
					log.Printf("Disk monitoring error: %v", err)
				}

			case <-netTicker.C:
				if err := m.netMonitor.Monitor(); err != nil {
					log.Printf("Network monitoring error: %v", err)
				}
			}
		}
	}()

	return nil
}

// Stop stops the monitor
func (m *Monitor) Stop() error {
	if !m.running {
		return nil
	}

	m.running = false
	close(m.stopChan)

	if m.ticker != nil {
		m.ticker.Stop()
	}

	return nil
}

// GetStats returns current monitoring statistics
func (m *Monitor) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"cpu_percent":    m.cpuMonitor.GetLastPercent(),
		"memory_percent": m.memMonitor.GetLastPercent(),
		"network_stats":  m.netMonitor.GetLastStats(),
	}
}

// IsRunning returns whether the monitor is currently running
func (m *Monitor) IsRunning() bool {
	return m.running
}
