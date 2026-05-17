package daemon

import (
	"fmt"
	"sync"
	"time"

	"github.com/mbahati3-wq/the-system-daemon/internal/logger"
)

// Alert represents a system alert
type Alert struct {
	ID        string
	Timestamp time.Time
	Level     AlertLevel
	Title     string
	Message   string
	Metric    string
	Value     float64
	Threshold float64
}

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	AlertLevelInfo    AlertLevel = "INFO"
	AlertLevelWarning AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
)

// AlertHandler is called when an alert is triggered
type AlertHandler func(alert *Alert) error

// AlertManager manages system alerts
type AlertManager struct {
	logger    *logger.Logger
	alerts    map[string]*Alert
	handlers  map[AlertLevel][]AlertHandler
	mu        sync.RWMutex
	alertChan chan *Alert
	stopChan  chan struct{}
}

// NewAlertManager creates a new alert manager
func NewAlertManager(log *logger.Logger) *AlertManager {
	return &AlertManager{
		logger:    log,
		alerts:    make(map[string]*Alert),
		handlers:  make(map[AlertLevel][]AlertHandler),
		alertChan: make(chan *Alert, 100),
		stopChan:  make(chan struct{}),
	}
}

// Start starts processing alerts
func (am *AlertManager) Start() {
	go func() {
		for {
			select {
			case alert := <-am.alertChan:
				am.processAlert(alert)
			case <-am.stopChan:
				return
			}
		}
	}()
}

// Stop stops the alert manager
func (am *AlertManager) Stop() {
	close(am.stopChan)
}

// RegisterHandler registers an alert handler for a specific level
func (am *AlertManager) RegisterHandler(level AlertLevel, handler AlertHandler) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.handlers[level] = append(am.handlers[level], handler)
}

// TriggerAlert triggers a new alert
func (am *AlertManager) TriggerAlert(level AlertLevel, title, message, metric string, value, threshold float64) {
	alert := &Alert{
		ID:        fmt.Sprintf("%s-%d", level, time.Now().UnixNano()),
		Timestamp: time.Now(),
		Level:     level,
		Title:     title,
		Message:   message,
		Metric:    metric,
		Value:     value,
		Threshold: threshold,
	}

	am.mu.Lock()
	am.alerts[alert.ID] = alert
	am.mu.Unlock()

	am.alertChan <- alert
}

// CheckCPUThreshold checks if CPU usage exceeds threshold
func (am *AlertManager) CheckCPUThreshold(usage float64, threshold float64) {
	if usage > threshold {
		am.TriggerAlert(
			AlertLevelWarning,
			"High CPU Usage",
			fmt.Sprintf("CPU usage is %.2f%%, exceeding threshold of %.2f%%", usage, threshold),
			"cpu",
			usage,
			threshold,
		)
	}
}

// CheckMemoryThreshold checks if memory usage exceeds threshold
func (am *AlertManager) CheckMemoryThreshold(usage float64, threshold float64) {
	if usage > threshold {
		am.TriggerAlert(
			AlertLevelWarning,
			"High Memory Usage",
			fmt.Sprintf("Memory usage is %.2f%%, exceeding threshold of %.2f%%", usage, threshold),
			"memory",
			usage,
			threshold,
		)
	}
}

// CheckDiskThreshold checks if disk usage exceeds threshold
func (am *AlertManager) CheckDiskThreshold(usage float64, threshold float64, path string) {
	if usage > threshold {
		am.TriggerAlert(
			AlertLevelWarning,
			"High Disk Usage",
			fmt.Sprintf("Disk usage on %s is %.2f%%, exceeding threshold of %.2f%%", path, usage, threshold),
			"disk",
			usage,
			threshold,
		)
	}
}

// CheckProcessAlert checks if a process is not running
func (am *AlertManager) CheckProcessAlert(processName string, running bool) {
	if !running {
		am.TriggerAlert(
			AlertLevelCritical,
			"Process Not Running",
			fmt.Sprintf("Process %s is not running", processName),
			"process",
			0,
			1,
		)
	}
}

// GetAlerts returns all active alerts
func (am *AlertManager) GetAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alerts := make([]*Alert, 0, len(am.alerts))
	for _, alert := range am.alerts {
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetAlertsByLevel returns alerts of a specific level
func (am *AlertManager) GetAlertsByLevel(level AlertLevel) []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*Alert
	for _, alert := range am.alerts {
		if alert.Level == level {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

// ClearAlert removes an alert
func (am *AlertManager) ClearAlert(id string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.alerts, id)
}

// ClearAlertsByLevel removes all alerts of a specific level
func (am *AlertManager) ClearAlertsByLevel(level AlertLevel) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for id, alert := range am.alerts {
		if alert.Level == level {
			delete(am.alerts, id)
		}
	}
}

// ClearAllAlerts removes all alerts
func (am *AlertManager) ClearAllAlerts() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.alerts = make(map[string]*Alert)
}

// processAlert processes an alert
func (am *AlertManager) processAlert(alert *Alert) {
	// Log the alert
	logMessage := fmt.Sprintf("[%s] %s: %s (Value: %.2f, Threshold: %.2f)",
		alert.Level, alert.Title, alert.Message, alert.Value, alert.Threshold)

	switch alert.Level {
	case AlertLevelInfo:
		am.logger.Info(logMessage)
	case AlertLevelWarning:
		am.logger.Warn(logMessage)
	case AlertLevelCritical:
		am.logger.Error(logMessage)
	}

	// Execute registered handlers
	am.mu.RLock()
	handlers := am.handlers[alert.Level]
	am.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(alert); err != nil {
			am.logger.Error(fmt.Sprintf("Alert handler error: %v", err))
		}
	}
}

// LogAlertHandler returns a handler that logs alerts
func LogAlertHandler(log *logger.Logger) AlertHandler {
	return func(alert *Alert) error {
		log.Info(fmt.Sprintf("Alert %s: %s", alert.ID, alert.Title))
		return nil
	}
}
