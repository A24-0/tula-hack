# TC-003: Негативный — неподдерживаемый формат

**ID:** TC-003
**Модуль:** Backend /upload
**Приоритет:** High
**Вид:** Негативное / чёрный ящик / ручное

## Предусловия
- Сервисы подняты. Файл `samples/not_audio.txt`.

## Шаги
1. Выполнить: `curl -i -F "audio=@samples/not_audio.txt" http://localhost:8080/upload`

## Ожидаемый результат
- HTTP 415, JSON `{"error":"unsupported audio format"}`.
- Файл НЕ создан в `UPLOAD_DIR`.
- Job НЕ создан.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:** HTTP 415, `{"error":"unsupported audio format"}`; `ls uploads/` не показал новых файлов; jobs store не изменился.
**Связанные баги:** нет
