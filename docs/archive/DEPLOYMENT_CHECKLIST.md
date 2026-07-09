# StardewPanel 部署检查清单

## 📋 部署前检查

### 环境准备
- [ ] Go 1.22+ 已安装（开发模式）
- [ ] Node.js 18+ 已安装（开发模式）
- [ ] Docker 和 Docker Compose 已安装（容器部署）
- [ ] 端口 8080 可用
- [ ] 端口 5173 可用（开发模式前端）

### 文件结构
- [ ] server/ 目录完整
- [ ] web/ 目录完整
- [ ] docker/ 目录完整
- [ ] config.yaml 存在
- [ ] README.md 完整

## 🔧 首次运行

### 方式一：开发模式

#### 后端
```bash
cd server
go mod download
go run main.go
```

**验证**：
- [ ] 访问 http://localhost:8080/health 返回 {"status":"ok"}
- [ ] 控制台无错误信息
- [ ] data/stardew-panel.db 文件已创建

#### 前端
```bash
cd web
npm install
npm run dev
```

**验证**：
- [ ] 访问 http://localhost:5173 能看到界面
- [ ] 控制台无错误信息
- [ ] 能正常访问所有页面

### 方式二：Docker 部署

```bash
cd docker
docker-compose up -d --build
```

**验证**：
- [ ] `docker ps` 能看到 stardew-panel 容器
- [ ] 容器状态为 Up
- [ ] 访问 http://localhost:8080 能看到界面
- [ ] `docker-compose logs` 无错误信息

## ✅ 功能测试

### 1. 安装测试
- [ ] 访问 /setup 页面
- [ ] 三种安装方式都能正常显示
- [ ] 上传文件功能正常（可选）
- [ ] 路径验证功能正常（可选）

### 2. 服务器管理
- [ ] 访问 / (Server) 页面
- [ ] 服务器状态正常显示
- [ ] 启动/停止按钮可用（需先安装游戏）

### 3. MOD 管理
- [ ] 访问 /mods 页面
- [ ] 能显示 MOD 列表（可能为空）
- [ ] 上传按钮可用

### 4. 玩家管理
- [ ] 访问 /players 页面
- [ ] 白名单列表正常显示
- [ ] 添加玩家功能正常

### 5. 存档管理
- [ ] 访问 /saves 页面
- [ ] 备份列表正常显示
- [ ] 创建备份按钮可用

### 6. 日志查看
- [ ] 访问 /logs 页面
- [ ] 日志区域正常显示
- [ ] 刷新按钮可用

## 🐛 常见问题排查

### 后端无法启动
**问题**：`listen tcp :8080: bind: address already in use`
**解决**：
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux
lsof -i :8080
kill -9 <PID>
```

### 前端无法连接后端
**问题**：`Network Error` 或 CORS 错误
**检查**：
1. 后端是否正常运行
2. `.env.development` 中 API_URL 是否正确
3. 浏览器控制台查看具体错误

### 数据库错误
**问题**：`database is locked` 或权限错误
**解决**：
```bash
# 检查文件权限
ls -la data/stardew-panel.db

# 修复权限
chmod 644 data/stardew-panel.db
```

### Docker 构建失败
**问题**：构建过程中断或失败
**解决**：
```bash
# 清理缓存重新构建
docker-compose down
docker system prune -a
docker-compose up -d --build
```

## 📊 性能检查

### 资源使用
- [ ] 后端内存使用 < 100MB（空闲时）
- [ ] 前端构建成功
- [ ] Docker 容器正常运行

### 响应时间
- [ ] API 响应时间 < 500ms
- [ ] 页面加载时间 < 2s
- [ ] 文件上传流畅

## 🚀 生产部署建议

### 安全配置
- [ ] 修改默认端口（可选）
- [ ] 配置反向代理（nginx/caddy）
- [ ] 添加 HTTPS 证书
- [ ] 配置防火墙规则
- [ ] 设置文件权限

### 性能优化
- [ ] 启用 gzip 压缩
- [ ] 配置静态文件缓存
- [ ] 数据库定期备份
- [ ] 日志文件轮转

### 监控
- [ ] 配置健康检查
- [ ] 设置日志收集
- [ ] 配置告警通知

## ✨ 完成

如果所有检查项都通过，恭喜你！StardewPanel 已经成功部署并可以使用了。

访问 http://localhost:8080 开始管理你的星露谷物语服务器吧！

---

**遇到问题？**
1. 查看 [开发文档](docs/development.md)
2. 查看 [实现总结](IMPLEMENTATION.md)
3. 提交 [Issue](https://github.com/LoganLazy/stardew-panel/issues)
