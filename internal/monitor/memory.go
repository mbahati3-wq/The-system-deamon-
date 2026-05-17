package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryMonitor struct {
	lastPercent float64
}

// GetMemoryUsage returns current memory usage statistics
func GetMemoryUsage() (*mem.VirtualMemoryStat, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}
	return vmem, nil
}

// GetMemoryPercent returns the current memory usage percentage
func GetMemoryPercent() (float64, error) {
	vmem, err := GetMemoryUsage()
	if err != nil {
		return 0, err
	}
	return vmem.UsedPercent, nil
}

// GetSwapUsage returns current swap usage statistics
func GetSwapUsage() (*mem.SwapMemoryStat, error) {
	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get swap usage: %w", err)
	}
	return swap, nil
}

// GetMemoryStats returns detailed memory statistics
func GetMemoryStats() (map[string]interface{}, error) {
	vmem, err := GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	swap, err := GetSwapUsage()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":        vmem.Total,
		"available":    vmem.Available,
		"used":         vmem.Used,
		"free":         vmem.Free,
		"used_percent": vmem.UsedPercent,
		"swap_total":   swap.Total,
		"swap_used":    swap.Used,
		"swap_free":    swap.Free,
	}, nil
}

// Monitor periodically monitors memory
func (mm *MemoryMonitor) Monitor() error {
	percent, err := GetMemoryPercent()
	if err != nil {
		return err
	}

	mm.lastPercent = percent
	return nil
}

// GetLastPercent returns the last recorded memory percentage
func (mm *MemoryMonitor) GetLastPercent() float64 {
	return mm.lastPercent
}
