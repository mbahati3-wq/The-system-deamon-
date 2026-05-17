#!/bin/bash

echo "Checking mydaemon status..."

# Check if systemd service exists
if command -v systemctl &> /dev/null; then
  echo ""
  echo "Systemd Service Status:"
  systemctl status mydaemon --no-pager || echo "Service not found"
fi

# Check if process is running
PIDFILE="/var/run/mydaemon.pid"
if [ -f "$PIDFILE" ]; then
  PID=$(cat "$PIDFILE")
  if ps -p "$PID" > /dev/null; then
    echo ""
    echo "Process Status: RUNNING"
    echo "PID: $PID"
    echo ""
    echo "Process Details:"
    ps aux | grep $PID | grep -v grep
  else
    echo ""
    echo "Process Status: NOT RUNNING (stale PID file)"
    rm -f "$PIDFILE"
  fi
else
  echo ""
  echo "Process Status: NOT RUNNING"
fi

# Check log file
LOGFILE="/var/log/mydaemon/mydaemon.log"
if [ -f "$LOGFILE" ]; then
  echo ""
  echo "Last 20 log entries:"
  tail -20 "$LOGFILE"
fi
