package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignals sets up signal handling for the daemon
func HandleSignals(d *Daemon) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGUSR1)

	go func() {
		for sig := range sigChan {
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				log.Println("Received termination signal, shutting down gracefully...")
				if err := d.Shutdown(); err != nil {
					log.Printf("Error during shutdown: %v", err)
				}
				os.Exit(0)

			case syscall.SIGHUP:
				log.Println("Received SIGHUP, reloading configuration...")
				// TODO: Implement configuration reload logic
				d.logger.Info("Configuration reload requested")

			case syscall.SIGUSR1:
				log.Println("Received SIGUSR1, printing stats...")
				// TODO: Implement stats printing logic
				d.logger.Info("Statistics requested")

			default:
				log.Printf("Received unexpected signal: %v", sig)
			}
		}
	}()
}

// WaitForSignal blocks until a signal is received
func WaitForSignal() os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	fmt.Printf("Signal received: %v\n", sig)
	return sig
}
