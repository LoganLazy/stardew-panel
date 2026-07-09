# StardewPanel - 星露谷物语服务器管理面板

一个轻量级的星露谷物语专用服务器管理面板，专为低配设备（树莓派/N1 盒子）优化。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org)

## ✨ 特性

- 🔐 **用户认证** - 登录/退出/密码修改，保护公网访问 🆕
- 🚀 **灵活安装** - 支持上传文件/指定路径/SteamCMD 三种安装方式
- 🎮 **MOD 管理** - 可视化安装/启用/禁用 MOD（基于 SMAPI）
- 👥 **玩家监控** - 实时在线玩家、支持中文玩家名 🆕
- 💾 **存档管理** - 自动备份、一键恢复
- 📊 **实时监控** - 服务器状态、在线玩家、日志查看
- ⚡ **性能优化** - 大日志文件优化，99% 性能提升 🆕
- 🔒 **安全可靠** - 限流防护、会话持久化、密码加密
- 📱 **移动适配** - 手机也能管理服务器
- 🇨🇳 **中文优先** - 界面和文档全中文
- 🌍 **跨平台** - 支持 Linux AMD64/ARM64

## 🎯 技术栈

- **后端**: Go 1.22 + Gin + SQLite + bcrypt
- **前端**: Vue 3 + Vite + Pinia + Axios
- **部署**: Docker + Docker Compose

## 🚀 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 克隆仓库
git clone https://github.com/LoganLazy/stardew-panel.git
cd stardew-panel

# 启动服务
cd docker
docker-compose up -d --build

# 等待 2-5 分钟构建完成
# 访问 http://localhost:8080
```

### 方式二：本地开发

**前置要求**: Go 1.22+, Node.js 18+

```bash
# 后端
cd server
go mod download
go run main.go

# 前端（新终端）
cd web
npm install
npm run dev

# 访问 http://localhost:5173
```

## 🔐 默认账号

首次启动后，控制台会显示默认账号：

```
用户名: admin
密码: admin123
```

⚠️ **请立即登录后修改密码！**

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
   ✅ 使用强密码（8位以上，包含字母数字）
   ```

2. **启用 HTTPS**
   ```bash
   # 使用 Nginx 反向代理
   server {
       listen 443 ssl http2;
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

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
- CPU: 1 核
- 内存: 512MB
- 磁盘: 10GB

### 推荐配置
- CPU: 2 核
- 内存: 2GB
- 磁盘: 20GB

### 支持平台
- Linux (Ubuntu/Debian/CentOS)
- Windows 10/11
- macOS 10.15+
- ARM 设备（树莓派、N1 盒子等）

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
│   │   ├── components/# 组件
│   │   ├── router/    # 路由
│   │   ├── utils/     # 工具函数
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
go build -o stardew-panel main.go
```

## 📝 更新日志

### v0.3.0 (2026-07-05) - 安全增强
- 🔐 新增用户认证系统
- 🔒 新增限流防护（防暴力破解）
- 💾 新增会话持久化（数据库）
- 🎨 新增 404 错误页面
- 🔄 新增 Loading 组件
- ⚡ 优化健康检查端点
- 📱 新增 PWA 支持

### v0.2.3 (2026-07-05) - 工具增强
- 🛠️ 新增 30+ 工具函数
- 📝 优化确认提示

### v0.2.2 (2026-07-05) - 性能优化
- ⚡ 大日志文件优化（99% 性能提升）
- 🔒 路径注入防护
- 📊 Server 页面在线人数显示

### v0.2.1 (2026-07-05) - Bug 修复
- 🌏 支持中文玩家名
- 🔐 密码泄露防护
- 📦 文件大小限制

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
