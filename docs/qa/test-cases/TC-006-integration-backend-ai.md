# TC-006: Интеграция backend ↔ ai-service

**ID:** TC-006
**Модуль:** backend/internal/service + ai-service/main.py
**Приоритет:** High
**Вид:** Интеграционное / серый ящик / комбинированное (ручное + авто)

## Предусловия
- Общий Docker volume `uploads` смонтирован в `/uploads` у обоих контейнеров.

## Шаги
1. Загрузить валидный файл через UI (получить `job_id`).
2. В момент `status=processing` выполнить:
   - `docker compose exec backend ls /uploads/`
   - `docker compose exec ai-service ls /uploads/`
3. Запустить `cd backend && go test ./...`.

## Ожидаемый результат
- Файл `<job_id>.mp3` виден в обоих контейнерах.
- После завершения обработки `<job_id>_redacted.mp3` также виден у обоих.
- `go test ./...` — PASS для `internal/api` и `internal/service`.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:** Файлы видны в обоих контейнерах; `go test ./...` — ok (handler_test + processor client тесты прошли).
**Связанные баги:** нет
