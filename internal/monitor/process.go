package monitor

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessMonitor struct {
	processes []*process.Process
}

// ProcessInfo holds information about a process
type ProcessInfo struct {
	PID     int32
	Name    string
	Status  string
	Memory  uint64
	CPUPercent float64
	NumThreads int32
}

// FindProcessesByName finds all processes matching the given name
func FindProcessesByName(name string) ([]*ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var results []*ProcessInfo

	for _, proc := range procs {
		procName, err := proc.Name()
		if err != nil {
			continue
		}

		if procName == name {
			info, err := getProcessInfo(proc)
			if err != nil {
				log.Printf("Warning: failed to get info for process %d: %v", proc.Pid, err)
				continue
			}
			results = append(results, info)
		}
	}

	return results, nil
}

// GetProcessByPID returns information about a specific process
func GetProcessByPID(pid int32) (*ProcessInfo, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get process %d: %w", pid, err)
	}

	return getProcessInfo(proc)
}

// GetAllProcesses returns information about all running processes
func GetAllProcesses() ([]*ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var results []*ProcessInfo

	for _, proc := range procs {
		info, err := getProcessInfo(proc)
		if err != nil {
			continue
		}
		results = append(results, info)
	}

	return results, nil
}

// getProcessInfo retrieves information about a process
func getProcessInfo(proc *process.Process) (*ProcessInfo, error) {
	name, err := proc.Name()
	if err != nil {
		return nil, err
	}

	status, err := proc.Status()
	if err != nil {
		status = []string{"unknown"}
	}

	memInfo, err := proc.MemoryInfo()
	var mem uint64
	if err == nil && memInfo != nil {
		mem = memInfo.RSS
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		cpuPercent = 0
	}

	numThreads, err := proc.NumThreads()
	if err != nil {
		numThreads = 0
	}

	return &ProcessInfo{
		PID:        proc.Pid,
		Name:       name,
		Status:     status[0],
		Memory:     mem,
		CPUPercent: cpuPercent,
		NumThreads: numThreads,
	}, nil
}

// Monitor monitors specific processes
func (pm *ProcessMonitor) Monitor(names []string) error {
	pm.processes = nil

	for _, name := range names {
		procs, err := FindProcessesByName(name)
		if err != nil {
			log.Printf("Warning: failed to find process %s: %v", name, err)
			continue
		}

		if len(procs) == 0 {
			log.Printf("Warning: process %s not found", name)
		}
	}

	return nil
}

// GetMonitoredProcesses returns information about monitored processes
func (pm *ProcessMonitor) GetMonitoredProcesses(names []string) (map[string][]*ProcessInfo, error) {
	result := make(map[string][]*ProcessInfo)

	for _, name := range names {
		procs, err := FindProcessesByName(name)
		if err != nil {
			log.Printf("Warning: failed to find process %s: %v", name, err)
			continue
		}
		result[name] = procs
	}

	return result, nil
}

// IsProcessRunning checks if a process with the given name is running
func IsProcessRunning(name string) bool {
	procs, err := FindProcessesByName(name)
	if err != nil {
		return false
	}
	return len(procs) > 0
}

// GetProcessCount returns the number of processes with the given name
func GetProcessCount(name string) (int, error) {
	procs, err := FindProcessesByName(name)
	if err != nil {
		return 0, err
	}
	return len(procs), nil
}
