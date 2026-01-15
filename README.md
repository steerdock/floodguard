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

```bash
# 下载二进制文件
wget https://github.com/steerdock/floodguard/releases/latest/download/floodguard-linux-amd64
chmod +x floodguard-linux-amd64
sudo mv floodguard-linux-amd64 /usr/local/bin/floodguard

# 或使用 Go 安装
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
go build -o floodguard cmd/floodguard/main.go

# 运行测试
go test ./...
```

## License

MIT License
