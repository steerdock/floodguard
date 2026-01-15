package firewall

import (
	"go.uber.org/zap"
)

type Dummy struct {
	logger      *zap.Logger
	blockedIPs  map[string]bool
}

func NewDummy(logger *zap.Logger) *Dummy {
	return &Dummy{
		logger:     logger,
		blockedIPs: make(map[string]bool),
	}
}

func (d *Dummy) IsAvailable() bool {
	return true
}

func (d *Dummy) Block(ip string) error {
	d.blockedIPs[ip] = true
	d.logger.Info("[DUMMY] Blocked IP", zap.String("ip", ip))
	return nil
}

func (d *Dummy) Unblock(ip string) error {
	delete(d.blockedIPs, ip)
	d.logger.Info("[DUMMY] Unblocked IP", zap.String("ip", ip))
	return nil
}

func (d *Dummy) List() ([]string, error) {
	var ips []string
	for ip := range d.blockedIPs {
		ips = append(ips, ip)
	}
	return ips, nil
}
