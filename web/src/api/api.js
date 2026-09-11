import api, { API_BASE_URL } from './index'

// 认证相关
export const login = (username, password) =>
  api.post('/auth/login', { username, password })

export const logout = () => api.post('/auth/logout')

export const changePassword = (oldPassword, newPassword) =>
  api.post('/auth/change-password', {
    old_password: oldPassword,
    new_password: newPassword
  })

// 检查安装状态（容器模式：Steam 已配置 + docker 可用即视为已安装）
// 获取安装向导就绪状态（Steam 账号/docker/游戏容器）
export const getSetupStatus = () => api.get('/install/setup')

// 服务器管理
export const getServerStatus = () => api.get('/server/status')
export const startServer = () => api.post('/server/start')
export const stopServer = () => api.post('/server/stop')
export const restartServer = () => api.post('/server/restart')
export const getInviteCode = () => api.get('/server/invite-code')
export const getServerSettings = () => api.get('/server/settings')
export const updateServerSetting = (key, value) =>
  api.put(`/server/settings/${encodeURIComponent(key)}`, { value })

// MOD 管理
export const listMods = () => api.get('/mods')
export const uploadMod = (formData) =>
  api.post('/mods/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
export const toggleMod = (id, enabled) =>
  api.put(`/mods/${id}/toggle`, { enabled })
export const deleteMod = (id) => api.delete(`/mods/${id}`)

// 玩家管理
export const listPlayers = () => api.get('/players')
export const refreshPlayers = () => api.post('/players/refresh')
export const cleanupPlayers = (days = 30) =>
  api.delete(`/players/cleanup?days=${days}`)

// 存档管理
export const listSaves = () => api.get('/saves')
export const createBackup = (saveName) =>
  api.post('/saves/backup', { save_name: saveName })
export const restoreBackup = (id) => api.post(`/saves/restore/${id}`)
export const deleteBackup = (id) => api.delete(`/saves/${id}`)

// 日志
export const getLogs = (lines = 100) => api.get(`/logs?lines=${lines}`)

// 性能监控（容器资源占用 + 面板自身指标）
export const getMetrics = () => api.get('/server/metrics')

// 实时日志流（SSE）。认证走 Bearer 头，原生 EventSource 带不了自定义头，
// 因此用 fetch 流式读取并手动解析 SSE 帧。
// 返回的 Promise 在流结束（服务端关闭/中断）时 resolve，调用方据此决定是否重连。
export async function streamLogs(lines, { onLine, onError, onReady, signal } = {}) {
  const token = localStorage.getItem('token')
  const resp = await fetch(`${API_BASE_URL}/logs/stream?lines=${lines}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    signal,
  })

  if (!resp.ok || !resp.body) {
    if (resp.status === 401) {
      // 与 axios 拦截器行为保持一致：清凭证并回登录页
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      window.location.href = '/login'
      throw new Error('登录已过期')
    }
    let message = `日志流连接失败 (${resp.status})`
    try {
      const body = await resp.json()
      if (body.error) message = body.error
    } catch { /* 非 JSON 响应，用默认提示 */ }
    throw new Error(message)
  }

  const handleBlock = (block) => {
    let event = 'message'
    const dataLines = []
    for (const rawLine of block.split('\n')) {
      if (rawLine.startsWith('event:')) event = rawLine.slice(6).trim()
      else if (rawLine.startsWith('data:')) dataLines.push(rawLine.slice(5).trim())
    }
    if (dataLines.length === 0) return
    let payload
    try {
      payload = JSON.parse(dataLines.join('\n'))
    } catch {
      return // 不完整的块，丢弃
    }
    if (event === 'log' && onLine) onLine(payload.line)
    else if (event === 'error' && onError) onError(payload.error)
    else if (event === 'log-ready' && onReady) onReady()
  }

  const reader = resp.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buffer += decoder.decode(value, { stream: true })
    let sep
    while ((sep = buffer.indexOf('\n\n')) !== -1) {
      const block = buffer.slice(0, sep)
      buffer = buffer.slice(sep + 2)
      handleBlock(block)
    }
  }
}
