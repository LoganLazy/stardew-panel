<template>
  <div class="players-view">
    <div class="header">
      <h1>在线<em>玩家</em></h1>
      <button class="btn btn-secondary" @click="refreshPlayers">
        刷新
      </button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="card stat-card">
        <div class="stat-label">在线玩家</div>
        <div class="stat-value">{{ stats.online_count || 0 }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-label">历史玩家</div>
        <div class="stat-value">{{ stats.total_count || 0 }}</div>
      </div>
    </div>

    <!-- 在线玩家列表 -->
    <div class="card">
      <h3>当前在线玩家</h3>

      <div class="players-list" v-if="players.length > 0">
        <div class="player-item" v-for="player in players" :key="player.id">
          <div class="player-info">
            <div class="player-name">
              <span class="badge badge-success">● 在线</span>
              {{ player.player_name }}
            </div>
            <div class="player-meta">
              <span>加入时间: {{ formatTime(player.connected_at) }}</span>
              <span>在线时长: {{ calculateDuration(player.connected_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="empty" v-else>
        <p>😴 当前没有玩家在线</p>
        <small>玩家加入游戏后会自动显示在这里</small>
      </div>
    </div>

    <!-- 管理操作 -->
    <div class="card">
      <h3>管理操作</h3>
      <div class="action-buttons">
        <button class="btn btn-secondary" @click="cleanupOldRecords">
          清理历史记录（保留30天）
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as api from '@/api/api'

const players = ref([])
const stats = ref({
  online_count: 0,
  total_count: 0
})

let refreshInterval = null

const loadPlayers = async () => {
  try {
    const { data } = await api.listPlayers()
    players.value = data.players || []
    stats.value = data.stats || { online_count: 0, total_count: 0 }
  } catch (error) {
    console.error('Failed to load players:', error)
  }
}

const refreshPlayers = async () => {
  try {
    await api.refreshPlayers()
    await loadPlayers()
  } catch (error) {
    console.error('Failed to refresh players:', error)
    alert('刷新失败: ' + (error.response?.data?.error || error.message))
  }
}

const cleanupOldRecords = async () => {
  if (!confirm('确定要清理30天前的历史记录吗？')) return

  try {
    await api.cleanupPlayers(30)
    alert('清理完成')
    await loadPlayers()
  } catch (error) {
    console.error('Failed to cleanup:', error)
    alert('清理失败: ' + (error.response?.data?.error || error.message))
  }
}

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const calculateDuration = (connectedAt) => {
  const start = new Date(connectedAt)
  const now = new Date()
  const diff = now - start

  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))

  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  }
  return `${minutes}分钟`
}

onMounted(() => {
  loadPlayers()
  // 每10秒自动刷新
  refreshInterval = setInterval(loadPlayers, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  text-align: center;
  padding: 2rem;
}

.stat-label {
  font-size: 0.9rem;
  color: #5C635D;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 500;
  color: #C4612F;
}

.card {
  margin-bottom: 2rem;
}

.card h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.player-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
  background: #FBF9F5;
  transition: all 0.2s;
}

.player-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(31, 36, 33, 0.08);
}

.player-info {
  flex: 1;
}

.player-name {
  font-size: 1.1rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.player-meta {
  display: flex;
  gap: 1.5rem;
  font-size: 0.9rem;
  color: #5C635D;
}

.empty {
  text-align: center;
  padding: 3rem 2rem;
}

.empty p {
  font-size: 1.1rem;
  color: #5C635D;
  margin-bottom: 0.5rem;
}

.empty small {
  color: #5C635D;
  font-size: 0.9rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
}
</style>
