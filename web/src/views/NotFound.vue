<template>
  <div class="error-page">
    <div class="error-content">
      <h1>404</h1>
      <h2>页面不存在</h2>
      <p>抱歉，您访问的页面不存在或已被移除。</p>
      <div class="actions">
        <router-link to="/" class="btn btn-primary">返回首页</router-link>
        <button @click="goBack" class="btn btn-secondary">返回上一页</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F7F4EF 0%, #FBF9F5 100%);
  padding: 2rem;
}

.error-content {
  text-align: center;
  max-width: 500px;
}

h1 {
  font-size: 8rem;
  color: #C4612F;
  margin: 0;
  font-weight: 700;
  line-height: 1;
}

h2 {
  font-size: 2rem;
  color: #1F2421;
  margin: 1rem 0;
  font-weight: 600;
}

p {
  color: #5C635D;
  font-size: 1.1rem;
  margin: 1.5rem 0 2rem 0;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 999px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
  display: inline-block;
}

.btn-primary {
  background: #C4612F;
  color: white;
}

.btn-primary:hover {
  background: #A94E22;
  transform: translateY(-2px);
}

.btn-secondary {
  background: #5C635D;
  color: white;
}

.btn-secondary:hover {
  background: #4a504b;
}
</style>
