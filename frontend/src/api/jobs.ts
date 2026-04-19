import { apiClient } from './client'
import type { Job } from '@/types/job'
import { adaptJob } from '@/utils/adapters'

export async function getJob(id: string): Promise<Job> {
  const response = await apiClient.get(`/jobs/${id}`)
  return adaptJob(response.data)
}
