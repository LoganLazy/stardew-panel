# StardewPanel 开发指南

## 环境要求

- Go 1.25+
- Node.js 20.19+（推荐 Node.js 22 LTS）
- Docker Engine + Docker Compose 插件（完整联调需要）
- 本地直接运行后端时还需要 CGO、C 编译器和 SQLite 开发库

项目使用 `go-sqlite3`。`CGO_ENABLED=0` 可以完成部分编译检查，但生成的程序无法打开数据库。
生产和完整联调优先使用 `docker/Dockerfile.panel`。

## 前端

```bash
cd web
npm ci
npm run dev
```

Vite 默认监听 `http://localhost:3000`，并把 `/api` 代理到 `http://localhost:9090`。

生产构建：

```bash
cd web
npm run build
```

## 后端

首次运行必须设置管理员初始密码：

```bash
cd server
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD='至少十二个字符的强密码'
CGO_ENABLED=1 go run .
```

Windows PowerShell：

```powershell
cd server
$env:ADMIN_USERNAME = 'admin'
$env:ADMIN_PASSWORD = '至少十二个字符的强密码'
$env:CGO_ENABLED = '1'
go run .
```

后端默认监听 `http://localhost:9090`，健康检查为 `GET /health`。

## 完整容器联调

```bash
cd docker
cp .env.example .env
# 填写 STEAM_USERNAME、STEAM_PASSWORD、ADMIN_PASSWORD、VNC_PASSWORD、API_KEY
docker compose config
docker compose build stardew-panel
docker compose up -d
```

首次 Steam Guard 验证仍需在命令行完成：

```bash
docker compose run --rm -it steam-auth setup
```

## 检查命令

```bash
cd server
go test ./...
go vet ./...

cd ../web
npm run build
```

涉及 SQLite 的集成测试必须在 `CGO_ENABLED=1` 且存在 C 编译器的环境中运行。

## 当前 API

所有 `/api/v1` 接口（登录除外）都要求 `Authorization: Bearer <token>`。

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/check`
- `POST /api/v1/auth/change-password`
- `GET /api/v1/install/check`
- `GET /api/v1/install/setup`
- `GET /api/v1/server/status`
- `POST /api/v1/server/start`
- `POST /api/v1/server/stop`
- `POST /api/v1/server/restart`
- `GET /api/v1/server/invite-code`
- `GET /api/v1/server/settings`
- `PUT /api/v1/server/settings/:key`
- `GET /api/v1/mods`
- `POST /api/v1/mods/upload`
- `PUT /api/v1/mods/:id/toggle`
- `DELETE /api/v1/mods/:id`
- `GET /api/v1/players`
- `POST /api/v1/players/refresh`
- `DELETE /api/v1/players/cleanup`
- `GET /api/v1/saves`
- `POST /api/v1/saves/backup`
- `POST /api/v1/saves/restore/:id`
- `DELETE /api/v1/saves/:id`
- `GET /api/v1/logs`

`POST /api/v1/players/kick` 仅为兼容旧客户端保留；上游游戏 API 未提供可靠的踢人接口，当前返回 `501`，前端不展示该操作。

## 目录约定

- `data/`：SQLite 数据库和上传临时文件
- `game/server/Mods`：与游戏容器 `game-data` 卷共享的 MOD 目录
- `game/saves/Saves`：游戏活动存档目录
- `game/saves/backups`：面板备份目录，不放入游戏活动存档
- `web/dist`：前端生产构建产物
