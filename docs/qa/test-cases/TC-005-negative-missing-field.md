# TC-005: Негативный — отсутствует поле `audio`

**ID:** TC-005
**Модуль:** Backend /upload
**Приоритет:** Medium
**Вид:** Негативное / чёрный ящик / ручное

## Шаги
1. Выполнить: `curl -i -F "file=@samples/ok_ru_30s.mp3" http://localhost:8080/upload` (поле названо `file`, а не `audio`).

## Ожидаемый результат
- HTTP 400, JSON `{"error":"field 'audio' is required"}`.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:** HTTP 400, `{"error":"field 'audio' is required"}`.
**Связанные баги:** нет
