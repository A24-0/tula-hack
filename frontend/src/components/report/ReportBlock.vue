<template>
  <section id="security" class="result-card glass-panel">
    <div class="panel-content">
      <div class="block-head">
        <div>
          <span class="eyebrow">Report</span>
          <h3>Найденные данные</h3>
        </div>
      </div>

      <div class="metrics">
        <article v-for="metric in metrics" :key="metric.label" class="metric">
          <span>{{ metric.label }}</span>
          <strong>{{ metric.value }}</strong>
        </article>
      </div>

      <div class="entity-grid">
        <article v-for="row in rows" :key="row.type" :class="['entity-card', entityCss[row.type]]">
          <div>
            <span>{{ entityLabels[row.type] }}</span>
            <strong>{{ row.count }}</strong>
          </div>
          <p>{{ row.examples.join(', ') || 'Примеры скрыты' }}</p>
        </article>
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Тип</th>
              <th>Значение</th>
              <th>Интервал</th>
              <th>Уверенность</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(entity, index) in data.entities" :key="`${entity.type}-${index}`">
              <td><span :class="['type-pill', entityCss[entity.type]]">{{ entity.type }}</span></td>
              <td>{{ entity.text || 'Скрыто' }}</td>
              <td>{{ interval(entity.start_sec, entity.end_sec) }}</td>
              <td>{{ confidence(entity.confidence) }}</td>
            </tr>
            <tr v-if="!data.entities.length">
              <td colspan="4" class="empty">Персональные данные не найдены.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TranscriptResult } from '@/types/transcript'
import { entityCss, entityLabels } from '@/utils/colors'
import { formatDuration, formatTime } from '@/utils/time'

const props = defineProps<{
  data: TranscriptResult
  duration: number
}>()

const totalRedacted = computed(() => props.data.entities.reduce((sum, entity) => {
  if (typeof entity.start_sec !== 'number' || typeof entity.end_sec !== 'number') return sum
  return sum + Math.max(0, entity.end_sec - entity.start_sec)
}, 0))

const metrics = computed(() => [
  { label: 'Всего сущностей', value: String(props.data.entities.length) },
  { label: 'Скрыто фрагментов', value: String(props.data.entities.filter((item) => item.start_sec !== undefined).length) },
  { label: 'Длительность скрытия', value: formatDuration(totalRedacted.value) },
  { label: 'Общая длительность', value: props.duration ? formatTime(props.duration) : 'Не определена' },
  { label: 'Доля скрытого', value: props.duration ? `${Math.min(100, (totalRedacted.value / props.duration) * 100).toFixed(1)}%` : 'Н/Д' },
])

const rows = computed(() => Object.entries(props.data.stats).map(([type, count]) => ({
  type: type as keyof typeof entityLabels,
  count: Number(count) || 0,
  examples: props.data.entities.filter((entity) => entity.type === type).slice(0, 3).map((entity) => entity.text).filter(Boolean),
})).filter((row) => row.count > 0))

function interval(start?: number, end?: number) {
  if (typeof start !== 'number' || typeof end !== 'number') return 'Н/Д'
  return `${formatTime(start)} - ${formatTime(end)}`
}

function confidence(value?: number) {
  if (typeof value !== 'number') return 'Н/Д'
  return `${Math.round(value * 100)}%`
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

.metrics {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
}

.metric {
  min-height: 112px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
}

.metric span {
  color: var(--text-muted);
  font-size: 12px;
}

.metric strong {
  display: block;
  margin-top: 18px;
  font-family: var(--font-mono);
  font-size: 25px;
  letter-spacing: -0.06em;
}

.entity-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-top: 14px;
}

.entity-card {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-soft);
  border-radius: 14px;
}

.entity-card div {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.entity-card span {
  color: var(--text-muted);
  font-size: 13px;
}

.entity-card strong {
  color: var(--entity-color);
  font-family: var(--font-mono);
  font-size: 24px;
}

.entity-card p {
  min-height: 40px;
  margin: 12px 0 0;
  color: var(--text-muted);
  line-height: 1.45;
}

.table-wrap {
  margin-top: 16px;
  overflow-x: auto;
  border: 1px solid var(--border-soft);
  border-radius: 16px;
}

table {
  width: 100%;
  min-width: 720px;
  border-collapse: collapse;
}

th,
td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-soft);
}

th {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

td {
  color: var(--text-secondary);
}

tr:last-child td {
  border-bottom: 0;
}

.type-pill {
  padding: 5px 8px;
  color: #050607;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 760;
  background: var(--entity-color);
  border-radius: 999px;
}

.empty {
  text-align: center;
}

@media (max-width: 980px) {
  .metrics,
  .entity-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 620px) {
  .metrics,
  .entity-grid {
    grid-template-columns: 1fr;
  }
}
</style>
