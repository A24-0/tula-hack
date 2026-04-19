<template>
  <section id="action" class="section container action-section">
    <div class="action-shell glass-panel">
      <div class="panel-content">
        <Transition name="section" mode="out-in">
          <div v-if="isCompact" class="summary">
            <div>
              <span class="eyebrow">Текущая обработка</span>
              <h2>{{ source?.name || 'Новая запись' }}</h2>
              <p>{{ summaryStatus }}<template v-if="entityCount"> · найдено сущностей: {{ entityCount }}</template></p>
            </div>
            <button class="button button-secondary" type="button" @click="processing.reset">Новая обработка</button>
          </div>

          <div v-else class="action-panel__full">
            <div class="action-panel__head">
              <span class="eyebrow">Start here</span>
              <h2>Передайте аудио в безопасный пайплайн</h2>
              <p>
                Запишите фрагмент с микрофона или загрузите готовый файл. После отправки сервис
                создаст транскрипт, отчёт по ПД и анонимизированную аудиоверсию.
              </p>
              <div class="format-row">
                <span>MP3</span>
                <span>WAV</span>
                <span>M4A</span>
                <span>OGG</span>
                <span>WEBM</span>
              </div>
            </div>

            <div class="action-workspace">
              <div class="action-grid">
                <AudioRecorder @recorded="onRecorded" />
                <UploadDropzone @select="onFile" />
              </div>

              <div v-if="source" class="selected-source">
                <div>
                  <span>{{ source.kind === 'recording' ? 'Новая запись' : 'Выбранный файл' }}</span>
                  <strong>{{ source.name }}</strong>
                </div>
                <button class="button button-primary" type="button" :disabled="isBusy" @click="processing.start">
                  Отправить на обработку
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import AudioRecorder from '@/components/recorder/AudioRecorder.vue'
import UploadDropzone from '@/components/upload/UploadDropzone.vue'
import { useProcessingStore } from '@/stores/processing'
import { useTranscriptStore } from '@/stores/transcript'

const processing = useProcessingStore()
const transcriptStore = useTranscriptStore()
const { source, status, isBusy, summaryStatus } = storeToRefs(processing)
const { entityCount } = storeToRefs(transcriptStore)
const isCompact = computed(() => status.value === 'processing' || status.value === 'pending' || status.value === 'done')

function onFile(file: File) {
  processing.setSource({ file, name: file.name, size: file.size, mime: file.type, kind: 'file' })
}

function onRecorded(file: File, url: string) {
  processing.setSource({ file, blobUrl: url, name: file.name, size: file.size, mime: file.type, kind: 'recording' })
}
</script>

<style scoped>
.action-section {
  padding-top: 76px;
}

.action-shell {
  padding: clamp(22px, 3vw, 34px);
}

.action-panel__full {
  display: grid;
  grid-template-columns: minmax(260px, 0.32fr) minmax(0, 1fr);
  gap: clamp(28px, 4vw, 56px);
  align-items: start;
}

.action-panel__head {
  position: sticky;
  top: 92px;
}

.action-panel__head h2 {
  max-width: 470px;
  margin: 16px 0 0;
  font-size: clamp(30px, 4vw, 54px);
  font-weight: 650;
  line-height: 1.02;
  letter-spacing: -0.055em;
}

.action-panel__head p {
  margin: 18px 0 0;
  color: var(--text-muted);
  font-size: 15px;
  line-height: 1.72;
}

.format-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 26px;
}

.format-row span {
  padding: 7px 10px;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-soft);
  border-radius: 999px;
}

.action-workspace {
  min-width: 0;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.selected-source {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-top: 14px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
}

.selected-source span,
.summary p {
  color: var(--text-muted);
  font-size: 13px;
}

.selected-source strong {
  display: block;
  margin-top: 3px;
  overflow-wrap: anywhere;
  font-size: 15px;
}

.summary {
  display: flex;
  min-height: 124px;
  align-items: center;
  justify-content: space-between;
  gap: 22px;
}

.summary h2 {
  max-width: 760px;
  margin: 14px 0 0;
  overflow-wrap: anywhere;
  font-size: clamp(24px, 3.6vw, 44px);
  font-weight: 650;
  line-height: 1.05;
  letter-spacing: -0.045em;
}

.summary p {
  margin: 10px 0 0;
}

@media (max-width: 980px) {
  .action-panel__full {
    grid-template-columns: 1fr;
  }

  .action-panel__head {
    position: static;
  }
}

@media (max-width: 720px) {
  .action-grid {
    grid-template-columns: 1fr;
  }

  .selected-source,
  .summary {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
