// utils/errors.js - 用户友好的错误处理

/**
 * 将技术错误转换为用户友好的消息
 * @param {Error} error - 错误对象
 * @returns {string} - 友好的错误消息
 */
export function getFriendlyErrorMessage(error) {
  // 如果有后端返回的错误消息，直接使用
  if (error.response?.data?.error) {
    return error.response.data.error
  }

  // 网络错误
  if (error.message === 'Network Error') {
    return '网络连接失败，请检查服务器是否正常运行'
  }

  // 超时错误
  if (error.code === 'ECONNABORTED') {
    return '请求超时，请稍后重试'
  }

  // HTTP 状态码错误
  if (error.response) {
    const status = error.response.status
    switch (status) {
      case 400:
        return '请求参数错误，请检查输入内容'
      case 401:
        return '未授权，请先登录'
      case 403:
        return '没有权限执行此操作'
      case 404:
        return '请求的资源不存在'
      case 500:
        return '服务器内部错误，请联系管理员'
      case 502:
        return '网关错误，服务器可能正在重启'
      case 503:
        return '服务暂时不可用，请稍后重试'
      default:
        return `服务器错误 (${status})`
    }
  }

  // 默认错误消息
  return error.message || '操作失败，请重试'
}

/**
 * 显示友好的错误提示
 * @param {Error} error - 错误对象
 * @param {string} action - 操作名称（如"上传文件"）
 */
export function showFriendlyError(error, action) {
  const message = getFriendlyErrorMessage(error)
  alert(`${action}失败：${message}`)
}

/**
 * 显示成功提示
 * @param {string} action - 操作名称
 */
export function showSuccess(action) {
  alert(`${action}成功！`)
}

/**
 * 显示确认对话框
 * @param {string} message - 确认消息
 * @param {string} warning - 警告信息（可选）
 * @returns {boolean} - 用户是否确认
 */
export function showConfirm(message, warning = null) {
  let fullMessage = message
  if (warning) {
    fullMessage += `\n\n⚠️ ${warning}`
  }
  return confirm(fullMessage)
}
