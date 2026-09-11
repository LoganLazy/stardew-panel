<template>
  <div class="server-view">
    <div class="hero">
      <div class="eyebrow">
        <span class="badge" :class="statusBadgeClass">● {{ statusLabel }}</span>
      </div>
      <h1>服务器<em>控制台</em></h1>
      <p>管理星露谷物语服务器的启动、停止和配置</p>
      <div class="actions">
        <button class="btn btn-primary" @click="startServer" v-if="canStart" :disabled="isTransitioning">
          启动服务器
        </button>
        <button class="btn btn-secondary" @click="stopServer" v-if="serverStatus === 'running'" :disabled="isTransitioning">
          停止服务器
        </button>
        <button class="btn btn-secondary" @click="restartServer" v-if="serverStatus === 'running'">
          重启服务器
        </button>
        <a class="btn btn-secondary" :href="vncUrl" target="_blank" v-if="serverStatus === 'running'">
          打开游戏画面
        </a>
      </div>
      <p class="operation-error" v-if="operationError">{{ operationError }}</p>
    </div>

    <!-- 邀请码：玩家联机要用 -->
    <div class="card invite-card" v-if="serverStatus === 'running'">
      <h3>联机邀请码</h3>
      <p class="invite-hint">把邀请码发给好友，他们在游戏「加入联机 → 输入邀请码」即可加入。</p>
      <div class="invite-box">
        <code class="invite-code">{{ inviteCode || '获取中…' }}</code>
        <button class="btn btn-secondary btn-sm" @click="copyInvite" :disabled="!inviteCode">
          {{ copied ? '已复制' : '复制' }}
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

    <!-- 性能监控：容器资源占用 + 面板进程自身指标 -->
    <div class="card">
      <h3>性能监控</h3>
      <template v-if="metrics.container">
        <div class="meter-row">
          <div class="meter-head">
            <span>CPU 占用</span>
            <span class="meter-value">{{ metrics.container.cpu_percent.toFixed(1) }}%</span>
          </div>
          <div class="meter">
            <div class="meter-fill" :style="{ width: cpuBarWidth }"></div>
          </div>
        </div>
        <div class="meter-row">
          <div class="meter-head">
            <span>内存占用</span>
            <span class="meter-value">
              {{ formatBytes(metrics.container.memory_bytes) }} / {{ formatBytes(metrics.container.memory_limit_bytes) }}
              （{{ metrics.container.memory_percent.toFixed(1) }}%）
            </span>
          </div>
          <div class="meter">
            <div class="meter-fill" :style="{ width: memBarWidth }"></div>
          </div>
        </div>
        <div class="net-row">
          <span>网络 ↓ {{ formatBytes(metrics.container.net_read_bytes) }}</span>
          <span>↑ {{ formatBytes(metrics.container.net_write_bytes) }}</span>
        </div>
      </template>
      <div class="metrics-empty" v-else>服务器未运行，暂无容器监控数据</div>
      <div class="panel-metrics">
        面板进程：内存 {{ formatBytes(metrics.panel.alloc_bytes) }} · 协程 {{ metrics.panel.goroutines }} 个 · 已运行 {{ formatUptime(metrics.panel.uptime_seconds) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import * as api from '@/api/api'

const serverStatus = ref('stopped')
const stats = ref({
  onlinePlayers: 0,
  version: '未知',
  uptime: '0分钟',
  modsCount: 0
})

const serverInfo = ref({
  host: window.location.hostname,
  port: '24642',
  saveName: '未加载',
  maxPlayers: 8
})

const inviteCode = ref('')
const copied = ref(false)
const operationError = ref('')
const vncPort = ref('5800')

const metrics = ref({
  container: null,
  panel: { alloc_bytes: 0, sys_bytes: 0, goroutines: 0, uptime_seconds: 0 }
})

const vncUrl = computed(() => `http://${window.location.hostname}:${vncPort.value}`)

// CPU/内存进度条宽度（百分比夹在 0-100，CPU 可短暂超过 100%）
const cpuBarWidth = computed(() => `${Math.min(100, metrics.value.container?.cpu_percent || 0)}%`)
const memBarWidth = computed(() => `${Math.min(100, metrics.value.container?.memory_percent || 0)}%`)
const isTransitioning = computed(() => ['starting', 'stopping', 'restarting'].includes(serverStatus.value))
const canStart = computed(() => ['stopped', 'error'].includes(serverStatus.value))
const statusLabel = computed(() => ({
  running: '运行中',
  starting: '启动中',
  stopping: '停止中',
  restarting: '重启中',
  error: '操作失败',
  stopped: '已停止'
}[serverStatus.value] || '未知状态'))
const statusBadgeClass = computed(() => ({
  running: 'badge-success',
  stopped: 'badge-error',
  error: 'badge-error'
}[serverStatus.value] || 'badge-warning'))

let refreshInterval = null

const formatUptime = (seconds) => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分钟`
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}小时${minutes}分钟`
}

const formatBytes = (bytes) => {
  if (!Number.isFinite(bytes) || bytes < 0) return '—'
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)))
  return `${(bytes / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

const fetchMetrics = async () => {
  try {
    const { data } = await api.getMetrics()
    metrics.value = data
  } catch (error) {
    console.error('获取性能指标失败:', error)
  }
}

const fetchStatus = async () => {
  try {
    const { data } = await api.getServerStatus()
    serverStatus.value = data.status
    stats.value.onlinePlayers = data.online_players || 0
    stats.value.version = data.version || '未知'
    stats.value.uptime = formatUptime(data.uptime || 0)
    serverInfo.value.host = window.location.hostname
    serverInfo.value.port = data.game_port || serverInfo.value.port
    serverInfo.value.saveName = data.farm_name || '未加载'
    serverInfo.value.maxPlayers = data.max_players || serverInfo.value.maxPlayers
    vncPort.value = data.vnc_port || vncPort.value
    operationError.value = data.error || ''

    // 运行中才拉邀请码；停止时清空
    if (data.status === 'running') {
      fetchInviteCode()
      fetchMetrics()
    } else {
      inviteCode.value = ''
      metrics.value.container = null
    }
  } catch (error) {
    console.error('获取服务器状态失败:', error)
  }
}

const fetchInviteCode = async () => {
  try {
    const { data } = await api.getInviteCode()
    inviteCode.value = data.invite_code || ''
  } catch (error) {
    // 游戏刚启动、API 未就绪时会失败，静默忽略
    inviteCode.value = ''
  }
}

const openVNC = () => {
  window.open(vncUrl, '_blank')
}

const copyInvite = async () => {
  if (!inviteCode.value) return
  try {
    await navigator.clipboard.writeText(inviteCode.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

const startServer = async () => {
  try {
    await api.startServer()
    serverStatus.value = 'starting'
    operationError.value = ''
    await fetchStatus()
  } catch (error) {
    alert('启动服务器失败: ' + (error.response?.data?.error || error.message))
  }
}

const stopServer = async () => {
  if (!confirm('确定要停止服务器吗？在线玩家将被断开连接。')) return

  try {
    const { data } = await api.stopServer()
    serverStatus.value = 'stopping'
    if (data.warning) alert('服务器已停止，但自动备份失败: ' + data.warning)
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
  fetchMetrics() // 停止状态下也有面板自身指标

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
.invite-card {
  text-align: center;
}

.invite-hint {
  color: #5C635D;
  margin-bottom: 1rem;
}

.invite-box {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
}

.invite-code {
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 1.25rem;
  letter-spacing: 0.05em;
  background: #F7F4EF;
  color: #C4612F;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  user-select: all;
}

.btn-sm {
  padding: 0.4rem 1rem;
  font-size: 0.9rem;
  margin-top: 0;
}

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

.hero .operation-error {
  color: #C0392B;
  background: #FFF3F3;
  border: 1px solid #FFCDD2;
  padding: 0.75rem 1rem;
  margin-top: 1rem;
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

.meter-row {
  margin-bottom: 1.25rem;
}

.meter-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.meter-value {
  color: #C4612F;
  font-variant-numeric: tabular-nums;
}

.meter {
  height: 10px;
  border-radius: 5px;
  background: rgba(196, 97, 47, 0.15);
  overflow: hidden;
}

.meter-fill {
  height: 100%;
  border-radius: 5px;
  background: #C4612F;
  transition: width 0.5s ease;
}

.net-row {
  display: flex;
  gap: 2rem;
  color: #5C635D;
  font-variant-numeric: tabular-nums;
  margin-bottom: 1.25rem;
}

.metrics-empty {
  color: #5C635D;
  padding: 1rem 0;
}

.panel-metrics {
  border-top: 1px solid rgba(92, 99, 93, 0.2);
  padding-top: 1rem;
  color: #5C635D;
  font-size: 0.9rem;
}
</style>
