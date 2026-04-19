import type { Entity, EntityType } from './entity'

export interface Word {
  word: string
  start: number
  end: number
  index: number
  entity_type?: EntityType
  redacted?: boolean
}

export interface TranscriptResult {
  id?: string
  original_text: string
  redacted_text: string
  words: Word[]
  entities: Entity[]
  stats: Partial<Record<EntityType, number>>
}
