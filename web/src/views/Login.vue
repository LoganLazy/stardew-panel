<template>
  <div class="login-container">
    <div class="login-box">
      <div class="logo">
        <h1>🌾 StardewPanel</h1>
        <p>星露谷物语服务器管理面板</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>用户名</label>
          <input
            v-model="username"
            type="text"
            placeholder="请输入用户名"
            autofocus
            required
          />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="password"
            type="password"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="error-message" v-if="errorMessage">
          ❌ {{ errorMessage }}
        </div>

        <button type="submit" class="btn-login" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <div class="info-box">
        <p>💡 <strong>首次登录默认账号：</strong></p>
        <p>用户名: <code>admin</code></p>
        <p>密码: <code>admin123</code></p>
        <p class="warning">⚠️ 登录后请立即修改密码！</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as api from '@/api/api'

const router = useRouter()

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const handleLogin = async () => {
  errorMessage.value = ''
  loading.value = true

  try {
    const { data } = await api.login(username.value, password.value)

    if (data.success) {
      // 保存 token
      localStorage.setItem('token', data.token)
      localStorage.setItem('username', username.value)

      // 跳转到首页
      router.push('/')
    } else {
      errorMessage.value = data.message || '登录失败'
    }
  } catch (error) {
    console.error('Login failed:', error)
    errorMessage.value = error.response?.data?.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F7F4EF 0%, #FBF9F5 100%);
  padding: 2rem;
}

.login-box {
  background: #FFFFFF;
  border-radius: 16px;
  padding: 3rem;
  max-width: 450px;
  width: 100%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.logo {
  text-align: center;
  margin-bottom: 2rem;
}

.logo h1 {
  font-size: 2rem;
  color: #1F2421;
  margin: 0 0 0.5rem 0;
  font-weight: 600;
}

.logo p {
  color: #5C635D;
  font-size: 0.95rem;
  margin: 0;
}

.login-form {
  margin-bottom: 2rem;
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
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #C4612F;
}

.error-message {
  background: #FFF3F3;
  border: 1px solid #FFCDD2;
  color: #C62828;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.btn-login {
  width: 100%;
  padding: 0.875rem;
  background: #C4612F;
  color: white;
  border: none;
  border-radius: 999px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-login:hover:not(:disabled) {
  background: #A94E22;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(196, 97, 47, 0.3);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.info-box {
  background: #FFFBF0;
  border: 1px solid #FFE69C;
  border-radius: 8px;
  padding: 1rem;
  font-size: 0.9rem;
}

.info-box p {
  margin: 0.5rem 0;
  color: #5C635D;
}

.info-box code {
  background: #FFF;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  color: #C4612F;
}

.info-box .warning {
  color: #C62828;
  font-weight: 500;
  margin-top: 0.5rem;
}
</style>
