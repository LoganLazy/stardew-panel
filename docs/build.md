# 构建说明

## 推荐方式：Docker

Dockerfile 已安装 `go-sqlite3` 所需的 C 编译器和 SQLite 运行库：

```bash
cd docker
docker compose build stardew-panel
```

多架构镜像：

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f docker/Dockerfile.panel \
  -t stardew-panel:latest .
```

## 本机构建

后端依赖 CGO，不能使用纯 Go 的 `CGO_ENABLED=0` 交叉编译。

Linux 需要 gcc、musl-dev/glibc-dev 和 SQLite 开发库，然后执行：

```bash
cd server
CGO_ENABLED=1 go build -o stardew-panel .
```

Windows 需要可供 Go 使用的 MinGW-w64 C 编译器：

```powershell
cd server
$env:CGO_ENABLED = '1'
go build -o stardew-panel.exe .
```

跨架构构建还需要对应目标架构的 C 交叉编译器；仅设置 `GOOS`/`GOARCH` 不足以构建 `go-sqlite3`。

前端构建：

```bash
cd web
npm ci
npm run build
```
