package firewall

import (
	"fmt"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type Iptables struct {
	logger *zap.Logger
	chain  string
}

func NewIptables(logger *zap.Logger) *Iptables {
	return &Iptables{
		logger: logger,
		chain:  "FIREGUARD",
	}
}

func (i *Iptables) IsAvailable() bool {
	cmd := exec.Command("iptables", "--version")
	return cmd.Run() == nil
}

func (i *Iptables) Block(ip string) error {
	i.ensureChain()

	cmd := exec.Command("iptables", "-I", i.chain, "-s", ip, "-j", "DROP")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to block IP %s: %w, output: %s", ip, err, output)
	}

	i.logger.Info("Blocked IP", zap.String("ip", ip))
	return nil
}

func (i *Iptables) Unblock(ip string) error {
	cmd := exec.Command("iptables", "-D", i.chain, "-s", ip, "-j", "DROP")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to unblock IP %s: %w, output: %s", ip, err, output)
	}

	i.logger.Info("Unblocked IP", zap.String("ip", ip))
	return nil
}

func (i *Iptables) List() ([]string, error) {
	cmd := exec.Command("iptables", "-L", i.chain, "-n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	var ips []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "DROP") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				ips = append(ips, fields[3])
			}
		}
	}

	return ips, nil
}

func (i *Iptables) ensureChain() error {
	cmd := exec.Command("iptables", "-N", i.chain)
	cmd.Run()

	cmd = exec.Command("iptables", "-C", "INPUT", "-j", i.chain)
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("iptables", "-I", "INPUT", "-j", i.chain)
		return cmd.Run()
	}

	return nil
}
