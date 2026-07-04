<template>
  <div class="setup-view">
    <div class="header">
      <h1>服务器<em>安装</em></h1>
    </div>

    <div class="card" v-if="!serverInstalled">
      <h3>🌾 配置星露谷服务端</h3>
      <p>选择一种方式提供游戏文件</p>

      <div class="setup-methods">
        <div class="method-card" :class="{ selected: setupMethod === 'upload' }" @click="setupMethod = 'upload'">
          <div class="method-icon">📦</div>
          <h4>上传游戏文件</h4>
          <p>上传已有的游戏压缩包</p>
          <span class="badge badge-success">最简单</span>
        </div>

        <div class="method-card" :class="{ selected: setupMethod === 'path' }" @click="setupMethod = 'path'">
          <div class="method-icon">📁</div>
          <h4>指定游戏路径</h4>
          <p>服务器上已有游戏文件</p>
          <span class="badge badge-success">最快速</span>
        </div>

        <div class="method-card" :class="{ selected: setupMethod === 'steamcmd' }" @click="setupMethod = 'steamcmd'">
          <div class="method-icon">⬇️</div>
          <h4>SteamCMD 下载</h4>
          <p>自动从 Steam 下载</p>
          <span class="badge badge-warning">需要时间</span>
        </div>
      </div>

      <!-- 上传方式 -->
      <div v-if="setupMethod === 'upload'" class="method-content">
        <div class="form-group">
          <label>上传游戏压缩包（.zip）</label>
          <div class="upload-area" @click="$refs.fileInput.click()">
            <div v-if="!uploadFile">
              <div class="upload-icon">📤</div>
              <p>点击选择文件或拖放到此处</p>
              <small>支持 .zip 格式，大约 500MB</small>
            </div>
            <div v-else>
              <div class="upload-icon">✅</div>
              <p>{{ uploadFile.name }}</p>
              <small>{{ formatFileSize(uploadFile.size) }}</small>
            </div>
          </div>
          <input ref="fileInput" type="file" accept=".zip" @change="handleFileSelect" style="display: none">
        </div>

        <div class="form-group">
          <label>
            <input type="checkbox" v-model="installSMAPI">
            安装 SMAPI（支持 MOD）
          </label>
        </div>

        <button class="btn btn-primary" @click="startUpload" :disabled="!uploadFile || installing">
          {{ installing ? '上传中...' : '开始上传' }}
        </button>
      </div>

      <!-- 指定路径方式 -->
      <div v-if="setupMethod === 'path'" class="method-content">
        <div class="form-group">
          <label>游戏文件路径</label>
          <input v-model="gamePath" placeholder="/home/user/.steam/steam/steamapps/common/Stardew Valley" />
          <small>填写服务器上已存在的游戏目录完整路径</small>
        </div>

        <div class="form-group">
          <label>
            <input type="checkbox" v-model="detectSMAPI">
            自动检测是否已安装 SMAPI
          </label>
        </div>

        <button class="btn btn-primary" @click="verifyPath" :disabled="!gamePath || installing">
          {{ installing ? '验证中...' : '验证路径' }}
        </button>
      </div>

      <!-- SteamCMD 方式 -->
      <div v-if="setupMethod === 'steamcmd'" class="method-content">
        <div class="form-group">
          <label>安装路径</label>
          <input v-model="installPath" placeholder="/opt/stardew-server" />
          <small>游戏将被下载到此目录</small>
        </div>

        <div class="form-group">
          <label>
            <input type="checkbox" v-model="installSMAPI">
            下载后自动安装 SMAPI
          </label>
        </div>

        <div class="alert">
          <strong>⚠️ 注意：</strong>
          <ul>
            <li>使用 SteamCMD 从官方服务器下载（免费）</li>
            <li>下载速度取决于网络环境，国内可能较慢</li>
            <li>建议使用代理或 VPN 以获得更好的下载速度</li>
          </ul>
        </div>

        <button class="btn btn-primary" @click="startSteamCMD" :disabled="!installPath || installing">
          {{ installing ? '下载中...' : '开始下载' }}
        </button>
      </div>

      <div class="install-log" v-if="installLog.length > 0">
        <div class="log-header">安装日志</div>
        <div class="log-content">
          <div v-for="(log, index) in installLog" :key="index">{{ log }}</div>
        </div>
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
      <h3>💡 配置说明</h3>
      <div class="help-section">
        <div class="help-item">
          <strong>📦 上传文件</strong>
          <p>从你的电脑上传游戏压缩包。适合已从 Steam 下载游戏的用户。</p>
          <p><em>游戏位置：</em> Steam 库 → Stardew Valley → 右键 → 管理 → 浏览本地文件</p>
        </div>
        <div class="help-item">
          <strong>📁 指定路径</strong>
          <p>服务器上已有游戏文件，只需指定路径。最快速的方式。</p>
          <p><em>常见路径：</em> <code>~/.steam/steam/steamapps/common/Stardew Valley</code></p>
        </div>
        <div class="help-item">
          <strong>⬇️ SteamCMD</strong>
          <p>使用 SteamCMD 从官方服务器下载（免费，无需购买游戏）。</p>
          <p><em>注意：</em> 国内下载速度可能较慢，建议使用代理。</p>
        </div>
      </div>

      <div class="help-item">
        <strong>🔌 关于 SMAPI</strong>
        <p>SMAPI 是星露谷物语的 MOD 加载器，安装后可以使用各种 MOD 增强游戏体验。如果不需要 MOD，可以取消勾选。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const serverInstalled = ref(false)
const setupMethod = ref('upload')
const gamePath = ref('')
const installPath = ref('/opt/stardew-server')
const installSMAPI = ref(true)
const detectSMAPI = ref(true)
const installing = ref(false)
const installLog = ref([])
const uploadFile = ref(null)

const serverInfo = ref({
  path: '/opt/stardew-server',
  type: 'SMAPI',
  version: '1.6.9'
})

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file && file.name.endsWith('.zip')) {
    uploadFile.value = file
  } else {
    alert('请选择 .zip 格式的文件')
  }
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const startUpload = async () => {
  installing.value = true
  installLog.value = []

  try {
    installLog.value.push('📤 上传游戏文件...')

    const formData = new FormData()
    formData.append('file', uploadFile.value)
    formData.append('installSMAPI', installSMAPI.value)

    // TODO: 实现上传 API
    // await axios.post('/api/v1/install/upload', formData, {
    //   onUploadProgress: (progressEvent) => {
    //     const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
    //     installLog.value.push(`上传进度: ${percentCompleted}%`)
    //   }
    // })

    // 模拟上传过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    installLog.value.push('📦 解压游戏文件...')

    await new Promise(resolve => setTimeout(resolve, 1500))
    installLog.value.push('🔧 配置服务器...')

    if (installSMAPI.value) {
      await new Promise(resolve => setTimeout(resolve, 1500))
      installLog.value.push('🔌 安装 SMAPI...')
    }

    await new Promise(resolve => setTimeout(resolve, 1000))
    installLog.value.push('✅ 配置完成！')

    serverInstalled.value = true
  } catch (error) {
    console.error('Upload failed:', error)
    installLog.value.push('❌ 上传失败：' + error.message)
  } finally {
    installing.value = false
  }
}

const verifyPath = async () => {
  installing.value = true
  installLog.value = []

  try {
    installLog.value.push('🔍 验证游戏路径...')

    // TODO: 调用后端 API
    // const { data } = await axios.post('/api/v1/install/verify', {
    //   path: gamePath.value,
    //   detectSMAPI: detectSMAPI.value
    // })

    // 模拟验证过程
    await new Promise(resolve => setTimeout(resolve, 1000))
    installLog.value.push('✅ 检测到 Stardew Valley 可执行文件')

    await new Promise(resolve => setTimeout(resolve, 500))
    if (detectSMAPI.value) {
      installLog.value.push('✅ 检测到 SMAPI 已安装')
    }

    await new Promise(resolve => setTimeout(resolve, 500))
    installLog.value.push('🔧 保存配置...')

    await new Promise(resolve => setTimeout(resolve, 500))
    installLog.value.push('✅ 配置完成！')

    serverInstalled.value = true
  } catch (error) {
    console.error('Verification failed:', error)
    installLog.value.push('❌ 验证失败：' + error.message)
  } finally {
    installing.value = false
  }
}

const startSteamCMD = async () => {
  installing.value = true
  installLog.value = []

  try {
    installLog.value.push('🚀 开始下载...')
    installLog.value.push('📥 检查 SteamCMD...')

    // TODO: 调用后端 API
    // const { data } = await axios.post('/api/v1/install/steamcmd', {
    //   path: installPath.value,
    //   installSMAPI: installSMAPI.value
    // })

    // 模拟下载过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    installLog.value.push('📥 使用 SteamCMD 下载游戏...')
    installLog.value.push('⏳ 这可能需要 5-15 分钟，请耐心等待...')

    await new Promise(resolve => setTimeout(resolve, 3000))
    installLog.value.push('📦 下载完成，验证文件完整性...')

    await new Promise(resolve => setTimeout(resolve, 1500))
    installLog.value.push('🔧 配置服务器...')

    if (installSMAPI.value) {
      await new Promise(resolve => setTimeout(resolve, 2000))
      installLog.value.push('🔌 下载并安装 SMAPI...')
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
  if (confirm('确定要重新配置吗？现有配置将被覆盖。')) {
    serverInstalled.value = false
    installLog.value = []
    uploadFile.value = null
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

.setup-methods {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
  margin: 1.5rem 0;
}

.method-card {
  padding: 1.5rem;
  border: 2px solid #E7E1D7;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.method-card:hover {
  border-color: #C4612F;
  transform: translateY(-2px);
}

.method-card.selected {
  border-color: #C4612F;
  background: #F2E3D6;
}

.method-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.method-card h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.15rem;
}

.method-card p {
  margin: 0 0 0.75rem 0;
  color: #5C635D;
  font-size: 0.9rem;
}

.method-content {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #E7E1D7;
}

.upload-area {
  border: 2px dashed #E7E1D7;
  border-radius: 12px;
  padding: 3rem 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #FBF9F5;
}

.upload-area:hover {
  border-color: #C4612F;
  background: #F2E3D6;
}

.upload-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.upload-area p {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  color: #1F2421;
}

.upload-area small {
  color: #5C635D;
  font-size: 0.85rem;
}

.alert {
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 8px;
  padding: 1rem;
  margin: 1rem 0;
}

.alert strong {
  display: block;
  margin-bottom: 0.5rem;
  color: #92400e;
}

.alert ul {
  margin: 0.5rem 0 0 0;
  padding-left: 1.5rem;
  color: #92400e;
  line-height: 1.6;
}

.alert ul li {
  margin-bottom: 0.25rem;
}

.install-log {
  margin-top: 1.5rem;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
  overflow: hidden;
}

.log-header {
  background: #5C635D;
  color: #F7F4EF;
  padding: 0.75rem 1rem;
  font-weight: 400;
}

.log-content {
  padding: 1rem;
  background: #1F2421;
  color: #E7E1D7;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.9rem;
  max-height: 300px;
  overflow-y: auto;
  line-height: 1.6;
}

.log-content div {
  margin-bottom: 0.25rem;
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

.form-group label input[type="checkbox"] {
  margin-right: 0.5rem;
  width: auto;
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

.btn {
  margin-top: 1rem;
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

.badge-warning {
  background: #fffbeb;
  color: #92400e;
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

.help-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.help-item {
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #E7E1D7;
}

.help-item:last-child {
  padding-bottom: 0;
  border-bottom: none;
}

.help-item strong {
  display: block;
  margin-bottom: 0.5rem;
  color: #1F2421;
  font-size: 1.05rem;
}

.help-item p {
  margin: 0 0 0.5rem 0;
  color: #5C635D;
  line-height: 1.6;
}

.help-item em {
  color: #5C635D;
  font-style: normal;
  font-size: 0.9rem;
}

.help-item code {
  background: #F7F4EF;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.85rem;
  color: #C4612F;
}
</style>
