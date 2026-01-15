package monitor

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yourusername/floodguard/internal/config"
	"github.com/yourusername/floodguard/internal/detector"
	"github.com/yourusername/floodguard/internal/firewall"
	"go.uber.org/zap"
)

type Monitor struct {
	config   *config.Config
	logger   *zap.Logger
	detector *detector.Detector
	firewall firewall.Firewall
	stopChan chan struct{}
}

func New(cfg *config.Config, logger *zap.Logger) *Monitor {
	fw := firewall.NewAuto(logger)
	det := detector.New(cfg, logger)

	return &Monitor{
		config:   cfg,
		logger:   logger,
		detector: det,
		firewall: fw,
		stopChan: make(chan struct{}),
	}
}

func (m *Monitor) Start() error {
	m.logger.Info("Starting FloodGuard protection")

	interval, err := time.ParseDuration(m.config.Monitor.Interval)
	if err != nil {
		interval = 10 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	m.logger.Info("FloodGuard is running", zap.Duration("interval", interval))

	for {
		select {
		case <-ticker.C:
			if err := m.check(); err != nil {
				m.logger.Error("Check failed", zap.Error(err))
			}
		case <-sigChan:
			m.logger.Info("Received shutdown signal")
			return nil
		case <-m.stopChan:
			return nil
		}
	}
}

func (m *Monitor) check() error {
	connections, err := m.getConnections()
	if err != nil {
		return fmt.Errorf("failed to get connections: %w", err)
	}

	maliciousIPs := m.detector.Analyze(connections)

	for _, ip := range maliciousIPs {
		if m.isWhitelisted(ip) {
			m.logger.Info("IP is whitelisted, skipping", zap.String("ip", ip))
			continue
		}

		m.logger.Warn("Blocking malicious IP", zap.String("ip", ip))
		if err := m.firewall.Block(ip); err != nil {
			m.logger.Error("Failed to block IP", zap.String("ip", ip), zap.Error(err))
		}
	}

	return nil
}

func (m *Monitor) getConnections() (map[string]int, error) {
	connections := make(map[string]int)

	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		remoteAddr := fields[2]
		parts := strings.Split(remoteAddr, ":")
		if len(parts) != 2 {
			continue
		}

		ip := m.hexToIP(parts[0])
		connections[ip]++
	}

	return connections, scanner.Err()
}

func (m *Monitor) hexToIP(hex string) string {
	return hex
}

func (m *Monitor) isWhitelisted(ip string) bool {
	for _, whiteIP := range m.config.Whitelist {
		if ip == whiteIP {
			return true
		}
	}
	return false
}

func (m *Monitor) Stop() {
	close(m.stopChan)
}
