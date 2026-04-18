import { createRouter, createWebHistory } from 'vue-router'
import UploadView from '@/views/UploadView.vue'
import ResultView from '@/views/ResultView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: UploadView },
    { path: '/result/:id', component: ResultView },
  ],
})
