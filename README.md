# OpsDash

Мини-дэшборд для мониторинга сервера (Vue 3 + Go) в стиле dark glass UI.  
Сейчас проект поднимается одной командой через Docker Compose и показывает **реальные CPU/RAM метрики** с хоста через **Prometheus + node_exporter**.

## Что внутри

- **UI (Vue)**: карточки метрик/сервисов/алертов, автообновление каждые 10 секунд
- **API (Go)**: `/api/*` эндпоинты для дэшборда
- **Monitoring**: `prometheus` + `node-exporter` (метрики хоста)
- **Reverse proxy**: nginx (раздаёт SPA и проксирует `/api/*` в backend)

## Быстрый старт (рекомендуется)

Требования: Docker + Docker Compose plugin.

```bash
docker compose up -d --build
```

После запуска:

- **UI**: `http://<server-ip>/`
- **API**: `http://<server-ip>/api/metrics`

Проверка:

```bash
curl -s http://<server-ip>/api/metrics
```

## Эндпоинты API

- `GET /api/metrics` — CPU/RAM (из Prometheus)
- `GET /api/services` — статусы сервисов (пока демо)
- `GET /api/alerts` — алерты (пока демо)
- `GET /healthz` — healthcheck backend

## Локальная разработка (без Docker)

### Frontend

```bash
npm install
npm run dev
```

По умолчанию фронт ходит на `/api`. В dev-режиме это проксируется на `http://localhost:8080` через `vite.config.ts`.

### Backend

```bash
cd backend
go run .
```

По умолчанию backend слушает `:8080` (можно переопределить через `PORT`).

## Реальные метрики: Prometheus + node_exporter

В `docker-compose.yml` добавлены сервисы:

- `node-exporter` — снимает метрики хоста
- `prometheus` — скрапит `node-exporter` и отдаёт API

Backend читает Prometheus через `PROMETHEUS_URL` и строит CPU/RAM через PromQL.

## Структура репозитория

```text
backend/          # Go API
nginx/            # nginx config (SPA + /api proxy)
prometheus/       # Prometheus config
src/              # Vue frontend
docker-compose.yml
Dockerfile
```

## Roadmap

- Подключить `services` и `alerts` к реальным данным (targets/alert rules)
- Добавить Grafana (дашборды + алертинг)
- HTTPS + домен (reverse proxy на хосте или в контейнере)
