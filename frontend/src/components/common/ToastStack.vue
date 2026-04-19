<template>
  <Teleport to="body">
    <TransitionGroup v-if="toasts.length" name="toast" tag="div" class="toast-stack">
      <article v-for="toast in toasts" :key="toast.id" :class="['toast', `toast--${toast.tone}`]">
        <button class="toast__close" type="button" aria-label="Закрыть" @click="app.dismissToast(toast.id)">×</button>
        <strong>{{ toast.title }}</strong>
        <span v-if="toast.text">{{ toast.text }}</span>
      </article>
    </TransitionGroup>
  </Teleport>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/stores/app'

const app = useAppStore()
const { toasts } = storeToRefs(app)
</script>

<style scoped>
.toast-stack {
  position: fixed;
  right: 22px;
  bottom: 22px;
  z-index: 100;
  display: grid;
  gap: 10px;
  width: min(390px, calc(100vw - 32px));
}

.toast {
  position: relative;
  display: grid;
  gap: 5px;
  padding: 16px 44px 16px 16px;
  background: rgba(13, 14, 17, 0.94);
  border: 1px solid var(--border-soft);
  border-left-color: var(--accent-blue);
  border-radius: 14px;
  box-shadow: var(--shadow-raised);
  backdrop-filter: blur(18px);
}

.toast--success { border-left-color: var(--accent-green); }
.toast--error { border-left-color: var(--accent-red); }

.toast strong {
  font-size: 14px;
}

.toast span {
  color: var(--text-muted);
  font-size: 13px;
  line-height: 1.45;
}

.toast__close {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 26px;
  height: 26px;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.045);
  border: 0;
  border-radius: 999px;
  cursor: pointer;
}

.toast-enter-active,
.toast-leave-active {
  transition: opacity var(--motion-base), transform var(--motion-base);
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
