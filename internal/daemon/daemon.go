package daemon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mbahati3-wq/the-system-daemon/internal/logger"
	"github.com/mbahati3-wq/the-system-daemon/internal/monitor"
)

type Daemon struct {
	config         *Config
	logger         *logger.Logger
	monitor        *monitor.Monitor
	alertManager   *AlertManager
	healthServer   *HealthServer
	metricsCollector *MetricsCollector
	ctx            context.Context
	cancel         context.CancelFunc
	sigChan        chan os.Signal
	reloadChan     chan struct{}
	mu             sync.RWMutex
	running        bool
}

// New creates a new daemon instance
func New(cfg *Config) (*Daemon, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize logger
	l, err := logger.New(cfg.LogPath, cfg.Debug)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize monitor
	mon, err := monitor.New(&monitor.Config{
		CPUInterval:      time.Duration(cfg.CPUMonitorInterval) * time.Second,
		MemoryInterval:   time.Duration(cfg.MemoryMonitorInterval) * time.Second,
		DiskInterval:     time.Duration(cfg.DiskMonitorInterval) * time.Second,
		NetworkInterval:  time.Duration(cfg.NetworkMonitorInterval) * time.Second,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize monitor: %w", err)
	}

	// Initialize alert manager
	alertMgr := NewAlertManager(l)

	// Initialize health server
	healthSrv := NewHealthServer(cfg.HealthCheckPort)

	// Initialize metrics collector
	var metricsColl *MetricsCollector
	if cfg.EnableMetrics {
		metricsColl = NewMetricsCollector(cfg.MetricsPort)
	}

	d := &Daemon{
		config:           cfg,
		logger:           l,
		monitor:          mon,
		alertManager:     alertMgr,
		healthServer:     healthSrv,
		metricsCollector: metricsColl,
		ctx:              ctx,
		cancel:           cancel,
		sigChan:          make(chan os.Signal, 1),
		reloadChan:       make(chan struct{}, 1),
	}

	return d, nil
}

// Start starts the daemon
func (d *Daemon) Start() error {
	d.mu.Lock()
	d.running = true
	d.mu.Unlock()

	d.logger.Info("Daemon starting...")

	// Setup signal handlers
	signal.Notify(d.sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	// Start alert manager
	d.alertManager.Start()
	d.logger.Info("Alert manager started")

	// Start health server
	if err := d.healthServer.Start(); err != nil {
		d.logger.Error(fmt.Sprintf("Failed to start health server: %v", err))
	} else {
		d.logger.Info(fmt.Sprintf("Health server started on port %d", d.config.HealthCheckPort))
		d.healthServer.SetAlive(true)
	}

	// Start metrics server if enabled
	if d.config.EnableMetrics && d.metricsCollector != nil {
		if err := d.metricsCollector.Start(); err != nil {
			d.logger.Error(fmt.Sprintf("Failed to start metrics server: %v", err))
		} else {
			d.logger.Info(fmt.Sprintf("Metrics server started on port %d", d.config.MetricsPort))
		}
	}

	// Start monitoring
	if err := d.monitor.Start(d.ctx); err != nil {
		return fmt.Errorf("failed to start monitor: %w", err)
	}

	// Register health checks
	d.healthServer.RegisterCheck("monitor", func() CheckStatus {
		if d.monitor.IsRunning() {
			return CheckStatus{Status: "ok", Message: "Monitor is running"}
		}
		return CheckStatus{Status: "error", Message: "Monitor is not running"}
	})

	d.healthServer.SetReady(true)
	d.logger.Info("Daemon started successfully")

	// Main loop
	return d.mainLoop()
}

// mainLoop is the main event loop of the daemon
func (d *Daemon) mainLoop() error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sig := <-d.sigChan:
			d.logger.Info(fmt.Sprintf("Received signal: %v", sig))
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				return d.Shutdown()
			case syscall.SIGHUP:
				d.logger.Info("Reloading configuration...")
				if err := d.reloadConfig(); err != nil {
					d.logger.Error(fmt.Sprintf("Configuration reload failed: %v", err))
				} else {
					d.logger.Info("Configuration reloaded successfully")
				}
			}

		case <-d.reloadChan:
			if err := d.reloadConfig(); err != nil {
				d.logger.Error(fmt.Sprintf("Configuration reload failed: %v", err))
			}

		case <-ticker.C:
			d.updateMetrics()

		case <-d.ctx.Done():
			return nil
		}
	}
}

// reloadConfig reloads the configuration from disk
func (d *Daemon) reloadConfig() error {
	cfg, err := LoadConfig(d.config.LogPath)
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	d.mu.Lock()
	d.config = cfg
	d.mu.Unlock()

	return nil
}

// updateMetrics updates metrics and checks thresholds
func (d *Daemon) updateMetrics() {
	if d.metricsCollector == nil {
		return
	}

	stats := d.monitor.GetStats()

	if cpu, ok := stats["cpu_percent"].(float64); ok {
		d.metricsCollector.UpdateCPUUsage(cpu)
		d.alertManager.CheckCPUThreshold(cpu, d.config.CPUThreshold)
	}

	if mem, ok := stats["memory_percent"].(float64); ok {
		d.metricsCollector.UpdateMemoryUsage(mem)
		d.alertManager.CheckMemoryThreshold(mem, d.config.MemoryThreshold)
	}

	// Check monitored processes
	if len(d.config.Processes) > 0 {
		for _, procName := range d.config.Processes {
			running := monitor.IsProcessRunning(procName)
			if !running {
				d.alertManager.CheckProcessAlert(procName, running)
			}
		}
	}

	d.mu.RLock()
	uptime := time.Since(time.Now())
	d.mu.RUnlock()

	d.metricsCollector.UpdateUptime(uptime)
}

// Shutdown gracefully shuts down the daemon
func (d *Daemon) Shutdown() error {
	d.mu.Lock()
	d.running = false
	d.mu.Unlock()

	d.logger.Info("Shutting down daemon...")

	// Stop health server
	if err := d.healthServer.Stop(); err != nil {
		d.logger.Error(fmt.Sprintf("Error stopping health server: %v", err))
	}

	// Stop metrics server
	if d.metricsCollector != nil {
		if err := d.metricsCollector.Stop(); err != nil {
			d.logger.Error(fmt.Sprintf("Error stopping metrics server: %v", err))
		}
	}

	// Stop alert manager
	d.alertManager.Stop()

	// Cancel context
	d.cancel()

	// Stop monitor
	if err := d.monitor.Stop(); err != nil {
		d.logger.Error(fmt.Sprintf("Error stopping monitor: %v", err))
	}

	// Close logger
	if err := d.logger.Close(); err != nil {
		log.Printf("Error closing logger: %v", err)
	}

	d.logger.Info("Daemon shutdown complete")
	return nil
}

// IsRunning returns whether the daemon is running
func (d *Daemon) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}

// GetConfig returns the current configuration
func (d *Daemon) GetConfig() *Config {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.config
}

// SetConfig updates the daemon configuration
func (d *Daemon) SetConfig(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	d.mu.Lock()
	d.config = cfg
	d.mu.Unlock()

	return nil
}

// TriggerReload triggers a configuration reload
func (d *Daemon) TriggerReload() {
	select {
	case d.reloadChan <- struct{}{}:
	default:
	}
}
