# 游戏服务器说明

游戏服务由 `sdvd/server` 容器负责，容器内已经包含 SMAPI、Xvfb、REST API 和 VNC。
面板不再直接启动 `StardewValley --server`，也不再执行 SteamCMD 或手动 SMAPI 安装。

## 首次启动

```bash
cd docker
cp .env.example .env
# 填写 Steam 账号、ADMIN_PASSWORD、VNC_PASSWORD、API_KEY
docker compose run --rm -it steam-auth setup
docker compose up -d
```

Steam Guard 必须在命令行交互完成。游戏文件、MOD 和存档分别通过 `game-data`、
`saves` 卷持久化。

## MOD 和存档

- MOD 通过面板上传 `.zip`，会安全解压到共享的 `game-data` 卷。
- 存档通过面板创建备份、恢复和删除。
- 恢复操作会先验证备份并暂存当前存档，避免损坏备份导致原存档丢失。

## 联机和画面

- 面板服务器页显示游戏容器生成的邀请码。
- VNC 画面入口为宿主机 `5800` 端口，必须设置 `VNC_PASSWORD`。
- 玩家连接使用 `GAME_PORT`（默认 UDP 24642）。

详细步骤和安全限制见 [DEPLOY.md](DEPLOY.md)。
