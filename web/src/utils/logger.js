// 日志工具 - 生产环境优化
const isDev = import.meta.env.DEV

export const logger = {
  log: (...args) => {
    if (isDev) {
      console.log(...args)
    }
  },

  error: (...args) => {
    console.error(...args)
  },

  warn: (...args) => {
    if (isDev) {
      console.warn(...args)
    }
  },

  info: (...args) => {
    if (isDev) {
      console.info(...args)
    }
  }
}

// 使用方法：
// import { logger } from '@/utils/logger'
// logger.log('Debug info')  // 只在开发环境输出
// logger.error('Error')      // 生产环境也输出
