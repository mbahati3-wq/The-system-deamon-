#!/bin/bash
set -e

echo "Uninstalling mydaemon..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
  echo "This script must be run as root"
  exit 1
fi

# Stop daemon if running
if command -v systemctl &> /dev/null; then
  if systemctl is-active --quiet mydaemon; then
    echo "Stopping systemd service..."
    systemctl stop mydaemon
  fi
  echo "Removing systemd service..."
  rm -f /etc/systemd/system/mydaemon.service
  systemctl daemon-reload
else
  if service mydaemon status &>/dev/null; then
    echo "Stopping init service..."
    service mydaemon stop
  fi
  echo "Removing init script..."
  rm -f /etc/init.d/mydaemon
  update-rc.d -f mydaemon remove || true
fi

# Remove binary
echo "Removing binary..."
rm -f /usr/local/bin/mydaemon

# Remove config
echo "Removing configuration files..."
rm -f /etc/mydaemon/mydaemon.yaml

# Optional: Ask about removing logs
read -p "Remove log files? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  rm -rf /var/log/mydaemon
  rm -rf /var/run/mydaemon
fi

echo "Uninstallation complete!"
