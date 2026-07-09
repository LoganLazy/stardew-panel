// utils/format.js - 格式化工具函数

/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string} - 格式化后的大小（如 "1.5 MB"）
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  if (!bytes) return '-'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

/**
 * 格式化日期时间
 * @param {string|Date} date - 日期
 * @returns {string} - 格式化后的日期（如 "2026-07-05 14:30"）
 */
export function formatDateTime(date) {
  if (!date) return '-'

  const d = new Date(date)
  if (isNaN(d.getTime())) return '-'

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hour = String(d.getHours()).padStart(2, '0')
  const minute = String(d.getMinutes()).padStart(2, '0')

  return `${year}-${month}-${day} ${hour}:${minute}`
}

/**
 * 格式化日期（仅日期）
 * @param {string|Date} date - 日期
 * @returns {string} - 格式化后的日期（如 "2026-07-05"）
 */
export function formatDate(date) {
  if (!date) return '-'

  const d = new Date(date)
  if (isNaN(d.getTime())) return '-'

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

/**
 * 格式化运行时长
 * @param {number} seconds - 秒数
 * @returns {string} - 格式化后的时长（如 "2小时30分钟"）
 */
export function formatDuration(seconds) {
  if (!seconds || seconds < 0) return '-'

  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)

  const parts = []
  if (days > 0) parts.push(`${days}天`)
  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分钟`)
  if (secs > 0 && days === 0 && hours === 0) parts.push(`${secs}秒`)

  return parts.length > 0 ? parts.join('') : '刚刚'
}

/**
 * 格式化相对时间
 * @param {string|Date} date - 日期
 * @returns {string} - 相对时间（如 "3分钟前"）
 */
export function formatRelativeTime(date) {
  if (!date) return '-'

  const d = new Date(date)
  if (isNaN(d.getTime())) return '-'

  const now = new Date()
  const diff = Math.floor((now - d) / 1000) // 秒

  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}天前`

  return formatDate(date)
}

/**
 * 截断长文本
 * @param {string} text - 文本
 * @param {number} maxLength - 最大长度
 * @returns {string} - 截断后的文本
 */
export function truncateText(text, maxLength = 50) {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

/**
 * 高亮关键字
 * @param {string} text - 文本
 * @param {string} keyword - 关键字
 * @returns {string} - HTML 字符串
 */
export function highlightKeyword(text, keyword) {
  if (!text || !keyword) return text
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}
