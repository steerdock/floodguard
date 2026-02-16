# FloodGuard

一个现代化的轻量级 Linux 防火墙工具，用于防御 CC 攻击和 DDoS 攻击。

## 特性

- 🚀 **轻量高效**：Go 语言编写，单二进制文件，资源占用少
- 🛡️ **智能防护**：多维度检测异常连接，自动封禁攻击 IP
- 🔧 **灵活配置**：支持 YAML 配置，可自定义各种阈值和策略
- 📊 **实时监控**：连接数统计、攻击日志、封禁记录
- 🔔 **通知告警**：支持 Webhook 通知（钉钉、企业微信、Slack）
- 🌐 **多后端支持**：自动适配 iptables、nftables、firewalld
- 📝 **详细日志**：结构化日志输出，支持多种格式

## 快速开始

### 安装

# 使用 Go 安装
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

**注意**：安装时会自动检测服务器的公网 IP 和本地网络 IP，并添加到白名单，防止误封。

### 使用

```bash
# 生成默认配置文件
sudo floodguard init

# 启动防护
sudo floodguard start

# 查看状态
sudo floodguard status

# 查看封禁列表
sudo floodguard list

# 解封 IP
sudo floodguard unban 1.2.3.4
```

## 配置说明

配置文件位于 `/etc/floodguard/config.yaml`

```yaml
# 监控设置
monitor:
  interval: 10s              # 检测间隔
  max_connections: 100       # 单 IP 最大连接数
  max_qps: 50                # 单 IP 最大 QPS
  
# 封禁策略
ban:
  duration: 3600            # 封禁时长（秒），0 为永久
  mode: "auto"              # auto/iptables/nftables/firewalld

# 白名单
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"
  
# 通知
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## 系统要求

- Linux 系统（内核 3.10+）
- root 权限
- iptables 或 nftables

## 开发

```bash
# 克隆项目
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# 安装依赖
go mod download

# 编译
go build -o build/floodguard cmd/floodguard/main.go

# 安装到系统
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# 修复 SELinux 上下文（RHEL/CentOS/Fedora 系统必需）
sudo restorecon -v /usr/local/bin/floodguard

# 初始化配置文件（重要：必须先执行此步骤！）
sudo /usr/local/bin/floodguard init

# 创建 systemd 服务
sudo tee /etc/systemd/system/floodguard.service > /dev/null <<EOF
[Unit]
Description=FloodGuard - DDoS 防护服务
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

# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard

# 查看服务状态
sudo systemctl status floodguard

# 运行测试
go test ./...
```

## 服务管理

```bash
# 启动服务
sudo systemctl start floodguard

# 停止服务
sudo systemctl stop floodguard

# 重启服务
sudo systemctl restart floodguard

# 查看状态
sudo systemctl status floodguard

# 查看日志
sudo journalctl -u floodguard -f

# 开机自启动
sudo systemctl enable floodguard

# 取消开机自启动
sudo systemctl disable floodguard
```

## License

MIT License
