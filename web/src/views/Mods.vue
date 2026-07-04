<template>
  <div class="mods-view">
    <div class="header">
      <h1>MOD <em>管理</em></h1>
      <button class="btn btn-primary" @click="showUpload = true">
        上传 MOD
      </button>
    </div>

    <div class="mods-list">
      <div class="card mod-item" v-for="mod in mods" :key="mod.id">
        <div class="mod-info">
          <h3>{{ mod.name }}</h3>
          <p>{{ mod.description }}</p>
          <div class="mod-meta">
            <span class="badge badge-success">v{{ mod.version }}</span>
            <span class="mod-author">作者: {{ mod.author }}</span>
          </div>
        </div>
        <div class="mod-actions">
          <label class="toggle">
            <input type="checkbox" :checked="mod.enabled" @change="toggleMod(mod.id)">
            <span class="toggle-slider"></span>
          </label>
          <button class="btn-icon" @click="deleteMod(mod.id)">删除</button>
        </div>
      </div>

      <div class="empty" v-if="mods.length === 0">
        <p>暂无已安装的 MOD</p>
        <button class="btn btn-primary" @click="showUpload = true">
          上传第一个 MOD
        </button>
      </div>
    </div>

    <!-- 上传弹窗（简化版） -->
    <div class="modal" v-if="showUpload" @click.self="showUpload = false">
      <div class="modal-content card">
        <h2>上传 MOD</h2>
        <p>将 MOD 文件拖放到此处或点击选择文件</p>
        <input type="file" accept=".zip" @change="uploadMod">
        <button class="btn btn-secondary" @click="showUpload = false">取消</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const mods = ref([])
const showUpload = ref(false)

const loadMods = async () => {
  try {
    const { data } = await axios.get('/api/v1/mods')
    mods.value = data.mods || []
  } catch (error) {
    console.error('Failed to load mods:', error)
  }
}

const toggleMod = async (id) => {
  try {
    await axios.put(`/api/v1/mods/${id}/toggle`)
    await loadMods()
  } catch (error) {
    console.error('Failed to toggle mod:', error)
  }
}

const deleteMod = async (id) => {
  if (!confirm('确定要删除这个 MOD 吗？')) return
  try {
    await axios.delete(`/api/v1/mods/${id}`)
    await loadMods()
  } catch (error) {
    console.error('Failed to delete mod:', error)
  }
}

const uploadMod = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    await axios.post('/api/v1/mods/upload', formData)
    showUpload.value = false
    await loadMods()
  } catch (error) {
    console.error('Failed to upload mod:', error)
  }
}

onMounted(loadMods)
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

.mods-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mod-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mod-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
}

.mod-info p {
  margin: 0 0 0.75rem 0;
  color: #5C635D;
  font-size: 0.95rem;
}

.mod-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.mod-author {
  font-size: 0.9rem;
  color: #5C635D;
}

.mod-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.toggle {
  position: relative;
  width: 48px;
  height: 26px;
  display: inline-block;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #E7E1D7;
  border-radius: 999px;
  transition: 0.3s;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background: white;
  border-radius: 50%;
  transition: 0.3s;
}

.toggle input:checked + .toggle-slider {
  background: #C4612F;
}

.toggle input:checked + .toggle-slider:before {
  transform: translateX(22px);
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

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(31, 36, 33, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  max-width: 480px;
  width: 90%;
  text-align: center;
}

.modal-content h2 {
  margin-top: 0;
}

.modal-content input[type="file"] {
  margin: 1.5rem 0;
  width: 100%;
}
</style>
