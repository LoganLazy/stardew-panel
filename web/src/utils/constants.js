// utils/constants.js - 常量定义

/**
 * API 刷新间隔（毫秒）
 */
export const REFRESH_INTERVALS = {
  SERVER_STATUS: 5000,    // 服务器状态：5秒
  PLAYERS: 10000,         // 在线玩家：10秒
  LOGS: 3000,             // 日志：3秒
  MODS: 30000,            // MOD 列表：30秒
  SAVES: 60000,           // 存档列表：60秒
}

/**
 * 文件上传限制（字节）
 */
export const FILE_SIZE_LIMITS = {
  MOD: 500 * 1024 * 1024,        // MOD：500MB
  GAME: 5 * 1024 * 1024 * 1024,  // 游戏文件：5GB
}

/**
 * 文件类型
 */
export const FILE_TYPES = {
  ZIP: '.zip',
  JSON: '.json',
}

/**
 * 服务器状态
 */
export const SERVER_STATUS = {
  RUNNING: 'running',
  STOPPED: 'stopped',
  STARTING: 'starting',
  STOPPING: 'stopping',
  RESTARTING: 'restarting',
  ERROR: 'error',
}

/**
 * 服务器状态显示文本
 */
export const SERVER_STATUS_TEXT = {
  [SERVER_STATUS.RUNNING]: '运行中',
  [SERVER_STATUS.STOPPED]: '已停止',
  [SERVER_STATUS.STARTING]: '启动中',
  [SERVER_STATUS.STOPPING]: '停止中',
  [SERVER_STATUS.RESTARTING]: '重启中',
  [SERVER_STATUS.ERROR]: '错误',
}

/**
 * 服务器状态颜色
 */
export const SERVER_STATUS_COLOR = {
  [SERVER_STATUS.RUNNING]: '#10b981',   // 绿色
  [SERVER_STATUS.STOPPED]: '#6b7280',   // 灰色
  [SERVER_STATUS.STARTING]: '#3b82f6',  // 蓝色
  [SERVER_STATUS.STOPPING]: '#f59e0b',  // 橙色
  [SERVER_STATUS.RESTARTING]: '#8b5cf6', // 紫色
  [SERVER_STATUS.ERROR]: '#ef4444',     // 红色
}

/**
 * 安装方法
 */
export const INSTALL_METHODS = {
  UPLOAD: 'upload',
  PATH: 'path',
  STEAMCMD: 'steamcmd',
}

/**
 * 安装方法显示文本
 */
export const INSTALL_METHOD_TEXT = {
  [INSTALL_METHODS.UPLOAD]: '上传文件',
  [INSTALL_METHODS.PATH]: '指定路径',
  [INSTALL_METHODS.STEAMCMD]: 'SteamCMD',
}

/**
 * 日志级别
 */
export const LOG_LEVELS = {
  DEBUG: 'debug',
  INFO: 'info',
  WARN: 'warn',
  ERROR: 'error',
}

/**
 * MOD 状态
 */
export const MOD_STATUS = {
  ENABLED: true,
  DISABLED: false,
}

/**
 * 默认配置
 */
export const DEFAULT_CONFIG = {
  LOG_LINES: 100,           // 默认日志行数
  BACKUP_RETENTION: 30,     // 备份保留天数
  MAX_PLAYERS: 8,           // 最大玩家数
}

/**
 * 路由路径
 */
export const ROUTES = {
  SERVER: '/',
  SETUP: '/setup',
  MODS: '/mods',
  PLAYERS: '/players',
  SAVES: '/saves',
  LOGS: '/logs',
}

/**
 * 本地存储键
 */
export const STORAGE_KEYS = {
  THEME: 'stardew_panel_theme',
  LANGUAGE: 'stardew_panel_language',
  LAST_PATH: 'stardew_panel_last_path',
}

/**
 * 错误消息
 */
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查服务器是否正常运行',
  TIMEOUT: '请求超时，请稍后重试',
  SERVER_ERROR: '服务器内部错误，请联系管理员',
  FILE_TOO_LARGE: '文件过大',
  INVALID_FILE_TYPE: '文件类型不正确',
  PERMISSION_DENIED: '没有权限执行此操作',
}

/**
 * 成功消息
 */
export const SUCCESS_MESSAGES = {
  SERVER_STARTED: '服务器启动成功',
  SERVER_STOPPED: '服务器停止成功',
  SERVER_RESTARTED: '服务器重启成功',
  MOD_UPLOADED: 'MOD 上传成功',
  MOD_ENABLED: 'MOD 已启用',
  MOD_DISABLED: 'MOD 已禁用',
  MOD_DELETED: 'MOD 已删除',
  BACKUP_CREATED: '备份创建成功',
  BACKUP_RESTORED: '存档恢复成功',
  BACKUP_DELETED: '备份删除成功',
}

/**
 * 确认消息
 */
export const CONFIRM_MESSAGES = {
  STOP_SERVER: '确定要停止服务器吗？在线玩家将被断开连接。',
  RESTART_SERVER: '确定要重启服务器吗？在线玩家将被断开连接。',
  DELETE_MOD: '确定要删除这个 MOD 吗？',
  RESTORE_BACKUP: '确定要恢复这个存档吗？当前存档将被覆盖。',
  DELETE_BACKUP: '确定要删除这个备份吗？',
  KICK_PLAYER: '确定要踢出这个玩家吗？',
}
