package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mbahati3-wq/the-system-daemon/internal/daemon"
	"github.com/mbahati3-wq/the-system-daemon/pkg/pidfile"
)

func main() {
	var (
		configPath = flag.String("config", "/etc/mydaemon/mydaemon.yaml", "Path to configuration file")
		pidPath    = flag.String("pidfile", "/var/run/mydaemon.pid", "Path to PID file")
		debug      = flag.Bool("debug", false, "Enable debug mode")
		validate   = flag.Bool("validate-config", false, "Validate configuration and exit")
	)
	flag.Parse()

	// Validate configuration if requested
	if *validate {
		cfg, err := daemon.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		if err := cfg.Validate(); err != nil {
			log.Fatalf("Configuration validation failed: %v", err)
		}

		fmt.Println("Configuration is valid")
		os.Exit(0)
	}

	// Create PID file
	pf, err := pidfile.New(*pidPath)
	if err != nil {
		log.Fatalf("Failed to create PID file: %v", err)
	}
	defer pf.Remove()

	// Load and validate configuration
	cfg, err := daemon.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Enable debug if requested
	if *debug {
		cfg.Debug = true
	}

	// Create and start daemon
	d, err := daemon.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create daemon: %v", err)
	}

	fmt.Println("Starting mydaemon...")
	if err := d.Start(); err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}
}
