# StardewPanel 开发指南

## 项目结构

```
stardew-panel/
├── server/          # Go 后端
│   ├── main.go      # 入口文件
│   ├── config/      # 配置管理
│   ├── handler/     # API 处理器
│   ├── router/      # 路由
│   └── service/     # 业务逻辑（待实现）
├── web/             # Vue 3 前端
│   ├── src/
│   │   ├── views/   # 页面
│   │   ├── stores/  # Pinia 状态管理
│   │   └── router/  # 路由
│   └── package.json
├── docker/          # Docker 配置
├── scripts/         # 安装脚本
└── docs/            # 文档
```

## 本地开发

### 后端开发

```bash
cd server
go mod download
go run main.go
```

后端将在 `http://localhost:8080` 启动。

### 前端开发

```bash
cd web
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动，API 请求会自动代理到 `http://localhost:8080`。

## API 设计

### 服务器管理
- `GET /api/v1/server/status` - 获取服务器状态
- `POST /api/v1/server/start` - 启动服务器
- `POST /api/v1/server/stop` - 停止服务器
- `POST /api/v1/server/restart` - 重启服务器

### MOD 管理
- `GET /api/v1/mods` - 获取 MOD 列表
- `POST /api/v1/mods/upload` - 上传 MOD
- `PUT /api/v1/mods/:id/toggle` - 启用/禁用 MOD
- `DELETE /api/v1/mods/:id` - 删除 MOD

### 玩家管理
- `GET /api/v1/players` - 获取玩家列表
- `POST /api/v1/players/whitelist` - 添加白名单
- `DELETE /api/v1/players/whitelist/:id` - 移除白名单
- `POST /api/v1/players/kick` - 踢出玩家

### 存档管理
- `GET /api/v1/saves` - 获取存档列表
- `POST /api/v1/saves/backup` - 创建备份
- `POST /api/v1/saves/restore/:id` - 恢复存档
- `DELETE /api/v1/saves/:id` - 删除存档

### 日志
- `GET /api/v1/logs` - 获取日志

## 下一步开发

1. **服务器管理核心**
   - 实现 Docker 容器管理
   - 星露谷服务端启停逻辑
   - 实时状态监控

2. **MOD 管理**
   - 文件上传处理
   - MOD 解析和验证
   - 启用/禁用机制

3. **存档管理**
   - 自动备份定时任务
   - 存档压缩和恢复
   - 云端备份集成（可选）

4. **实时通信**
   - WebSocket 支持
   - 日志实时推送
   - 服务器状态实时更新

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 代码规范

- Go: 使用 `gofmt` 格式化
- Vue: 使用 ESLint + Prettier
- 提交信息: 遵循 [Conventional Commits](https://www.conventionalcommits.org/)
