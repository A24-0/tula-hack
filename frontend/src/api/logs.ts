import { apiClient } from './client'
import type { ProcessingLog } from '@/types/log'
import { adaptLogs } from '@/utils/adapters'

export async function getLogs(id: string): Promise<ProcessingLog[]> {
  const response = await apiClient.get(`/logs/${id}`)
  return adaptLogs(response.data)
}
