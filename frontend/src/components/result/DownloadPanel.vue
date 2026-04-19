<template>
  <aside class="downloads glass-panel">
    <div class="panel-content">
      <span class="eyebrow">Export</span>
      <h3>Пакет готов</h3>
      <p>Скачайте аудио, транскрипты или JSON-отчёт для внутреннего аудита.</p>
      <div class="downloads__grid">
        <button class="button button-primary" type="button" @click="downloadFromUrl('voice-redaction-redacted-audio.webm', redactedUrl)">Анонимизированное аудио</button>
        <button class="button button-secondary" type="button" @click="downloadText('original-transcript.txt', data.original_text)">Исходный транскрипт</button>
        <button class="button button-secondary" type="button" @click="downloadText('redacted-transcript.txt', data.redacted_text)">Безопасный транскрипт</button>
        <button class="button button-secondary" type="button" @click="downloadJson('voice-redaction-report.json', report)">Лог и отчёт</button>
        <button v-if="sourceUrl" class="button button-secondary" type="button" @click="downloadFromUrl(sourceName, sourceUrl)">Исходное аудио</button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ProcessingLog } from '@/types/log'
import type { TranscriptResult } from '@/types/transcript'
import { downloadFromUrl, downloadJson, downloadText } from '@/utils/download'

const props = defineProps<{
  data: TranscriptResult
  logs: ProcessingLog[]
  redactedUrl: string
  sourceUrl?: string
  sourceName: string
}>()

const report = computed(() => ({
  transcript: props.data,
  logs: props.logs,
}))
</script>

<style scoped>
.downloads {
  height: 100%;
  padding: 22px;
}

h3 {
  margin: 12px 0 0;
  font-size: 30px;
  font-weight: 650;
  letter-spacing: -0.045em;
}

p {
  margin: 12px 0 0;
  color: var(--text-muted);
  line-height: 1.6;
}

.downloads__grid {
  display: grid;
  gap: 10px;
  margin-top: 24px;
}

.button {
  width: 100%;
}
</style>
