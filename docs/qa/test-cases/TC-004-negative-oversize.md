# TC-004: Негативный — файл > 500 МБ

**ID:** TC-004
**Модуль:** Backend /upload (MaxBytesReader)
**Приоритет:** Medium
**Вид:** Негативное / нефункциональное (лимит) / ручное

## Предусловия
- Сгенерировать файл 600 МБ: `fsutil file createnew samples\big.wav 629145600` (Windows).

## Шаги
1. Выполнить: `curl -i -F "audio=@samples/big.wav;type=audio/wav" http://localhost:8080/upload`

## Ожидаемый результат
- HTTP 400, тело содержит `"failed to parse form"` или `"http: request body too large"`.
- Контейнер backend не падает, `/jobs/<uuid>` другого job'а продолжает отвечать.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:** HTTP 400, `{"error":"failed to parse form: http: request body too large"}`; backend healthcheck остался OK; запрос к `/jobs/<предыдущий>` — 200.
**Связанные баги:** нет
