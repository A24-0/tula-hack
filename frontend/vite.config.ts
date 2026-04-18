import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    port: 3000,
    proxy: {
      '/upload': 'http://localhost:8080',
      '/jobs': 'http://localhost:8080',
      '/transcript': 'http://localhost:8080',
      '/audio': 'http://localhost:8080',
      '/logs': 'http://localhost:8080',
    },
  },
})
