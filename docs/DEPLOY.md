# StardewPanel 部署指南

在一台 Linux 服务器上，用 StardewPanel（中文管理面板）+ sdvd/server（社区无头游戏容器）
跑起一个星露谷物语联机服。照着做即可。

---

## 一、前置条件

| 项 | 要求 |
|----|------|
| 服务器 | Linux（x86_64 或 arm64），**内存建议 4G**，2G 吃紧 |
| Docker | 已装 Docker Engine + Compose 插件（`docker compose version` 能出版本） |
| Steam 账号 | **拥有《星露谷物语》的正版账号**，用来下载游戏文件（绕不过去） |
| 端口 | 面板 9090（TCP）、VNC 5800（TCP）、联机 24642（UDP）默认放行 |

> 为什么必须正版 Steam 账号：星露谷没有免费的专用服务器，游戏文件只能用
> 拥有该游戏的账号通过 Steam 下载。这是游戏本身的限制，不是本面板的。

检查 Docker：

```bash
docker compose version   # 出版本号即可；报错说明 Compose 插件没装
```

---

## 二、拿到项目

```bash
# 把 stardew-panel 目录放到服务器上（git clone 或 scp 上传均可）
cd stardew-panel/docker
```

后续所有命令都在 `stardew-panel/docker/` 目录下执行。

---

## 三、填配置

```bash
cp .env.example .env
nano .env          # 或用你顺手的编辑器
```

**至少填这几项：**

```ini
# 你的正版 Steam 账号
STEAM_USERNAME="你的steam用户名"
STEAM_PASSWORD="你的steam密码"

# VNC 网页密码（看游戏画面用），公网部署必须设
VNC_PASSWORD="自己设一个"

# 游戏 API 密钥，公网部署强烈建议设。生成一个：
#   openssl rand -base64 32
API_KEY="把生成的贴这里"

# 面板首次管理员账号（必须设置，密码至少 12 个字符）
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="一条随机的强密码"
```

小内存服务器建议顺手打开（去掉行首 `#`）：

```ini
SERVER_TPS=30      # 降 CPU 占用，NPC 移动略卡但省资源
```

> 安全提醒：`.env` 里有 Steam 密码，别提交到 git、别泄露。
> 项目 `.gitignore` 已忽略 `.env`。

---

## 四、首次登录 Steam（过 Steam Guard）

游戏文件由 `steam-auth` 容器下载，**首次要交互式登录一次**（如果账号开了
Steam Guard 两步验证，需要在这一步输入手机/邮箱验证码）。这一步是交互式的，
**面板网页替代不了**，必须在服务器命令行跑：

```bash
docker compose run --rm -it steam-auth setup
```

按提示登录。成功后会话（refresh token）会持久化到 `steam-session` 卷，
以后重启不用再登。这一步也会开始下载游戏文件（几个 G，视网速几分钟到几十分钟）。

---

## 五、启动

```bash
docker compose up -d
```

首次会拉镜像 + 下载游戏，耐心等。查看进度：

```bash
docker compose logs -f server        # 看游戏容器日志
docker compose ps                    # 看各容器状态
```

当 `server` 容器状态正常、日志出现游戏已加载/正在 host 的字样，就绪。

---

## 六、访问面板

浏览器打开：

```
http://<服务器IP>:9090
```

使用 `.env` 中的 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 登录；首次登录后建议在「设置」页再次修改密码。

面板里能做：
- **服务器页**：启动/停止/重启游戏容器、看在线玩家、拿**联机邀请码**、打开 VNC 看画面
- **设置页**：改最大玩家数、自动过天等游戏参数（改动直达游戏容器）
- **存档页**：备份/恢复存档

---

## 七、玩家怎么进服

1. 在面板「服务器」页复制**邀请码**
2. 发给好友
3. 好友在游戏里「加入联机游戏 → 输入邀请码」即可

如果设了 `SERVER_PASSWORD`，进服后还需在游戏聊天框输入 `!login <密码>`。

---

## 八、日常运维

```bash
# 停服（保留数据）
docker compose stop server

# 全部停并移除容器（不删数据卷，存档/游戏文件都还在）
docker compose down

# 重启
docker compose restart server

# 升级 sdvd 游戏镜像到新版
docker compose pull
docker compose up -d
```

数据都在 Docker 命名卷里，`down` 不会丢：
- `saves` — 存档
- `game-data` — 下载的游戏文件
- `steam-session` — Steam 登录态
- `panel-data` — 面板数据库

---

## 九、安全须知（重要）

1. **面板挂载了宿主的 `docker.sock`**（为了能起停游戏容器）。这等于给了面板
   宿主机 root 级的 docker 权限。**绝对不要把 9090 端口直接裸露到公网**。
   要远程访问，走以下任一方式：
   - Tailscale / WireGuard 等私有网络（推荐）
   - Nginx 反代 + HTTPS + 额外鉴权
   - SSH 隧道

2. **公网可达时，`VNC_PASSWORD` 和 `API_KEY` 必须设**，否则任何人都能看你游戏画面/调 API。

3. 不要把 `ADMIN_PASSWORD` 提交到 git；修改密码后可从 `.env` 中移除初始密码配置。

---

## 十、常见问题

**Q: 点「启动服务器」网页转半天没反应？**
首次启动要下载几个 G 游戏文件，面板已做成异步——点完立即返回，状态会显示
「启动中」。去 `docker compose logs -f server` 看真实进度，就绪后面板状态自动变「运行中」。

**Q: 邀请码显示"获取中…"一直不出来？**
游戏容器还没完全就绪（还在加载存档）。等日志显示 host 完成后刷新即可。

**Q: 2G 服务器能跑吗？**
勉强。游戏本体加载存档就要 1G+，叠加 Docker/Xvfb 会紧张。建议设 `SERVER_TPS=30`、
别装太多 MOD、盯着点内存。4G 才舒服。

**Q: 面板设置页显示"服务器未运行"？**
游戏参数存在游戏容器里，容器没跑时读不到。先在服务器页启动，再回设置页。

---

## 架构参考

```
浏览器 → 面板(9090, 本项目)
              ├── docker compose 起停 ──→ server 容器(sdvd/server)
              └── 代理 REST API ─────────→   ├── 游戏 + SMAPI + Xvfb
                                             ├── VNC 网页(5800)
                                             └── 联机(24642/udp)
                              steam-auth 容器 ── 下载游戏文件
```

游戏容器用的是社区项目
[stardew-valley-dedicated-server/server](https://github.com/stardew-valley-dedicated-server/server)，
无头联机、SMAPI、虚拟显示这些都由它负责；本面板负责编排和中文管理界面。
