<template>
  <div class="server-view">
    <div class="hero">
      <div class="eyebrow">
        <span class="badge badge-success" v-if="serverStatus === 'running'">● 运行中</span>
        <span class="badge badge-error" v-else>● 已停止</span>
      </div>
      <h1>服务器<em>控制台</em></h1>
      <p>管理星露谷物语服务器的启动、停止和配置</p>
      <div class="actions">
        <button class="btn btn-primary" @click="startServer" v-if="serverStatus !== 'running'">
          启动服务器
        </button>
        <button class="btn btn-secondary" @click="stopServer" v-else>
          停止服务器
        </button>
        <button class="btn btn-secondary" @click="restartServer" v-if="serverStatus === 'running'">
          重启服务器
        </button>
      </div>
    </div>

    <div class="stats">
      <div class="card stat-card">
        <div class="stat-label">在线玩家</div>
        <div class="stat-value">{{ stats.onlinePlayers }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-label">服务器版本</div>
        <div class="stat-value">{{ stats.version }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-label">运行时长</div>
        <div class="stat-value">{{ stats.uptime }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-label">已安装 MOD</div>
        <div class="stat-value">{{ stats.modsCount }}</div>
      </div>
    </div>

    <div class="card">
      <h3>服务器信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">服务器地址</span>
          <span class="info-value">{{ serverInfo.host }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">端口</span>
          <span class="info-value">{{ serverInfo.port }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">存档名称</span>
          <span class="info-value">{{ serverInfo.saveName }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">最大玩家数</span>
          <span class="info-value">{{ serverInfo.maxPlayers }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useServerStore } from '@/stores/server'

const serverStore = useServerStore()

const serverStatus = ref('stopped')
const stats = ref({
  onlinePlayers: 0,
  version: '1.6.9',
  uptime: '0分钟',
  modsCount: 0
})

const serverInfo = ref({
  host: 'localhost',
  port: '24642',
  saveName: '未加载',
  maxPlayers: 8
})

const startServer = async () => {
  await serverStore.startServer()
  serverStatus.value = 'running'
}

const stopServer = async () => {
  await serverStore.stopServer()
  serverStatus.value = 'stopped'
}

const restartServer = async () => {
  await serverStore.restartServer()
}

onMounted(async () => {
  const status = await serverStore.getStatus()
  serverStatus.value = status.status
})
</script>

<style scoped>
.hero {
  text-align: center;
  padding: 3rem 0;
  max-width: 640px;
  margin: 0 auto;
}

.eyebrow {
  margin-bottom: 1rem;
}

.hero h1 {
  font-size: 3rem;
  margin: 0 0 1rem 0;
  line-height: 1.1;
}

.hero em {
  color: #C4612F;
  font-style: italic;
}

.hero p {
  font-size: 1.1rem;
  color: #5C635D;
  margin-bottom: 2rem;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin: 3rem 0;
}

.stat-card {
  text-align: center;
}

.stat-label {
  font-size: 0.9rem;
  color: #5C635D;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 500;
  color: #C4612F;
}

.card h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-label {
  font-size: 0.9rem;
  color: #5C635D;
}

.info-value {
  font-size: 1.1rem;
  font-weight: 400;
}
</style>
