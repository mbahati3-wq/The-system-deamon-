package test

import (
	"testing"

	"github.com/mbahati3-wq/the-system-daemon/internal/monitor"
)

// BenchmarkCPUUsage benchmarks CPU usage retrieval
func BenchmarkCPUUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := monitor.GetCPUUsage()
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// BenchmarkMemoryUsage benchmarks memory usage retrieval
func BenchmarkMemoryUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := monitor.GetMemoryUsage()
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// BenchmarkDiskUsage benchmarks disk usage retrieval
func BenchmarkDiskUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := monitor.GetDiskUsage("/")
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// BenchmarkNetworkStats benchmarks network stats retrieval
func BenchmarkNetworkStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := monitor.GetNetworkStats(false)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}
