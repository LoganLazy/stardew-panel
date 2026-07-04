# StardewPanel 安装指南

## 系统要求

- Linux (amd64/arm64)
- Docker 20.10+
- Docker Compose 2.0+
- 最低 1GB RAM / 10GB 磁盘空间

## 快速安装

### 一键安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/stardew-panel/main/scripts/install.sh | bash
```

### 手动安装

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/stardew-panel.git
cd stardew-panel

# 2. 启动服务
cd docker
docker-compose up -d

# 3. 访问面板
# 浏览器打开 http://your-server-ip:8080
```

## 配置

### 修改端口

编辑 `docker/docker-compose.yml`：

```yaml
services:
  stardew-panel:
    ports:
      - "你的端口:8080"  # 修改左边的端口
```

### 数据持久化

所有数据存储在 Docker 卷中：
- `panel-data`: 面板数据（数据库、配置）
- `game-data`: 游戏数据（存档、MOD）

备份命令：
```bash
docker run --rm -v stardew-panel_game-data:/data -v $(pwd):/backup alpine tar czf /backup/backup.tar.gz /data
```

## 常见问题

### 1. 端口被占用

修改 `docker-compose.yml` 中的端口映射。

### 2. 权限问题

确保当前用户在 docker 组：
```bash
sudo usermod -aG docker $USER
```

### 3. N1/树莓派部署

确认架构后拉取对应镜像：
```bash
docker pull --platform linux/arm64 your-image
```

## 卸载

```bash
cd /opt/stardew-panel/docker
docker-compose down -v  # -v 会删除所有数据
```

## 下一步

- [配置星露谷服务器](server-setup.md)
- [安装 MOD](mod-installation.md)
- [常见问题](faq.md)
