# MyDaemon API Documentation

## Overview

MyDaemon provides a comprehensive API for system monitoring. This document describes the internal API interfaces and their usage.

## Daemon API

### Daemon Interface

```go
type Daemon struct {
    config  *Config
    logger  *logger.Logger
    monitor *monitor.Monitor
    ctx     context.Context
    cancel  context.CancelFunc
}
```

#### Methods

##### `New(cfg *Config) (*Daemon, error)`
Creates a new daemon instance.

**Parameters:**
- `cfg`: Configuration object

**Returns:**
- `*Daemon`: Daemon instance
- `error`: Error if initialization fails

**Example:**
```go
cfg := daemon.DefaultConfig()
d, err := daemon.New(cfg)
if err != nil {
    log.Fatal(err)
}
```

##### `Start() error`
Starts the daemon and blocks until shutdown.

**Returns:**
- `error`: Error if startup fails

##### `Shutdown() error`
Gracefully shuts down the daemon.

**Returns:**
- `error`: Error during shutdown

## Configuration API

### Config Structure

```go
type Config struct {
    AppName                string
    Debug                  bool
    LogPath                string
    MaxLogSize             int
    MaxLogBackups          int
    CPUMonitorInterval     int
    MemoryMonitorInterval  int
    DiskMonitorInterval    int
    NetworkMonitorInterval int
    CPUThreshold           float64
    MemoryThreshold        float64
}
```

#### Functions

##### `DefaultConfig() *Config`
Returns default configuration.

##### `LoadConfig(path string) (*Config, error)`
Loads configuration from YAML file.

**Parameters:**
- `path`: Path to YAML configuration file

**Returns:**
- `*Config`: Loaded configuration
- `error`: Error if loading fails

##### `(c *Config) SaveConfig(path string) error`
Saves configuration to YAML file.

**Parameters:**
- `path`: Target path for YAML file

**Returns:**
- `error`: Error if saving fails

## Monitor API

### Monitor Interface

```go
type Monitor struct {
    config    *Config
    cpuMonitor    *CPUMonitor
    memMonitor    *MemoryMonitor
    diskMonitor   *DiskMonitor
    netMonitor    *NetworkMonitor
}
```

#### Methods

##### `New(cfg *Config) (*Monitor, error)`
Creates a new monitor instance.

##### `Start(ctx context.Context) error`
Starts monitoring.

**Parameters:**
- `ctx`: Context for cancellation

**Returns:**
- `error`: Error if startup fails

##### `Stop() error`
Stops monitoring.

##### `GetStats() map[string]interface{}`
Returns current statistics.

**Returns:**
- `map[string]interface{}`: Statistics dictionary

##### `IsRunning() bool`
Checks if monitor is running.

### CPU Monitoring

#### Functions

##### `GetCPUUsage() (float64, error)`
Returns current CPU usage percentage.

**Returns:**
- `float64`: Percentage (0-100)
- `error`: Error if retrieval fails

##### `GetCPUCount() (int, error)`
Returns number of CPU cores.

##### `GetCPUStats() (map[string]interface{}, error)`
Returns detailed CPU statistics.

### Memory Monitoring

#### Functions

##### `GetMemoryPercent() (float64, error)`
Returns current memory usage percentage.

##### `GetMemoryUsage() (*mem.VirtualMemoryStat, error)`
Returns detailed memory usage.

##### `GetSwapUsage() (*mem.SwapMemoryStat, error)`
Returns swap usage statistics.

##### `GetMemoryStats() (map[string]interface{}, error)`
Returns comprehensive memory statistics.

### Disk Monitoring

#### Functions

##### `GetDiskUsage(path string) (*disk.UsageStat, error)`
Returns disk usage for specified path.

**Parameters:**
- `path`: Mount point path

##### `GetAllPartitions() ([]disk.PartitionStat, error)`
Returns all disk partitions.

##### `GetDiskStats() (map[string]interface{}, error)`
Returns comprehensive disk statistics.

### Network Monitoring

#### Functions

##### `GetNetworkStats(pernic bool) ([]net.IOCountersStat, error)`
Returns network I/O statistics.

**Parameters:**
- `pernic`: If true, returns per-interface stats

##### `GetConnections() ([]net.ConnectionStat, error)`
Returns active network connections.

##### `GetNetworkInterfaces() ([]net.InterfaceStat, error)`
Returns network interface information.

##### `GetDetailedNetworkStats() (map[string]interface{}, error)`
Returns comprehensive network statistics.

## Logger API

### Logger Interface

```go
type Logger struct {
    file       *os.File
    logger     *log.Logger
    logPath    string
    maxSize    int64
    maxBackups int
    debug      bool
}
```

#### Methods

##### `New(logPath string, debug bool) (*Logger, error)`
Creates a new logger instance.

##### `Info(msg string)`
Logs an info message.

##### `Error(msg string)`
Logs an error message.

##### `Warn(msg string)`
Logs a warning message.

##### `Debug(msg string)`
Logs a debug message (only if debug is enabled).

##### `Close() error`
Closes the logger.

##### `Rotate() error`
Performs log rotation.

##### `RotateBySize(maxSize int64) error`
Rotates log if size exceeds limit.

##### `RotateByAge(maxAge int) error`
Removes log files older than specified days.

## PID File API

### PIDFile Interface

```go
type PIDFile struct {
    path string
    pid  int
}
```

#### Methods

##### `New(path string) (*PIDFile, error)`
Creates a new PID file.

**Behavior:**
- Checks if process already running
- Writes current process ID to file

##### `Read(path string) (int, error)`
Reads PID from file.

##### `Remove() error`
Removes the PID file.

##### `GetPID() int`
Returns stored PID.

##### `Path() string`
Returns path to PID file.

##### `IsProcessRunning(pid int) bool`
Checks if process with given PID is running.

## Usage Examples

### Starting the Daemon

```go
package main

import (
    "log"
    "github.com/mbahati3-wq/the-system-daemon/internal/daemon"
)

func main() {
    cfg, err := daemon.LoadConfig("/etc/mydaemon/mydaemon.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    d, err := daemon.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    if err := d.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Getting System Stats

```go
usage, err := monitor.GetCPUUsage()
if err != nil {
    log.Fatal(err)
}
log.Printf("CPU Usage: %.2f%%", usage)

mem, err := monitor.GetMemoryStats()
if err != nil {
    log.Fatal(err)
}
log.Printf("Memory: %v", mem)
```

### Logging

```go
logger, err := logger.New("/var/log/mydaemon.log", false)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

logger.Info("Daemon started")
logger.Warn("High CPU usage detected")
logger.Error("Failed to read config")
```
