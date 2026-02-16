package monitor

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
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

	if len(connections) == 0 {
		m.logger.Debug("No connections found to analyze")
		return nil
	}

	m.logger.Debug("Analyzing connections", zap.Int("total_connections", len(connections)))

	maliciousIPs := m.detector.Analyze(connections)
	if len(maliciousIPs) == 0 {
		m.logger.Debug("No malicious IPs detected")
		return nil
	}

	m.logger.Info("Detected malicious IPs", zap.Int("count", len(maliciousIPs)))

	blockedCount := 0
	whitelistedCount := 0
	failedCount := 0

	for _, ip := range maliciousIPs {
		if m.isWhitelisted(ip) {
			whitelistedCount++
			m.logger.Info("IP is whitelisted, skipping", 
				zap.String("ip", ip),
				zap.Int("connections", connections[ip]))
			continue
		}

		m.logger.Warn("Blocking malicious IP", 
			zap.String("ip", ip),
			zap.Int("connections", connections[ip]))
		
		if err := m.firewall.Block(ip); err != nil {
			failedCount++
			m.logger.Error("Failed to block IP", 
				zap.String("ip", ip), 
				zap.Error(err),
				zap.String("error_type", "firewall_block_failed"))
		} else {
			blockedCount++
			m.logger.Info("Successfully blocked IP", zap.String("ip", ip))
		}
	}

	m.logger.Info("Check completed", 
		zap.Int("total_malicious", len(maliciousIPs)),
		zap.Int("blocked", blockedCount),
		zap.Int("whitelisted", whitelistedCount),
		zap.Int("failed", failedCount))

	if failedCount > 0 {
		return fmt.Errorf("failed to block %d out of %d malicious IPs", failedCount, len(maliciousIPs))
	}

	return nil
}

func (m *Monitor) getConnections() (map[string]int, error) {
	connections := make(map[string]int)

	// Process IPv4 connections
	if err := m.processConnectionFile("/proc/net/tcp", connections); err != nil {
		m.logger.Warn("Failed to process IPv4 connections", zap.Error(err))
	}

	// Process IPv6 connections
	if err := m.processConnectionFile("/proc/net/tcp6", connections); err != nil {
		m.logger.Warn("Failed to process IPv6 connections", zap.Error(err))
	}

	return connections, nil
}

func (m *Monitor) processConnectionFile(filepath string, connections map[string]int) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", filepath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header line
	if !scanner.Scan() {
		return fmt.Errorf("failed to read header from %s", filepath)
	}

	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			m.logger.Debug("Skipping malformed line", 
				zap.String("file", filepath), 
				zap.Int("line", lineCount),
				zap.String("content", line))
			continue
		}

		remoteAddr := fields[2]
		parts := strings.Split(remoteAddr, ":")
		if len(parts) != 2 {
			m.logger.Debug("Invalid remote address format", 
				zap.String("file", filepath),
				zap.String("address", remoteAddr))
			continue
		}

		ip := m.hexToIP(parts[0])
		if ip != parts[0] { // Only count if conversion was successful
			connections[ip]++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filepath, err)
	}

	m.logger.Debug("Processed connection file", 
		zap.String("file", filepath), 
		zap.Int("lines", lineCount))
	
	return nil
}

func (m *Monitor) hexToIP(hex string) string {
	// Handle both IPv4 and IPv6 hex formats
	if len(hex) == 8 {
		// IPv4: 8 hex chars = 4 bytes
		return m.hexToIPv4(hex)
	} else if len(hex) == 32 {
		// IPv6: 32 hex chars = 16 bytes
		return m.hexToIPv6(hex)
	}
	
	// If format is unknown, return original for debugging
	m.logger.Warn("Unknown hex IP format", zap.String("hex", hex))
	return hex
}

func (m *Monitor) hexToIPv4(hex string) string {
	if len(hex) != 8 {
		return hex
	}
	
	var ip [4]byte
	for i := 0; i < 4; i++ {
		// Parse each byte (2 hex chars) in little-endian format
		val, err := strconv.ParseUint(hex[i*2:(i+1)*2], 16, 8)
		if err != nil {
			m.logger.Warn("Failed to parse hex byte", zap.String("hex", hex), zap.Error(err))
			return hex
		}
		// /proc/net/tcp stores IP in little-endian format
		ip[3-i] = byte(val)
	}
	
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (m *Monitor) hexToIPv6(hex string) string {
	if len(hex) != 32 {
		return hex
	}
	
	var ip [16]byte
	for i := 0; i < 16; i++ {
		val, err := strconv.ParseUint(hex[i*2:(i+1)*2], 16, 8)
		if err != nil {
			m.logger.Warn("Failed to parse hex byte for IPv6", zap.String("hex", hex), zap.Error(err))
			return hex
		}
		ip[i] = byte(val)
	}
	
	// Convert to standard IPv6 format
	ipv6 := net.IP(ip[:])
	return ipv6.String()
}

func (m *Monitor) isWhitelisted(ip string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		m.logger.Warn("Invalid IP address format", zap.String("ip", ip))
		return false
	}
	
	for _, whiteEntry := range m.config.Whitelist {
		// Check if it's a CIDR notation
		if strings.Contains(whiteEntry, "/") {
			if m.isIPInCIDR(clientIP, whiteEntry) {
				return true
			}
		} else {
			// Simple IP comparison
			whiteIP := net.ParseIP(whiteEntry)
			if whiteIP != nil && clientIP.Equal(whiteIP) {
				return true
			}
		}
	}
	return false
}

func (m *Monitor) isIPInCIDR(ip net.IP, cidr string) bool {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		m.logger.Warn("Invalid CIDR format in whitelist", zap.String("cidr", cidr), zap.Error(err))
		return false
	}
	
	return network.Contains(ip)
}

func (m *Monitor) Stop() {
	close(m.stopChan)
}
