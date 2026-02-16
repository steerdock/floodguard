# FloodGuard

A modern, lightweight Linux firewall tool for defending against CC and DDoS attacks.

## Features

- 🚀 **Lightweight & Efficient**: Written in Go, single binary, minimal resource usage
- 🛡️ **Smart Protection**: Multi-dimensional anomaly detection, automatic IP blocking
- 🔧 **Flexible Configuration**: YAML-based config with customizable thresholds
- 📊 **Real-time Monitoring**: Connection statistics, attack logs, ban records
- 🔔 **Alert Notifications**: Webhook support (DingTalk, WeCom, Slack)
- 🌐 **Multiple Backends**: Auto-detect iptables, nftables, firewalld
- 📝 **Detailed Logging**: Structured logs with multiple formats

## Quick Start

### Installation

**Note**: During installation, the server's public IP and local network IPs will be automatically detected and added to the whitelist to prevent accidental blocking.

### Usage

```bash
# Generate default config
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
  max_qps: 50               # Max QPS per IP
  
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
go build -o build/floodguard cmd/floodguard/main.go

# Install to system
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# Fix SELinux context (required for RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# Initialize configuration (IMPORTANT: Run this first!)
sudo /usr/local/bin/floodguard init

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

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard

# Check service status
sudo systemctl status floodguard
```

## Service Management

```bash
# Start service
sudo systemctl start floodguard

# Stop service
sudo systemctl stop floodguard

# Restart service
sudo systemctl restart floodguard

# Check status
sudo systemctl status floodguard

# View logs
sudo journalctl -u floodguard -f

# Enable auto-start on boot
sudo systemctl enable floodguard

# Disable auto-start
sudo systemctl disable floodguard
```

## License

MIT License
