package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName                string   `yaml:"app_name"`
	Debug                  bool     `yaml:"debug"`
	LogPath                string   `yaml:"log_path"`
	MaxLogSize             int      `yaml:"max_log_size"`
	MaxLogBackups          int      `yaml:"max_log_backups"`
	CPUMonitorInterval     int      `yaml:"cpu_monitor_interval"`
	MemoryMonitorInterval  int      `yaml:"memory_monitor_interval"`
	DiskMonitorInterval    int      `yaml:"disk_monitor_interval"`
	NetworkMonitorInterval int      `yaml:"network_monitor_interval"`
	CPUThreshold           float64  `yaml:"cpu_threshold"`
	MemoryThreshold        float64  `yaml:"memory_threshold"`
	HealthCheckPort        int      `yaml:"health_check_port"`
	EnableMetrics          bool     `yaml:"enable_metrics"`
	MetricsPort            int      `yaml:"metrics_port"`
	Processes              []string `yaml:"processes"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		AppName:                "mydaemon",
		Debug:                  false,
		LogPath:                "/var/log/mydaemon/mydaemon.log",
		MaxLogSize:             100, // MB
		MaxLogBackups:          10,
		CPUMonitorInterval:     5,
		MemoryMonitorInterval:  5,
		DiskMonitorInterval:    10,
		NetworkMonitorInterval: 10,
		CPUThreshold:           80.0,
		MemoryThreshold:        80.0,
		HealthCheckPort:        8080,
		EnableMetrics:          false,
		MetricsPort:            9090,
		Processes:              []string{},
	}
}

// LoadConfig loads configuration from YAML file
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Use defaults if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// SaveConfig saves configuration to YAML file
func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("app_name cannot be empty")
	}

	if c.LogPath == "" {
		return fmt.Errorf("log_path cannot be empty")
	}

	logDir := filepath.Dir(c.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("cannot create log directory %s: %w", logDir, err)
	}

	if c.MaxLogSize < 1 {
		return fmt.Errorf("max_log_size must be >= 1 MB, got %d", c.MaxLogSize)
	}

	if c.MaxLogBackups < 0 {
		return fmt.Errorf("max_log_backups must be >= 0, got %d", c.MaxLogBackups)
	}

	if c.CPUMonitorInterval < 1 {
		return fmt.Errorf("cpu_monitor_interval must be >= 1 second, got %d", c.CPUMonitorInterval)
	}

	if c.MemoryMonitorInterval < 1 {
		return fmt.Errorf("memory_monitor_interval must be >= 1 second, got %d", c.MemoryMonitorInterval)
	}

	if c.DiskMonitorInterval < 1 {
		return fmt.Errorf("disk_monitor_interval must be >= 1 second, got %d", c.DiskMonitorInterval)
	}

	if c.NetworkMonitorInterval < 1 {
		return fmt.Errorf("network_monitor_interval must be >= 1 second, got %d", c.NetworkMonitorInterval)
	}

	if c.CPUThreshold < 0 || c.CPUThreshold > 100 {
		return fmt.Errorf("cpu_threshold must be between 0 and 100, got %.2f", c.CPUThreshold)
	}

	if c.MemoryThreshold < 0 || c.MemoryThreshold > 100 {
		return fmt.Errorf("memory_threshold must be between 0 and 100, got %.2f", c.MemoryThreshold)
	}

	if c.HealthCheckPort < 1 || c.HealthCheckPort > 65535 {
		return fmt.Errorf("health_check_port must be between 1 and 65535, got %d", c.HealthCheckPort)
	}

	if c.EnableMetrics && (c.MetricsPort < 1 || c.MetricsPort > 65535) {
		return fmt.Errorf("metrics_port must be between 1 and 65535, got %d", c.MetricsPort)
	}

	return nil
}
