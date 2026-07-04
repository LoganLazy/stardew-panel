import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Server',
    component: () => import('@/views/Server.vue')
  },
  {
    path: '/mods',
    name: 'Mods',
    component: () => import('@/views/Mods.vue')
  },
  {
    path: '/players',
    name: 'Players',
    component: () => import('@/views/Players.vue')
  },
  {
    path: '/saves',
    name: 'Saves',
    component: () => import('@/views/Saves.vue')
  },
  {
    path: '/logs',
    name: 'Logs',
    component: () => import('@/views/Logs.vue')
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/Setup.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
