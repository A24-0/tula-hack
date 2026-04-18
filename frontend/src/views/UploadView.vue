<template>
  <div class="container">
    <div class="card">
      <h1 class="page-title">Загрузка записи</h1>

      <DropZone @change="onFileSelected" />

      <div v-if="status !== 'idle'" class="progress-wrap">
        <div class="progress-bar">
          <div class="progress-bar__fill" :style="{ width: progress + '%' }" />
        </div>
        <span class="progress-label">{{ statusText }}</span>
      </div>

      <button
        class="btn btn-primary submit-btn"
        :disabled="!selectedFile || status === 'uploading' || status === 'processing'"
        @click="start"
      >
        {{ status === 'uploading' || status === 'processing' ? 'Обрабатывается...' : 'Анализировать' }}
      </button>

      <p v-if="status === 'error'" class="error-msg">{{ errorMsg }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import DropZone from '@/components/DropZone.vue'
import { useJobStore } from '@/stores/job'

const router = useRouter()
const store = useJobStore()
const { status, progress, errorMsg, jobId } = storeToRefs(store)
const selectedFile = ref<File | null>(null)

const statusText = computed(() => {
  if (status.value === 'uploading') return `Загрузка... ${progress.value}%`
  if (status.value === 'processing') return 'Обработка...'
  return ''
})

function onFileSelected(f: File) {
  selectedFile.value = f
  store.reset()
}

async function start() {
  if (!selectedFile.value) return
  await store.upload(selectedFile.value)
}

watch(status, (val) => {
  if (val === 'done' && jobId.value) {
    router.push(`/result/${jobId.value}`)
  }
})
</script>

<style scoped>
.card {
  max-width: 560px;
  margin: 0 auto;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
}

.progress-wrap {
  margin-top: 16px;
}

.progress-bar {
  height: 4px;
  background: #1e293b;
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar__fill {
  height: 100%;
  background: #2563eb;
  transition: width 0.3s ease;
}

.progress-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 6px;
  display: block;
}

.submit-btn {
  margin-top: 16px;
  width: 100%;
  justify-content: center;
  padding: 11px;
}

.error-msg {
  color: #f87171;
  font-size: 13px;
  margin-top: 10px;
  text-align: center;
}
</style>
