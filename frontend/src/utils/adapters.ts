import type { Entity, EntityType } from '@/types/entity'
import type { Job } from '@/types/job'
import type { ProcessingLog } from '@/types/log'
import type { TranscriptResult, Word } from '@/types/transcript'
import { entityTypes } from '@/types/entity'
import { nowIso } from './time'

const legacyMap: Record<string, EntityType> = {
  person: 'PERSON',
  fio: 'PERSON',
  phone: 'PHONE',
  email: 'EMAIL',
  address: 'ADDRESS',
  inn: 'INN',
  snils: 'SNILS',
  passport: 'PASSPORT',
}

export function normalizeEntityType(value: unknown): EntityType {
  const raw = String(value ?? '').trim()
  const upper = raw.toUpperCase()
  if (entityTypes.includes(upper as EntityType)) return upper as EntityType
  return legacyMap[raw.toLowerCase()] ?? 'PERSON'
}

export function adaptUploadResponse(data: unknown): string {
  const value = data as { id?: string; job_id?: string }
  const id = value.id ?? value.job_id
  if (!id) throw new Error('Backend не вернул идентификатор задачи')
  return id
}

export function adaptJob(data: unknown): Job {
  const raw = data as Record<string, unknown>
  const status = String(raw.status ?? 'processing')
  return {
    id: String(raw.id ?? raw.job_id ?? ''),
    status: status === 'done' || status === 'error' || status === 'pending' ? status : 'processing',
    progress: typeof raw.progress === 'number' ? raw.progress : undefined,
    stage: typeof raw.stage === 'string' ? raw.stage : undefined,
    error: typeof raw.error === 'string' ? raw.error : null,
    created_at: typeof raw.created_at === 'string' ? raw.created_at : undefined,
    updated_at: typeof raw.updated_at === 'string' ? raw.updated_at : undefined,
  }
}

export function adaptTranscript(data: unknown, id?: string): TranscriptResult {
  const raw = data as Record<string, any>
  const source = raw.result && typeof raw.result === 'object' ? raw.result : raw
  const original = String(source.original_text ?? source.transcript ?? '')
  const redacted = String(source.redacted_text ?? source.redacted_transcript ?? '')
  const entities = adaptEntities(source.entities ?? source.pii_events ?? [])
  const words = Array.isArray(source.words) && source.words.length
    ? source.words.map((word: any, index: number) => ({
        word: String(word.word ?? ''),
        start: Number(word.start ?? word.start_sec ?? 0),
        end: Number(word.end ?? word.end_sec ?? word.start ?? 0),
        index: Number(word.index ?? index),
        entity_type: word.entity_type ? normalizeEntityType(word.entity_type) : undefined,
        redacted: Boolean(word.redacted),
      }))
    : buildWords(original, entities)

  return {
    id: String(source.id ?? raw.id ?? id ?? ''),
    original_text: original,
    redacted_text: redacted || redactText(original, entities),
    words,
    entities,
    stats: source.stats ? adaptStats(source.stats) : buildStats(entities),
  }
}

export function adaptLogs(data: unknown): ProcessingLog[] {
  const raw = data as any
  const events = Array.isArray(raw) ? raw : Array.isArray(raw?.events) ? raw.events : []
  return events.map((event: any) => ({
    timestamp: String(event.timestamp ?? nowIso()),
    type: String(event.type ?? 'EVENT'),
    message: String(event.message ?? 'Событие обработки'),
    entity_type: event.entity_type ? normalizeEntityType(event.entity_type) : undefined,
    start_sec: typeof event.start_sec === 'number' ? event.start_sec : undefined,
    end_sec: typeof event.end_sec === 'number' ? event.end_sec : undefined,
  }))
}

export function fallbackLogs(transcript: TranscriptResult): ProcessingLog[] {
  const logs: ProcessingLog[] = [
    { timestamp: nowIso(), type: 'TRANSCRIBE_DONE', message: 'Транскрипция завершена' },
  ]
  transcript.entities.forEach((entity) => {
    logs.push({
      timestamp: nowIso(),
      type: 'PII_FOUND',
      message: `Найден ${entity.type}`,
      entity_type: entity.type,
      start_sec: entity.start_sec,
      end_sec: entity.end_sec,
    })
  })
  logs.push({ timestamp: nowIso(), type: 'AUDIO_REDACTED', message: 'Аудио анонимизировано' })
  logs.push({ timestamp: nowIso(), type: 'FILES_READY', message: 'Файл готов к скачиванию' })
  return logs
}

function adaptEntities(items: any[]): Entity[] {
  if (!Array.isArray(items)) return []
  return items.map((item) => ({
    type: normalizeEntityType(item.type ?? item.entity_type),
    text: String(item.text ?? item.original ?? ''),
    start_char: numberOrUndefined(item.start_char),
    end_char: numberOrUndefined(item.end_char),
    start_sec: numberOrUndefined(item.start_sec ?? item.start),
    end_sec: numberOrUndefined(item.end_sec ?? item.end),
    confidence: numberOrUndefined(item.confidence),
  }))
}

function adaptStats(raw: Record<string, unknown>) {
  return Object.fromEntries(Object.entries(raw).map(([key, value]) => [normalizeEntityType(key), Number(value) || 0]))
}

function buildStats(entities: Entity[]) {
  return entities.reduce<Record<string, number>>((acc, entity) => {
    acc[entity.type] = (acc[entity.type] ?? 0) + 1
    return acc
  }, {})
}

function buildWords(text: string, entities: Entity[]): Word[] {
  const chunks = text.match(/\S+/g) ?? []
  const duration = Math.max(12, chunks.length * 0.42)
  return chunks.map((word, index) => {
    const start = Number(((index / Math.max(chunks.length, 1)) * duration).toFixed(2))
    const end = Number((start + 0.34).toFixed(2))
    const entity = entities.find((item) => item.text && word.includes(item.text.split(/\s+/)[0]))
    return {
      word,
      index,
      start: entity?.start_sec ?? start,
      end: entity?.end_sec ?? end,
      entity_type: entity?.type,
      redacted: Boolean(entity),
    }
  })
}

function redactText(text: string, entities: Entity[]) {
  return entities.reduce((acc, entity) => {
    if (!entity.text) return acc
    return acc.replaceAll(entity.text, `[${entity.type}]`)
  }, text)
}

function numberOrUndefined(value: unknown) {
  return typeof value === 'number' && Number.isFinite(value) ? value : undefined
}
