# TC-011: Whitebox — regex PII и юнит-тесты backend

**ID:** TC-011
**Модуль:** ai-service/ner.py + backend/internal/api/handler_test.go
**Приоритет:** Medium
**Вид:** Whitebox / автоматизированное

## Шаги
1. `cd backend && go test ./... -v`
2. `cd ai-service && python -c "from ner import extract_entities; print(extract_entities('Мой телефон +7 900 123-45-67, ИНН 7707083893'))"`
3. Проверить, что `PHONE` и `INN` извлечены.

## Ожидаемый результат
- `go test` — PASS.
- Python-скрипт возвращает список с entity типами `PHONE` и `INN` с корректными значениями и offset'ами.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:** `go test ./...` → ok; python: извлечены `PHONE` («+7 900 123-45-67») и `INN` («7707083893») с корректными offset'ами.
**Связанные баги:** нет
