<template>
  <div class="setup-view">
    <div class="header">
      <h1>服务器<em>安装</em></h1>
    </div>

    <div class="card" v-if="!serverInstalled">
      <h3>🌾 安装星露谷服务端</h3>
      <p>首次使用需要安装星露谷物语服务端</p>

      <div class="install-options">
        <div class="option-card" :class="{ selected: installMode === 'official' }" @click="installMode = 'official'">
          <h4>官方服务端</h4>
          <p>纯净服务器，不支持 MOD</p>
          <span class="badge badge-success">推荐新手</span>
        </div>
        <div class="option-card" :class="{ selected: installMode === 'smapi' }" @click="installMode = 'smapi'">
          <h4>SMAPI 服务端</h4>
          <p>支持安装 MOD，功能更强大</p>
          <span class="badge badge-success">推荐进阶</span>
        </div>
      </div>

      <div class="form-group">
        <label>安装路径</label>
        <input v-model="installPath" placeholder="/opt/stardew-server" />
        <small>建议使用默认路径</small>
      </div>

      <button class="btn btn-primary" @click="startInstall" :disabled="installing">
        {{ installing ? '安装中...' : '开始安装' }}
      </button>

      <div class="install-log" v-if="installLog.length > 0">
        <div v-for="(log, index) in installLog" :key="index">{{ log }}</div>
      </div>
    </div>

    <div class="card" v-else>
      <h3>✅ 服务端已安装</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">安装路径</span>
          <span class="info-value">{{ serverInfo.path }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">服务端类型</span>
          <span class="info-value">{{ serverInfo.type }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">版本</span>
          <span class="info-value">{{ serverInfo.version }}</span>
        </div>
      </div>
      <button class="btn btn-secondary" @click="reinstall">重新安装</button>
    </div>

    <div class="card">
      <h3>💡 安装说明</h3>
      <ul>
        <li><strong>官方服务端：</strong>使用 SteamCMD 下载，无需购买游戏</li>
        <li><strong>SMAPI 服务端：</strong>在官方基础上自动安装 SMAPI 框架</li>
        <li><strong>安装时间：</strong>根据网络情况，大约需要 5-15 分钟</li>
        <li><strong>磁盘空间：</strong>需要至少 2GB 可用空间</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const serverInstalled = ref(false)
const installMode = ref('smapi')
const installPath = ref('/opt/stardew-server')
const installing = ref(false)
const installLog = ref([])

const serverInfo = ref({
  path: '/opt/stardew-server',
  type: 'SMAPI',
  version: '1.6.9'
})

const startInstall = async () => {
  installing.value = true
  installLog.value = []

  try {
    installLog.value.push('🚀 开始安装...')
    installLog.value.push('📥 下载 SteamCMD...')

    // TODO: 调用后端 API
    // const { data } = await axios.post('/api/v1/server/install', {
    //   mode: installMode.value,
    //   path: installPath.value
    // })

    // 模拟安装过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    installLog.value.push('📥 下载星露谷服务端...')

    await new Promise(resolve => setTimeout(resolve, 2000))
    installLog.value.push('🔧 配置服务器...')

    if (installMode.value === 'smapi') {
      await new Promise(resolve => setTimeout(resolve, 1500))
      installLog.value.push('🔌 安装 SMAPI...')
    }

    await new Promise(resolve => setTimeout(resolve, 1000))
    installLog.value.push('✅ 安装完成！')

    serverInstalled.value = true
  } catch (error) {
    console.error('Installation failed:', error)
    installLog.value.push('❌ 安装失败：' + error.message)
  } finally {
    installing.value = false
  }
}

const reinstall = () => {
  if (confirm('确定要重新安装吗？现有配置将被保留，但服务端文件会被覆盖。')) {
    serverInstalled.value = false
    installLog.value = []
  }
}
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

.card {
  margin-bottom: 2rem;
}

.card h3 {
  margin-top: 0;
}

.install-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
  margin: 1.5rem 0;
}

.option-card {
  padding: 1.5rem;
  border: 2px solid #E7E1D7;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.option-card:hover {
  border-color: #C4612F;
  transform: translateY(-2px);
}

.option-card.selected {
  border-color: #C4612F;
  background: #F2E3D6;
}

.option-card h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
}

.option-card p {
  margin: 0 0 0.75rem 0;
  color: #5C635D;
  font-size: 0.95rem;
}

.form-group {
  margin: 1.5rem 0;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 400;
  color: #1F2421;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.95rem;
}

.form-group small {
  display: block;
  margin-top: 0.5rem;
  color: #5C635D;
  font-size: 0.85rem;
}

.install-log {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #1F2421;
  color: #E7E1D7;
  border-radius: 8px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.9rem;
  max-height: 300px;
  overflow-y: auto;
}

.install-log div {
  margin-bottom: 0.25rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
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

.card ul {
  margin: 0;
  padding-left: 1.5rem;
  color: #5C635D;
  line-height: 1.8;
}

.card ul li {
  margin-bottom: 0.5rem;
}
</style>
