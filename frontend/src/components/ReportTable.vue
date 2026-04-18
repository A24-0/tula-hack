<template>
  <div>
    <h2 class="section-title">Удалённые данные</h2>

    <p v-if="!rows.length" class="empty">Персональных данных не найдено</p>

    <table v-else class="table">
      <thead>
        <tr>
          <th>Тип</th>
          <th>Количество</th>
          <th>Примеры</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.type">
          <td>
            <span :class="['type-badge', `pii--${row.type}`]">
              {{ typeLabel[row.type] ?? row.type }}
            </span>
          </td>
          <td class="count">{{ row.count }}</td>
          <td>
            <span v-for="(ex, i) in row.examples" :key="i" class="example">{{ ex }}</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PIIEvent } from '@/types'

const props = defineProps<{ events: PIIEvent[] }>()

const typeLabel: Record<string, string> = {
  passport: 'Паспортные данные',
  inn: 'ИНН',
  snils: 'СНИЛС',
  phone: 'Телефон',
  email: 'Email',
  address: 'Адрес',
}

const rows = computed(() => {
  const map = new Map<string, { count: number; examples: string[] }>()
  for (const ev of props.events) {
    const entry = map.get(ev.type) ?? { count: 0, examples: [] }
    entry.count++
    if (entry.examples.length < 3 && !entry.examples.includes(ev.original)) {
      entry.examples.push(ev.original)
    }
    map.set(ev.type, entry)
  }
  return [...map.entries()].map(([type, d]) => ({ type, ...d }))
})
</script>

<style scoped>
.section-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 14px;
}

.empty {
  color: #475569;
  font-size: 14px;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.table th {
  text-align: left;
  padding: 8px 12px;
  font-size: 11px;
  font-weight: 600;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid #334155;
}

.table td {
  padding: 10px 12px;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
}

.table tr:last-child td {
  border-bottom: none;
}

.type-badge {
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.count {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.example {
  display: inline-block;
  background: #162032;
  border: 1px solid #334155;
  border-radius: 4px;
  padding: 1px 7px;
  font-size: 12px;
  color: #94a3b8;
  margin-right: 5px;
  margin-bottom: 3px;
  font-family: monospace;
}
</style>
