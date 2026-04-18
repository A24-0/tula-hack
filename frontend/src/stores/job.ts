import { ref } from 'vue'
import { defineStore } from 'pinia'
import { uploadAudio, getJobStatus } from '@/api'

export const useJobStore = defineStore('job', () => {
  const jobId = ref<string | null>(null)
  const status = ref<'idle' | 'uploading' | 'processing' | 'done' | 'error'>('idle')
  const progress = ref(0)
  const errorMsg = ref('')

  let pollTimer: ReturnType<typeof setInterval> | null = null

  async function upload(file: File) {
    status.value = 'uploading'
    progress.value = 0
    errorMsg.value = ''

    try {
      const id = await uploadAudio(file, (pct) => {
        progress.value = pct
      })
      jobId.value = id
      status.value = 'processing'
      progress.value = 90
      startPolling(id)
    } catch {
      status.value = 'error'
      errorMsg.value = 'Ошибка при загрузке файла'
    }
  }

  function startPolling(id: string) {
    pollTimer = setInterval(async () => {
      try {
        const job = await getJobStatus(id)
        if (job.status === 'done') {
          status.value = 'done'
          progress.value = 100
          stopPolling()
        } else if (job.status === 'error') {
          status.value = 'error'
          errorMsg.value = job.error ?? 'Ошибка обработки'
          stopPolling()
        }
      } catch {
        // сеть временно недоступна, попробуем в следующий раз
      }
    }, 2000)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  function reset() {
    stopPolling()
    jobId.value = null
    status.value = 'idle'
    progress.value = 0
    errorMsg.value = ''
  }

  return { jobId, status, progress, errorMsg, upload, reset }
})
