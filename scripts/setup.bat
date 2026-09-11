@echo off
chcp 65001 >nul
echo 🌾 StardewPanel 快速启动脚本
echo ================================
echo.

REM 检查 Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 未安装 Go，请先安装 Go 1.25+
    pause
    exit /b 1
)
echo ✅ Go 已安装

REM 检查 Node.js
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 未安装 Node.js，请先安装 Node.js 20.19+
    pause
    exit /b 1
)
echo ✅ Node.js 已安装

echo.
echo 📁 创建目录...
if not exist "data" mkdir data
if not exist "data\temp" mkdir data\temp
if not exist "game\server" mkdir game\server
if not exist "game\server\Mods" mkdir game\server\Mods
if not exist "game\saves\Saves" mkdir game\saves\Saves
if not exist "game\saves\backups" mkdir game\saves\backups

echo.
echo 📦 下载 Go 依赖...
cd server
go mod download
if %errorlevel% neq 0 (
    echo ❌ 下载 Go 依赖失败
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo 📦 安装前端依赖...
cd web
call npm ci
if %errorlevel% neq 0 (
    echo ❌ 安装前端依赖失败
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo ✅ 环境配置完成！
echo.
echo 启动方式：
echo   1. 开发模式（前后端分离）：
echo      - 后端: cd server ^&^& go run main.go
echo      - 前端: cd web ^&^& npm run dev
echo.
echo   2. Docker 部署：
echo      cd docker ^&^& docker compose up -d
echo.
echo 访问地址：
echo   - 前端开发: http://localhost:3000
echo   - 后端 API: http://localhost:9090
echo   - 健康检查: http://localhost:9090/health
echo.
pause
