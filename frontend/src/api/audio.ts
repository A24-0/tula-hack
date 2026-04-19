export function redactedAudioUrl(id: string) {
  return `/audio/${id}/redacted`
}

export function originalAudioUrl(id: string) {
  return `/audio/${id}/original`
}

export async function originalAudioExists(id: string) {
  try {
    const response = await fetch(originalAudioUrl(id), { method: 'HEAD' })
    return response.ok
  } catch {
    return false
  }
}
