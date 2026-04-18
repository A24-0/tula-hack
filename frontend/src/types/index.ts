export type JobStatus = 'pending' | 'done' | 'error'

export interface PIIEvent {
  type: 'passport' | 'inn' | 'snils' | 'phone' | 'email' | 'address'
  original: string
  start_sec: number
  end_sec: number
}

export interface TranscriptResult {
  transcript: string
  redacted_transcript: string
  redacted_audio_path: string | null
  pii_events: PIIEvent[]
}

export interface Job {
  id: string
  status: JobStatus
  error?: string
  created_at: string
  updated_at: string
}
