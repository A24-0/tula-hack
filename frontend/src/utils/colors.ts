import type { EntityType } from '@/types/entity'

export const entityLabels: Record<EntityType, string> = {
  PERSON: 'ФИО',
  PHONE: 'Телефон',
  EMAIL: 'Email',
  ADDRESS: 'Адрес',
  INN: 'ИНН',
  SNILS: 'СНИЛС',
  PASSPORT: 'Паспорт',
}

export const entityCss: Record<EntityType, string> = {
  PERSON: 'entity-person',
  PHONE: 'entity-phone',
  EMAIL: 'entity-email',
  ADDRESS: 'entity-address',
  INN: 'entity-inn',
  SNILS: 'entity-snils',
  PASSPORT: 'entity-passport',
}
