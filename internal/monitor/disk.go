package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskMonitor struct {
	lastPartitions []disk.PartitionStat
}

// GetDiskUsage returns disk usage for the specified path
func GetDiskUsage(path string) (*disk.UsageStat, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage for %s: %w", path, err)
	}
	return usage, nil
}

// GetAllPartitions returns all disk partitions
func GetAllPartitions() ([]disk.PartitionStat, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %w", err)
	}
	return partitions, nil
}

// GetDiskStats returns detailed disk statistics
func GetDiskStats() (map[string]interface{}, error) {
	partitions, err := GetAllPartitions()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	partitionStats := []map[string]interface{}{}

	for _, partition := range partitions {
		usage, err := GetDiskUsage(partition.Mountpoint)
		if err != nil {
			continue
		}

		partitionStats = append(partitionStats, map[string]interface{}{
			"device":       partition.Device,
			"mountpoint":   partition.Mountpoint,
			"fstype":       partition.Fstype,
			"total":        usage.Total,
			"used":         usage.Used,
			"free":         usage.Free,
			"used_percent": usage.UsedPercent,
		})
	}

	stats["partitions"] = partitionStats
	return stats, nil
}

// Monitor periodically monitors disk
func (dm *DiskMonitor) Monitor() error {
	partitions, err := GetAllPartitions()
	if err != nil {
		return err
	}

	dm.lastPartitions = partitions
	return nil
}

// GetLastPartitions returns the last recorded partitions
func (dm *DiskMonitor) GetLastPartitions() []disk.PartitionStat {
	return dm.lastPartitions
}
