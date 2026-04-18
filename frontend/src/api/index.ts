import axios from 'axios'
import type { Job, TranscriptResult } from '@/types'

export async function uploadAudio(
  file: File,
  onProgress: (pct: number) => void,
): Promise<string> {
  const form = new FormData()
  form.append('audio', file)

  const res = await axios.post<{ job_id: string }>('/upload', form, {
    onUploadProgress(e) {
      const pct = e.total ? Math.round((e.loaded / e.total) * 90) : 0
      onProgress(pct)
    },
  })

  return res.data.job_id
}

export async function getJobStatus(id: string): Promise<Job> {
  const res = await axios.get<Job>(`/jobs/${id}`)
  return res.data
}

export async function getTranscript(id: string): Promise<TranscriptResult> {
  const res = await axios.get<TranscriptResult>(`/transcript/${id}`)
  return res.data
}

export function redactedAudioUrl(id: string): string {
  return `/audio/${id}/redacted`
}
