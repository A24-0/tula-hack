import { ref } from 'vue'
import { defineStore } from 'pinia'

export const usePlayerStore = defineStore('player', () => {
  const currentTime = ref(0)
  const duration = ref(0)
  const isPlaying = ref(false)

  function setTime(value: number) {
    currentTime.value = value
  }

  function setDuration(value: number) {
    duration.value = value
  }

  function setPlaying(value: boolean) {
    isPlaying.value = value
  }

  return { currentTime, duration, isPlaying, setTime, setDuration, setPlaying }
})
