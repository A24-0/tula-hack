import type { EntityType } from './entity'

export interface ProcessingLog {
  timestamp: string
  type: string
  message: string
  entity_type?: EntityType
  start_sec?: number
  end_sec?: number
}
