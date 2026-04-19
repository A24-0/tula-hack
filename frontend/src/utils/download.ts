export function downloadText(filename: string, text: string, type = 'text/plain;charset=utf-8') {
  const blob = new Blob([text], { type })
  downloadBlob(filename, blob)
}

export function downloadJson(filename: string, value: unknown) {
  downloadText(filename, JSON.stringify(value, null, 2), 'application/json;charset=utf-8')
}

export function downloadBlob(filename: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

export async function downloadFromUrl(filename: string, url: string) {
  const response = await fetch(url)
  if (!response.ok) throw new Error('Не удалось скачать файл')
  const blob = await response.blob()
  downloadBlob(filename, blob)
}
