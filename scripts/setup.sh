#!/bin/bash

echo "🌾 StardewPanel 快速启动脚本"
echo "================================"

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ 未安装 Go，请先安装 Go 1.22+"
    exit 1
fi
echo "✅ Go 版本: $(go version)"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 未安装 Node.js，请先安装 Node.js 18+"
    exit 1
fi
echo "✅ Node.js 版本: $(node --version)"

# 创建必要的目录
echo ""
echo "📁 创建目录..."
mkdir -p data
mkdir -p data/temp
mkdir -p game/server
mkdir -p game/mods
mkdir -p game/saves

# 下载后端依赖
echo ""
echo "📦 下载 Go 依赖..."
cd server
go mod download
if [ $? -ne 0 ]; then
    echo "❌ 下载 Go 依赖失败"
    exit 1
fi
cd ..

# 安装前端依赖
echo ""
echo "📦 安装前端依赖..."
cd web
npm install
if [ $? -ne 0 ]; then
    echo "❌ 安装前端依赖失败"
    exit 1
fi
cd ..

echo ""
echo "✅ 环境配置完成！"
echo ""
echo "启动方式："
echo "  1. 开发模式（前后端分离）："
echo "     - 后端: cd server && go run main.go"
echo "     - 前端: cd web && npm run dev"
echo ""
echo "  2. Docker 部署："
echo "     cd docker && docker-compose up -d"
echo ""
echo "访问地址："
echo "  - 前端开发: http://localhost:5173"
echo "  - 后端 API: http://localhost:8080"
echo "  - 健康检查: http://localhost:8080/health"
