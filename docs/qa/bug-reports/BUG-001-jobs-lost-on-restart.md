# BUG-001: Job'ы теряются после рестарта backend

**ID:** BUG-001
**Заголовок:** После `docker compose restart backend` все завершённые и активные job'ы становятся недоступны (404)

**Приоритет:** High
**Серьёзность:** Major
**Компонент:** backend/internal/service (Store)
**Окружение:** Windows 11, Docker Desktop 4.x, Chrome 124, backend: образ из `docker-compose.yml`

## Описание
Хранилище job'ов реализовано как `map` в оперативной памяти ([backend/internal/service/store.go](../../../backend/internal/service/store.go)). При рестарте контейнера все job'ы теряются; запросы `GET /jobs/<id>` возвращают 404, что ломает клиентский polling и делает невозможным повторное получение уже обработанной транскрипции.

## Шаги воспроизведения
1. `docker compose up -d`
2. Загрузить `samples/ok_ru_30s.mp3`, получить `job_id`.
3. Дождаться `status=done`.
4. `docker compose restart backend`.
5. `curl -i http://localhost:8080/jobs/<job_id>`

## Ожидаемый результат
200 OK, тело — информация о job'е со статусом `done`.

## Фактический результат
404 Not Found, `{"error":"job not found"}`.

## Вложения
- `docs/qa/evidence/TC-007/before_restart.png`
- `docs/qa/evidence/TC-007/after_restart.png`
- `docs/qa/evidence/TC-007/backend.log`

## Дополнительная информация
Частота: 100%. Уже задокументировано в `CLAUDE.md` как известное ограничение, но для промышленной эксплуатации требуется персистентное хранилище (SQLite/Postgres/Redis).
