#!/bin/bash

# StardewPanel 一键安装脚本

set -e

echo "🌾 StardewPanel 安装脚本"
echo "=========================="

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

echo "✅ Docker 环境检查通过"

# 创建安装目录
INSTALL_DIR="/opt/stardew-panel"
echo "📁 创建安装目录: $INSTALL_DIR"
sudo mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

# 下载项目文件
echo "📥 下载项目文件..."
# TODO: 替换为实际的 GitHub 仓库地址
# git clone https://github.com/yourusername/stardew-panel.git .

# 启动服务
echo "🚀 启动服务..."
cd docker
docker-compose up -d

echo ""
echo "✅ 安装完成！"
echo ""
echo "📝 访问地址: http://localhost:8080"
echo "📖 文档: https://github.com/yourusername/stardew-panel"
echo ""
echo "常用命令:"
echo "  启动: cd $INSTALL_DIR/docker && docker-compose up -d"
echo "  停止: cd $INSTALL_DIR/docker && docker-compose down"
echo "  日志: cd $INSTALL_DIR/docker && docker-compose logs -f"
echo ""
