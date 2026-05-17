# MyDaemon - System Monitoring Daemon

A lightweight, production-ready system monitoring daemon written in Go. MyDaemon provides comprehensive monitoring of system resources including CPU, memory, disk, and network statistics.

## Features

- 🖥️ **System Monitoring**: Real-time CPU, memory, disk, and network statistics
- 🔄 **Graceful Signal Handling**: SIGTERM, SIGINT, SIGHUP support for clean shutdown and config reload
- 📝 **Structured Logging**: Comprehensive logging with automatic rotation
- 🔒 **Single Instance**: PID file management ensures only one daemon instance runs
- 🐳 **Containerization**: Docker and Kubernetes deployment ready
- 📊 **Performance**: Minimal resource footprint with configurable monitoring intervals
- 🛠️ **Easy Installation**: Systemd service, SysV init, and manual installation options
- 📚 **Documentation**: Comprehensive architecture, API, and debugging guides

## Quick Start

### Prerequisites

- Go 1.21 or later
- Linux operating system
- Root/sudo access for system-wide installation

### Building

```bash
# Clone the repository
git clone https://github.com/mbahati3-wq/The-system-deamon-.git
cd The-system-deamon-

# Build the daemon
make build

# Run in debug mode
make debug

# Run tests
make test
```

### Installation (Linux)

```bash
# Install as system service
sudo make install

# Start the daemon
sudo systemctl start mydaemon

# Enable auto-start on boot
sudo systemctl enable mydaemon

# Check status
sudo systemctl status mydaemon
```

### Development

```bash
# Build binary
make build

# Run unit tests
make test

# Run benchmarks
make bench

# Clean build artifacts
make clean
```

## Project Structure

```
├── cmd/mydaemon/           # Main application entry point
├── internal/
│   ├── daemon/            # Core daemon logic
│   ├── monitor/           # System monitoring modules
│   └── logger/            # Logging functionality
├── pkg/pidfile/           # PID file management
├── configs/               # Configuration files
├── deployments/           # Docker & Kubernetes configs
├── scripts/               # Installation & management scripts
├── test/                  # Integration & benchmark tests
└── docs/                  # Documentation
```

## Configuration

The daemon is configured via YAML file at `/etc/mydaemon/mydaemon.yaml`:

```yaml
app_name: mydaemon
debug: false
log_path: /var/log/mydaemon/mydaemon.log
max_log_size: 100  # MB
max_log_backups: 10

# Monitoring intervals in seconds
cpu_monitor_interval: 5
memory_monitor_interval: 5
disk_monitor_interval: 10
network_monitor_interval: 10

# Alert thresholds (percentages)
cpu_threshold: 80.0
memory_threshold: 80.0
```

## Usage

### Command Line Options

```bash
mydaemon [flags]

Flags:
  -config string     Path to configuration file (default "/etc/mydaemon/mydaemon.yaml")
  -pidfile string    Path to PID file (default "/var/run/mydaemon.pid")
  -debug             Enable debug mode
  -h, -help          Show help message
```

### Systemd Management

```bash
# Start daemon
sudo systemctl start mydaemon

# Stop daemon
sudo systemctl stop mydaemon

# Restart daemon
sudo systemctl restart mydaemon

# Reload configuration
sudo systemctl reload mydaemon

# View logs
sudo journalctl -u mydaemon -f

# Check status
sudo systemctl status mydaemon
```

### Manual Control

```bash
# Start daemon
sudo ./build/mydaemon -config configs/mydaemon.yaml

# Check if running
bash scripts/status.sh

# View recent logs
tail -f /var/log/mydaemon/mydaemon.log
```

## Signal Handling

The daemon responds to the following signals:

| Signal | Action |
|--------|--------|
| SIGTERM | Graceful shutdown |
| SIGINT | Immediate shutdown |
| SIGHUP | Reload configuration |
| SIGUSR1 | Print statistics |

```bash
# Send signals to daemon
kill -TERM $(cat /var/run/mydaemon.pid)    # Graceful shutdown
kill -HUP $(cat /var/run/mydaemon.pid)     # Reload config
kill -USR1 $(cat /var/run/mydaemon.pid)    # Print stats
```

## Deployment

### Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run

# Or manually
docker build -f deployments/docker/Dockerfile -t mydaemon:latest .
docker run -d --name mydaemon mydaemon:latest
```

### Kubernetes

```bash
# Deploy to Kubernetes
kubectl apply -f deployments/kubernetes/deployment.yaml

# Check deployment
kubectl get pods -l app=mydaemon
kubectl logs -f pod/mydaemon-<pod-id>
```

## Monitoring

### Check Daemon Status

```bash
bash scripts/status.sh
```

Output includes:
- Service status
- Process PID
- Recent log entries

### View Logs

```bash
# Live logs
tail -f /var/log/mydaemon/mydaemon.log

# Search for errors
grep ERROR /var/log/mydaemon/mydaemon.log

# Last 100 entries
tail -100 /var/log/mydaemon/mydaemon.log
```

### Performance Monitoring

```bash
# Monitor resource usage
top -p $(pidof mydaemon)

# Check memory usage
ps -eo pid,vsz,rss,comm | grep mydaemon

# Monitor CPU usage
ps aux | grep mydaemon
```

## Testing

### Unit Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Benchmarks

```bash
# Run benchmarks
make bench

# Specific benchmark
go test -bench=BenchmarkCPUUsage -benchmem ./test
```

### Integration Tests

```bash
bash scripts/test.sh
```

## Troubleshooting

### Daemon Won't Start

**Check permissions:**
```bash
ls -la /var/run/mydaemon
chmod 755 /var/run/mydaemon
```

**Check for orphaned process:**
```bash
cat /var/run/mydaemon.pid
ps -p <pid>
rm /var/run/mydaemon.pid  # Only if process not running
```

### High Memory Usage

**Check log size:**
```bash
du -h /var/log/mydaemon/mydaemon.log
```

**Rotate logs manually:**
```bash
kill -HUP $(cat /var/run/mydaemon.pid)
```

### Configuration Issues

**Verify config file:**
```bash
ls -la /etc/mydaemon/mydaemon.yaml
cat /etc/mydaemon/mydaemon.yaml
```

**Reload configuration:**
```bash
systemctl reload mydaemon
```

See [debugging.md](docs/debugging.md) for comprehensive troubleshooting guide.

## Documentation

- [Architecture](docs/architecture.md) - System design and components
- [API Reference](docs/api.md) - Internal API documentation
- [Debugging Guide](docs/debugging.md) - Troubleshooting and diagnostics

## Uninstallation

```bash
# Remove daemon from system
sudo make uninstall

# Or manually
sudo bash scripts/uninstall.sh
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or suggestions:

1. Check the [debugging guide](docs/debugging.md)
2. Review [existing issues](https://github.com/mbahati3-wq/The-system-deamon-/issues)
3. Create a new issue with details

## Changelog

### Version 1.0.0 (Initial Release)

- Initial project structure and implementation
- Core daemon functionality
- System monitoring (CPU, memory, disk, network)
- Structured logging with rotation
- Docker and Kubernetes support
- Comprehensive documentation
