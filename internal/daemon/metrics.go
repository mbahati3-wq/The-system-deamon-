package daemon

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Metrics holds system metrics
type Metrics struct {
	CPUUsage           float64
	MemoryUsage        float64
	DiskUsage          float64
	NetworkBytesRecv   uint64
	NetworkBytesSent   uint64
	ProcessCount       int
	Uptime             time.Duration
	AlertsTriggered    int64
	AlertsResolved     int64
	LastUpdateTime     time.Time
}

// MetricsCollector collects and stores metrics
type MetricsCollector struct {
	metrics    *Metrics
	mu         sync.RWMutex
	server     *http.Server
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(port int) *MetricsCollector {
	mc := &MetricsCollector{
		metrics: &Metrics{
			LastUpdateTime: time.Now(),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", mc.handleMetrics)

	mc.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return mc
}

// Start starts the metrics server
func (mc *MetricsCollector) Start() error {
	go func() {
		if err := mc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Metrics server error: %v\n", err)
		}
	}()
	return nil
}

// Stop stops the metrics server
func (mc *MetricsCollector) Stop() error {
	if mc.server != nil {
		return mc.server.Close()
	}
	return nil
}

// UpdateCPUUsage updates CPU usage metric
func (mc *MetricsCollector) UpdateCPUUsage(usage float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.CPUUsage = usage
	mc.metrics.LastUpdateTime = time.Now()
}

// UpdateMemoryUsage updates memory usage metric
func (mc *MetricsCollector) UpdateMemoryUsage(usage float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.MemoryUsage = usage
	mc.metrics.LastUpdateTime = time.Now()
}

// UpdateDiskUsage updates disk usage metric
func (mc *MetricsCollector) UpdateDiskUsage(usage float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.DiskUsage = usage
	mc.metrics.LastUpdateTime = time.Now()
}

// UpdateNetworkMetrics updates network metrics
func (mc *MetricsCollector) UpdateNetworkMetrics(bytesRecv, bytesSent uint64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.NetworkBytesRecv = bytesRecv
	mc.metrics.NetworkBytesSent = bytesSent
	mc.metrics.LastUpdateTime = time.Now()
}

// UpdateProcessCount updates process count metric
func (mc *MetricsCollector) UpdateProcessCount(count int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.ProcessCount = count
	mc.metrics.LastUpdateTime = time.Now()
}

// UpdateUptime updates uptime metric
func (mc *MetricsCollector) UpdateUptime(uptime time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.Uptime = uptime
	mc.metrics.LastUpdateTime = time.Now()
}

// IncrementAlerts increments alerts triggered counter
func (mc *MetricsCollector) IncrementAlerts() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.AlertsTriggered++
	mc.metrics.LastUpdateTime = time.Now()
}

// IncrementResolved increments alerts resolved counter
func (mc *MetricsCollector) IncrementResolved() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics.AlertsResolved++
	mc.metrics.LastUpdateTime = time.Now()
}

// GetMetrics returns a copy of current metrics
func (mc *MetricsCollector) GetMetrics() *Metrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	m := &Metrics{
		CPUUsage:         mc.metrics.CPUUsage,
		MemoryUsage:      mc.metrics.MemoryUsage,
		DiskUsage:        mc.metrics.DiskUsage,
		NetworkBytesRecv: mc.metrics.NetworkBytesRecv,
		NetworkBytesSent: mc.metrics.NetworkBytesSent,
		ProcessCount:     mc.metrics.ProcessCount,
		Uptime:           mc.metrics.Uptime,
		AlertsTriggered:  mc.metrics.AlertsTriggered,
		AlertsResolved:   mc.metrics.AlertsResolved,
		LastUpdateTime:   mc.metrics.LastUpdateTime,
	}
	return m
}

// handleMetrics handles the /metrics endpoint
func (mc *MetricsCollector) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	m := mc.GetMetrics()
	output := formatPrometheus(m)
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, output)
}

// formatPrometheus formats metrics in Prometheus text format
func formatPrometheus(m *Metrics) string {
	output := "# HELP daemon_cpu_usage_percent CPU usage percentage\n"
	output += "# TYPE daemon_cpu_usage_percent gauge\n"
	output += fmt.Sprintf("daemon_cpu_usage_percent %.2f\n\n", m.CPUUsage)

	output += "# HELP daemon_memory_usage_percent Memory usage percentage\n"
	output += "# TYPE daemon_memory_usage_percent gauge\n"
	output += fmt.Sprintf("daemon_memory_usage_percent %.2f\n\n", m.MemoryUsage)

	output += "# HELP daemon_disk_usage_percent Disk usage percentage\n"
	output += "# TYPE daemon_disk_usage_percent gauge\n"
	output += fmt.Sprintf("daemon_disk_usage_percent %.2f\n\n", m.DiskUsage)

	output += "# HELP daemon_network_bytes_received Total network bytes received\n"
	output += "# TYPE daemon_network_bytes_received counter\n"
	output += fmt.Sprintf("daemon_network_bytes_received %d\n\n", m.NetworkBytesRecv)

	output += "# HELP daemon_network_bytes_sent Total network bytes sent\n"
	output += "# TYPE daemon_network_bytes_sent counter\n"
	output += fmt.Sprintf("daemon_network_bytes_sent %d\n\n", m.NetworkBytesSent)

	output += "# HELP daemon_process_count Total number of processes\n"
	output += "# TYPE daemon_process_count gauge\n"
	output += fmt.Sprintf("daemon_process_count %d\n\n", m.ProcessCount)

	output += "# HELP daemon_uptime_seconds Daemon uptime in seconds\n"
	output += "# TYPE daemon_uptime_seconds gauge\n"
	output += fmt.Sprintf("daemon_uptime_seconds %.2f\n\n", m.Uptime.Seconds())

	output += "# HELP daemon_alerts_triggered Total alerts triggered\n"
	output += "# TYPE daemon_alerts_triggered counter\n"
	output += fmt.Sprintf("daemon_alerts_triggered %d\n\n", m.AlertsTriggered)

	output += "# HELP daemon_alerts_resolved Total alerts resolved\n"
	output += "# TYPE daemon_alerts_resolved counter\n"
	output += fmt.Sprintf("daemon_alerts_resolved %d\n", m.AlertsResolved)

	return output
}
