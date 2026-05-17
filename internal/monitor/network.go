package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

type NetworkMonitor struct {
	lastStats []net.IOCountersStat
}

// GetNetworkStats returns network I/O statistics
func GetNetworkStats(pernic bool) ([]net.IOCountersStat, error) {
	stats, err := net.IOCounters(pernic)
	if err != nil {
		return nil, fmt.Errorf("failed to get network stats: %w", err)
	}
	return stats, nil
}

// GetConnections returns active network connections
func GetConnections() ([]net.ConnectionStat, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}
	return conns, nil
}

// GetNetworkInterfaces returns network interface information
func GetNetworkInterfaces() ([]net.InterfaceStat, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}
	return interfaces, nil
}

// GetDetailedNetworkStats returns comprehensive network statistics
func GetDetailedNetworkStats() (map[string]interface{}, error) {
	stats, err := GetNetworkStats(true)
	if err != nil {
		return nil, err
	}

	interfaces, err := GetNetworkInterfaces()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"io_stats":   stats,
		"interfaces": interfaces,
	}, nil
}

// Monitor periodically monitors network
func (nm *NetworkMonitor) Monitor() error {
	stats, err := GetNetworkStats(true)
	if err != nil {
		return err
	}

	nm.lastStats = stats
	return nil
}

// GetLastStats returns the last recorded network statistics
func (nm *NetworkMonitor) GetLastStats() []net.IOCountersStat {
	return nm.lastStats
}
