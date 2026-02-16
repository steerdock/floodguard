#!/bin/bash
# FloodGuard Uninstallation Script
# FloodGuard 卸载脚本

set -e

INSTALL_DIR="/usr/bin"
CONFIG_DIR="/etc/floodguard"
LOG_DIR="/var/log/floodguard"
SERVICE_FILE="/etc/systemd/system/floodguard.service"

echo "Uninstalling FloodGuard..."
echo "正在卸载 FloodGuard..."

# Check root privileges / 检查 root 权限
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root"
    echo "请使用 root 权限运行"
    exit 1
fi

# Stop and disable service / 停止并禁用服务
if systemctl is-active --quiet floodguard; then
    systemctl stop floodguard
    echo "Service stopped"
    echo "服务已停止"
fi

if systemctl is-enabled --quiet floodguard; then
    systemctl disable floodguard
    echo "Service disabled"
    echo "服务已禁用"
fi

# Remove service file / 删除服务文件
if [ -f "$SERVICE_FILE" ]; then
    rm "$SERVICE_FILE"
    systemctl daemon-reload
    echo "Service file removed"
    echo "服务文件已删除"
fi

# Remove binary file / 删除二进制文件
if [ -f "$INSTALL_DIR/floodguard" ]; then
    rm "$INSTALL_DIR/floodguard"
    echo "Binary removed"
    echo "二进制文件已删除"
fi

# Ask whether to remove configuration and logs / 询问是否删除配置和日志
echo ""
read -p "Remove configuration and logs? / 是否删除配置和日志? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$CONFIG_DIR"
    rm -rf "$LOG_DIR"
    echo "Configuration and logs removed"
    echo "配置和日志已删除"
fi

echo ""
echo "Uninstallation complete! / 卸载完成！"
