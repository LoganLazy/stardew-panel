import api from './index'

// 认证相关
export const login = (username, password) =>
  api.post('/auth/login', { username, password })

export const logout = () => api.post('/auth/logout')

export const checkAuth = () => api.get('/auth/check')

export const changePassword = (oldPassword, newPassword) =>
  api.post('/auth/change-password', {
    old_password: oldPassword,
    new_password: newPassword
  })

// 检查安装状态
export const checkInstallation = () => api.get('/install/check')

// 验证游戏路径
export const verifyGamePath = (path, detectSMAPI) =>
  api.post('/install/verify', { path, detect_smapi: detectSMAPI })

// 上传游戏文件
export const uploadGameFiles = (formData, onProgress) =>
  api.post('/install/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress,
  })

// SteamCMD 安装
export const installViaSteamCMD = (path, username, password, installSMAPI) =>
  api.post('/install/steamcmd', {
    path,
    steam_username: username,
    steam_password: password,
    install_smapi: installSMAPI
  })

// 获取安装状态
export const getInstallStatus = () => api.get('/install/status')

// 服务器管理
export const getServerStatus = () => api.get('/server/status')
export const startServer = () => api.post('/server/start')
export const stopServer = () => api.post('/server/stop')
export const restartServer = () => api.post('/server/restart')

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
export const kickPlayer = (playerName) =>
  api.post('/players/kick', { player_name: playerName })
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
