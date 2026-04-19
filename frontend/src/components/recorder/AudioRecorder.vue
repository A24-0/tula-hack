<template>
  <div :class="['recorder', { 'recorder--active': isRecording, 'recorder--ready': audioUrl }]">
    <div class="recorder__meta">
      <span>record</span>
      <strong>{{ title }}</strong>
      <em>{{ timer }}</em>
    </div>

    <button class="recorder__mic" type="button" :aria-label="isRecording ? 'Остановить запись' : 'Записать аудио'" @click="toggleRecording">
      <span v-if="isRecording" class="recorder__ring"></span>
      <svg v-if="!isRecording" viewBox="0 0 24 24"><path d="M12 14a4 4 0 0 0 4-4V7a4 4 0 0 0-8 0v3a4 4 0 0 0 4 4z"/><path d="M5 10a7 7 0 0 0 14 0"/><path d="M12 17v4"/><path d="M8 21h8"/></svg>
      <svg v-else viewBox="0 0 24 24"><path d="M8 8h8v8H8z"/></svg>
    </button>

    <canvas ref="canvas" class="recorder__wave" width="640" height="96" />
    <audio v-if="audioUrl" class="recorder__preview" :src="audioUrl" controls />
    <button v-if="audioUrl" class="button button-secondary recorder__reset" type="button" @click.stop="resetRecording">Перезаписать</button>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useAppStore } from '@/stores/app'
import { formatTime } from '@/utils/time'

const emit = defineEmits<{ recorded: [file: File, url: string] }>()
const app = useAppStore()
const canvas = ref<HTMLCanvasElement | null>(null)
const isRecording = ref(false)
const seconds = ref(0)
const audioUrl = ref('')
const chunks: BlobPart[] = []
let mediaRecorder: MediaRecorder | null = null
let stream: MediaStream | null = null
let audioContext: AudioContext | null = null
let analyser: AnalyserNode | null = null
let raf = 0
let timerId = 0

const title = computed(() => {
  if (isRecording.value) return 'Идёт запись'
  if (audioUrl.value) return 'Запись готова'
  return 'Записать голос'
})
const timer = computed(() => formatTime(seconds.value))

async function toggleRecording() {
  if (isRecording.value) {
    stopRecording()
    return
  }
  await startRecording()
}

async function startRecording() {
  try {
    resetRecording()
    stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    const mimeType = supportedMime()
    mediaRecorder = new MediaRecorder(stream, mimeType ? { mimeType } : {})
    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size) chunks.push(event.data)
    }
    mediaRecorder.onstop = finishRecording
    mediaRecorder.start()
    isRecording.value = true
    timerId = window.setInterval(() => seconds.value += 1, 1000)
    setupAnalyser(stream)
    drawWave()
  } catch {
    app.pushToast({ tone: 'error', title: 'Не удалось начать запись', text: 'Проверьте доступ к микрофону в браузере.' })
  }
}

function stopRecording() {
  mediaRecorder?.stop()
  isRecording.value = false
  window.clearInterval(timerId)
  stream?.getTracks().forEach((track) => track.stop())
  cancelAnimationFrame(raf)
}

function finishRecording() {
  const mime = mediaRecorder?.mimeType || 'audio/webm'
  const blob = new Blob(chunks, { type: mime })
  const ext = mime.startsWith('audio/mp4') ? 'mp4' : mime.startsWith('audio/ogg') ? 'ogg' : 'webm'
  const name = `voice-redaction-recording-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')}.${ext}`
  const file = new File([blob], name, { type: blob.type })
  audioUrl.value = URL.createObjectURL(blob)
  emit('recorded', file, audioUrl.value)
}

function resetRecording() {
  chunks.splice(0)
  seconds.value = 0
  if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
  audioUrl.value = ''
  clearCanvas()
}

function setupAnalyser(nextStream: MediaStream) {
  audioContext = new AudioContext()
  const source = audioContext.createMediaStreamSource(nextStream)
  analyser = audioContext.createAnalyser()
  analyser.fftSize = 1024
  source.connect(analyser)
}

function drawWave() {
  const el = canvas.value
  if (!el || !analyser) return
  const ctx = el.getContext('2d')
  if (!ctx) return
  const data = new Uint8Array(analyser.frequencyBinCount)
  analyser.getByteTimeDomainData(data)
  ctx.clearRect(0, 0, el.width, el.height)
  ctx.fillStyle = 'rgba(255,255,255,0.018)'
  ctx.fillRect(0, 0, el.width, el.height)
  ctx.strokeStyle = isRecording.value ? '#ff8b8b' : '#f6f6f3'
  ctx.lineWidth = 2
  ctx.beginPath()
  data.forEach((value, index) => {
    const x = (index / data.length) * el.width
    const y = (value / 255) * el.height
    if (index === 0) ctx.moveTo(x, y)
    else ctx.lineTo(x, y)
  })
  ctx.stroke()
  raf = requestAnimationFrame(drawWave)
}

function clearCanvas() {
  const el = canvas.value
  const ctx = el?.getContext('2d')
  if (!el || !ctx) return
  ctx.clearRect(0, 0, el.width, el.height)
}

function supportedMime() {
  const candidates = ['audio/webm;codecs=opus', 'audio/webm', 'audio/ogg;codecs=opus', 'audio/mp4']
  return candidates.find(m => MediaRecorder.isTypeSupported(m)) ?? ''
}

onBeforeUnmount(() => {
  if (isRecording.value) stopRecording()
  if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
  audioContext?.close()
})
</script>

<style scoped>
.recorder {
  position: relative;
  display: grid;
  min-height: 310px;
  grid-template-columns: 1fr auto;
  gap: 20px;
  align-content: space-between;
  padding: 26px;
  overflow: hidden;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.06), rgba(255, 255, 255, 0.022)),
    rgba(255, 255, 255, 0.018);
  border: 1px solid var(--border-soft);
  border-radius: 18px;
  transition: border-color var(--motion-base), transform var(--motion-base), background var(--motion-base);
}

.recorder:hover {
  transform: translateY(-3px);
  border-color: rgba(154, 164, 255, 0.32);
}

.recorder--active {
  border-color: rgba(255, 139, 139, 0.48);
}

.recorder__meta {
  display: grid;
  align-content: start;
}

.recorder__meta span {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  text-transform: uppercase;
}

.recorder__meta strong {
  margin-top: 20px;
  font-size: 27px;
  font-weight: 650;
  letter-spacing: -0.04em;
}

.recorder__meta em {
  margin-top: 8px;
  color: var(--accent-blue);
  font-family: var(--font-mono);
  font-size: 13px;
  font-style: normal;
}

.recorder__mic {
  position: relative;
  display: grid;
  width: 72px;
  height: 72px;
  place-items: center;
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.055);
  border: 1px solid var(--border-soft);
  border-radius: 18px;
  cursor: pointer;
  box-shadow: none;
}

.recorder--active .recorder__mic {
  background: rgba(255, 139, 139, 0.13);
  border-color: rgba(255, 139, 139, 0.44);
}

.recorder__ring {
  position: absolute;
  inset: -8px;
  border: 1px solid rgba(255, 139, 139, 0.74);
  border-radius: 24px;
  animation: pulseRing 1.8s infinite;
}

svg {
  width: 30px;
  height: 30px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.recorder__wave {
  width: 100%;
  height: 96px;
  grid-column: 1 / -1;
  border: 1px solid var(--border-soft);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.025);
}

.recorder__preview {
  width: 100%;
  grid-column: 1 / -1;
  accent-color: var(--accent-blue);
}

.recorder__reset {
  width: max-content;
  grid-column: 1 / -1;
}

@media (max-width: 620px) {
  .recorder {
    grid-template-columns: 1fr;
  }

  .recorder__mic {
    width: 66px;
    height: 66px;
  }
}
</style>
