# Voice Anonymization Service

Сервис анонимизации голосовых данных: транскрипция → детекция ПДн → редактирование текста и аудио.

**Стек:** Vue 3 (frontend) → Go/Chi (backend) → Python/FastAPI + Whisper + GigaChat (AI service).

## Запуск

### Dev (локально, со сборкой образов):

```bash
cp .env.example .env
docker compose up --build
# UI → http://localhost:3000
```

### Prod (готовые образы из реестра):

```bash
cp .env.example .env
TAG=1.0.0 REGISTRY=ghcr.io/a24-0 \
  docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
# UI → http://<host>/
```

Подробнее об архитектуре и эндпоинтах — в `CLAUDE.md`.
