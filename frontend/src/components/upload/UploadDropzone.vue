<template>
  <div
    :class="['dropzone', { 'dropzone--over': isOver, 'dropzone--filled': Boolean(fileName) }]"
    @dragenter.prevent="isOver = true"
    @dragover.prevent
    @dragleave.prevent="isOver = false"
    @drop.prevent="onDrop"
    @click="input?.click()"
  >
    <input ref="input" class="dropzone__input" type="file" :accept="accept" @change="onInput" />

    <template v-if="fileName">
      <div class="dropzone__file">
        <span class="dropzone__icon">
          <svg viewBox="0 0 24 24"><path d="M7 3h7l4 4v14H7z"/><path d="M14 3v5h5"/></svg>
        </span>
        <div>
          <strong>{{ fileName }}</strong>
          <span>{{ fileSize }}</span>
        </div>
      </div>
    </template>

    <template v-else>
      <span class="dropzone__kicker">upload</span>
      <span class="dropzone__icon">
        <svg viewBox="0 0 24 24"><path d="M12 16V4"/><path d="M7 9l5-5 5 5"/><path d="M5 20h14"/></svg>
      </span>
      <strong>Загрузить аудио</strong>
      <p>Перетащите запись сюда или выберите файл вручную.</p>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAppStore } from '@/stores/app'

const emit = defineEmits<{ select: [file: File] }>()
const app = useAppStore()
const input = ref<HTMLInputElement | null>(null)
const isOver = ref(false)
const selected = ref<File | null>(null)
const accept = '.mp3,.wav,.m4a,.ogg,.webm,audio/mpeg,audio/wav,audio/x-wav,audio/mp4,audio/ogg,audio/webm'
const allowedExtensions = ['mp3', 'wav', 'm4a', 'ogg', 'webm']

const fileName = computed(() => selected.value?.name)
const fileSize = computed(() => selected.value ? formatSize(selected.value.size) : '')

function onDrop(event: DragEvent) {
  isOver.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) choose(file)
}

function onInput(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) choose(file)
}

function choose(file: File) {
  const extension = file.name.split('.').pop()?.toLowerCase()
  if (!extension || !allowedExtensions.includes(extension)) {
    app.pushToast({ tone: 'error', title: 'Формат не поддерживается', text: 'Используйте MP3, WAV, M4A, OGG или WEBM.' })
    return
  }
  selected.value = file
  emit('select', file)
}

function formatSize(bytes: number) {
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} КБ`
  return `${(bytes / 1024 / 1024).toFixed(1)} МБ`
}
</script>

<style scoped>
.dropzone {
  position: relative;
  display: grid;
  min-height: 310px;
  align-content: center;
  justify-items: start;
  padding: 26px;
  overflow: hidden;
  cursor: pointer;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.055), rgba(255, 255, 255, 0.022)),
    rgba(255, 255, 255, 0.018);
  border: 1px dashed rgba(255, 255, 255, 0.16);
  border-radius: 18px;
  transition: border-color var(--motion-base), transform var(--motion-base), background var(--motion-base);
}

.dropzone::after {
  position: absolute;
  right: -70px;
  bottom: -80px;
  width: 190px;
  height: 190px;
  pointer-events: none;
  content: "";
  background: radial-gradient(circle, rgba(143, 155, 255, 0.18), transparent 68%);
}

.dropzone:hover,
.dropzone--over {
  transform: translateY(-3px);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.026)),
    rgba(255, 255, 255, 0.025);
  border-color: rgba(154, 164, 255, 0.5);
}

.dropzone__input {
  display: none;
}

.dropzone__kicker {
  margin-bottom: 22px;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  text-transform: uppercase;
}

.dropzone__icon {
  display: grid;
  width: 58px;
  height: 58px;
  place-items: center;
  margin-bottom: 24px;
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.055);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
}

svg {
  width: 25px;
  height: 25px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

strong {
  display: block;
  max-width: 360px;
  font-size: 27px;
  font-weight: 650;
  letter-spacing: -0.04em;
}

p,
.dropzone__file span {
  max-width: 330px;
  margin: 10px 0 0;
  color: var(--text-muted);
  line-height: 1.6;
}

.dropzone__file {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  max-width: 100%;
  gap: 16px;
  text-align: left;
}

.dropzone__file .dropzone__icon {
  flex: 0 0 auto;
  margin: 0;
}

.dropzone__file strong {
  overflow-wrap: anywhere;
}
</style>
