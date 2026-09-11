# 贡献指南

感谢参与 StardewPanel。

## 开发检查

提交前请至少运行：

```bash
cd web
npm ci
npm audit --audit-level=high
npm run build

cd ../server
go test ./...
go test -race ./...
go vet ./...
```

涉及容器配置时还应运行：

```bash
docker compose --env-file docker/.env.example -f docker/docker-compose.yml config --quiet
docker build -f docker/Dockerfile.panel -t stardew-panel:local .
```

不要提交真实凭据、数据库、游戏文件、存档、日志、`node_modules` 或前端构建产物。
