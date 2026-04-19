# Voice Redaction Frontend

Frontend MVP сервиса Voice Redaction: запись или загрузка аудио, отправка на обработку, polling статуса, просмотр результата, WaveSurfer-плеер, синхронизированный транскрипт, отчёт, журнал и скачивание материалов.

## Стек

- Vue 3
- TypeScript
- Vite
- Vue Router
- Pinia
- Axios
- WaveSurfer.js

## Установка

```bash
npm install
```

## Запуск в dev

```bash
npm run dev
```

По умолчанию Vite стартует на `http://localhost:3000` и проксирует API на `http://localhost:8080`.

## Запуск через Docker

Только frontend:

```bash
docker compose -f frontend/docker-compose.yml up --build
```

Открыть:

```text
http://localhost:3000
```

Из папки `frontend` можно короче:

```bash
docker compose up --build
```

## Сборка

```bash
npm run build
```

## API

Фронтенд ожидает endpoints:

- `POST /upload` с multipart-полем `file`; дополнительно отправляется поле `audio` для совместимости.
- `GET /jobs/{id}`
- `GET /transcript/{id}`
- `GET /audio/{id}/redacted`
- `GET /logs/{id}`

`VITE_API_BASE_URL` можно задать в `.env`, если API находится не на том же origin.
