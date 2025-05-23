import { createMemoryHistory, createRouter } from 'vue-router'

import Mine from './Mine.vue'
import Login from './Login.vue'
import Translation from './Translation.vue'


const routes = [
  { path: '/Mine', component: Mine },
  { path: '/', component: Translation },
  { path: '/login', component: Login},
]

export  const router = createRouter({
  history: createMemoryHistory(),
  routes,
})