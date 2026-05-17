#!/bin/bash
set -e

echo "Running test suite..."

# Run unit tests
echo "Running unit tests..."
go test -v -race -coverprofile=coverage.out ./...

# Run integration tests
echo "Running integration tests..."
go test -v ./test/...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report: coverage.html"

# Run benchmarks
echo "Running benchmarks..."
go test -bench=. -benchmem ./test/... | tee benchmark.txt

echo "All tests completed!"
