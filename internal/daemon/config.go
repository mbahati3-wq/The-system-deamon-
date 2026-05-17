package daemon

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName                string `yaml:"app_name"`
	Debug                  bool   `yaml:"debug"`
	LogPath                string `yaml:"log_path"`
	MaxLogSize             int    `yaml:"max_log_size"`
	MaxLogBackups          int    `yaml:"max_log_backups"`
	CPUMonitorInterval     int    `yaml:"cpu_monitor_interval"`
	MemoryMonitorInterval  int    `yaml:"memory_monitor_interval"`
	DiskMonitorInterval    int    `yaml:"disk_monitor_interval"`
	NetworkMonitorInterval int    `yaml:"network_monitor_interval"`
	CPUThreshold           float64 `yaml:"cpu_threshold"`
	MemoryThreshold        float64 `yaml:"memory_threshold"`
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
