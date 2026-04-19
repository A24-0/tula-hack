import re
from dataclasses import dataclass


@dataclass
class Entity:
    type: str
    value: str
    start_char: int
    end_char: int


PATTERNS = {
    "phone": re.compile(
        r"(?:\+7|8)[\s\-](?:[\s\-]*\d){10}"
    ),
    # spoken: "номер телефона плюс 8 плюс 1 2 3" / "плюс 7 .829"
    "phone_spoken": re.compile(
        r"(?:(?:номер\s+)?телефон[аеуои]?\s+)(?:(?:плюс\s+)?[\d][\d\s.,]*)+",
        re.IGNORECASE | re.UNICODE,
    ),
    "email": re.compile(
        r"[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}"
    ),
    "email_spoken": re.compile(
        r"[\w\s]+собака\s+[\w]+\s*\.[\w]+",
        re.IGNORECASE | re.UNICODE,
    ),
    "inn": re.compile(
        r"\b\d{10}\b|\b\d{12}\b"
        r"|\b(?:\d[\s]?){12}\b"
    ),
    "snils": re.compile(
        r"\b\d{3}[\s\-]\d{3}[\s\-]\d{3}[\s\-]\d{2}\b"
    ),
    # spoken: "снилс 798" / "снил 185" (Whisper часто теряет финальный "с")
    "snils_spoken": re.compile(
        r"снил[сc]?\s+[\d][\d\s.,]+",
        re.IGNORECASE | re.UNICODE,
    ),
    "passport": re.compile(
        r"\b\d{4}[\s\-]\d{6}\b"
        r"|\b\d{2}[\s\-]\d{2}[\s\-]\d{2}[\s\-]\d{2}[\s\-]\d{2}\b"
    ),
    # spoken: "паспортные данные 182 .3457"
    "passport_spoken": re.compile(
        r"паспорт(?:ные\s+данные?)?\s+[\d][\d\s.,]+",
        re.IGNORECASE | re.UNICODE,
    ),
    "address": re.compile(
        r"(?:ул(?:ица)?\.?\s+[А-ЯЁа-яё\w]+|пр(?:оспект)?\.?\s+[А-ЯЁа-яё\w]+|"
        r"д(?:ом)?\.?\s*\d+|кв(?:артира)?\.?\s*\d+)",
        re.IGNORECASE,
    ),
}

LABELS = {
    "phone": "[ТЕЛЕФОН]",
    "phone_spoken": "[ТЕЛЕФОН]",
    "email": "[EMAIL]",
    "inn": "[ИНН]",
    "snils": "[СНИЛС]",
    "snils_spoken": "[СНИЛС]",
    "passport": "[ПАСПОРТ]",
    "passport_spoken": "[ПАСПОРТ]",
    "address": "[АДРЕС]",
}


def find_entities(text: str) -> list[Entity]:
    entities = []
    for entity_type, pattern in PATTERNS.items():
        base_type = entity_type.removesuffix("_spoken")
        for match in pattern.finditer(text):
            entities.append(Entity(
                type=base_type,
                value=match.group(),
                start_char=match.start(),
                end_char=match.end(),
            ))
    # remove entities fully contained within a larger entity of same type
    entities.sort(key=lambda e: (e.start_char, -(e.end_char - e.start_char)))
    deduped: list[Entity] = []
    for e in entities:
        if deduped and deduped[-1].type == e.type and e.start_char >= deduped[-1].start_char and e.end_char <= deduped[-1].end_char:
            continue
        deduped.append(e)
    return deduped


def redact_text(text: str, entities: list[Entity]) -> str:
    result = text
    for entity in sorted(entities, key=lambda e: e.start_char, reverse=True):
        label = LABELS.get(entity.type, "[ДАННЫЕ]")
        result = result[:entity.start_char] + label + result[entity.end_char:]
    return result