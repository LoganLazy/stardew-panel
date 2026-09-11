# 安装指南

项目当前使用 Docker 编排 `stardew-panel`、`sdvd/server` 和 `steam-auth` 三个服务。
旧版 SteamCMD 下载、手动启动 StardewValley、8080 面板的流程已移除。

## 快速安装

```bash
git clone https://github.com/LoganLazy/stardew-panel.git
cd stardew-panel/docker
cp .env.example .env
nano .env
docker compose config
docker compose build
docker compose run --rm -it steam-auth setup
docker compose up -d
```

`.env` 至少需要填写：

- `STEAM_USERNAME` / `STEAM_PASSWORD`
- `ADMIN_PASSWORD`（至少 12 个字符）
- `VNC_PASSWORD`
- `API_KEY`

面板地址为 `http://<服务器IP>:9090`，VNC 地址为 `http://<服务器IP>:5800`。

## 数据和端口

- `panel-data`：面板数据库
- `game-data`：游戏文件和 Mods
- `saves`：存档
- `steam-session`：Steam 登录态
- `server-settings`：游戏服务器设置

默认端口：面板 TCP 9090、VNC TCP 5800、联机 UDP 24642、查询 UDP 27015。
修改宿主端口请编辑 `.env`，不要修改容器内部端口。

## 卸载

```bash
docker compose down       # 保留数据卷
docker compose down -v    # 同时删除数据卷，谨慎执行
```

完整部署说明见 [DEPLOY.md](DEPLOY.md)。
