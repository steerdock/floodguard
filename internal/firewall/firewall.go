package firewall

import (
	"go.uber.org/zap"
)

type Firewall interface {
	Block(ip string) error
	Unblock(ip string) error
	List() ([]string, error)
	IsAvailable() bool
}

func NewAuto(logger *zap.Logger) Firewall {
	if nft := NewNftables(logger); nft.IsAvailable() {
		logger.Info("Using nftables backend")
		return nft
	}

	if ipt := NewIptables(logger); ipt.IsAvailable() {
		logger.Info("Using iptables backend")
		return ipt
	}

	logger.Warn("No firewall backend available, using dummy")
	return NewDummy(logger)
}
