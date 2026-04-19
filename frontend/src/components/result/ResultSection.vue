<template>
  <Transition name="section">
    <section v-if="data && jobId" id="result" class="section container result-section">
      <div class="result-section__head">
        <span class="eyebrow">Result</span>
        <h2 class="section-title">Материалы готовы к проверке и скачиванию</h2>
        <p class="section-copy">
          Сравните аудио и транскрипт, проверьте найденные сущности и выгрузите отчёт для аудита.
        </p>
      </div>

      <div class="result-grid">
        <WaveAudioPlayer ref="audioPlayer" :src="redactedUrl" :entities="data.entities" />
        <DownloadPanel
          :data="data"
          :logs="logs"
          :redacted-url="redactedUrl"
          :source-url="source?.blobUrl"
          :source-name="source?.name || 'voice-redaction-original-audio.webm'"
        />
      </div>

      <TranscriptViewer :data="data" :current-time="currentTime" @seek="seek" />
      <ReportBlock :data="data" :duration="duration" />
      <LogsBlock :events="logs" />

      <div class="bottom-cta">
        <button class="button button-primary" type="button" @click="processing.reset">Обработать новый файл</button>
      </div>
    </section>
  </Transition>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import DownloadPanel from '@/components/result/DownloadPanel.vue'
import LogsBlock from '@/components/logs/LogsBlock.vue'
import WaveAudioPlayer from '@/components/player/WaveAudioPlayer.vue'
import ReportBlock from '@/components/report/ReportBlock.vue'
import TranscriptViewer from '@/components/transcript/TranscriptViewer.vue'
import { redactedAudioUrl } from '@/api/audio'
import { usePlayerStore } from '@/stores/player'
import { useProcessingStore } from '@/stores/processing'
import { useTranscriptStore } from '@/stores/transcript'

const processing = useProcessingStore()
const transcriptStore = useTranscriptStore()
const player = usePlayerStore()
const { jobId, source } = storeToRefs(processing)
const { transcript: data, logs } = storeToRefs(transcriptStore)
const { currentTime, duration } = storeToRefs(player)
const audioPlayer = ref<InstanceType<typeof WaveAudioPlayer> | null>(null)
const redactedUrl = computed(() => jobId.value ? redactedAudioUrl(jobId.value) : '')

function seek(seconds: number) {
  audioPlayer.value?.seekTo(seconds)
}
</script>

<style scoped>
.result-section {
  display: grid;
  gap: 16px;
  padding-top: 72px;
}

.result-section__head {
  max-width: 780px;
  margin-bottom: 16px;
}

.result-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 370px;
  gap: 14px;
  align-items: stretch;
}

.bottom-cta {
  display: flex;
  justify-content: center;
  padding: 24px 0 18px;
}

@media (max-width: 980px) {
  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>
