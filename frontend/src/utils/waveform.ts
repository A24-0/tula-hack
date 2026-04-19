import type { Entity } from '@/types/entity'

export function redactionIntervals(entities: Entity[]) {
  return entities
    .filter((item) => typeof item.start_sec === 'number' && typeof item.end_sec === 'number')
    .map((item) => ({
      start: item.start_sec ?? 0,
      end: Math.max(item.end_sec ?? 0, (item.start_sec ?? 0) + 0.08),
      type: item.type,
    }))
}
