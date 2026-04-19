export function formatTime(seconds = 0): string {
  if (!Number.isFinite(seconds)) return '0:00'
  const total = Math.max(0, Math.floor(seconds))
  const min = Math.floor(total / 60)
  const sec = total % 60
  return `${min}:${sec.toString().padStart(2, '0')}`
}

export function formatDuration(seconds = 0): string {
  if (!Number.isFinite(seconds) || seconds <= 0) return '0 сек'
  if (seconds < 60) return `${seconds.toFixed(1)} сек`
  const min = Math.floor(seconds / 60)
  const sec = Math.round(seconds % 60)
  return `${min} мин ${sec} сек`
}

export function nowIso(): string {
  return new Date().toISOString()
}
