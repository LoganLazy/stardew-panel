<template>
  <div class="saves-view">
    <div class="header">
      <h1>存档<em>管理</em></h1>
      <button class="btn btn-primary" @click="createBackup">
        创建备份
      </button>
    </div>

    <div class="saves-list">
      <div class="card save-item" v-for="save in saves" :key="save.id">
        <div class="save-info">
          <h3>{{ save.name }}</h3>
          <p>{{ save.date }}</p>
          <span class="badge badge-success">{{ save.size }}</span>
        </div>
        <div class="save-actions">
          <button class="btn btn-secondary" @click="restoreSave(save.id)">
            恢复
          </button>
          <button class="btn-icon" @click="deleteSave(save.id)">删除</button>
        </div>
      </div>

      <div class="empty" v-if="saves.length === 0">
        <p>暂无存档备份</p>
        <button class="btn btn-primary" @click="createBackup">
          创建第一个备份
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const saves = ref([])

const loadSaves = async () => {
  try {
    const { data } = await axios.get('/api/v1/saves')
    saves.value = data.saves || []
  } catch (error) {
    console.error('Failed to load saves:', error)
  }
}

const createBackup = async () => {
  try {
    await axios.post('/api/v1/saves/backup')
    await loadSaves()
  } catch (error) {
    console.error('Failed to create backup:', error)
  }
}

const restoreSave = async (id) => {
  if (!confirm('确定要恢复这个存档吗？当前存档将被覆盖。')) return
  try {
    await axios.post(`/api/v1/saves/restore/${id}`)
    await loadSaves()
  } catch (error) {
    console.error('Failed to restore save:', error)
  }
}

const deleteSave = async (id) => {
  if (!confirm('确定要删除这个备份吗？')) return
  try {
    await axios.delete(`/api/v1/saves/${id}`)
    await loadSaves()
  } catch (error) {
    console.error('Failed to delete save:', error)
  }
}

onMounted(loadSaves)
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

.saves-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.save-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.save-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
}

.save-info p {
  margin: 0 0 0.5rem 0;
  color: #5C635D;
  font-size: 0.9rem;
}

.save-actions {
  display: flex;
  gap: 1rem;
}

.btn-icon {
  background: transparent;
  color: #5C635D;
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.btn-icon:hover {
  color: #C4612F;
}

.empty {
  text-align: center;
  padding: 4rem 2rem;
  color: #5C635D;
}

.empty p {
  margin-bottom: 1.5rem;
}
</style>
