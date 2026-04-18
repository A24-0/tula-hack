<template>
  <div
    class="dropzone"
    :class="{ over: isDragging, filled: !!file }"
    @dragenter.prevent="isDragging = true"
    @dragleave.prevent="isDragging = false"
    @dragover.prevent
    @drop.prevent="onDrop"
    @click="inputEl?.click()"
  >
    <input
      ref="inputEl"
      type="file"
      accept=".mp3,.wav,.ogg,.m4a,.flac,audio/*"
      style="display:none"
      @change="onFileChange"
    />

    <div v-if="!file" class="dropzone__idle">
      <p class="dropzone__label">Перетащите файл или нажмите для выбора</p>
      <p class="dropzone__formats">mp3 · wav · ogg · m4a · flac</p>
    </div>

    <div v-else class="dropzone__filled">
      <span class="dropzone__name">{{ file.name }}</span>
      <span class="dropzone__size">{{ formatSize(file.size) }}</span>
      <button class="dropzone__clear" @click.stop="clear">Удалить</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{ change: [file: File] }>()

const inputEl = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const file = ref<File | null>(null)

function onDrop(e: DragEvent) {
  isDragging.value = false
  const f = e.dataTransfer?.files[0]
  if (f) setFile(f)
}

function onFileChange(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0]
  if (f) setFile(f)
}

function setFile(f: File) {
  file.value = f
  emit('change', f)
}

function clear() {
  file.value = null
  if (inputEl.value) inputEl.value.value = ''
}

function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(0) + ' КБ'
  return (bytes / (1024 * 1024)).toFixed(1) + ' МБ'
}
</script>

<style scoped>
.dropzone {
  border: 1px dashed #334155;
  border-radius: 8px;
  padding: 40px 24px;
  text-align: center;
  cursor: pointer;
  background: #162032;
  transition: border-color 0.2s, background 0.2s;
}

.dropzone.over {
  border-color: #3b82f6;
  background: #172a45;
}

.dropzone.filled {
  border-style: solid;
  border-color: #2563eb;
  padding: 16px 24px;
}

.dropzone__label {
  font-size: 14px;
  color: #94a3b8;
  margin-bottom: 6px;
}

.dropzone__formats {
  font-size: 12px;
  color: #475569;
}

.dropzone__filled {
  display: flex;
  align-items: center;
  gap: 12px;
  text-align: left;
}

.dropzone__name {
  font-size: 14px;
  color: #e2e8f0;
  flex: 1;
  word-break: break-all;
}

.dropzone__size {
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
}

.dropzone__clear {
  background: none;
  border: none;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}

.dropzone__clear:hover {
  color: #94a3b8;
  background: #1e293b;
}
</style>
