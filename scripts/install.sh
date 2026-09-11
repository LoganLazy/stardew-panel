#!/bin/sh

set -eu

echo "StardewPanel Docker 部署"
echo "========================"

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker 未安装，请先安装 Docker Engine。" >&2
  exit 1
fi
if ! docker compose version >/dev/null 2>&1; then
  echo "Docker Compose 插件未安装，请先安装 Compose v2。" >&2
  exit 1
fi

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
PROJECT_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
cd "$PROJECT_DIR/docker"

if [ ! -f .env ]; then
  cp .env.example .env
  echo "已创建 docker/.env，请填写 Steam 账号、ADMIN_PASSWORD、VNC_PASSWORD 和 API_KEY 后重新运行。" >&2
  exit 1
fi

docker compose config >/dev/null
docker compose up -d --build

echo "部署完成：http://<服务器IP>:9090"
echo "首次 Steam Guard 验证：docker compose run --rm -it steam-auth setup"
