import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { uploadAudio } from '@/api/upload'
import { getJob } from '@/api/jobs'
import { getTranscript } from '@/api/transcript'
import { getLogs } from '@/api/logs'
import { fallbackLogs } from '@/utils/adapters'
import { useAppStore } from './app'
import { useTranscriptStore } from './transcript'
import type { JobStage, JobStatus } from '@/types/job'

interface SourceAudio {
  file?: File
  blobUrl?: string
  name: string
  size?: number
  mime?: string
  kind: 'file' | 'recording'
}

const stageOrder: JobStage[] = ['uploading', 'queued', 'transcribing', 'detecting', 'redacting_text', 'redacting_audio', 'packaging', 'done']

export const stageLabels: Record<string, string> = {
  uploading: 'Загрузка аудио',
  queued: 'Постановка в очередь',
  transcribing: 'Распознавание речи',
  detecting: 'Поиск персональных данных',
  redacting_text: 'Анонимизация текста',
  redacting_audio: 'Маскирование аудио',
  packaging: 'Подготовка материалов',
  done: 'Обработка завершена',
  error: 'Ошибка обработки',
}

export const useProcessingStore = defineStore('processing', () => {
  const app = useAppStore()
  const transcriptStore = useTranscriptStore()
  const source = ref<SourceAudio | null>(null)
  const jobId = ref<string | null>(null)
  const status = ref<JobStatus>('idle')
  const stage = ref<JobStage>('queued')
  const progress = ref(0)
  const error = ref('')
  const startedAt = ref<number | null>(null)
  const pollStartedAt = ref<number | null>(null)
  let pollTimer: number | null = null
  let smoothTimer: number | null = null

  const isBusy = computed(() => ['pending', 'processing'].includes(status.value))
  const hasSource = computed(() => Boolean(source.value))
  const summaryStatus = computed(() => {
    if (status.value === 'done') return 'Готово'
    if (status.value === 'error') return 'Нужна повторная обработка'
    if (isBusy.value) return stageLabels[stage.value] ?? 'Обработка'
    return 'Готово к отправке'
  })

  function setSource(next: SourceAudio) {
    if (source.value?.blobUrl && source.value.blobUrl !== next.blobUrl) URL.revokeObjectURL(source.value.blobUrl)
    source.value = next
    status.value = 'idle'
    progress.value = 0
    error.value = ''
    jobId.value = null
    transcriptStore.reset()
  }

  async function start() {
    if (!source.value?.file) {
      app.pushToast({ tone: 'error', title: 'Файл не выбран', text: 'Загрузите файл или запишите аудио.' })
      return
    }

    clearTimers()
    status.value = 'pending'
    stage.value = 'uploading'
    progress.value = 1
    error.value = ''
    startedAt.value = Date.now()
    startSmoothing()

    try {
      const id = await uploadAudio(source.value.file, source.value.name, (value) => {
        progress.value = Math.max(progress.value, value)
      })
      jobId.value = id
      status.value = 'processing'
      stage.value = 'queued'
      progress.value = Math.max(progress.value, 18)
      pollStartedAt.value = Date.now()
      await pollOnce(id)
      pollTimer = window.setInterval(() => pollOnce(id), 1800)
    } catch (err) {
      fail(err instanceof Error ? err.message : 'Не удалось отправить аудио')
    }
  }

  async function pollOnce(id: string) {
    try {
      const job = await getJob(id)
      if (job.stage) stage.value = normalizeStage(job.stage)
      if (typeof job.progress === 'number') progress.value = Math.max(progress.value, Math.min(job.progress, 99))

      if (job.status === 'done') {
        stage.value = 'packaging'
        progress.value = Math.max(progress.value, 96)
        await loadResult(id)
        complete()
      } else if (job.status === 'error') {
        fail(job.error || 'Backend вернул ошибку обработки')
      } else {
        status.value = 'processing'
      }

      if (pollStartedAt.value && Date.now() - pollStartedAt.value > 20 * 60 * 1000) {
        fail('Превышено время ожидания обработки')
      }
    } catch (err) {
      if (!pollStartedAt.value || Date.now() - pollStartedAt.value < 15000) return
      fail(err instanceof Error ? err.message : 'Не удалось получить статус задачи')
    }
  }

  async function loadResult(id: string) {
    const transcript = await getTranscript(id)
    let logs = []
    try {
      logs = await getLogs(id)
    } catch {
      logs = fallbackLogs(transcript)
      app.pushToast({ tone: 'info', title: 'Журнал собран локально', text: 'Backend не вернул /logs, показаны события по результату.' })
    }
    transcriptStore.setResult(transcript, logs.length ? logs : fallbackLogs(transcript))
  }

  function complete() {
    clearTimers()
    status.value = 'done'
    stage.value = 'done'
    progress.value = 100
    app.pushToast({ tone: 'success', title: 'Обработка завершена', text: 'Материалы готовы к просмотру и скачиванию.' })
  }

  function fail(message: string) {
    clearTimers()
    status.value = 'error'
    stage.value = 'error'
    error.value = message
    app.pushToast({ tone: 'error', title: 'Ошибка обработки', text: message })
  }

  function reset() {
    clearTimers()
    if (source.value?.blobUrl) URL.revokeObjectURL(source.value.blobUrl)
    source.value = null
    jobId.value = null
    status.value = 'idle'
    stage.value = 'queued'
    progress.value = 0
    error.value = ''
    startedAt.value = null
    pollStartedAt.value = null
    transcriptStore.reset()
  }

  function startSmoothing() {
    smoothTimer = window.setInterval(() => {
      if (!isBusy.value || progress.value >= 94) return
      const stageIndex = Math.max(0, stageOrder.indexOf(stage.value))
      const stageCap = Math.min(94, 16 + stageIndex * 12)
      progress.value = Math.min(stageCap, progress.value + Math.random() * 1.8)
    }, 650)
  }

  function clearTimers() {
    if (pollTimer) window.clearInterval(pollTimer)
    if (smoothTimer) window.clearInterval(smoothTimer)
    pollTimer = null
    smoothTimer = null
  }

  function normalizeStage(value: string): JobStage {
    return value in stageLabels ? value as JobStage : 'queued'
  }

  return {
    source,
    jobId,
    status,
    stage,
    progress,
    error,
    isBusy,
    hasSource,
    summaryStatus,
    setSource,
    start,
    reset,
  }
})
