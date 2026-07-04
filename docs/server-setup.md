# 星露谷物语服务器配置指南

## 服务端安装

StardewPanel 支持两种服务端模式：

### 方式 1：官方 Dedicated Server（推荐）

星露谷 1.5.5+ 版本自带独立服务器。

**Linux 安装：**

```bash
# 1. 安装 SteamCMD
sudo apt update
sudo apt install steamcmd

# 2. 下载星露谷服务端
steamcmd +login anonymous +app_update 1026920 validate +quit

# 3. 服务端位置
# ~/.steam/steamapps/common/StardewValley/
```

**配置路径：**
- 服务端：`~/.steam/steamapps/common/StardewValley/`
- 启动文件：`StardewValley` (可执行文件)

### 方式 2：SMAPI + MOD 支持

如果需要 MOD 支持，使用 SMAPI（Stardew Modding API）。

**安装 SMAPI：**

```bash
# 1. 下载 SMAPI
wget https://github.com/Pathoschild/SMAPI/releases/latest/download/SMAPI-3.18.6-installer.zip
unzip SMAPI-3.18.6-installer.zip

# 2. 安装
cd SMAPI-3.18.6-installer
chmod +x install.sh
./install.sh

# 3. 服务端启动改用 SMAPI
# 启动文件变为：StardewValley/StardewModdingAPI
```

## StardewPanel 配置

编辑 `config.yaml`：

```yaml
game:
  # 服务端路径
  server_path: "/home/user/.steam/steamapps/common/StardewValley"
  
  # MOD 路径（SMAPI）
  mods_path: "/home/user/.steam/steamapps/common/StardewValley/Mods"
  
  # 存档路径
  saves_path: "/home/user/.config/StardewValley/Saves"
  
  # 启动命令
  # 官方服务端：
  start_command: "./StardewValley"
  # SMAPI 服务端：
  # start_command: "./StardewModdingAPI"
  
  # 服务器端口
  port: 24642
```

## 一键安装脚本（即将支持）

StardewPanel 将提供一键安装功能：

1. 在面板中点击"安装服务器"
2. 选择安装模式（官方 / SMAPI）
3. 自动下载和配置
4. 完成后即可启动

## Docker 部署（推荐）

使用 Docker 最简单，已包含服务端：

```bash
cd docker
docker-compose up -d
```

Docker 镜像已预装：
- 星露谷服务端
- SMAPI（可选启用）
- 常用 MOD

## 手动启动测试

```bash
# 进入服务端目录
cd ~/.steam/steamapps/common/StardewValley

# 启动服务器
./StardewValley
```

首次启动会生成配置文件 `startup_preferences`。

## 常见问题

### 1. Linux 上没有 Steam 怎么办？

使用 SteamCMD（无需 Steam 客户端）：
```bash
sudo apt install steamcmd
```

### 2. ARM64 设备（N1/树莓派）能运行吗？

星露谷官方服务端只支持 x86_64，ARM 需要通过以下方式：
- 使用 Box64 模拟器（性能损耗较大）
- 使用 Windows ARM 版本 + Wine

**推荐方案：** N1 做管理面板，另一台 x86 设备跑游戏服务器。

### 3. 需要正版游戏吗？

服务端需要 Steam 账号下载（可以匿名），但**不需要购买游戏**。客户端玩家需要正版。

## 下一步

安装完成后：
1. 在 StardewPanel 中配置路径
2. 启动服务器
3. 安装 MOD（可选）
4. 邀请玩家加入
