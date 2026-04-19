import { apiClient } from './client'
import type { TranscriptResult } from '@/types/transcript'
import { adaptTranscript } from '@/utils/adapters'

export async function getTranscript(id: string): Promise<TranscriptResult> {
  const response = await apiClient.get(`/transcript/${id}`)
  return adaptTranscript(response.data, id)
}
