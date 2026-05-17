# Features Implementation Guide

## Features Implemented

This document describes all the features implemented in MyDaemon and how to use them.

### 1. Configuration Validation

**Status**: ✅ Complete

The daemon now validates all configuration parameters on startup. Invalid configurations will prevent the daemon from starting.

**Validation Rules**:
- `app_name`: Cannot be empty
- `log_path`: Cannot be empty, directory must be writable
- `max_log_size`: Must be >= 1 MB
- `max_log_backups`: Must be >= 0
- Monitor intervals: Must be >= 1 second
- Thresholds: Must be between 0 and 100
- Port numbers: Must be between 1 and 65535

**Usage**:
```bash
# Validate configuration without starting daemon
./build/mydaemon -validate-config -config configs/mydaemon.yaml

# Start with automatic validation
./build/mydaemon -config configs/mydaemon.yaml
```

### 2. Dynamic Configuration Reload

**Status**: ✅ Complete

Configuration can be reloaded without restarting the daemon using SIGHUP signal.

**Usage**:
```bash
# Send SIGHUP to reload configuration
kill -HUP $(cat /var/run/mydaemon.pid)

# Or with systemctl
systemctl reload mydaemon
```

**Log Output**:
```
[INFO] Received signal: hangup
[INFO] Reloading configuration...
[INFO] Configuration reloaded successfully
```

### 3. Process Monitoring

**Status**: ✅ Complete

Monitor specific system processes for availability and resource usage.

**Configuration**:
```yaml
processes:
  - nginx
  - redis-server
  - postgresql
```

**API Functions**:
```go
// Find processes by name
procs, err := monitor.FindProcessesByName("nginx")

// Get specific process by PID
proc, err := monitor.GetProcessByPID(1234)

// Get all processes
procs, err := monitor.GetAllProcesses()

// Check if process is running
running := monitor.IsProcessRunning("nginx")

// Count processes with given name
count, err := monitor.GetProcessCount("nginx")
```

**Alert**: Automatically triggers CRITICAL alert if monitored process stops running.

### 4. HTTP Health Endpoint

**Status**: ✅ Complete

Provides HTTP endpoints for health checking, used by Kubernetes, Docker, and monitoring systems.

**Configuration**:
```yaml
health_check_port: 8080
```

**Endpoints**:

#### `/health` - Full Health Status
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "uptime": "1h23m45s",
  "start_time": "2024-05-17T10:30:00Z",
  "checks": {
    "monitor": {
      "status": "ok",
      "message": "Monitor is running"
    }
  },
  "is_ready": true,
  "is_alive": true
}
```

#### `/ready` - Readiness Probe (Kubernetes)
```bash
curl http://localhost:8080/ready
```

Returns 200 when daemon is ready to serve requests, 503 otherwise.

#### `/live` - Liveness Probe (Kubernetes)
```bash
curl http://localhost:8080/live
```

Returns 200 when daemon is alive, 503 if critical failure detected.

**Kubernetes Configuration**:
```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### 5. Alerting System

**Status**: ✅ Complete

Comprehensive alert management with automatic threshold checking.

**Configuration**:
```yaml
cpu_threshold: 80.0
memory_threshold: 80.0
```

**Alert Types**:
- `INFO`: Informational alerts
- `WARNING`: Warning alerts (threshold exceeded)
- `CRITICAL`: Critical alerts (process missing, critical failure)

**Automatic Alerts**:
- CPU usage exceeds threshold → WARNING
- Memory usage exceeds threshold → WARNING
- Monitored process stops running → CRITICAL
- Disk usage exceeds threshold → WARNING

**API Usage**:
```go
// Get all alerts
alerts := daemon.alertManager.GetAlerts()

// Get alerts by level
warnings := daemon.alertManager.GetAlertsByLevel(daemon.AlertLevelWarning)

// Clear specific alert
daemon.alertManager.ClearAlert(alertID)

// Clear all alerts of a level
daemon.alertManager.ClearAlertsByLevel(daemon.AlertLevelCritical)

// Trigger custom alert
daemon.alertManager.TriggerAlert(
    daemon.AlertLevelWarning,
    "High Temperature",
    "CPU temperature is 95°C",
    "temperature",
    95.0,
    90.0,
)
```

### 6. Prometheus Metrics Export

**Status**: ✅ Complete

Export system metrics in Prometheus text format for monitoring and alerting.

**Configuration**:
```yaml
enable_metrics: true
metrics_port: 9090
```

**Metrics Endpoint**: `http://localhost:9090/metrics`

**Available Metrics**:
- `daemon_cpu_usage_percent` - CPU usage percentage
- `daemon_memory_usage_percent` - Memory usage percentage
- `daemon_disk_usage_percent` - Disk usage percentage
- `daemon_network_bytes_received` - Network bytes received
- `daemon_network_bytes_sent` - Network bytes sent
- `daemon_process_count` - Number of running processes
- `daemon_uptime_seconds` - Daemon uptime in seconds
- `daemon_alerts_triggered` - Total alerts triggered
- `daemon_alerts_resolved` - Total alerts resolved

**Prometheus Configuration Example**:
```yaml
scrape_configs:
  - job_name: 'mydaemon'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 15s
```

**Example Output**:
```
# HELP daemon_cpu_usage_percent CPU usage percentage
# TYPE daemon_cpu_usage_percent gauge
daemon_cpu_usage_percent 45.32

# HELP daemon_memory_usage_percent Memory usage percentage
# TYPE daemon_memory_usage_percent gauge
daemon_memory_usage_percent 62.15

# HELP daemon_uptime_seconds Daemon uptime in seconds
# TYPE daemon_uptime_seconds gauge
daemon_uptime_seconds 3665.42
```

### 7. Complete Log Rotation

**Status**: ✅ Partial (Enhanced rotate.go)

Automatic log file rotation with size, age, and backup count management.

**Configuration**:
```yaml
max_log_size: 100  # MB
max_log_backups: 10
```

**API Usage**:
```go
// Rotate if size exceeded
logger.RotateBySize(100 * 1024 * 1024)

// Rotate by age (days)
logger.RotateByAge(30)

// Clean up old logs
logger.CleanupOldLogs(10)
```

**Automatic Rotation**:
- Triggered when log file exceeds `max_log_size`
- Old logs named with timestamp: `mydaemon.log.2024-05-17_14-30-45`
- Only `max_log_backups` most recent backups kept

## Configuration Example with All Features

```yaml
app_name: mydaemon
debug: false
log_path: /var/log/mydaemon/mydaemon.log
max_log_size: 100
max_log_backups: 10

# Monitoring
cpu_monitor_interval: 5
memory_monitor_interval: 5
disk_monitor_interval: 10
network_monitor_interval: 10

# Thresholds
cpu_threshold: 80.0
memory_threshold: 85.0

# Health checks
health_check_port: 8080

# Prometheus metrics
enable_metrics: true
metrics_port: 9090

# Process monitoring
processes:
  - nginx
  - redis-server
  - postgresql
```

## Testing Features

### Test Configuration Validation
```bash
# Valid config
./build/mydaemon -validate-config

# Invalid config (will show validation errors)
sed 's/max_log_size: 100/max_log_size: 0/' configs/mydaemon.yaml > /tmp/bad.yaml
./build/mydaemon -validate-config -config /tmp/bad.yaml
```

### Test Health Endpoints
```bash
# In separate terminal, run daemon
./build/mydaemon -debug -config configs/mydaemon.yaml

# Test health endpoint
curl -v http://localhost:8080/health
curl -v http://localhost:8080/ready
curl -v http://localhost:8080/live

# Test metrics (if enabled)
curl http://localhost:9090/metrics | grep daemon_
```

### Test Dynamic Reload
```bash
# Modify config (increase threshold)
sed -i 's/cpu_threshold: 80.0/cpu_threshold: 70.0/' configs/mydaemon.yaml

# Trigger reload
kill -HUP $(cat /var/run/mydaemon.pid)

# Check logs
tail -f /var/log/mydaemon/mydaemon.log
```

### Test Process Monitoring
```yaml
# Add to config
processes:
  - bash
```

Then check if alerts trigger when process stops.

## Integration with Docker and Kubernetes

### Docker Health Check
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

### Kubernetes Deployment
See `deployments/kubernetes/deployment.yaml` for complete example with health probes.

## Command-Line Enhancements

```bash
# Validate config and exit
./build/mydaemon -validate-config

# Run with debug output
./build/mydaemon -debug

# Custom config path
./build/mydaemon -config /etc/custom/config.yaml

# Custom PID file location
./build/mydaemon -pidfile /tmp/mydaemon.pid
```

## Performance Impact

- **Memory**: < 50 MB typical usage
- **CPU**: < 2% when idle
- **Network**: Minimal impact
- **Disk I/O**: Only on log rotation

## Security Considerations

- Health/metrics endpoints listen on localhost by default
- Configuration files should be readable only by daemon user
- Log files contain sensitive information
- PID file should be in /var/run with proper permissions

## Troubleshooting

### Configuration Won't Validate
```bash
# Check validation errors
./build/mydaemon -validate-config -debug

# Fix issues and try again
```

### Health Endpoint Not Responding
```bash
# Check if daemon is running
ps aux | grep mydaemon

# Verify port is not in use
lsof -i :8080

# Check logs
tail -f /var/log/mydaemon/mydaemon.log
```

### Metrics Not Updating
```bash
# Enable metrics in config
sed -i 's/enable_metrics: false/enable_metrics: true/' configs/mydaemon.yaml

# Reload daemon
kill -HUP $(cat /var/run/mydaemon.pid)

# Check metrics endpoint
curl http://localhost:9090/metrics
```
