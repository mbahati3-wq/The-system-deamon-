# MyDaemon Debugging Guide

## Running in Debug Mode

### Command Line

Run the daemon with debug mode enabled:

```bash
./build/mydaemon -debug -config configs/mydaemon.yaml
```

Or use the provided script:

```bash
bash scripts/debug.sh
```

### Debug Output

In debug mode, the daemon will:
- Output all logs to stderr and log file
- Include detailed debugging information
- Show signal handling events
- Display monitoring data

## Checking Daemon Status

### Using the Status Script

```bash
bash scripts/status.sh
```

This will show:
- Systemd service status (if available)
- Process status and PID
- Recent log entries

### Manual Status Check

Check if the daemon is running:

```bash
ps aux | grep mydaemon
```

Check if the service is active (systemd):

```bash
systemctl status mydaemon
```

## Viewing Logs

### Log Location

Production: `/var/log/mydaemon/mydaemon.log`
Development: `./logs/mydaemon.log`

### View Recent Logs

```bash
tail -f /var/log/mydaemon/mydaemon.log
```

### View Last 100 Entries

```bash
tail -100 /var/log/mydaemon/mydaemon.log
```

### Search for Errors

```bash
grep ERROR /var/log/mydaemon/mydaemon.log
```

## Common Issues

### Daemon Won't Start

**Problem:** `failed to create PID file`

**Solution:** Check permissions on /var/run directory:
```bash
ls -la /var/run/ | grep mydaemon
chmod 755 /var/run/mydaemon
```

**Problem:** `process already running`

**Solution:** Check for orphaned PID file:
```bash
cat /var/run/mydaemon.pid
ps -p <pid>
rm /var/run/mydaemon.pid  # Only if process not running
```

### High Memory Usage

**Problem:** Daemon consuming excessive memory

**Solution:** Check log file size:
```bash
du -h /var/log/mydaemon/mydaemon.log
```

Manually rotate logs:
```bash
kill -HUP $(cat /var/run/mydaemon.pid)
```

### Configuration Not Loading

**Problem:** Custom config ignored

**Solution:** Verify config file permissions:
```bash
ls -la /etc/mydaemon/mydaemon.yaml
cat /etc/mydaemon/mydaemon.yaml
```

Reload configuration:
```bash
systemctl reload mydaemon
```

## Performance Debugging

### CPU Usage

Monitor CPU consumption:

```bash
top -p $(cat /var/run/mydaemon.pid)
```

View CPU profiling data:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

### Memory Usage

Monitor memory consumption:

```bash
ps -eo pid,vsz,rss,comm | grep mydaemon
```

View memory profiling data:

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

## Testing

### Run Unit Tests

```bash
make test
```

### Run Integration Tests

```bash
bash scripts/test.sh
```

### Run Benchmarks

```bash
make bench
```

## Signal Handling

### Graceful Shutdown

```bash
kill -TERM $(cat /var/run/mydaemon.pid)
```

Or using systemctl:

```bash
systemctl stop mydaemon
```

### Reload Configuration

```bash
kill -HUP $(cat /var/run/mydaemon.pid)
```

Or using systemctl:

```bash
systemctl reload mydaemon
```

### Print Statistics

```bash
kill -USR1 $(cat /var/run/mydaemon.pid)
```

## Docker Debugging

### Run in Docker with Debug Output

```bash
docker run -it mydaemon:latest mydaemon -debug
```

### View Docker Logs

```bash
docker logs -f <container_id>
```

### Execute Commands in Container

```bash
docker exec -it <container_id> sh
```

## Kubernetes Debugging

### View Pod Logs

```bash
kubectl logs -f pod/mydaemon-<pod_id>
```

### Describe Pod

```bash
kubectl describe pod mydaemon-<pod_id>
```

### Port Forward for Debugging

```bash
kubectl port-forward pod/mydaemon-<pod_id> 6060:6060
```

### View Events

```bash
kubectl events --for=pod/mydaemon-<pod_id>
```

## Environment Variables

Set for debugging:

```bash
export GODEBUG=gctrace=1  # GC tracing
export GODEBUG=madvdontneed=0  # Memory debugging
export GOMAXPROCS=2  # CPU debugging
```

## Useful Commands

### Build for Debugging

```bash
go build -gcflags="all=-N -l" -o build/mydaemon ./cmd/mydaemon
```

### Run with Debugger (Delve)

```bash
dlv exec ./build/mydaemon -- -config configs/mydaemon.yaml
```

### Generate Core Dump on Crash

```bash
ulimit -c unlimited
./build/mydaemon
```

Analyze core dump:

```bash
gdb ./build/mydaemon ./core
```

## Troubleshooting Checklist

- [ ] Check daemon is running: `systemctl status mydaemon`
- [ ] Check logs for errors: `tail /var/log/mydaemon/mydaemon.log`
- [ ] Verify configuration: `cat /etc/mydaemon/mydaemon.yaml`
- [ ] Check permissions: `ls -la /var/log/mydaemon`
- [ ] Check disk space: `df -h /var/log`
- [ ] Test manually: `./build/mydaemon -debug`
- [ ] Check systemd status: `journalctl -u mydaemon -50`
- [ ] Monitor resource usage: `top -p $(pidof mydaemon)`
