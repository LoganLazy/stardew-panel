// utils/validators.js - 验证工具函数

/**
 * 验证文件大小
 * @param {File} file - 文件对象
 * @param {number} maxSize - 最大大小（字节）
 * @returns {boolean} - 是否有效
 */
export function validateFileSize(file, maxSize) {
  if (!file) return false
  return file.size <= maxSize
}

/**
 * 验证文件类型
 * @param {File} file - 文件对象
 * @param {string[]} allowedTypes - 允许的类型（如 ['.zip', '.json']）
 * @returns {boolean} - 是否有效
 */
export function validateFileType(file, allowedTypes) {
  if (!file || !allowedTypes || allowedTypes.length === 0) return false

  const fileName = file.name.toLowerCase()
  return allowedTypes.some(type => fileName.endsWith(type.toLowerCase()))
}

/**
 * 验证路径
 * @param {string} path - 路径
 * @returns {boolean} - 是否有效
 */
export function validatePath(path) {
  if (!path || typeof path !== 'string') return false

  // 不允许包含 ..
  if (path.includes('..')) return false

  // 不允许空白路径
  if (path.trim() === '') return false

  return true
}

/**
 * 验证 Steam 用户名
 * @param {string} username - 用户名
 * @returns {boolean} - 是否有效
 */
export function validateSteamUsername(username) {
  if (!username || typeof username !== 'string') return false

  // 长度检查
  if (username.length < 3 || username.length > 32) return false

  // 只允许字母、数字、下划线、连字符
  const regex = /^[a-zA-Z0-9_-]+$/
  return regex.test(username)
}

/**
 * 验证密码强度（简单检查）
 * @param {string} password - 密码
 * @returns {object} - { valid: boolean, message: string }
 */
export function validatePassword(password) {
  if (!password || typeof password !== 'string') {
    return { valid: false, message: '密码不能为空' }
  }

  if (password.length < 6) {
    return { valid: false, message: '密码至少需要 6 个字符' }
  }

  return { valid: true, message: '' }
}

/**
 * 验证玩家名称
 * @param {string} name - 玩家名称
 * @returns {boolean} - 是否有效
 */
export function validatePlayerName(name) {
  if (!name || typeof name !== 'string') return false

  // 长度检查
  if (name.trim().length === 0) return false
  if (name.length > 32) return false

  return true
}

/**
 * 验证端口号
 * @param {number|string} port - 端口号
 * @returns {boolean} - 是否有效
 */
export function validatePort(port) {
  const portNum = typeof port === 'string' ? parseInt(port, 10) : port

  if (isNaN(portNum)) return false
  if (portNum < 1 || portNum > 65535) return false

  return true
}

/**
 * 验证 IP 地址
 * @param {string} ip - IP 地址
 * @returns {boolean} - 是否有效
 */
export function validateIP(ip) {
  if (!ip || typeof ip !== 'string') return false

  // IPv4 正则
  const ipv4Regex = /^(\d{1,3}\.){3}\d{1,3}$/
  if (!ipv4Regex.test(ip)) return false

  // 检查每个部分是否在 0-255 之间
  const parts = ip.split('.')
  return parts.every(part => {
    const num = parseInt(part, 10)
    return num >= 0 && num <= 255
  })
}

/**
 * 验证 URL
 * @param {string} url - URL
 * @returns {boolean} - 是否有效
 */
export function validateURL(url) {
  if (!url || typeof url !== 'string') return false

  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 验证邮箱
 * @param {string} email - 邮箱地址
 * @returns {boolean} - 是否有效
 */
export function validateEmail(email) {
  if (!email || typeof email !== 'string') return false

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * 验证表单数据
 * @param {object} data - 表单数据
 * @param {object} rules - 验证规则
 * @returns {object} - { valid: boolean, errors: object }
 */
export function validateForm(data, rules) {
  const errors = {}
  let valid = true

  for (const [field, rule] of Object.entries(rules)) {
    const value = data[field]

    // 必填检查
    if (rule.required && (!value || value.toString().trim() === '')) {
      errors[field] = rule.message || `${field} 不能为空`
      valid = false
      continue
    }

    // 自定义验证函数
    if (rule.validator && typeof rule.validator === 'function') {
      const result = rule.validator(value)
      if (!result) {
        errors[field] = rule.message || `${field} 验证失败`
        valid = false
      }
    }

    // 最小长度
    if (rule.minLength && value && value.length < rule.minLength) {
      errors[field] = rule.message || `${field} 至少需要 ${rule.minLength} 个字符`
      valid = false
    }

    // 最大长度
    if (rule.maxLength && value && value.length > rule.maxLength) {
      errors[field] = rule.message || `${field} 不能超过 ${rule.maxLength} 个字符`
      valid = false
    }

    // 正则匹配
    if (rule.pattern && value && !rule.pattern.test(value)) {
      errors[field] = rule.message || `${field} 格式不正确`
      valid = false
    }
  }

  return { valid, errors }
}
