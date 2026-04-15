# OpsDash Frontend

Мини-дэшборд для мониторинга сервера на Vue 3.  
Проект отображает метрики, статусы сервисов и алерты в стиле dark glass UI.

## Что реализовано

- Карточка метрик CPU/RAM с прогресс-барами
- Карточка статусов сервисов (`up`, `down`, `degraded`) с цветными бейджами
- Карточка алертов с подсветкой critical-состояний
- Mock API-данные из `public/mock/*.json`
- Загрузка данных через `Axios`
- Автообновление данных каждые 10 секунд
- Адаптив и dev-переключатель режимов: `Auto / Phone / Tablet / Desktop`

## Стек

- Vue 3 + Vite
- TypeScript
- Axios
- Pinia (установлен, можно использовать на следующих этапах)
- Vue Router (установлен, пока без маршрутов)

## Структура проекта

```text
src/
  components/
    AlertsCard.vue
    MetricsCard.vue
    ServicesCard.vue
  services/
    dashboardApi.ts
  types/
    dashboard.ts
  App.vue
  main.ts

public/
  mock/
    alerts.json
    metrics.json
    services.json
```

## Запуск

Установка зависимостей:

```bash
npm install
```

Режим разработки:

```bash
npm run dev
```

Проверка типов:

```bash
npm run type-check
```

Production build:

```bash
npm run build
```

## Откуда берутся данные

Сейчас приложение работает на моках:

- `/mock/metrics.json`
- `/mock/services.json`
- `/mock/alerts.json`

Запросы выполняются через `src/services/dashboardApi.ts`.

## Следующие шаги (roadmap)

- Добавить `VITE_API_BASE_URL` для гибкого переключения источника данных
- Подключить Go backend (`/api/metrics`, `/api/services`, `/api/alerts`)
- Подготовить Dockerfile и docker-compose
- Интегрировать Prometheus/node_exporter
