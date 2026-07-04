<template>
  <div class="players-view">
    <h1>玩家<em>管理</em></h1>

    <div class="card">
      <h3>在线玩家</h3>
      <div class="players-list" v-if="onlinePlayers.length > 0">
        <div class="player-item" v-for="player in onlinePlayers" :key="player.id">
          <div>{{ player.name }}</div>
          <button class="btn btn-secondary" @click="kickPlayer(player.id)">踢出</button>
        </div>
      </div>
      <div class="empty" v-else>
        <p>当前没有在线玩家</p>
      </div>
    </div>

    <div class="card">
      <h3>白名单</h3>
      <div class="whitelist-add">
        <input v-model="newPlayerName" placeholder="输入玩家名称" />
        <button class="btn btn-primary" @click="addToWhitelist">添加</button>
      </div>
      <div class="players-list" v-if="whitelist.length > 0">
        <div class="player-item" v-for="player in whitelist" :key="player.id">
          <div>{{ player.name }}</div>
          <button class="btn-icon" @click="removeFromWhitelist(player.id)">移除</button>
        </div>
      </div>
      <div class="empty" v-else>
        <p>白名单为空</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const onlinePlayers = ref([])
const whitelist = ref([])
const newPlayerName = ref('')

const loadPlayers = async () => {
  try {
    const { data } = await axios.get('/api/v1/players')
    onlinePlayers.value = data.online || []
    whitelist.value = data.whitelist || []
  } catch (error) {
    console.error('Failed to load players:', error)
  }
}

const kickPlayer = async (id) => {
  try {
    await axios.post('/api/v1/players/kick', { id })
    await loadPlayers()
  } catch (error) {
    console.error('Failed to kick player:', error)
  }
}

const addToWhitelist = async () => {
  if (!newPlayerName.value.trim()) return
  try {
    await axios.post('/api/v1/players/whitelist', { name: newPlayerName.value })
    newPlayerName.value = ''
    await loadPlayers()
  } catch (error) {
    console.error('Failed to add to whitelist:', error)
  }
}

const removeFromWhitelist = async (id) => {
  try {
    await axios.delete(`/api/v1/players/whitelist/${id}`)
    await loadPlayers()
  } catch (error) {
    console.error('Failed to remove from whitelist:', error)
  }
}

onMounted(loadPlayers)
</script>

<style scoped>
h1 {
  margin-bottom: 2rem;
  font-size: 2.5rem;
}

h1 em {
  color: #C4612F;
  font-style: italic;
}

.card {
  margin-bottom: 2rem;
}

.card h3 {
  margin-top: 0;
}

.whitelist-add {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.whitelist-add input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.95rem;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.player-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #FBF9F5;
  border: 1px solid #E7E1D7;
  border-radius: 8px;
}

.empty {
  text-align: center;
  padding: 2rem;
  color: #5C635D;
}

.btn-icon {
  background: transparent;
  color: #5C635D;
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
}

.btn-icon:hover {
  color: #C4612F;
}
</style>
