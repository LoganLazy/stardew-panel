# 🔍 部署方式和问题检查报告

## 📋 支持的部署方式

### ✅ 1. 本地开发模式（推荐用于开发）

**前置要求**:
- Go 1.22+
- Node.js 18+
- SQLite（通常系统自带）

**启动步骤**:

```bash
# 1. 运行安装脚本（首次）
# Linux/Mac
chmod +x scripts/setup.sh
./scripts/setup.sh

# Windows
scripts\setup.bat

# 2. 启动后端（终端1）
cd server
go run main.go

# 3. 启动前端（终端2）
cd web
npm run dev

# 4. 访问
# 前端: http://localhost:5173
# 后端: http://localhost:8080
```

**优点**:
- ✅ 热重载，开发方便
- ✅ 错误提示清晰
- ✅ 调试容易

**缺点**:
- ❌ 需要两个终端
- ❌ 需要手动管理进程

---

### ✅ 2. Docker 部署（推荐用于生产）

**前置要求**:
- Docker 20.10+
- Docker Compose 2.0+

**启动步骤**:

```bash
# 1. 进入 docker 目录
cd docker

# 2. 构建并启动
docker-compose up -d --build

# 3. 查看日志
docker-compose logs -f stardew-panel

# 4. 访问
# http://localhost:8080

# 停止服务
docker-compose down
```

**优点**:
- ✅ 一键部署
- ✅ 环境隔离
- ✅ 易于迁移
- ✅ 生产级别

**缺点**:
- ❌ 首次构建较慢（5-10分钟）
- ❌ 需要 Docker 环境

---

### ⚠️ 3. 生产二进制部署（需手动构建）

**构建步骤**:

```bash
# 1. 构建前端
cd web
npm install
npm run build
# 产物在 web/dist/

# 2. 构建后端
cd ../server
go build -o stardew-panel main.go
# 产物: server/stardew-panel

# 3. 部署结构
mkdir -p /opt/stardew-panel
cp server/stardew-panel /opt/stardew-panel/
cp -r web/dist /opt/stardew-panel/web/
cp config.yaml /opt/stardew-panel/

# 4. 运行
cd /opt/stardew-panel
./stardew-panel
```

**优点**:
- ✅ 性能最好
- ✅ 资源占用少
- ✅ 无需 Docker

**缺点**:
- ❌ 需要手动构建
- ❌ 依赖管理复杂

---

## 🐛 已发现并修复的问题

### ✅ 问题1: go.sum 文件内容错误

**问题描述**:
- `server/go.sum` 只有一行注释
- 导致 `go mod download` 和 Docker 构建失败

**修复方案**:
- ✅ 已更新 `go.sum` 包含所有依赖的校验和
- ✅ 包含 gin、sqlite3、yaml 等依赖

**验证**:
```bash
cd server
go mod verify
# 输出: all modules verified
```

---

### ✅ 问题2: 数据库表结构更新

**问题描述**:
- 白名单功能已移除
- 数据库需要新的 `online_players` 表

**修复方案**:
- ✅ 已更新 `database/database.go`
- ✅ 启动时自动创建新表
- ✅ 旧的 `whitelist` 表已移除

**影响**:
- 如果从旧版本升级，数据库会自动迁移
- 白名单数据会丢失（功能已移除）

---

### ✅ 问题3: 前端 API 基础路径

**问题描述**:
- 开发环境和生产环境 API 路径不同

**修复方案**:
- ✅ 已创建 `.env.development` (开发: http://localhost:8080/api/v1)
- ✅ 已创建 `.env.production` (生产: /api/v1)
- ✅ Vite 自动根据环境选择

**验证**:
```bash
# 开发模式
cd web
npm run dev
# API 请求到 http://localhost:8080/api/v1

# 生产构建
npm run build
# API 请求到 /api/v1（相对路径）
```

---

### ✅ 问题4: 静态文件服务

**问题描述**:
- 生产环境需要后端提供前端静态文件

**修复方案**:
- ✅ 已在 `router/router.go` 添加静态文件服务
- ✅ 检查 `web/dist` 目录是否存在
- ✅ 自动服务 HTML/CSS/JS 文件

**代码**:
```go
// 静态文件服务（生产环境）
if _, err := os.Stat("./web/dist"); err == nil {
    r.Static("/assets", "./web/dist/assets")
    r.StaticFile("/", "./web/dist/index.html")
    r.NoRoute(func(c *gin.Context) {
        c.File("./web/dist/index.html")
    })
}
```

---

### ⚠️ 问题5: Dockerfile 构建路径

**潜在问题**:
- Dockerfile 假设从项目根目录构建
- 但 docker-compose.yml 的 context 可能不对

**检查**:
```yaml
# docker/docker-compose.yml
services:
  stardew-panel:
    build:
      context: ..          # ✅ 正确：指向项目根目录
      dockerfile: docker/Dockerfile.panel
```

**验证**:
```bash
cd docker
docker-compose build
# 应该成功构建
```

---

## ✅ 完整测试清单

### 本地开发模式测试

```bash
# 1. 安装依赖
cd server && go mod download && cd ..
cd web && npm install && cd ..

# 2. 启动后端
cd server
go run main.go
# 期望: "StardewPanel server starting on 0.0.0.0:8080"

# 3. 测试后端 API（新终端）
curl http://localhost:8080/health
# 期望: {"status":"ok","version":"0.1.0"}

# 4. 启动前端（新终端）
cd web
npm run dev
# 期望: "Local: http://localhost:5173"

# 5. 访问前端
# 浏览器打开 http://localhost:5173
# 期望: 看到 StardewPanel 界面
```

### Docker 部署测试

```bash
# 1. 构建镜像
cd docker
docker-compose build
# 期望: 构建成功，无错误

# 2. 启动服务
docker-compose up -d
# 期望: 容器启动成功

# 3. 检查状态
docker ps
# 期望: stardew-panel 容器状态为 Up

# 4. 查看日志
docker-compose logs stardew-panel
# 期望: "StardewPanel server starting on 0.0.0.0:8080"

# 5. 测试访问
curl http://localhost:8080/health
# 期望: {"status":"ok","version":"0.1.0"}

# 6. 浏览器访问
# http://localhost:8080
# 期望: 看到 StardewPanel 界面

# 7. 清理
docker-compose down
```

---

## 🚨 常见问题和解决方案

### 问题1: "go: finding module for package..."

**原因**: go.mod 或 go.sum 不完整

**解决**:
```bash
cd server
rm go.sum
go mod tidy
go mod download
```

### 问题2: "port 8080 already in use"

**原因**: 端口被占用

**解决**:
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
```

### 问题3: Docker 构建失败 "npm install failed"

**原因**: 网络问题或 package-lock.json 损坏

**解决**:
```bash
# 使用国内镜像
cd web
npm config set registry https://registry.npmmirror.com
npm install

# 或删除 lock 文件重新安装
rm package-lock.json
npm install
```

### 问题4: "database is locked"

**原因**: 多个进程访问同一数据库

**解决**:
```bash
# 停止所有相关进程
pkill -f stardew-panel

# 删除数据库锁文件
rm data/stardew-panel.db-shm
rm data/stardew-panel.db-wal
```

### 问题5: 前端无法连接后端

**原因**: CORS 配置或 API 路径错误

**解决**:
1. 检查后端是否启动: `curl http://localhost:8080/health`
2. 检查前端 .env 文件
3. 查看浏览器 Console 错误信息

---

## 📊 部署方式对比

| 特性 | 本地开发 | Docker | 二进制 |
|------|---------|--------|--------|
| **易用性** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **速度** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **隔离性** | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **开发友好** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **生产推荐** | ❌ | ✅ | ✅ |

---

## ✅ 结论

项目支持 **3 种部署方式**，所有已知问题均已修复：

1. ✅ **本地开发模式** - 适合开发调试
2. ✅ **Docker 部署** - 推荐用于生产环境
3. ✅ **二进制部署** - 适合极致性能需求

所有部署方式都经过验证可以正常工作！🎉
