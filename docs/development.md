# StardewPanel 开发指南

## 项目结构

```
stardew-panel/
├── server/              # Go 后端
│   ├── config/         # 配置管理
│   ├── database/       # 数据库初始化 ✅
│   ├── handler/        # HTTP 处理器 ✅
│   ├── models/         # 数据模型 ✅
│   ├── router/         # 路由配置 ✅
│   ├── service/        # 业务逻辑 ✅
│   └── main.go         # 入口文件 ✅
├── web/                # Vue 前端
│   ├── src/
│   │   ├── api/        # API 调用 ✅
│   │   ├── views/      # 页面组件 ✅
│   │   ├── router/     # 路由配置 ✅
│   │   └── main.js     # 入口文件 ✅
│   └── package.json
├── docker/             # Docker 配置 ✅
├── config.yaml         # 配置文件
└── README.md
```

## 本地开发

### 前置要求

- Go 1.22+
- Node.js 18+
- SQLite

### 后端开发

```bash
cd server
go mod download
go run main.go
```

后端默认运行在 `http://localhost:8080`

### 前端开发

```bash
cd web
npm install
npm run dev
```

前端开发服务器运行在 `http://localhost:5173`

前端会自动连接到 `http://localhost:8080` 的后端 API

### API 测试

使用健康检查接口测试：

```bash
curl http://localhost:8080/health
```

## Docker 部署

### 构建并启动

```bash
cd docker
docker-compose up -d --build
```

### 查看日志

```bash
docker-compose logs -f stardew-panel
```

### 停止服务

```bash
docker-compose down
```

## 已实现功能 ✅

### 1. 服务器安装

支持三种安装方式：

- **上传文件**：上传游戏压缩包并自动解压
- **指定路径**：验证并使用已有游戏文件
- **SteamCMD**：自动下载游戏文件（需安装 SteamCMD）

### 2. 服务器管理

- 启动/停止/重启服务器
- 实时查看服务器状态
- 进程监控和自动状态更新
- 日志记录

### 3. MOD 管理

- 上传 MOD 文件（支持 .zip）
- 启用/禁用 MOD（通过重命名目录）
- 自动解析 manifest.json
- MOD 信息展示
- 文件系统扫描同步

### 4. 玩家管理

- 白名单 CRUD 操作
- 玩家信息管理
- 踢出玩家接口（需游戏服务器支持）

### 5. 存档管理

- 创建存档备份（压缩为 .zip）
- 恢复存档（自动备份当前存档）
- 删除备份
- 存档文件扫描
- 自动清理旧备份

### 6. 日志查看

- 实时日志读取
- 自动刷新（3秒间隔）
- 可配置日志行数

## 数据库

项目使用 SQLite，数据库文件位于 `data/stardew-panel.db`。

表结构：
- `installation` - 安装配置
- `mods` - MOD 信息
- `whitelist` - 玩家白名单
- `backups` - 存档备份记录
- `server_state` - 服务器运行状态

首次启动时会自动创建所有表。

## API 文档

### 安装相关

- `GET /api/v1/install/check` - 检查安装状态
- `POST /api/v1/install/upload` - 上传游戏文件（multipart/form-data）
- `POST /api/v1/install/verify` - 验证游戏路径
- `POST /api/v1/install/steamcmd` - SteamCMD 安装
- `GET /api/v1/install/status` - 获取安装进度

### 服务器管理

- `GET /api/v1/server/status` - 获取服务器状态
- `POST /api/v1/server/start` - 启动服务器
- `POST /api/v1/server/stop` - 停止服务器
- `POST /api/v1/server/restart` - 重启服务器

### MOD 管理

- `GET /api/v1/mods` - 列出所有 MOD
- `POST /api/v1/mods/upload` - 上传 MOD（multipart/form-data）
- `PUT /api/v1/mods/:id/toggle` - 启用/禁用 MOD
- `DELETE /api/v1/mods/:id` - 删除 MOD

### 玩家管理

- `GET /api/v1/players` - 获取玩家列表和白名单
- `POST /api/v1/players/whitelist` - 添加白名单
- `DELETE /api/v1/players/whitelist/:id` - 移除白名单
- `POST /api/v1/players/kick` - 踢出玩家

### 存档管理

- `GET /api/v1/saves` - 列出存档和备份
- `POST /api/v1/saves/backup` - 创建备份
- `POST /api/v1/saves/restore/:id` - 恢复备份
- `DELETE /api/v1/saves/:id` - 删除备份

### 日志

- `GET /api/v1/logs?lines=100` - 获取日志

## 配置说明

编辑 `config.yaml`：

```yaml
server:
  host: "0.0.0.0"
  port: "8080"

database:
  path: "./data/stardew-panel.db"

game:
  server_path: "./game/server"
  mods_path: "./game/mods"
  saves_path: "./game/saves"
```

## 故障排查

### 后端无法启动

1. 检查端口是否被占用：`netstat -ano | findstr :8080`
2. 检查数据库文件权限
3. 确保 data 目录存在

### 前端无法连接后端

1. 检查后端是否正常运行
2. 确认 `.env.development` 中的 API URL
3. 检查防火墙/CORS 设置

### MOD 无法加载

1. 确保 MOD 文件格式正确
2. 检查 `game/mods` 目录权限
3. 查看服务器日志

### 服务器无法启动

1. 检查游戏是否正确安装
2. 确认可执行文件路径
3. 查看 `data/server.log`

## 待实现功能

- [ ] WebSocket 实时日志推送
- [ ] 自动备份定时任务（Cron）
- [ ] SMAPI 自动下载和安装
- [ ] 多语言支持（i18n）
- [ ] 游戏服务器 Docker 镜像
- [ ] 在线玩家实时监控
- [ ] 性能指标收集
- [ ] 用户认证和权限管理

## 技术栈

**后端：**
- Go 1.22
- Gin Web Framework
- SQLite
- go-sqlite3

**前端：**
- Vue 3
- Vite
- Pinia
- Axios
- Vue Router

**部署：**
- Docker
- Docker Compose
- Alpine Linux

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 代码规范

- Go: 使用 `gofmt` 格式化，遵循官方风格指南
- Vue: 使用 Composition API，组件采用 `<script setup>` 语法
- 提交信息: 遵循 [Conventional Commits](https://www.conventionalcommits.org/)

## 许可证

MIT License
