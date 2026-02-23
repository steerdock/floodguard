package detector

import (
	"testing"
	"time"

	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

func newTestDetector(maxConn, maxQPS int) *Detector {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		Monitor: config.MonitorConfig{
			MaxConnections: maxConn,
			MaxQPS:         maxQPS,
		},
	}
	return New(cfg, logger)
}

// containsIP checks whether ip exists in slice.
func containsIP(slice []string, ip string) bool {
	for _, s := range slice {
		if s == ip {
			return true
		}
	}
	return false
}

// TestAnalyze_ConnectionThreshold verifies that IPs exceeding MaxConnections are detected.
func TestAnalyze_ConnectionThreshold(t *testing.T) {
	d := newTestDetector(100, 0)

	connections := map[string]int{
		"1.2.3.4": 150, // exceeds threshold
		"5.6.7.8": 50,  // normal
	}
	result := d.Analyze(connections)

	if !containsIP(result, "1.2.3.4") {
		t.Errorf("expected 1.2.3.4 to be flagged for excessive connections")
	}
	if containsIP(result, "5.6.7.8") {
		t.Errorf("expected 5.6.7.8 NOT to be flagged")
	}
}

// TestAnalyze_FirstRun_NoQPS verifies that QPS detection is skipped on the first run.
func TestAnalyze_FirstRun_NoQPS(t *testing.T) {
	d := newTestDetector(1000, 10)

	// Even though MaxQPS is set, there is no previous snapshot yet.
	connections := map[string]int{
		"1.2.3.4": 500,
	}
	result := d.Analyze(connections)

	if len(result) != 0 {
		t.Errorf("expected no IPs on first run (no QPS baseline), got %v", result)
	}
}

// TestAnalyze_QPSThreshold verifies that IPs exceeding MaxQPS are detected on subsequent runs.
func TestAnalyze_QPSThreshold(t *testing.T) {
	d := newTestDetector(10000, 50) // high conn limit so only QPS triggers

	// First run — establishes the baseline snapshot.
	first := map[string]int{"1.2.3.4": 100}
	d.Analyze(first)

	// Manually backdate lastCheckTime to simulate a 2-second interval.
	d.mu.Lock()
	d.lastCheckTime = d.lastCheckTime.Add(-2 * time.Second)
	d.mu.Unlock()

	// Second run: delta = 300-100 = 200, elapsed ≈ 2s → QPS ≈ 100, exceeds 50.
	second := map[string]int{"1.2.3.4": 300}
	result := d.Analyze(second)

	if !containsIP(result, "1.2.3.4") {
		t.Errorf("expected 1.2.3.4 to be flagged for excessive QPS")
	}
}

// TestAnalyze_QPSDisabled verifies that setting max_qps=0 disables QPS detection.
func TestAnalyze_QPSDisabled(t *testing.T) {
	d := newTestDetector(10000, 0) // QPS detection disabled

	first := map[string]int{"1.2.3.4": 0}
	d.Analyze(first)

	d.mu.Lock()
	d.lastCheckTime = d.lastCheckTime.Add(-1 * time.Second)
	d.mu.Unlock()

	// Large increase, but QPS check is disabled.
	second := map[string]int{"1.2.3.4": 9999}
	result := d.Analyze(second)

	if len(result) != 0 {
		t.Errorf("expected no IPs when max_qps=0, got %v", result)
	}
}

// TestAnalyze_NegativeDelta verifies that a decreasing connection count does not trigger QPS detection.
func TestAnalyze_NegativeDelta(t *testing.T) {
	d := newTestDetector(10000, 10) // low QPS threshold

	first := map[string]int{"1.2.3.4": 500}
	d.Analyze(first)

	d.mu.Lock()
	d.lastCheckTime = d.lastCheckTime.Add(-2 * time.Second)
	d.mu.Unlock()

	// Connections dropped — delta is negative, should NOT trigger.
	second := map[string]int{"1.2.3.4": 100}
	result := d.Analyze(second)

	if len(result) != 0 {
		t.Errorf("expected no IPs when connection count drops, got %v", result)
	}
}

// TestAnalyze_IntervalTooShort verifies that QPS check is skipped when elapsed < 1s.
func TestAnalyze_IntervalTooShort(t *testing.T) {
	d := newTestDetector(10000, 1) // extremely low QPS threshold

	first := map[string]int{"1.2.3.4": 0}
	d.Analyze(first)

	// Do NOT backdate — lastCheckTime is just now, elapsed < 1s.
	second := map[string]int{"1.2.3.4": 9999}
	result := d.Analyze(second)

	if len(result) != 0 {
		t.Errorf("expected QPS check to be skipped when interval < 1s, got %v", result)
	}
}

// TestAnalyze_BothThresholdsExceeded verifies an IP is only returned once even if both checks trigger.
func TestAnalyze_BothThresholdsExceeded(t *testing.T) {
	d := newTestDetector(50, 10)

	first := map[string]int{"1.2.3.4": 0}
	d.Analyze(first)

	d.mu.Lock()
	d.lastCheckTime = d.lastCheckTime.Add(-2 * time.Second)
	d.mu.Unlock()

	// count=200 exceeds MaxConnections(50) AND QPS=(200/2s=100)>10.
	second := map[string]int{"1.2.3.4": 200}
	result := d.Analyze(second)

	count := 0
	for _, ip := range result {
		if ip == "1.2.3.4" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1.2.3.4 to appear exactly once, got %d times", count)
	}
}
