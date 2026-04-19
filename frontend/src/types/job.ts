export type JobStatus = 'idle' | 'pending' | 'processing' | 'done' | 'error'

export type JobStage =
  | 'uploading'
  | 'queued'
  | 'transcribing'
  | 'detecting'
  | 'redacting_text'
  | 'redacting_audio'
  | 'packaging'
  | 'done'
  | 'error'

export interface Job {
  id: string
  status: Exclude<JobStatus, 'idle'>
  progress?: number
  stage?: JobStage | string
  error?: string | null
  created_at?: string
  updated_at?: string
}
