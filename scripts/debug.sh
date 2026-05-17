#!/bin/bash

echo "Running mydaemon in debug mode..."

# Build if not exists
if [ ! -f "build/mydaemon" ]; then
  echo "Building daemon first..."
  make build
fi

# Create log directory
mkdir -p logs

# Run in debug mode with output to console
./build/mydaemon -config configs/mydaemon.yaml -debug -pidfile /tmp/mydaemon.pid

