package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthStatus represents the health status of the daemon
type HealthStatus struct {
	Status        string                 `json:"status"`
	Uptime        string                 `json:"uptime"`
	StartTime     time.Time              `json:"start_time"`
	Checks        map[string]CheckStatus `json:"checks"`
	IsReady       bool                   `json:"is_ready"`
	IsAlive       bool                   `json:"is_alive"`
}

// CheckStatus represents the status of a health check
type CheckStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthServer manages health check endpoints
type HealthServer struct {
	server      *http.Server
	startTime   time.Time
	checks      map[string]CheckFunc
	mu          sync.RWMutex
	isReady     bool
	isAlive     bool
}

// CheckFunc is a function that performs a health check
type CheckFunc func() CheckStatus

// NewHealthServer creates a new health server
func NewHealthServer(port int) *HealthServer {
	hs := &HealthServer{
		startTime: time.Now(),
		checks:    make(map[string]CheckFunc),
		isReady:   false,
		isAlive:   true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", hs.handleHealth)
	mux.HandleFunc("/ready", hs.handleReady)
	mux.HandleFunc("/live", hs.handleLive)

	hs.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return hs
}

// Start starts the health server
func (hs *HealthServer) Start() error {
	go func() {
		if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Health server error: %v\n", err)
		}
	}()
	return nil
}

// Stop stops the health server
func (hs *HealthServer) Stop() error {
	if hs.server != nil {
		return hs.server.Close()
	}
	return nil
}

// SetReady marks the daemon as ready
func (hs *HealthServer) SetReady(ready bool) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.isReady = ready
}

// SetAlive marks the daemon as alive
func (hs *HealthServer) SetAlive(alive bool) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.isAlive = alive
}

// RegisterCheck registers a new health check
func (hs *HealthServer) RegisterCheck(name string, check CheckFunc) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.checks[name] = check
}

// handleHealth handles the /health endpoint
func (hs *HealthServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hs.mu.RLock()
	uptime := time.Since(hs.startTime)
	isReady := hs.isReady
	isAlive := hs.isAlive

	checks := make(map[string]CheckStatus)
	for name, check := range hs.checks {
		checks[name] = check()
	}
	hs.mu.RUnlock()

	status := "unhealthy"
	if isAlive && isReady {
		status = "healthy"
	} else if isAlive {
		status = "starting"
	}

	health := HealthStatus{
		Status:    status,
		Uptime:    uptime.String(),
		StartTime: hs.startTime,
		Checks:    checks,
		IsReady:   isReady,
		IsAlive:   isAlive,
	}

	w.Header().Set("Content-Type", "application/json")

	if !isAlive || !isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(health)
}

// handleReady handles the /ready endpoint (Kubernetes readiness probe)
func (hs *HealthServer) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hs.mu.RLock()
	isReady := hs.isReady
	hs.mu.RUnlock()

	response := map[string]bool{"ready": isReady}
	w.Header().Set("Content-Type", "application/json")

	if !isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

// handleLive handles the /live endpoint (Kubernetes liveness probe)
func (hs *HealthServer) handleLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hs.mu.RLock()
	isAlive := hs.isAlive
	hs.mu.RUnlock()

	response := map[string]bool{"alive": isAlive}
	w.Header().Set("Content-Type", "application/json")

	if !isAlive {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
