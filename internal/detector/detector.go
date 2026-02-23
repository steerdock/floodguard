package detector

import (
	"sync"
	"time"

	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

// Detector analyzes per-IP connection counts to identify malicious traffic.
// It supports two detection strategies:
//   - Max Connections: flags IPs whose current connection count exceeds the threshold.
//   - Max QPS: flags IPs whose request rate (connections per second) exceeds the threshold,
//     calculated by comparing consecutive snapshots.
type Detector struct {
	config        *config.Config
	logger        *zap.Logger
	mu            sync.Mutex
	lastSnapshot  map[string]int // connection counts from the previous check
	lastCheckTime time.Time      // timestamp of the previous check
}

func New(cfg *config.Config, logger *zap.Logger) *Detector {
	return &Detector{
		config:       cfg,
		logger:       logger,
		lastSnapshot: make(map[string]int),
	}
}

// Analyze inspects the current connection snapshot and returns a list of IPs
// that exceed either the max_connections or max_qps threshold.
func (d *Detector) Analyze(connections map[string]int) []string {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	var maliciousIPs []string

	// Compute elapsed time since the last check for QPS calculation.
	// Skip QPS detection on the first run (lastCheckTime is zero) or if the
	// interval is too short to produce a reliable rate.
	elapsed := 0.0
	qpsEnabled := d.config.Monitor.MaxQPS > 0 && !d.lastCheckTime.IsZero()
	if qpsEnabled {
		elapsed = now.Sub(d.lastCheckTime).Seconds()
		if elapsed < 1.0 {
			// Interval too short — disable QPS check this round to avoid noise.
			qpsEnabled = false
		}
	}

	for ip, count := range connections {
		blocked := false

		// --- Connection count check ---
		if d.config.Monitor.MaxConnections > 0 && count > d.config.Monitor.MaxConnections {
			d.logger.Warn("Detected excessive connections",
				zap.String("ip", ip),
				zap.Int("connections", count),
				zap.Int("threshold", d.config.Monitor.MaxConnections),
			)
			blocked = true
		}

		// --- QPS check ---
		if !blocked && qpsEnabled {
			prev := d.lastSnapshot[ip]
			delta := count - prev
			if delta > 0 {
				qps := float64(delta) / elapsed
				if qps > float64(d.config.Monitor.MaxQPS) {
					d.logger.Warn("Detected excessive QPS",
						zap.String("ip", ip),
						zap.Float64("qps", qps),
						zap.Int("threshold", d.config.Monitor.MaxQPS),
						zap.Int("connections_delta", delta),
						zap.Float64("elapsed_seconds", elapsed),
					)
					blocked = true
				}
			}
		}

		if blocked {
			maliciousIPs = append(maliciousIPs, ip)
		}
	}

	// Update state for the next round.
	d.lastSnapshot = make(map[string]int, len(connections))
	for ip, count := range connections {
		d.lastSnapshot[ip] = count
	}
	d.lastCheckTime = now

	return maliciousIPs
}
