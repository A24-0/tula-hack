<template>
  <div class="player glass-panel">
    <div class="panel-content">
      <div class="player__top">
        <div>
          <span class="eyebrow">Redacted audio</span>
          <h3>Анонимизированная версия</h3>
        </div>
        <button class="player__play" type="button" :aria-label="isPlaying ? 'Пауза' : 'Воспроизвести'" @click="toggle">
          <svg v-if="!isPlaying" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
          <svg v-else viewBox="0 0 24 24"><path d="M8 5v14M16 5v14"/></svg>
        </button>
      </div>

      <div ref="waveform" class="waveform"></div>

      <div class="player__meta">
        <span>{{ formatTime(currentTime) }}</span>
        <span>{{ formatTime(duration) }}</span>
      </div>

      <div v-if="audioError" class="player__error">{{ audioError }}</div>
      <div v-if="intervals.length" class="player__legend">
        <span></span>
        Замаскированные интервалы отмечены на waveform
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import WaveSurfer from 'wavesurfer.js'
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js'
import type { Entity } from '@/types/entity'
import { usePlayerStore } from '@/stores/player'
import { formatTime } from '@/utils/time'
import { redactionIntervals } from '@/utils/waveform'

const props = defineProps<{
  src: string
  entities: Entity[]
}>()

const player = usePlayerStore()
const { currentTime, duration, isPlaying } = storeToRefs(player)
const waveform = ref<HTMLElement | null>(null)
const audioError = ref('')
const intervals = computed(() => redactionIntervals(props.entities))
let wavesurfer: WaveSurfer | null = null
let regions: any = null

onMounted(async () => {
  await nextTick()
  createPlayer()
})

onBeforeUnmount(() => {
  wavesurfer?.destroy()
})

function createPlayer() {
  if (!waveform.value) return
  regions = RegionsPlugin.create()
  wavesurfer = WaveSurfer.create({
    container: waveform.value,
    url: props.src,
    height: 156,
    waveColor: 'rgba(246, 246, 243, 0.22)',
    progressColor: '#9aa4ff',
    cursorColor: '#f6f6f3',
    cursorWidth: 2,
    barWidth: 3,
    barGap: 3,
    barRadius: 8,
    dragToSeek: true,
    normalize: true,
    plugins: [regions],
  })

  wavesurfer.on('ready', () => {
    const total = wavesurfer?.getDuration() ?? 0
    player.setDuration(total)
    addRegions(total)
  })
  wavesurfer.on('timeupdate', (time) => player.setTime(time))
  wavesurfer.on('audioprocess', (time) => player.setTime(time))
  wavesurfer.on('play', () => player.setPlaying(true))
  wavesurfer.on('pause', () => player.setPlaying(false))
  wavesurfer.on('finish', () => player.setPlaying(false))
  wavesurfer.on('error', () => {
    audioError.value = 'Не удалось загрузить анонимизированное аудио. Текстовые материалы доступны ниже.'
  })
}

function addRegions(total: number) {
  intervals.value.forEach((interval) => {
    const start = Math.min(interval.start, Math.max(0, total - 0.1))
    const end = Math.min(Math.max(interval.end, start + 0.08), total)
    regions?.addRegion({
      start,
      end,
      color: 'rgba(255, 139, 139, 0.24)',
      drag: false,
      resize: false,
    })
  })
}

function toggle() {
  wavesurfer?.playPause()
}

function seekTo(seconds: number) {
  if (!wavesurfer || !duration.value) return
  wavesurfer.seekTo(Math.min(1, Math.max(0, seconds / duration.value)))
  player.setTime(seconds)
}

defineExpose({ seekTo })
</script>

<style scoped>
.player {
  padding: 24px;
}

.player__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 22px;
}

h3 {
  margin: 12px 0 0;
  font-size: clamp(27px, 4vw, 44px);
  font-weight: 650;
  line-height: 1.02;
  letter-spacing: -0.05em;
}

.player__play {
  display: grid;
  width: 70px;
  height: 70px;
  flex: 0 0 auto;
  place-items: center;
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.055);
  border: 1px solid var(--border-soft);
  border-radius: 18px;
  box-shadow: none;
  cursor: pointer;
  transition: transform var(--motion-base), background var(--motion-base), border-color var(--motion-base);
}

.player__play:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.075);
  border-color: var(--border-strong);
}

svg {
  width: 30px;
  height: 30px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.waveform {
  min-height: 186px;
  padding: 14px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.025);
  border: 1px solid var(--border-soft);
  border-radius: 16px;
}

.player__meta {
  display: flex;
  justify-content: space-between;
  margin-top: 12px;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.player__legend,
.player__error {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-top: 14px;
  color: var(--text-muted);
  font-size: 13px;
}

.player__legend span {
  width: 28px;
  height: 10px;
  background: rgba(255, 139, 139, 0.28);
  border: 1px solid rgba(255, 139, 139, 0.36);
  border-radius: 999px;
}

.player__error {
  color: var(--accent-red);
}
</style>
