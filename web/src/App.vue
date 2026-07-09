<template>
  <div class="app">
    <!-- 只在非登录页显示导航 -->
    <nav class="nav" v-if="$route.path !== '/login'">
      <div class="nav-brand">
        <h1>🌾 StardewPanel</h1>
      </div>
      <div class="nav-links">
        <router-link to="/">服务器</router-link>
        <router-link to="/mods">MOD</router-link>
        <router-link to="/players">玩家</router-link>
        <router-link to="/saves">存档</router-link>
        <router-link to="/logs">日志</router-link>
        <router-link to="/setup">安装</router-link>
        <router-link to="/settings">⚙️ 设置</router-link>
      </div>
    </nav>
    <main class="main" :class="{ 'no-nav': $route.path === '/login' }">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.app {
  min-height: 100vh;
  background: #F7F4EF;
}

.nav {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #E7E1D7;
  padding: 0 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.nav-brand h1 {
  font-size: 1.25rem;
  font-weight: 500;
  color: #1F2421;
  margin: 0;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-links a {
  color: #5C635D;
  text-decoration: none;
  font-weight: 300;
  transition: color 0.2s;
  font-size: 0.95rem;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: #C4612F;
}

.main {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
}

.main.no-nav {
  max-width: none;
  padding: 0;
}
</style>
