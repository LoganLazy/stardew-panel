<template>
  <div class="setup-view">
    <div class="header">
      <h1>服务器<em>安装</em></h1>
      <p>本面板通过编排 sdvd/server 容器运行星露谷无头服务器</p>
    </div>

    <!-- 就绪状态总览 -->
    <div class="card">
      <h3>🔍 就绪状态</h3>
      <div class="status-list">
        <div class="status-item">
          <span class="status-dot" :class="setup.docker_available ? 'ok' : 'bad'"></span>
          <div class="status-text">
            <strong>Docker 环境</strong>
            <small>{{ setup.docker_available ? '已就绪，面板可控制游戏容器' : '不可用，请在宿主安装 Docker 且确保面板挂载了 docker.sock' }}</small>
          </div>
        </div>

        <div class="status-item">
          <span class="status-dot" :class="setup.steam_configured ? 'ok' : 'bad'"></span>
          <div class="status-text">
            <strong>Steam 账号</strong>
            <small v-if="setup.steam_configured">已配置：{{ setup.steam_username }}（用于下载正版游戏文件）</small>
            <small v-else>未配置，请在 .env 填入 STEAM_USERNAME / STEAM_PASSWORD</small>
          </div>
        </div>

        <div class="status-item">
          <span class="status-dot" :class="setup.game_ready ? 'ok' : 'wait'"></span>
          <div class="status-text">
            <strong>游戏文件</strong>
            <small v-if="setup.game_ready">steam-auth 容器已就绪</small>
            <small v-else>尚未就绪，需首次运行下方的下载命令</small>
          </div>
        </div>
      </div>

      <button class="btn btn-secondary" @click="refresh" :disabled="loading">
        {{ loading ? '检查中...' : '重新检查' }}
      </button>
    </div>

    <!-- 部署引导 -->
    <div class="card">
      <h3>📖 部署步骤</h3>
      <p class="intro">
        星露谷没有官方无头服务器，本面板编排社区方案
        <a href="https://github.com/stardew-valley-dedicated-server/server" target="_blank">sdvd/server</a>
        （SMAPI + 无头 MOD + 虚拟显示都封装在容器内）。
        游戏文件需用你的<strong>正版 Steam 账号</strong>下载，首次要在服务器命令行操作一次。
      </p>

      <ol class="steps">
        <li>
          <strong>填写配置</strong>
          <p>在项目 <code>docker/</code> 目录：</p>
          <pre>cp .env.example .env</pre>
          <p>编辑 <code>.env</code>，填入 Steam 账号、VNC 密码、API_KEY（用 <code>openssl rand -base64 32</code> 生成）。</p>
        </li>
        <li>
          <strong>首次登录 Steam 并下载游戏</strong>
          <p>这一步是交互式的（要过 Steam Guard 验证码），网页替代不了，需在服务器上跑：</p>
          <pre>docker compose run --rm -it steam-auth setup</pre>
          <p>按提示输入验证码，完成后会话会持久化，之后无需重复。</p>
        </li>
        <li>
          <strong>启动全部服务</strong>
          <pre>docker compose up -d</pre>
          <p>之后回到本面板「服务器」页点启动即可，无需再碰命令行。</p>
        </li>
      </ol>

      <div class="alert alert-warning">
        <strong>⚠️ 注意</strong>
        <ul>
          <li>需要拥有星露谷物语的正版 Steam 账号，游戏文件由 Steam 官方下载</li>
          <li>首次下载游戏约几 GB，耗时取决于网络；建议 4G 内存以上服务器</li>
          <li>面板端口请勿裸露公网（它持有 docker 控制权限），建议配合 Tailscale 或反代加认证</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as api from '@/api/api'

const loading = ref(false)
const setup = ref({
  docker_available: false,
  steam_configured: false,
  steam_username: '',
  game_ready: false,
})

const refresh = async () => {
  loading.value = true
  try {
    const { data } = await api.getSetupStatus()
    setup.value = data
  } catch (error) {
    console.error('获取就绪状态失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(refresh)
</script>

<style scoped>
.header {
  margin-bottom: 2rem;
}

.header h1 {
  margin: 0;
  font-size: 2.5rem;
}

.header em {
  color: #C4612F;
  font-style: italic;
}

.header p {
  color: #5C635D;
  margin-top: 0.5rem;
}

.card {
  margin-bottom: 2rem;
}

.card h3 {
  margin-top: 0;
}

.status-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin: 1.5rem 0;
}

.status-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-top: 0.35rem;
  flex-shrink: 0;
}

.status-dot.ok {
  background: #4C9A5D;
}

.status-dot.bad {
  background: #C0392B;
}

.status-dot.wait {
  background: #E0A93F;
}

.status-text {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.status-text strong {
  color: #1F2421;
}

.status-text small {
  color: #5C635D;
  font-size: 0.85rem;
}

.intro {
  color: #5C635D;
  line-height: 1.7;
  margin-bottom: 1.5rem;
}

.intro a {
  color: #C4612F;
}

.steps {
  padding-left: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.steps li {
  line-height: 1.6;
}

.steps li strong {
  color: #1F2421;
  font-size: 1.05rem;
}

.steps li p {
  margin: 0.5rem 0;
  color: #5C635D;
}

.steps code {
  background: #F7F4EF;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.85rem;
  color: #C4612F;
}

.steps pre {
  background: #1F2421;
  color: #E7E1D7;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.85rem;
  overflow-x: auto;
  margin: 0.5rem 0;
}

.alert {
  background: #FFF3CD;
  border: 1px solid #FFE69C;
  border-radius: 8px;
  padding: 1rem;
  margin-top: 1.5rem;
}

.alert strong {
  display: block;
  margin-bottom: 0.5rem;
  color: #92400e;
}

.alert ul {
  margin: 0;
  padding-left: 1.5rem;
  color: #92400e;
  line-height: 1.7;
}

.btn {
  margin-top: 1rem;
}
</style>
