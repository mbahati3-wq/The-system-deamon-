#!/bin/bash
set -e

echo "Installing mydaemon..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
  echo "This script must be run as root"
  exit 1
fi

# Build the daemon
echo "Building daemon..."
make build

# Create user and group
echo "Creating daemon user and group..."
if ! id "daemon" &>/dev/null; then
  useradd -r -s /bin/false daemon 2>/dev/null || true
fi

# Install binary
echo "Installing binary to /usr/local/bin..."
install -m 0755 build/mydaemon /usr/local/bin/

# Create config directory
echo "Creating config directory..."
mkdir -p /etc/mydaemon
mkdir -p /var/log/mydaemon
mkdir -p /var/run/mydaemon

# Install config files
echo "Installing configuration files..."
install -m 0644 configs/mydaemon.yaml /etc/mydaemon/mydaemon.yaml.default || true

# Set permissions
chown -R daemon:daemon /var/log/mydaemon
chown -R daemon:daemon /var/run/mydaemon
chmod 755 /etc/mydaemon

# Install systemd service (if systemd is available)
if command -v systemctl &> /dev/null; then
  echo "Installing systemd service..."
  install -m 0644 configs/mydaemon.service /etc/systemd/system/
  systemctl daemon-reload
  echo "To start the daemon, run: systemctl start mydaemon"
  echo "To enable at boot, run: systemctl enable mydaemon"
else
  echo "Installing SysV init script..."
  install -m 0755 configs/mydaemon.init /etc/init.d/mydaemon
  echo "To start the daemon, run: service mydaemon start"
  echo "To enable at boot, run: update-rc.d mydaemon defaults"
fi

echo "Installation complete!"
