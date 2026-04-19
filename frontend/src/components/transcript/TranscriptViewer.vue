<template>
  <section class="result-card glass-panel">
    <div class="panel-content">
      <div class="block-head">
        <div>
          <span class="eyebrow">Transcript</span>
          <h3>Текст с синхронизацией</h3>
        </div>
        <div class="toggle" role="group" aria-label="Режим транскрипта">
          <button :class="{ active: mode === 'original' }" type="button" @click="mode = 'original'">С ПД</button>
          <button :class="{ active: mode === 'redacted' }" type="button" @click="mode = 'redacted'">Безопасный</button>
        </div>
      </div>

      <div class="transcript">
        <template v-if="mode === 'original'">
          <button
            v-for="word in data.words"
            :key="word.index"
            :class="['word', word.entity_type && entityCss[word.entity_type], { 'word--entity': word.entity_type, 'word--active': isCurrent(word) }]"
            type="button"
            @click="$emit('seek', word.start)"
          >
            {{ word.word }}
          </button>
        </template>
        <template v-else>
          <span v-for="(part, index) in redactedParts" :key="`${part.text}-${index}`" :class="['redacted-part', part.type && entityCss[part.type]]">
            {{ part.text }}
          </span>
        </template>
      </div>

      <div v-if="presentTypes.length" class="legend">
        <span v-for="type in presentTypes" :key="type" :class="['legend__pill', entityCss[type]]">{{ entityLabels[type] }}</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { EntityType } from '@/types/entity'
import type { TranscriptResult, Word } from '@/types/transcript'
import { entityCss, entityLabels } from '@/utils/colors'

const props = defineProps<{
  data: TranscriptResult
  currentTime: number
}>()

defineEmits<{ seek: [seconds: number] }>()
const mode = ref<'original' | 'redacted'>('original')
const presentTypes = computed(() => [...new Set(props.data.entities.map((entity) => entity.type))])

const redactedParts = computed(() => {
  const matches = [...props.data.redacted_text.matchAll(/\[(PERSON|PHONE|EMAIL|ADDRESS|INN|SNILS|PASSPORT)\]/g)]
  if (!matches.length) return [{ text: props.data.redacted_text }]
  const parts: Array<{ text: string; type?: EntityType }> = []
  let cursor = 0
  matches.forEach((match) => {
    if (match.index > cursor) parts.push({ text: props.data.redacted_text.slice(cursor, match.index) })
    parts.push({ text: match[0], type: match[1] as EntityType })
    cursor = match.index + match[0].length
  })
  if (cursor < props.data.redacted_text.length) parts.push({ text: props.data.redacted_text.slice(cursor) })
  return parts
})

function isCurrent(word: Word) {
  return props.currentTime >= word.start && props.currentTime <= word.end
}
</script>

<style scoped>
.result-card {
  padding: 24px;
}

.block-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 20px;
}

h3 {
  margin: 12px 0 0;
  font-size: clamp(27px, 3.6vw, 44px);
  font-weight: 650;
  line-height: 1.02;
  letter-spacing: -0.05em;
}

.toggle {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.035);
  border: 1px solid var(--border-soft);
  border-radius: 999px;
}

.toggle button {
  min-height: 36px;
  padding: 0 14px;
  color: var(--text-muted);
  background: transparent;
  border: 0;
  border-radius: 999px;
  cursor: pointer;
  transition: background var(--motion-base), color var(--motion-base);
}

.toggle button.active {
  color: var(--ink);
  background: var(--paper);
}

.transcript {
  min-height: 220px;
  padding: clamp(18px, 3vw, 28px);
  color: var(--text-secondary);
  font-size: 18px;
  line-height: 2.05;
  background: rgba(255, 255, 255, 0.025);
  border: 1px solid var(--border-soft);
  border-radius: 16px;
}

.word {
  display: inline;
  margin: 0 3px 0 0;
  padding: 2px 5px;
  color: inherit;
  background: transparent;
  border: 0;
  border-radius: 7px;
  cursor: pointer;
  transition: background var(--motion-fast), color var(--motion-fast), box-shadow var(--motion-fast);
}

.word:hover,
.word--active {
  color: var(--text-primary);
  background: rgba(143, 155, 255, 0.12);
  box-shadow: 0 0 0 1px rgba(143, 155, 255, 0.24);
}

.word--entity {
  color: #050607;
  background: color-mix(in srgb, var(--entity-color) 76%, #f6f6f3);
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--entity-color) 48%, transparent);
}

.word--entity.word--active {
  box-shadow: 0 0 0 2px rgba(246, 246, 243, 0.72);
}

.redacted-part {
  white-space: pre-wrap;
}

.redacted-part[class*='entity-'] {
  display: inline-flex;
  align-items: center;
  margin: 0 4px;
  padding: 2px 10px;
  color: #050607;
  font-weight: 760;
  background: color-mix(in srgb, var(--entity-color) 78%, #f6f6f3);
  border-radius: 999px;
}

.legend {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.legend__pill {
  padding: 6px 10px;
  color: #050607;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 760;
  background: var(--entity-color);
  border-radius: 999px;
}

@media (max-width: 760px) {
  .block-head {
    flex-direction: column;
  }

  .toggle {
    width: 100%;
  }

  .toggle button {
    flex: 1;
    padding: 0 8px;
    font-size: 13px;
  }

  .transcript {
    font-size: 16px;
  }
}
</style>
