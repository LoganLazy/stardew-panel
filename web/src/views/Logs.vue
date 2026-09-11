<template>
  <div class="logs-view">
    <h1>服务器<em>日志</em></h1>

    <div class="card log-container">
      <div class="log-header">
        <button class="btn" :class="live ? 'btn-primary' : 'btn-secondary'" @click="toggleLive">
          {{ live ? '⏸ 暂停实时' : '▶ 开启实时' }}
        </button>
        <span class="conn-state" :class="`conn-${conn}`">{{ connLabel }}</span>
        <button class="btn btn-secondary" @click="loadLogs" :disabled="live">刷新</button>
        <button class="btn btn-secondary" @click="clearLogs">清空</button>
        <label class="autoscroll">
          <input type="checkbox" v-model="autoScroll" />
          自动滚动
        </label>
      </div>
      <div class="log-content" ref="logContent">
        <div class="log-line" v-for="(log, index) in logs" :key="index">
          {{ log }}
        </div>
        <div class="empty" v-if="logs.length === 0">
          <p>{{ live ? '等待日志推送…' : '暂无日志' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import * as api from '@/api/api'

// 页面最多保留的日志行数，防止长会话内存膨胀
const MAX_LOG_LINES = 1000

const logs = ref([])
const live = ref(false)
const autoScroll = ref(true)
const conn = ref('idle') // idle | connecting | connected | reconnecting
const logContent = ref(null)

let abortController = null
let reconnectTimer = null
let retryDelay = 1000
// 流代数：停止实时后旧流的收尾逻辑不能再触发重连
let streamGeneration = 0

const connLabel = ({
  idle: '未连接',
  connecting: '连接中…',
  connected: '实时推送中',
  reconnecting: '连接断开，自动重连中…'
}[conn] || '')

const appendLine = (line) => {
  logs.value.push(line)
  if (logs.value.length > MAX_LOG_LINES) {
    logs.value.splice(0, logs.value.length - MAX_LOG_LINES)
  }
  if (autoScroll.value) {
    nextTick(() => {
      const el = logContent.value
      if (el) el.scrollTop = el.scrollHeight
    })
  }
}

const startStream = async () => {
  const generation = streamGeneration
  abortController = new AbortController()
  conn.value = 'connecting'

  try {
    await api.streamLogs(200, {
      onLine: appendLine,
      onError: (message) => console.error('日志流错误:', message),
      onReady: () => {
        if (generation === streamGeneration) {
          conn.value = 'connected'
          retryDelay = 1000 // 恢复正常后重置退避
        }
      },
      signal: abortController.signal,
    })
  } catch (error) {
    if (error.name === 'AbortError') return // 主动停止，不需要重连
    console.error('日志流连接失败:', error)
  }

  // 流结束（服务端关闭或异常断开）：仍在实时模式则退避重连
  if (generation === streamGeneration && live.value) {
    conn.value = 'reconnecting'
    reconnectTimer = setTimeout(() => {
      if (live.value) startStream()
    }, retryDelay)
    retryDelay = Math.min(retryDelay * 2, 10000)
  }
}

const stopStream = () => {
  streamGeneration++
  live.value = false
  conn.value = 'idle'
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (abortController) {
    abortController.abort()
    abortController = null
  }
}

const toggleLive = () => {
  if (live.value) {
    stopStream()
  } else {
    live.value = true
    startStream()
  }
}

const loadLogs = async () => {
  try {
    const { data } = await api.getLogs(200)
    logs.value = data.logs || []
    if (autoScroll.value) {
      nextTick(() => {
        const el = logContent.value
        if (el) el.scrollTop = el.scrollHeight
      })
    }
  } catch (error) {
    console.error('Failed to load logs:', error)
  }
}

const clearLogs = () => {
  if (confirm('确定要清空日志显示吗？（不会删除服务器日志文件）')) {
    logs.value = []
  }
}

onMounted(() => {
  // 进入页面即开启实时推送（对齐旧版自动刷新行为），无需手动点按钮
  live.value = true
  startStream()
})

onUnmounted(stopStream)
</script>

<style scoped>
h1 {
  margin-bottom: 2rem;
  font-size: 2.5rem;
}

h1 em {
  color: #C4612F;
  font-style: italic;
}

.log-container {
  background: #1F2421;
  color: #F7F4EF;
  padding: 0;
  overflow: hidden;
}

.log-header {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: rgba(247, 244, 239, 0.05);
  border-bottom: 1px solid rgba(231, 225, 215, 0.1);
  align-items: center;
  flex-wrap: wrap;
}

.conn-state {
  font-size: 0.9rem;
  color: #E7E1D7;
}

.conn-connected {
  color: #7FB069;
}

.conn-connecting,
.conn-reconnecting {
  color: #E0A458;
}

.autoscroll {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  margin-left: auto;
  font-size: 0.9rem;
  color: #E7E1D7;
  cursor: pointer;
  user-select: none;
}

.log-content {
  padding: 1rem;
  height: 600px;
  overflow-y: auto;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
}

.log-line {
  margin-bottom: 0.25rem;
  color: #E7E1D7;
  white-space: pre-wrap;
  word-break: break-all;
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #5C635D;
}
</style>
