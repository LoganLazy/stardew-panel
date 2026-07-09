<template>
  <div class="settings-container">
    <h2>⚙️ 系统设置</h2>

    <div class="card">
      <h3>🔐 修改密码</h3>
      <p class="description">定期修改密码可以提高账户安全性</p>

      <form @submit.prevent="handleChangePassword" class="password-form">
        <div class="form-group">
          <label>当前密码</label>
          <input
            v-model="oldPassword"
            type="password"
            placeholder="请输入当前密码"
            required
          />
        </div>

        <div class="form-group">
          <label>新密码</label>
          <input
            v-model="newPassword"
            type="password"
            placeholder="请输入新密码（至少6位）"
            required
          />
        </div>

        <div class="form-group">
          <label>确认新密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            required
          />
        </div>

        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '修改中...' : '修改密码' }}
        </button>
      </form>
    </div>

    <div class="card">
      <h3>👤 账户信息</h3>
      <div class="info-row">
        <span class="label">用户名：</span>
        <span class="value">{{ username }}</span>
      </div>
      <div class="info-row">
        <span class="label">登录状态：</span>
        <span class="value status-online">● 已登录</span>
      </div>
      <button @click="handleLogout" class="btn btn-secondary">退出登录</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import * as api from '@/api/api'

const router = useRouter()

const username = ref('')
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)

onMounted(() => {
  username.value = localStorage.getItem('username') || 'admin'
})

const handleChangePassword = async () => {
  if (newPassword.value !== confirmPassword.value) {
    alert('两次输入的新密码不一致')
    return
  }

  if (newPassword.value.length < 6) {
    alert('新密码至少需要6个字符')
    return
  }

  if (!confirm('确定要修改密码吗？修改后需要重新登录。')) return

  loading.value = true

  try {
    await api.changePassword(oldPassword.value, newPassword.value)
    alert('密码修改成功！请重新登录。')

    // 清除登录信息
    localStorage.removeItem('token')
    localStorage.removeItem('username')

    // 跳转到登录页
    router.push('/login')
  } catch (error) {
    console.error('Change password failed:', error)
    alert('修改密码失败：' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const handleLogout = async () => {
  if (!confirm('确定要退出登录吗？')) return

  try {
    await api.logout()
  } catch (error) {
    console.error('Logout failed:', error)
  }

  // 清除登录信息
  localStorage.removeItem('token')
  localStorage.removeItem('username')

  // 跳转到登录页
  router.push('/login')
}
</script>

<style scoped>
.settings-container {
  max-width: 800px;
  margin: 0 auto;
}

h2 {
  color: #1F2421;
  margin-bottom: 2rem;
}

.card {
  background: #FBF9F5;
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 1.5rem;
  border: 1px solid #E7E1D7;
}

.card h3 {
  color: #1F2421;
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
}

.description {
  color: #5C635D;
  margin: 0 0 1.5rem 0;
  font-size: 0.9rem;
}

.password-form {
  max-width: 400px;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #1F2421;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #C4612F;
}

.info-row {
  display: flex;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #E7E1D7;
}

.info-row:last-of-type {
  border-bottom: none;
  margin-bottom: 1rem;
}

.info-row .label {
  font-weight: 500;
  color: #5C635D;
  width: 120px;
}

.info-row .value {
  color: #1F2421;
}

.status-online {
  color: #10b981;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 999px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #C4612F;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #A94E22;
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #5C635D;
  color: white;
}

.btn-secondary:hover {
  background: #4a504b;
}
</style>
