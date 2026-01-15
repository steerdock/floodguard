#!/bin/bash
# FloodGuard Installation Script
# FloodGuard 安装脚本

set -e

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/floodguard"
LOG_DIR="/var/log/floodguard"
SERVICE_FILE="/etc/systemd/system/floodguard.service"

echo "Installing FloodGuard..."
echo "正在安装 FloodGuard..."

# Check root privileges / 检查 root 权限
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root"
    echo "请使用 root 权限运行"
    exit 1
fi

# Create necessary directories / 创建必要的目录
mkdir -p "$CONFIG_DIR"
mkdir -p "$LOG_DIR"

# Copy binary file / 复制二进制文件
if [ -f "build/floodguard" ]; then
    cp build/floodguard "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/floodguard"
    echo "Binary installed to $INSTALL_DIR/floodguard"
    echo "二进制文件已安装到 $INSTALL_DIR/floodguard"
else
    echo "Error: build/floodguard not found. Please run 'make build' first."
    echo "错误：未找到 build/floodguard，请先运行 'make build'"
    exit 1
fi

# Detect server's public IP / 检测服务器的公网 IP
echo ""
echo "Detecting server's public IP address..."
echo "正在检测服务器的公网 IP 地址..."

# Try multiple IP detection services / 尝试多个 IP 检测服务
PUBLIC_IP=$(curl -s --max-time 5 ifconfig.me 2>/dev/null || \
            curl -s --max-time 5 icanhazip.com 2>/dev/null || \
            curl -s --max-time 5 ipinfo.io/ip 2>/dev/null || \
            curl -s --max-time 5 api.ipify.org 2>/dev/null || \
            echo "")

# Also detect local network IPs / 同时检测本地网络 IP
LOCAL_IPS=$(hostname -I 2>/dev/null | tr ' ' '\n' | grep -v '^$' || ip addr show | grep 'inet ' | awk '{print $2}' | cut -d/ -f1 | grep -v '127.0.0.1')

if [ -n "$PUBLIC_IP" ]; then
    echo "Detected server public IP: $PUBLIC_IP"
    echo "检测到服务器公网 IP: $PUBLIC_IP"
    echo "This IP will be added to the whitelist automatically."
    echo "此 IP 将自动添加到白名单。"
else
    echo "Warning: Could not detect server's public IP address."
    echo "警告：无法检测服务器的公网 IP 地址。"
    PUBLIC_IP=""
fi

if [ -n "$LOCAL_IPS" ]; then
    echo ""
    echo "Detected local network IPs:"
    echo "检测到本地网络 IP:"
    echo "$LOCAL_IPS"
    echo "These IPs will also be added to the whitelist."
    echo "这些 IP 也将添加到白名单。"
fi

# Generate default configuration / 生成默认配置
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    # Pass both public and local IPs / 传递公网和本地 IP
    ALL_IPS="$PUBLIC_IP"
    if [ -n "$LOCAL_IPS" ]; then
        ALL_IPS="$ALL_IPS $LOCAL_IPS"
    fi
    
    floodguard init --public-ip="$ALL_IPS"
    echo "Configuration file created at $CONFIG_DIR/config.yaml"
    echo "配置文件已创建于 $CONFIG_DIR/config.yaml"
else
    echo "Configuration file already exists, skipping..."
    echo "配置文件已存在，跳过..."
fi

# Create systemd service / 创建 systemd 服务
cat > "$SERVICE_FILE" << EOF
[Unit]
Description=FloodGuard - Lightweight Firewall Protection
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/floodguard start
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd / 重载 systemd
systemctl daemon-reload

echo ""
echo "Installation complete! / 安装完成！"
echo ""
echo "Next steps / 下一步操作:"
echo "  1. Edit configuration / 编辑配置: nano $CONFIG_DIR/config.yaml"
echo "  2. Start service / 启动服务: systemctl start floodguard"
echo "  3. Enable on boot / 开机自启: systemctl enable floodguard"
echo "  4. Check status / 查看状态: systemctl status floodguard"
