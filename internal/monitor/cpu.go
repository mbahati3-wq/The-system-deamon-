package monitor

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUMonitor struct {
	lastPercent float64
}

// GetCPUUsage returns the current CPU usage percentage
func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, nil
}

// GetCPUCount returns the number of CPU cores
func GetCPUCount() (int, error) {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU count: %w", err)
	}
	return count, nil
}

// GetCPUStats returns detailed CPU statistics
func GetCPUStats() (map[string]interface{}, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	stats := make(map[string]interface{})
	if len(info) > 0 {
		stats["model_name"] = info[0].ModelName
		stats["cores"] = info[0].Cores
		stats["mhz"] = info[0].Mhz
		stats["family"] = info[0].Family
		stats["vendor_id"] = info[0].VendorID
	}

	usage, err := GetCPUUsage()
	if err != nil {
		log.Printf("Warning: failed to get CPU usage: %v", err)
	} else {
		stats["usage_percent"] = usage
	}

	return stats, nil
}

// Monitor periodically monitors CPU
func (cm *CPUMonitor) Monitor() error {
	usage, err := GetCPUUsage()
	if err != nil {
		return err
	}

	cm.lastPercent = usage
	return nil
}

// GetLastPercent returns the last recorded CPU percentage
func (cm *CPUMonitor) GetLastPercent() float64 {
	return cm.lastPercent
}
