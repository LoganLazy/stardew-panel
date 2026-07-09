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
import { ref, onMounted, onUnmounted } from 'vue'
import * as api from '@/api/api'

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

let refreshInterval = null

const formatUptime = (seconds) => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分钟`
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}小时${minutes}分钟`
}

const fetchStatus = async () => {
  try {
    const { data } = await api.getServerStatus()
    serverStatus.value = data.status
    stats.value.onlinePlayers = data.online_players || 0
    stats.value.version = data.version || '1.6.9'
    stats.value.uptime = formatUptime(data.uptime || 0)
  } catch (error) {
    console.error('获取服务器状态失败:', error)
  }
}

const startServer = async () => {
  try {
    await api.startServer()
    serverStatus.value = 'starting'
    await fetchStatus()
  } catch (error) {
    alert('启动服务器失败: ' + (error.response?.data?.error || error.message))
  }
}

const stopServer = async () => {
  if (!confirm('确定要停止服务器吗？在线玩家将被断开连接。')) return

  try {
    await api.stopServer()
    serverStatus.value = 'stopping'
    await fetchStatus()
  } catch (error) {
    alert('停止服务器失败: ' + (error.response?.data?.error || error.message))
  }
}

const restartServer = async () => {
  if (!confirm('确定要重启服务器吗？在线玩家将被断开连接。')) return

  try {
    await api.restartServer()
    serverStatus.value = 'restarting'
    await fetchStatus()
  } catch (error) {
    alert('重启服务器失败: ' + (error.response?.data?.error || error.message))
  }
}

onMounted(async () => {
  await fetchStatus()

  // 每5秒刷新一次状态
  refreshInterval = setInterval(fetchStatus, 5000)

  // 获取MOD数量
  try {
    const { data } = await api.listMods()
    stats.value.modsCount = data.mods?.length || 0
  } catch (error) {
    console.error('获取MOD列表失败:', error)
  }
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
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
