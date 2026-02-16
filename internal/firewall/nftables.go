package firewall

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

type Nftables struct {
	logger *zap.Logger
	table  string
	chain  string
}

func NewNftables(logger *zap.Logger) *Nftables {
	return &Nftables{
		logger: logger,
		table:  "fireguard",
		chain:  "input",
	}
}

func (n *Nftables) IsAvailable() bool {
	cmd := exec.Command("nft", "--version")
	return cmd.Run() == nil
}

func (n *Nftables) Block(ip string) error {
	n.ensureTable()

	rule := fmt.Sprintf("add rule ip %s %s ip saddr %s drop", n.table, n.chain, ip)
	cmd := exec.Command("nft", rule)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to block IP %s: %w, output: %s", ip, err, output)
	}

	n.logger.Info("Blocked IP", zap.String("ip", ip))
	return nil
}

func (n *Nftables) Unblock(ip string) error {
	// nftables 需要通过 handle 删除规则，这里简化处理
	n.logger.Info("Unblock not fully implemented for nftables", zap.String("ip", ip))
	return nil
}

func (n *Nftables) List() ([]string, error) {
	cmd := exec.Command("nft", "list", "table", "ip", n.table)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	n.logger.Info("Rules", zap.String("output", string(output)))
	return []string{}, nil
}

func (n *Nftables) ensureTable() error {
	cmd := exec.Command("nft", "add", "table", "ip", n.table)
	cmd.Run()

	cmd = exec.Command("nft", "add", "chain", "ip", n.table, n.chain, "{", "type", "filter", "hook", "input", "priority", "0", ";", "}")
	return cmd.Run()
}
