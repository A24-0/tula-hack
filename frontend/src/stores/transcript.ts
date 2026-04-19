import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { TranscriptResult } from '@/types/transcript'
import type { ProcessingLog } from '@/types/log'

export const useTranscriptStore = defineStore('transcript', () => {
  const transcript = ref<TranscriptResult | null>(null)
  const logs = ref<ProcessingLog[]>([])

  const entityCount = computed(() => transcript.value?.entities.length ?? 0)
  const hasResult = computed(() => Boolean(transcript.value))

  function setResult(result: TranscriptResult, nextLogs: ProcessingLog[]) {
    transcript.value = result
    logs.value = nextLogs
  }

  function reset() {
    transcript.value = null
    logs.value = []
  }

  return { transcript, logs, entityCount, hasResult, setResult, reset }
})
