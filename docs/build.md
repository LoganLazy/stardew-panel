# 跨平台编译脚本

## Linux AMD64
```bash
GOOS=linux GOARCH=amd64 go build -o stardew-panel-linux-amd64 main.go
```

## Linux ARM64 (N1/树莓派)
```bash
GOOS=linux GOARCH=arm64 go build -o stardew-panel-linux-arm64 main.go
```

## Windows AMD64
```bash
GOOS=windows GOARCH=amd64 go build -o stardew-panel-windows-amd64.exe main.go
```

## 多架构 Docker 镜像构建
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t stardew-panel:latest .
```
