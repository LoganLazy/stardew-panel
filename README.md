# StardewPanel - 星露谷物语服务器管理面板

一个轻量级的星露谷物语专用服务器管理面板，专为低配设备（树莓派/N1 盒子）优化。

## ✨ 特性

- 🚀 **灵活安装** - 支持上传文件/指定路径/SteamCMD 三种安装方式
- 🎮 **MOD 管理** - 可视化安装/启用/禁用 MOD（基于 SMAPI）
- 👥 **玩家管理** - 白名单、踢人、权限控制
- 💾 **存档管理** - 自动备份、一键恢复
- 📊 **实时监控** - 服务器状态、在线玩家、日志查看
- 📱 **移动适配** - 手机也能管理服务器
- 🇨🇳 **中文优先** - 界面和文档全中文
- 🌍 **跨平台** - 支持 Linux AMD64/ARM64

## 🎯 技术栈

- **后端**: Go 1.22 + Gin
- **前端**: Vue 3 + Vite + Pinia
- **数据库**: SQLite
- **容器**: Docker + Docker Compose

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 克隆仓库
git clone https://github.com/yourusername/stardew-panel.git
cd stardew-panel

# 启动服务
cd docker
docker-compose up -d

# 访问面板
# 浏览器打开 http://localhost:8080
```

### 本地开发

**后端：**
```bash
cd server
go mod download
go run main.go
```

**前端：**
```bash
cd web
npm install
npm run dev
```

## 📖 安装方式

StardewPanel 提供三种灵活的游戏文件配置方式：

### 1. 📦 上传游戏文件
- 从 Steam 导出游戏文件压缩包
- 上传到服务器
- 自动解压和配置

### 2. 📁 指定游戏路径
- 服务器已有游戏文件
- 直接指定路径
- 最快速的方式

### 3. ⬇️ SteamCMD 下载
- 使用 SteamCMD 自动下载
- 免费，无需购买游戏
- 国内可能需要代理

## 🔧 系统要求

- **最低配置**: 1GB RAM / 10GB 磁盘
- **推荐配置**: 2GB RAM / 20GB 磁盘
- **支持系统**: Linux (amd64/arm64)
- **依赖**: Docker 20.10+ / Docker Compose 2.0+

## 📱 跨平台支持

**管理面板**：
- ✅ Linux AMD64 (x86_64 服务器)
- ✅ Linux ARM64 (树莓派/N1 盒子)

**游戏服务端**：
- ✅ Linux AMD64 (官方支持)
- ⚠️ Linux ARM64 (需要模拟器，性能较差)

**推荐部署方案**：
- N1/树莓派运行管理面板
- x86 服务器运行游戏服务端
- 通过网络连接管理

## 📚 文档

- [安装指南](docs/installation.md)
- [服务器配置](docs/server-setup.md)
- [开发指南](docs/development.md)
- [跨平台构建](docs/build.md)

## 🛣️ 路线图

- [x] 项目框架搭建
- [x] 三种安装方式
- [x] 基础 UI 设计
- [ ] Docker 容器化游戏服务端
- [ ] 服务器启停管理
- [ ] MOD 上传和管理
- [ ] 玩家白名单功能
- [ ] 存档自动备份
- [ ] WebSocket 实时日志
- [ ] 一键安装脚本

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 开源协议

MIT License

## 🙏 鸣谢

- [Stardew Valley](https://www.stardewvalley.net/) - 游戏本体
- [SMAPI](https://smapi.io/) - MOD 加载器
- [ValleyServer](https://github.com/Lixeer/ValleyServer) - 参考项目

---

**Star ⭐ 这个项目以支持开发！**
