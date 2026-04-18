<template>
  <div>
    <div class="tb-header">
      <h2 class="tb-title">Транскрипт</h2>
      <div class="tb-tabs">
        <button :class="['tb-tab', { active: mode === 'original' }]" @click="mode = 'original'">
          С персональными данными
        </button>
        <button :class="['tb-tab', { active: mode === 'redacted' }]" @click="mode = 'redacted'">
          Анонимизированный
        </button>
      </div>
    </div>

    <div class="tb-text" v-html="displayText" />

    <div v-if="mode === 'original' && presentTypes.length" class="tb-legend">
      <span v-for="type in presentTypes" :key="type" :class="['legend-tag', `pii--${type}`]">
        {{ typeLabel[type] }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TranscriptResult } from '@/types'

const props = defineProps<{ data: TranscriptResult }>()
const mode = ref<'original' | 'redacted'>('original')

const typeLabel: Record<string, string> = {
  passport: 'Паспорт',
  inn: 'ИНН',
  snils: 'СНИЛС',
  phone: 'Телефон',
  email: 'Email',
  address: 'Адрес',
}

const presentTypes = computed(() => [...new Set(props.data.pii_events.map(e => e.type))])

const displayText = computed(() => {
  if (mode.value === 'redacted') return escapeHtml(props.data.redacted_transcript)

  let text = props.data.transcript
  const sorted = [...props.data.pii_events].sort((a, b) => b.original.length - a.original.length)
  for (const ev of sorted) {
    text = text.replace(
      ev.original,
      `<mark class="pii pii--${ev.type}" title="${typeLabel[ev.type] ?? ev.type}">${escapeHtml(ev.original)}</mark>`,
    )
  }
  return text
})

function escapeHtml(s: string) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}
</script>

<style scoped>
.tb-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 14px;
}

.tb-title {
  font-size: 15px;
  font-weight: 600;
}

.tb-tabs {
  display: flex;
  gap: 4px;
}

.tb-tab {
  padding: 5px 12px;
  font-size: 13px;
  border: 1px solid #334155;
  border-radius: 5px;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: background 0.15s;
}

.tb-tab.active {
  background: #2563eb;
  color: #fff;
  border-color: #2563eb;
}

.tb-text {
  background: #162032;
  border: 1px solid #334155;
  border-radius: 6px;
  padding: 16px;
  font-size: 14px;
  line-height: 1.9;
  white-space: pre-wrap;
  min-height: 100px;
  color: #cbd5e1;
}

.tb-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
}

.legend-tag {
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}
</style>

<style>
.pii { border-radius: 3px; padding: 1px 4px; cursor: default; }
.pii--phone    { background: #1e3a5f; color: #93c5fd; }
.pii--email    { background: #2e1065; color: #c4b5fd; }
.pii--passport { background: #451a03; color: #fcd34d; }
.pii--inn      { background: #052e16; color: #86efac; }
.pii--snils    { background: #052e16; color: #86efac; }
.pii--address  { background: #450a0a; color: #fca5a5; }
</style>
