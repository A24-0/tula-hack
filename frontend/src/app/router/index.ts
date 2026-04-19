import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'

export default createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    return { top: 0, behavior: 'smooth' }
  },
  routes: [
    { path: '/', name: 'home', component: HomePage },
  ],
})
