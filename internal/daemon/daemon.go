package daemon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbahati3-wq/the-system-daemon/internal/logger"
	"github.com/mbahati3-wq/the-system-daemon/internal/monitor"
)

type Daemon struct {
	config  *Config
	logger  *logger.Logger
	monitor *monitor.Monitor
	ctx     context.Context
	cancel  context.CancelFunc
	sigChan chan os.Signal
}

// New creates a new daemon instance
func New(cfg *Config) (*Daemon, error) {
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

	d := &Daemon{
		config:  cfg,
		logger:  l,
		monitor: mon,
		ctx:     ctx,
		cancel:  cancel,
		sigChan: make(chan os.Signal, 1),
	}

	return d, nil
}

// Start starts the daemon
func (d *Daemon) Start() error {
	d.logger.Info("Daemon starting...")

	// Setup signal handlers
	setupSignalHandlers(d.sigChan, d.cancel)

	// Start monitoring
	if err := d.monitor.Start(d.ctx); err != nil {
		return fmt.Errorf("failed to start monitor: %w", err)
	}

	d.logger.Info("Daemon started successfully")

	// Wait for signals
	for {
		select {
		case sig := <-d.sigChan:
			d.logger.Info(fmt.Sprintf("Received signal: %v", sig))
			return d.Shutdown()
		case <-d.ctx.Done():
			return nil
		}
	}
}

// Shutdown gracefully shuts down the daemon
func (d *Daemon) Shutdown() error {
	d.logger.Info("Shutting down daemon...")
	d.cancel()

	if err := d.monitor.Stop(); err != nil {
		d.logger.Error(fmt.Sprintf("Error stopping monitor: %v", err))
	}

	if err := d.logger.Close(); err != nil {
		log.Printf("Error closing logger: %v", err)
	}

	d.logger.Info("Daemon shutdown complete")
	return nil
}

func setupSignalHandlers(sigChan chan<- os.Signal, cancel context.CancelFunc) {
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
}
