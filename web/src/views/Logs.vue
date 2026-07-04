<template>
  <div class="logs-view">
    <h1>服务器<em>日志</em></h1>

    <div class="card log-container">
      <div class="log-header">
        <button class="btn btn-secondary" @click="loadLogs">刷新</button>
        <button class="btn btn-secondary" @click="clearLogs">清空</button>
      </div>
      <div class="log-content">
        <div class="log-line" v-for="(log, index) in logs" :key="index">
          {{ log }}
        </div>
        <div class="empty" v-if="logs.length === 0">
          <p>暂无日志</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const logs = ref([])

const loadLogs = async () => {
  try {
    const { data } = await axios.get('/api/v1/logs')
    logs.value = data.logs || []
  } catch (error) {
    console.error('Failed to load logs:', error)
  }
}

const clearLogs = () => {
  logs.value = []
}

onMounted(loadLogs)
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
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #5C635D;
}
</style>
