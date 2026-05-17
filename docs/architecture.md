# MyDaemon Architecture

## Overview

MyDaemon is a lightweight system monitoring daemon written in Go. It provides comprehensive system resource monitoring including CPU, memory, disk, and network statistics.

## Architecture Diagram

```
┌─────────────────────────────────────────┐
│          Main Application               │
│         (cmd/mydaemon/main.go)          │
└────────────────┬────────────────────────┘
                 │
    ┌────────────┴────────────┐
    │                         │
    v                         v
┌─────────────┐       ┌──────────────┐
│   Daemon    │       │   PID File   │
│  (daemon/)  │       │  Management  │
└─────────────┘       └──────────────┘
    │
    ├─────────────────────┬─────────────────────┐
    │                     │                     │
    v                     v                     v
┌────────────┐    ┌──────────────┐    ┌──────────────┐
│  Monitor   │    │    Logger    │    │   Config     │
│ (monitor/) │    │  (logger/)   │    │  Management  │
└────────────┘    └──────────────┘    └──────────────┘
    │
    ├────────────┬─────────────┬──────────────┐
    │            │             │              │
    v            v             v              v
┌──────────┐ ┌────────┐ ┌──────────┐ ┌──────────────┐
│ CPU Mon  │ │Memory  │ │ Disk     │ │ Network      │
│          │ │Monitor │ │ Monitor  │ │ Monitor      │
└──────────┘ └────────┘ └──────────┘ └──────────────┘
```

## Components

### 1. **Daemon (`internal/daemon/`)**
- Core daemon logic and lifecycle management
- Signal handling (SIGTERM, SIGINT, SIGHUP)
- Configuration management
- Graceful shutdown

### 2. **Monitor (`internal/monitor/`)**
- Unified monitoring interface
- CPU monitoring
- Memory monitoring
- Disk monitoring
- Network monitoring
- Configurable monitoring intervals

### 3. **Logger (`internal/logger/`)**
- Structured logging with timestamps
- File-based logging with rotation
- Debug mode support
- Log size management

### 4. **Configuration**
- YAML-based configuration files
- Default configuration values
- Runtime configuration management

### 5. **PID File Management (`pkg/pidfile/`)**
- Prevents multiple instances
- Process tracking
- Clean process removal

## Data Flow

```
┌─────────┐
│ Daemon  │
└────┬────┘
     │
     ├──► Monitor CPU ─────────────────┐
     │                                 │
     ├──► Monitor Memory ──────────────┤
     │                                 │
     ├──► Monitor Disk ────────────────┤──► Logger ──► Log File
     │                                 │
     └──► Monitor Network ─────────────┘
```

## Configuration

Configuration is managed through `configs/mydaemon.yaml`:

```yaml
app_name: mydaemon
debug: false
log_path: /var/log/mydaemon/mydaemon.log
cpu_monitor_interval: 5
memory_monitor_interval: 5
cpu_threshold: 80.0
memory_threshold: 80.0
```

## Deployment Options

### 1. **Systemd Service**
- Primary method on modern Linux distributions
- Service file: `configs/mydaemon.service`
- Automatic restart on failure

### 2. **SysV Init**
- Fallback for older systems
- Script: `configs/mydaemon.init`

### 3. **Docker**
- Containerized deployment
- Multi-stage build for minimal image size

### 4. **Kubernetes**
- Deployment manifest with ConfigMap
- Resource limits and health checks
- Security context configuration

## Signal Handling

- **SIGTERM/SIGINT**: Graceful shutdown
- **SIGHUP**: Configuration reload
- **SIGUSR1**: Statistics reporting

## Error Handling

- Graceful error handling throughout
- Non-blocking error propagation
- Comprehensive logging

## Security

- Non-root user execution
- Read-only root filesystem in containers
- Capability dropping
- Privilege escalation prevention
