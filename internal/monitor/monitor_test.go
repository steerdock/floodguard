package monitor

import (
	"testing"

	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

func TestHexToIPv4(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	m := &Monitor{logger: logger}

	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "localhost",
			hex:      "0100007F", // 127.0.0.1 in little-endian hex
			expected: "127.0.0.1",
		},
		{
			name:     "google dns",
			hex:      "08080808", // 8.8.8.8 in little-endian hex
			expected: "8.8.8.8",
		},
		{
			name:     "private network",
			hex:      "0101A8C0", // 192.168.1.1 in little-endian hex
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.hexToIPv4(tt.hex)
			if result != tt.expected {
				t.Errorf("hexToIPv4(%s) = %s, want %s", tt.hex, result, tt.expected)
			}
		})
	}
}

func TestIsWhitelisted(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		Whitelist: []string{
			"127.0.0.1",
			"192.168.0.0/16",
			"10.0.0.0/8",
			"::1",
			"2001:db8::/32",
		},
	}
	m := &Monitor{
		config: cfg,
		logger: logger,
	}

	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "exact match localhost",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "CIDR match private network",
			ip:       "192.168.1.100",
			expected: true,
		},
		{
			name:     "CIDR match 10.x network",
			ip:       "10.0.0.1",
			expected: true,
		},
		{
			name:     "not whitelisted",
			ip:       "8.8.8.8",
			expected: false,
		},
		{
			name:     "IPv6 localhost",
			ip:       "::1",
			expected: true,
		},
		{
			name:     "IPv6 CIDR match",
			ip:       "2001:db8::1",
			expected: true,
		},
		{
			name:     "IPv6 not whitelisted",
			ip:       "2001:4860:4860::8888",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.isWhitelisted(tt.ip)
			if result != tt.expected {
				t.Errorf("isWhitelisted(%s) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}
