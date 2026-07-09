<template>
  <div class="saves-view">
    <div class="header">
      <h1>存档<em>管理</em></h1>
      <button class="btn btn-primary" @click="createBackup">
        创建备份
      </button>
    </div>

    <div class="saves-list">
      <div class="card save-item" v-for="backup in backups" :key="backup.id">
        <div class="save-info">
          <h3>{{ backup.save_name }}</h3>
          <p>{{ formatDate(backup.created_at) }}</p>
          <span class="badge badge-success">{{ formatSize(backup.size) }}</span>
        </div>
        <div class="save-actions">
          <button class="btn btn-secondary" @click="restoreSave(backup.id)">
            恢复
          </button>
          <button class="btn-icon" @click="deleteSave(backup.id)">删除</button>
        </div>
      </div>

      <div class="empty" v-if="backups.length === 0">
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
import * as api from '@/api/api'

const saves = ref([])
const backups = ref([])
const selectedSave = ref('')

const loadSaves = async () => {
  try {
    const { data } = await api.listSaves()
    saves.value = data.saves || []
    backups.value = data.backups || []
  } catch (error) {
    console.error('Failed to load saves:', error)
  }
}

const createBackup = async () => {
  if (saves.value.length === 0) {
    alert('没有可备份的存档')
    return
  }

  // 如果只有一个存档，直接备份
  const saveName = saves.value.length === 1 ? saves.value[0] : selectedSave.value || saves.value[0]

  try {
    await api.createBackup(saveName)
    alert('备份创建成功！')
    await loadSaves()
  } catch (error) {
    console.error('Failed to create backup:', error)
    alert('创建备份失败: ' + (error.response?.data?.error || error.message))
  }
}

const restoreSave = async (id) => {
  if (!confirm('确定要恢复这个存档吗？当前存档将被覆盖。')) return
  try {
    await api.restoreBackup(id)
    alert('存档恢复成功！')
    await loadSaves()
  } catch (error) {
    console.error('Failed to restore save:', error)
    alert('恢复失败: ' + (error.response?.data?.error || error.message))
  }
}

const deleteSave = async (id) => {
  if (!confirm('确定要删除这个备份吗？')) return
  try {
    await api.deleteBackup(id)
    await loadSaves()
  } catch (error) {
    console.error('Failed to delete save:', error)
    alert('删除失败: ' + (error.response?.data?.error || error.message))
  }
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
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
