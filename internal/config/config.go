package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Monitor      MonitorConfig      `yaml:"monitor"`
	Ban          BanConfig          `yaml:"ban"`
	Whitelist    []string           `yaml:"whitelist"`
	Blacklist    []string           `yaml:"blacklist"`
	Notification NotificationConfig `yaml:"notification"`
	Log          LogConfig          `yaml:"log"`
}

type MonitorConfig struct {
	Interval       string `yaml:"interval"`
	MaxConnections int    `yaml:"max_connections"`
	MaxQPS         int    `yaml:"max_qps"`
	Ports          []int  `yaml:"ports"`
	CheckHTTP      bool   `yaml:"check_http"`
}

type BanConfig struct {
	Duration int    `yaml:"duration"`
	Mode     string `yaml:"mode"`
}

type NotificationConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Type       string `yaml:"type"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	File   string `yaml:"file"`
	Format string `yaml:"format"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func InitDefault(publicIP string) error {
	// Build whitelist / 构建白名单
	whitelist := []string{
		"127.0.0.1",
		"::1",
	}
	
	// Add public and local IPs if provided / 如果提供了公网和本地 IP，添加到白名单
	if publicIP != "" {
		// Split multiple IPs by space / 按空格分割多个 IP
		ips := strings.Fields(publicIP)
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && ip != "127.0.0.1" && ip != "::1" {
				whitelist = append(whitelist, ip)
			}
		}
	}
	
	defaultConfig := &Config{
		Monitor: MonitorConfig{
			Interval:       "10s",
			MaxConnections: 100,
			MaxQPS:         50,
			Ports:          []int{80, 443},
			CheckHTTP:      true,
		},
		Ban: BanConfig{
			Duration: 3600,
			Mode:     "auto",
		},
		Whitelist: whitelist,
		Blacklist: []string{},
		Notification: NotificationConfig{
			Enabled:    false,
			WebhookURL: "",
			Type:       "dingtalk",
		},
		Log: LogConfig{
			Level:  "info",
			File:   "/var/log/floodguard/floodguard.log",
			Format: "json",
		},
	}

	configDir := "/etc/floodguard"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Configuration file created at: %s\n", configPath)
	if publicIP != "" {
		ips := strings.Fields(publicIP)
		if len(ips) > 0 {
			fmt.Printf("The following IPs have been added to the whitelist:\n")
			fmt.Printf("以下 IP 已添加到白名单：\n")
			for _, ip := range ips {
				if ip != "127.0.0.1" && ip != "::1" {
					fmt.Printf("  - %s\n", ip)
				}
			}
		}
	}
	return nil
}
