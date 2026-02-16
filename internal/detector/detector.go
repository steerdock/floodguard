package detector

import (
	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

type Detector struct {
	config *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Detector {
	return &Detector{
		config: cfg,
		logger: logger,
	}
}

func (d *Detector) Analyze(connections map[string]int) []string {
	var maliciousIPs []string

	for ip, count := range connections {
		if count > d.config.Monitor.MaxConnections {
			d.logger.Warn("Detected excessive connections",
				zap.String("ip", ip),
				zap.Int("connections", count),
				zap.Int("threshold", d.config.Monitor.MaxConnections),
			)
			maliciousIPs = append(maliciousIPs, ip)
		}
	}

	return maliciousIPs
}
