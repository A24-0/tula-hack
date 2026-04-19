<template>
  <section class="result-card glass-panel">
    <div class="panel-content">
      <div class="block-head">
        <div>
          <span class="eyebrow">Audit log</span>
          <h3>События системы</h3>
        </div>
      </div>

      <div class="logs">
        <article v-for="(event, index) in events" :key="`${event.type}-${index}`" class="log-line">
          <time>{{ time(event.timestamp) }}</time>
          <span class="log-line__type">{{ event.type }}</span>
          <p>{{ event.message }}</p>
          <span v-if="event.entity_type" :class="['log-line__entity', entityCss[event.entity_type]]">{{ event.entity_type }}</span>
          <span v-if="event.start_sec !== undefined" class="log-line__interval">{{ interval(event.start_sec, event.end_sec) }}</span>
        </article>
        <div v-if="!events.length" class="empty">Журнал пока пуст.</div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { ProcessingLog } from '@/types/log'
import { entityCss } from '@/utils/colors'
import { formatTime } from '@/utils/time'

defineProps<{ events: ProcessingLog[] }>()

function time(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function interval(start?: number, end?: number) {
  if (typeof start !== 'number') return ''
  return `${formatTime(start)} - ${formatTime(end ?? start)}`
}
</script>

<style scoped>
.result-card {
  padding: 24px;
}

.block-head {
  margin-bottom: 20px;
}

h3 {
  margin: 12px 0 0;
  font-size: clamp(27px, 3.6vw, 44px);
  font-weight: 650;
  line-height: 1.02;
  letter-spacing: -0.05em;
}

.logs {
  display: grid;
  gap: 8px;
}

.log-line {
  display: grid;
  grid-template-columns: 92px 160px minmax(0, 1fr) auto auto;
  gap: 12px;
  align-items: center;
  padding: 13px 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
}

time,
.log-line__type,
.log-line__interval {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  font-variant-numeric: tabular-nums;
}

.log-line__type {
  color: var(--accent-blue);
}

p {
  margin: 0;
  color: var(--text-secondary);
}

.log-line__entity {
  padding: 4px 8px;
  color: #050607;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 760;
  background: var(--entity-color);
  border-radius: 999px;
}

.empty {
  padding: 18px;
  color: var(--text-muted);
  text-align: center;
}

@media (max-width: 860px) {
  .log-line {
    grid-template-columns: 1fr;
    gap: 7px;
  }
}
</style>
