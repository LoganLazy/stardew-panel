import { defineStore } from 'pinia'
import axios from 'axios'

export const useServerStore = defineStore('server', {
  state: () => ({
    status: 'stopped',
    stats: {
      onlinePlayers: 0,
      version: '1.6.9',
      uptime: 0,
      modsCount: 0
    }
  }),

  actions: {
    async getStatus() {
      try {
        const { data } = await axios.get('/api/v1/server/status')
        this.status = data.status
        return data
      } catch (error) {
        console.error('Failed to get server status:', error)
        return { status: 'stopped' }
      }
    },

    async startServer() {
      try {
        await axios.post('/api/v1/server/start')
        this.status = 'running'
      } catch (error) {
        console.error('Failed to start server:', error)
      }
    },

    async stopServer() {
      try {
        await axios.post('/api/v1/server/stop')
        this.status = 'stopped'
      } catch (error) {
        console.error('Failed to stop server:', error)
      }
    },

    async restartServer() {
      try {
        await axios.post('/api/v1/server/restart')
      } catch (error) {
        console.error('Failed to restart server:', error)
      }
    }
  }
})
