<template>
  <div class="container">
    <div v-if="loading" class="state-msg">Загрузка...</div>

    <div v-else-if="error" class="card">
      <p class="error-msg">{{ error }}</p>
      <RouterLink to="/" class="btn btn-secondary" style="margin-top:14px">Назад</RouterLink>
    </div>

    <template v-else-if="data">
      <div class="top-bar">
        <RouterLink to="/" class="btn btn-secondary">Новая запись</RouterLink>
        <a
          v-if="data.redacted_audio_path"
          :href="`/audio/${jobId}/redacted`"
          class="btn btn-primary"
          download
        >
          Скачать анонимизированное аудио
        </a>
      </div>

      <div class="card">
        <TranscriptBlock :data="data" />
      </div>

      <div class="card" style="margin-top:12px">
        <ReportTable :events="data.pii_events" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import TranscriptBlock from '@/components/TranscriptBlock.vue'
import ReportTable from '@/components/ReportTable.vue'
import { getTranscript } from '@/api'
import type { TranscriptResult } from '@/types'

const route = useRoute()
const jobId = route.params.id as string

const loading = ref(true)
const error = ref('')
const data = ref<TranscriptResult | null>(null)

onMounted(async () => {
  try {
    data.value = await getTranscript(jobId)
  } catch {
    error.value = 'Не удалось загрузить результат'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.state-msg {
  text-align: center;
  padding: 60px;
  color: #475569;
  font-size: 14px;
}

.error-msg {
  color: #f87171;
  font-size: 14px;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  gap: 10px;
}
</style>
