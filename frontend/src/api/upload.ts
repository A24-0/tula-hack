import { apiClient } from './client'
import { adaptUploadResponse } from '@/utils/adapters'

export async function uploadAudio(file: File | Blob, filename: string, onProgress: (progress: number) => void) {
  const form = new FormData()
  form.append('file', file, filename)
  form.append('audio', file, filename)

  const response = await apiClient.post('/upload', form, {
    onUploadProgress(event) {
      if (!event.total) return
      onProgress(Math.min(92, Math.round((event.loaded / event.total) * 88)))
    },
  })

  return adaptUploadResponse(response.data)
}
