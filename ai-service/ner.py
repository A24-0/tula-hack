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
    "passport": re.compile(
        r"\b\d{4}[\s\-]\d{6}\b"
        r"|\b\d{2}[\s\-]\d{2}[\s\-]\d{2}[\s\-]\d{2}[\s\-]\d{2}\b"
    ),
    "address": re.compile(
        r"(?:ул(?:ица)?\.?\s+[А-ЯЁа-яё\w]+|пр(?:оспект)?\.?\s+[А-ЯЁа-яё\w]+|"
        r"д(?:ом)?\.?\s*\d+|кв(?:артира)?\.?\s*\d+)",
        re.IGNORECASE,
    ),
}

LABELS = {
    "phone": "[ТЕЛЕФОН]",
    "email": "[EMAIL]",
    "inn": "[ИНН]",
    "snils": "[СНИЛС]",
    "passport": "[ПАСПОРТ]",
    "address": "[АДРЕС]",
}


def find_entities(text: str) -> list[Entity]:
    entities = []
    for entity_type, pattern in PATTERNS.items():
        for match in pattern.finditer(text):
            entities.append(Entity(
                type="email" if entity_type == "email_spoken" else entity_type,
                value=match.group(),
                start_char=match.start(),
                end_char=match.end(),
            ))
    entities.sort(key=lambda e: e.start_char)
    return entities


def redact_text(text: str, entities: list[Entity]) -> str:
    result = text
    for entity in sorted(entities, key=lambda e: e.start_char, reverse=True):
        label = LABELS.get(entity.type, "[ДАННЫЕ]")
        result = result[:entity.start_char] + label + result[entity.end_char:]
    return result