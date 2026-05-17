.PHONY: build test clean install uninstall run debug

BINARY_NAME=mydaemon
BINARY_PATH=./build/$(BINARY_NAME)
GO=go
GOFLAGS=-v

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p build
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) ./cmd/mydaemon

test:
	@echo "Running tests..."
	$(GO) test -v ./...

bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./test

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/
	@$(GO) clean

install: build
	@echo "Installing daemon..."
	@bash scripts/install.sh

uninstall:
	@echo "Uninstalling daemon..."
	@bash scripts/uninstall.sh

run: build
	@echo "Running daemon..."
	$(BINARY_PATH)

debug: build
	@echo "Running daemon in debug mode..."
	@bash scripts/debug.sh

status:
	@echo "Checking daemon status..."
	@bash scripts/status.sh

docker-build:
	@echo "Building Docker image..."
	@docker build -f deployments/docker/Dockerfile -t mydaemon:latest .

docker-run: docker-build
	@echo "Running daemon in Docker..."
	@docker run -d --name mydaemon mydaemon:latest

help:
	@echo "Available targets:"
	@echo "  build           - Build the daemon binary"
	@echo "  test            - Run tests"
	@echo "  bench           - Run benchmarks"
	@echo "  clean           - Clean build artifacts"
	@echo "  install         - Install daemon to system"
	@echo "  uninstall       - Remove daemon from system"
	@echo "  run             - Run daemon locally"
	@echo "  debug           - Run daemon with debug output"
	@echo "  status          - Check daemon status"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run daemon in Docker"
