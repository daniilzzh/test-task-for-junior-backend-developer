# Task Service с периодическими задачами

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## Добавленная функциональность

### Типы периодичности

1. **daily** - каждый N-й день
2. **monthly** - в определенное число месяца
3. **specific_dates** - на конкретные даты
4. **odd_even** - по четным/нечетным дням

### Примеры запросов

#### Создание ежедневной задачи
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Обзвон пациентов",
    "recurrence_type": "daily",
    "recurrence_rule": {"interval": 1}
  }'
```

#### Создание еженедельной задачи (каждые 7 дней)
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Инвентаризация",
    "recurrence_type": "daily",
    "recurrence_rule": {
      "interval": 7,
      "start_date": "2026-04-21"
    }
  }'
```

#### Создание задачи на конкретные даты
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Праздничное дежурство",
    "recurrence_type": "specific_dates",
    "recurrence_rule": {
      "dates": ["2026-12-25", "2026-12-31", "2027-01-01"]
    }
  }'
```

#### Получение задач на сегодня
```bash
curl "http://localhost:8080/api/v1/tasks/date"
```

#### Получение задач на конкретную дату
```bash
curl "http://localhost:8080/api/v1/tasks/date?date=2026-12-25"
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/date`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`
