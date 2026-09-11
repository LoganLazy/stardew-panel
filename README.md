# StardewPanel - 星露谷物语服务器管理面板

一个轻量级的星露谷物语专用服务器管理面板，专为低配设备（树莓派/N1 盒子）优化。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js)](https://vuejs.org)

> **架构说明（v0.4 起）**：本面板不再自己运行游戏进程，而是作为**中文编排层**，
> 通过 Docker 编排经过验证的社区无头服务器 [`sdvd/server`](https://github.com/stardew-valley-dedicated-server/server)
> （自带 SMAPI + Xvfb + REST API + VNC），并代理其数据。游戏的下载/无头运行由该容器负责，
> 面板负责起停、状态展示、存档、设置代理。**详细部署见 [docs/DEPLOY.md](docs/DEPLOY.md)**。

## ✨ 特性

- 🔐 **用户认证** - 登录/退出/密码修改，保护公网访问
- 🐳 **容器编排** - 一键起停 sdvd/server 游戏容器，无需手动折腾 SMAPI/Xvfb
- 🎮 **MOD 管理** - 可视化管理 MOD（基于 SMAPI）
- 👥 **玩家监控** - 代理游戏 API 拿实时在线玩家，支持中文玩家名
- 🔗 **邀请码** - 面板直接显示联机邀请码，一键复制发好友
- 🖥️ **VNC 画面** - 一键打开游戏画面（需开启渲染）
- 💾 **存档管理** - 停服自动备份、每日安全检查、一键恢复
- ⚙️ **设置代理** - 最大玩家数、自动过天等游戏设置直接在面板改
- 📊 **性能监控** - 容器 CPU/内存/网络占用与面板进程指标实时展示
- 📡 **实时日志** - 游戏容器日志 SSE 推送，断线自动重连
- 🔒 **安全可靠** - 限流防护、会话持久化、密码加密
- 🇨🇳 **中文优先** - 界面和文档全中文

## 🎯 技术栈

- **后端**: Go 1.25 + Gin + SQLite + bcrypt
- **前端**: Vue 3 + Vite + Vue Router + Axios
- **部署**: Docker + Docker Compose

## 🚀 快速开始

**前置要求**：一台 4G 内存的 Linux 服务器（2G 吃紧）、Docker + Docker Compose、**拥有《星露谷物语》的正版 Steam 账号**。

完整步骤见 **[docs/DEPLOY.md](docs/DEPLOY.md)**，概览：

```bash
cd docker
cp .env.example .env          # 填 Steam 账号、面板初始密码、VNC 密码、API_KEY

# 首次：交互式登录 Steam 并下载游戏（过 Steam Guard 验证）
docker compose run --rm -it steam-auth setup

# 起全部服务
docker compose up -d

# 访问面板 http://<服务器IP>:9090
```

> ⚠️ Steam Guard 首次验证是交互式的，只能在命令行完成（`steam-auth setup` 那步），
> 面板无法代替。之后 refresh token 会持久化，不用重复验证。

## 🔐 首次登录

首次启动前必须在 `docker/.env` 设置 `ADMIN_PASSWORD`（至少 12 个字符），
可选设置 `ADMIN_USERNAME`。面板不会再使用固定默认密码。

修改密码步骤：
1. 登录后点击右上角 "⚙️ 设置"
2. 填写修改密码表单
3. 提交后自动退出，使用新密码重新登录

## 📖 功能说明

### 服务器管理
- 启动/停止/重启服务器
- 实时查看服务器状态
- 在线人数统计
- 运行时长显示

### MOD 管理
- 上传 MOD 压缩包（.zip）
- 自动解析 manifest.json
- 一键启用/禁用 MOD
- 删除 MOD

### 玩家监控
- 实时在线玩家列表
- 支持中文玩家名
- 在线时长统计
- 历史玩家记录

### 存档管理
- 创建存档备份
- 一键恢复存档
- 自动压缩为 .zip
- 显示备份大小和时间

### 日志查看
- 实时查看服务器日志
- 自动刷新（3秒）
- 可配置显示行数
- 终端样式显示

### 认证系统
- 用户登录/退出
- 修改密码
- Token 认证（24小时）
- 限流防护（防暴力破解）
- 会话持久化

## 🔒 安全建议

### 生产环境必做

1. **立即修改默认密码**
   ```
   ❌ 不要使用 admin123
   ✅ 使用强密码（12位以上，包含字母数字）
   ```

2. **启用 HTTPS**
   ```bash
   # 使用 Nginx 反向代理
   server {
       listen 443 ssl http2;
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://localhost:9090;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

   使用反向代理时，在 `docker/.env` 将代理的实际 IP 或 CIDR 填入
   `TRUSTED_PROXIES`；默认不信任客户端提供的转发头。不要配置 `0.0.0.0/0`。

3. **限制访问 IP（可选）**
   ```nginx
   # Nginx 配置
   allow 192.168.1.0/24;
   deny all;
   ```

4. **定期备份数据库**
   ```bash
   # 备份数据库文件
   cp data/stardew-panel.db data/backup/stardew-panel-$(date +%Y%m%d).db
   ```

5. **使用防火墙**
   ```bash
   # 只允许特定端口
   ufw allow 443/tcp
   ufw enable
   ```

## 📊 系统要求

### 最低要求
- CPU: 2 核
- 内存: 2GB（仅适合少量 MOD，并建议将 `SERVER_TPS` 调低到 30）
- 磁盘: 20GB

### 推荐配置
- CPU: 2 核
- 内存: 4GB 或更多
- 磁盘: 30GB 或更多

### 支持平台
- 生产部署：Linux x86_64 / arm64，Docker Engine + Compose v2
- 本地开发：Linux、Windows 10/11、macOS（需 Go、Node.js 和 CGO 工具链）
- 低功耗 ARM 设备可运行，但建议至少 4GB 内存

## 🛠️ 开发

### 项目结构

```
stardew-panel/
├── server/              # Go 后端
│   ├── config/         # 配置管理
│   ├── database/       # 数据库
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── router/         # 路由
│   ├── service/        # 业务逻辑
│   └── main.go         # 入口
├── web/                # Vue 前端
│   ├── src/
│   │   ├── api/       # API 调用
│   │   ├── router/    # 路由
│   │   └── views/     # 页面
│   └── index.html
└── docker/             # Docker 配置
```

### 构建

```bash
# 前端构建
cd web
npm run build

# 后端构建
cd server
# go-sqlite3 需要 CGO 和 C 编译器；生产环境推荐使用 Dockerfile 构建
CGO_ENABLED=1 go build -o stardew-panel main.go
```

## 📝 更新日志

版本变更记录见 [CHANGELOG.md](CHANGELOG.md)。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [Stardew Valley](https://www.stardewvalley.net/) - 游戏本体
- [SMAPI](https://smapi.io/) - MOD 加载器
- [Gin](https://gin-gonic.com/) - Go Web 框架
- [Vue.js](https://vuejs.org/) - 前端框架

## 📮 联系方式

- GitHub: [LoganLazy/stardew-panel](https://github.com/LoganLazy/stardew-panel)
- Issues: [提交问题](https://github.com/LoganLazy/stardew-panel/issues)

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**
