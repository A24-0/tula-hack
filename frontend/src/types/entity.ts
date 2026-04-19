export type EntityType = 'PERSON' | 'PHONE' | 'EMAIL' | 'ADDRESS' | 'INN' | 'SNILS' | 'PASSPORT'

export interface Entity {
  type: EntityType
  text: string
  start_char?: number
  end_char?: number
  start_sec?: number
  end_sec?: number
  confidence?: number
}

export const entityTypes: EntityType[] = ['PERSON', 'PHONE', 'EMAIL', 'ADDRESS', 'INN', 'SNILS', 'PASSPORT']
