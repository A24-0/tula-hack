<template>
  <Transition name="section">
    <section v-if="visible" id="pipeline" class="section container">
      <div :class="['processing glass-panel', { 'processing--error': status === 'error' }]">
        <div class="panel-content">
          <div class="processing__header">
            <div>
              <span class="eyebrow">{{ status === 'error' ? 'Action needed' : 'Pipeline' }}</span>
              <h2>{{ title }}</h2>
              <p>{{ copy }}</p>
            </div>
            <div class="processing__percent">{{ Math.round(progress) }}%</div>
          </div>

          <div v-if="status !== 'error'" class="progress">
            <div class="progress__fill" :style="{ width: `${progress}%` }"></div>
          </div>

          <div v-if="status !== 'error'" class="steps">
            <div v-for="item in steps" :key="item.key" :class="['step', stepState(item.key)]">
              <span>{{ item.index }}</span>
              <strong>{{ item.label }}</strong>
            </div>
          </div>

          <div v-else class="error-card">
            <strong>{{ error || 'Обработка не завершилась' }}</strong>
            <p>Состояние сохранено: можно повторить отправку или начать новую обработку с другим файлом.</p>
            <div>
              <button class="button button-primary" type="button" @click="processing.start">Повторить</button>
              <button class="button button-secondary" type="button" @click="processing.reset">Начать заново</button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { stageLabels, useProcessingStore } from '@/stores/processing'

const processing = useProcessingStore()
const { status, stage, progress, error } = storeToRefs(processing)
const visible = computed(() => status.value !== 'idle')
const title = computed(() => status.value === 'done' ? 'Обработка завершена' : status.value === 'error' ? 'Не удалось обработать аудио' : stageLabels[stage.value] ?? 'Идёт обработка')
const copy = computed(() => {
  if (status.value === 'done') return 'Материалы готовы к просмотру, проверке и скачиванию.'
  if (status.value === 'error') return 'Интерфейс сохранил текущий контекст, чтобы быстро повторить действие.'
  return 'Аудио проходит этапы распознавания, поиска ПД и маскирования чувствительных фрагментов.'
})

const steps = [
  { key: 'queued', label: 'Очередь', index: '01' },
  { key: 'transcribing', label: 'Речь', index: '02' },
  { key: 'detecting', label: 'ПД', index: '03' },
  { key: 'redacting_text', label: 'Текст', index: '04' },
  { key: 'redacting_audio', label: 'Аудио', index: '05' },
  { key: 'packaging', label: 'Файлы', index: '06' },
]

function stepState(key: string) {
  const current = steps.findIndex((item) => item.key === stage.value)
  const index = steps.findIndex((item) => item.key === key)
  if (status.value === 'done' || index < current) return 'step--done'
  if (index === current) return 'step--active'
  return ''
}
</script>

<style scoped>
.processing {
  padding: clamp(22px, 3vw, 34px);
}

.processing__header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 24px;
  align-items: start;
}

.processing h2 {
  max-width: 760px;
  margin: 14px 0 0;
  font-size: clamp(30px, 4.8vw, 58px);
  font-weight: 650;
  line-height: 1.02;
  letter-spacing: -0.055em;
}

.processing p {
  max-width: 650px;
  margin: 16px 0 0;
  color: var(--text-muted);
  line-height: 1.68;
}

.processing__percent {
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: clamp(36px, 6vw, 66px);
  font-weight: 700;
  letter-spacing: -0.08em;
}

.progress {
  height: 8px;
  margin-top: 30px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-soft);
  border-radius: 999px;
}

.progress__fill {
  height: 100%;
  min-width: 16px;
  background: linear-gradient(90deg, var(--accent-blue), #f4f4f2);
  border-radius: inherit;
  transition: width 700ms var(--ease-soft);
}

.steps {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.step {
  min-height: 88px;
  padding: 14px;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
  transition: border-color var(--motion-base), background var(--motion-base), color var(--motion-base);
}

.step span {
  display: block;
  margin-bottom: 18px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-variant-numeric: tabular-nums;
}

.step strong {
  font-size: 14px;
}

.step--active {
  color: var(--text-primary);
  background: rgba(143, 155, 255, 0.11);
  border-color: rgba(143, 155, 255, 0.38);
  box-shadow: 0 0 46px rgba(143, 155, 255, 0.08) inset;
}

.step--done {
  color: var(--accent-green);
  border-color: rgba(150, 216, 173, 0.24);
}

.error-card {
  display: grid;
  gap: 14px;
  margin-top: 26px;
  padding: 20px;
  background: rgba(255, 139, 139, 0.09);
  border: 1px solid rgba(255, 139, 139, 0.3);
  border-radius: 16px;
}

.error-card div {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 900px) {
  .steps {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 620px) {
  .processing__header {
    grid-template-columns: 1fr;
  }
}
</style>
