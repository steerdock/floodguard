# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md)

A modern, lightweight Linux firewall tool for defending against CC and DDoS attacks.

## Features

- 🚀 **Lightweight & Fast**: Written in Go, single binary, minimal resource usage
- 🛡️ **Smart Protection**: Multi-dimensional detection of abnormal connections, automatic IP banning
- 🔧 **Flexible Configuration**: YAML-based config with customizable thresholds and policies
- 📊 **Real-time Monitoring**: Connection statistics, attack logs, and ban records
- 🔔 **Alert Notifications**: Webhook support (DingTalk, WeCom, Slack)
- 🌐 **Multi-backend Support**: Auto-detects iptables, nftables, firewalld
- 📝 **Detailed Logging**: Structured log output in multiple formats

## Quick Start

### Installation

```bash
# Install via Go
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **Note**: During `init`, FloodGuard automatically detects your server's public and local IPs and adds them to the whitelist to prevent accidental self-blocking.

### Usage

```bash
# Generate default config file
sudo floodguard init

# Start protection
sudo floodguard start

# Check status
sudo floodguard status

# List banned IPs
sudo floodguard list

# Unban an IP
sudo floodguard unban 1.2.3.4
```

## Configuration

Config file location: `/etc/floodguard/config.yaml`

```yaml
# Monitor settings
monitor:
  interval: 10s              # Check interval
  max_connections: 100       # Max connections per IP
  max_qps: 50                # Max QPS per IP

# Ban policy
ban:
  duration: 3600            # Ban duration (seconds), 0 for permanent
  mode: "auto"              # auto/iptables/nftables/firewalld

# Whitelist
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# Notifications
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## System Requirements

- Linux (kernel 3.10+)
- Root privileges
- iptables or nftables

## Development

```bash
# Clone repository
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

## Deployment (systemd)

```bash
# Install binary
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# Fix SELinux context (RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# Initialize configuration (run this first!)
sudo floodguard init

# Create systemd service
sudo tee /etc/systemd/system/floodguard.service > /dev/null <<EOF
[Unit]
Description=FloodGuard - DDoS Protection Service
After=network.target

[Service]
Type=exec
ExecStart=/usr/local/bin/floodguard start --config /etc/floodguard/config.yaml
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## Service Management

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## License

MIT License
