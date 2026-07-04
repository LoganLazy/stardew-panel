# StardewPanel

一个轻量级的星露谷物语专用服务器管理面板，专为低配设备（树莓派/N1 盒子）优化。

## 特性

- 🚀 **一键部署** - 单条命令完成服务器安装
- 🎮 **MOD 管理** - 可视化安装/启用/禁用 MOD
- 👥 **玩家管理** - 白名单、踢人、权限控制
- 💾 **存档管理** - 自动备份、一键恢复
- 📊 **实时监控** - 服务器状态、在线玩家、日志查看
- 📱 **移动适配** - 手机也能管理服务器
- 🇨🇳 **中文优先** - 界面和文档全中文

## 技术栈

- **后端**: Go + Gin
- **前端**: Vue 3 + Vite
- **数据库**: SQLite
- **容器**: Docker + Docker Compose

## 快速开始

```bash
# 一键安装（开发中）
curl -fsSL https://raw.githubusercontent.com/yourusername/stardew-panel/main/install.sh | bash
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/yourusername/stardew-panel.git
cd stardew-panel

# 启动后端
cd server
go run main.go

# 启动前端
cd web
npm install
npm run dev
```

## 系统要求

- **最低配置**: 1GB RAM / 10GB 磁盘
- **推荐配置**: 2GB RAM / 20GB 磁盘
- **支持系统**: Linux (amd64/arm64)

## 开源协议

MIT License

## 路线图

- [x] 项目初始化
- [ ] Docker 容器化星露谷服务端
- [ ] Go 后端 API
- [ ] Vue 前端界面
- [ ] MOD 管理功能
- [ ] 玩家管理功能
- [ ] 存档备份功能
- [ ] 一键安装脚本

## 贡献

欢迎提交 Issue 和 Pull Request！

## 鸣谢

- [Stardew Valley](https://www.stardewvalley.net/)
- [SMAPI](https://smapi.io/)
