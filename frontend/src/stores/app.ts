import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export interface Toast {
  id: number
  title: string
  text?: string
  tone: 'info' | 'success' | 'error'
}

export const useAppStore = defineStore('app', () => {
  const toasts = ref<Toast[]>([])
  let nextId = 1

  function pushToast(toast: Omit<Toast, 'id'>) {
    const id = nextId++
    toasts.value.push({ id, ...toast })
    window.setTimeout(() => dismissToast(id), 5200)
  }

  function dismissToast(id: number) {
    toasts.value = toasts.value.filter((toast) => toast.id !== id)
  }

  const hasToasts = computed(() => toasts.value.length > 0)

  return { toasts, hasToasts, pushToast, dismissToast }
})
