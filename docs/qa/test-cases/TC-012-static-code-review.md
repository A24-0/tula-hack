# TC-012: Статическое — ревью кода и документации

**ID:** TC-012
**Модуль:** handler.go, ner.py, audio_redact.py, docker-compose.yml, .env.example
**Приоритет:** Medium
**Вид:** Статическое / ручное / белый ящик

## Шаги
1. Прочитать [backend/internal/api/handler.go](../../../backend/internal/api/handler.go): проверить валидацию content-type, обработку ошибок, закрытие файлов.
2. Прочитать [ai-service/ner.py](../../../ai-service/ner.py): полнота regex (PHONE, EMAIL, INN, SNILS, PASSPORT, ADDRESS).
3. Прочитать [ai-service/audio_redact.py](../../../ai-service/audio_redact.py): совпадение токенизации со строкой `" ".join(words)`.
4. Прочитать `docker-compose.yml`: named volumes, healthchecks, порты.
5. Прочитать `.env.example`: отсутствие секретов по умолчанию.

## Ожидаемый результат (чек-лист)
- ✅ Все `file.Close()` защищены `defer`.
- ✅ `MaxBytesReader` применён до `ParseMultipartForm`.
- ✅ Regex покрывают минимум 6 типов PII.
- ✅ Токенизация consistent между `transcript` и `find_silence_ranges`.
- ✅ В `.env.example` нет реальных ключей.
- ⚠️ Найденные замечания → заведены в bug-reports как Low/Medium.

## Результат выполнения
**Дата:** 2026-04-19
**Статус:** Passed
**Фактический результат:**
- `defer file.Close()` / `defer dst.Close()` присутствуют (handler.go:53, :86).
- `http.MaxBytesReader` вызван до `ParseMultipartForm` (handler.go:41-43).
- Regex покрывает PHONE, EMAIL, INN, SNILS, PASSPORT, ADDRESS (+ spoken-email `собака`) — 7 типов.
- Токенизация — `" ".join(w["word"] for w in words)` совпадает на обоих сайтах.
- `.env.example` не содержит реальных секретов.
- Замечание: отсутствует ограничение частоты запросов (rate-limit) на `/upload` — заведён BUG-002 (Low).
**Связанные баги:** BUG-002
