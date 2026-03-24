# go-link-checker

Высокопроизводительный CLI-инструмент для массовой проверки доступности URL-адресов, написанный на Go.

## Возможности

- **Конкурентная обработка** — пул воркеров (по умолчанию 100) для параллельной проверки ссылок
- **Rate limiting** — ограничение 5 запросов/сек на каждый хост (in-memory или через Redis)
- **Дедупликация** — автоматический пропуск повторяющихся URL
- **Повторные попытки** — 3 ретрая с экспоненциальным backoff и jitter
- **Хранение результатов** — вывод в консоль или сохранение в Redis
- **Цветной отчёт** — статистика по статусам (2xx/3xx/4xx/5xx/таймауты) с таблицей результатов
- **Graceful shutdown** — корректное завершение по Ctrl+C

## Архитектура

Проект построен на паттерне **Producer-Consumer**:

```
[Файл с URL] → Producer (дедупликация) → Job Queue → Worker Pool → Consumer → Output (Console/Redis)
```

```
├── cmd/cli/              # Точка входа
├── internal/
│   ├── cli/              # Команды CLI (check, report, clear)
│   ├── config/           # Конфигурация (.env)
│   ├── domain/           # Доменная модель Link
│   ├── handler/          # HTTP-обработчик с ретраями
│   ├── limiter/          # Rate limiter (memory / redis)
│   ├── deduplicator/     # Дедупликатор (memory / redis)
│   ├── output/           # Вывод результатов (console / redis)
│   ├── report/           # Генерация отчёта
│   ├── service/          # Бизнес-логика
│   └── timer/            # Замер времени выполнения
└── pkg/
    ├── worker/           # Пул воркеров
    ├── producer/         # Продюсер URL
    └── consumer/         # Консьюмер результатов
```

## Быстрый старт

### Требования

- Go 1.25+
- Redis (опционально, для режима `redis`)

### Установка и запуск

```bash
git clone https://github.com/0gl04q/go-link-checker.git
cd go-link-checker
```

**Запуск Redis** (если нужен режим `redis`):

```bash
docker compose up -d
```

**Проверка ссылок** (вывод в консоль):

```bash
go run ./cmd/cli/ check -f links.example.txt
```

**Проверка ссылок** (сохранение в Redis):

```bash
go run ./cmd/cli/ check -f links.example.txt -o redis
```

**Просмотр отчёта**:

```bash
go run ./cmd/cli/ report
```

**Очистка результатов в Redis**:

```bash
go run ./cmd/cli/ clear
```

### Сборка

```bash
make build        # Собрать бинарник в bin/worker
make run ARGS="check -f links.example.txt"
make bench        # Запустить бенчмарки
```

## Команды CLI

| Команда   | Описание                                      |
|-----------|-----------------------------------------------|
| `check`   | Проверить доступность URL из файла             |
| `report`  | Показать отчёт по результатам из Redis         |
| `clear`   | Удалить все результаты проверок из Redis        |

### Флаги команды `check`

| Флаг                    | По умолчанию       | Описание                        |
|-------------------------|---------------------|---------------------------------|
| `-f`, `--file`          | `links.example.txt` | Файл со списком URL             |
| `-w`, `--worker-pool-size` | `100`            | Количество воркеров             |
| `-o`, `--output`        | `console`           | Тип вывода: `console` или `redis` |

## Конфигурация

Переменные окружения (или `.env` файл):

| Переменная       | По умолчанию     | Описание              |
|------------------|------------------|-----------------------|
| `LOG_LEVEL`      | `debug`          | Уровень логирования   |
| `REDIS_ADDR`     | `localhost:6379` | Адрес Redis           |
| `REDIS_PASSWORD` | —                | Пароль Redis          |
| `REDIS_DB`       | `0`              | Номер БД Redis        |

## Режимы работы

- **Console** — вывод в терминал, in-memory rate limiter и дедупликатор. Подходит для быстрой одноразовой проверки.
- **Redis** — сохранение результатов в Redis, распределённый rate limiter и дедупликатор. Подходит для повторного анализа и работы в несколько инстансов.

## Параметры HTTP-клиента

| Параметр              | Значение          |
|-----------------------|-------------------|
| Таймаут запроса       | 10 сек            |
| Общий таймаут         | 30 сек            |
| Количество ретраев    | 3                 |
| Backoff               | 500ms → 1s → 2s  |
| Jitter                | 0–1000ms          |
| Keep-Alive            | включён           |

## Зависимости

- [cobra](https://github.com/spf13/cobra) — CLI-фреймворк
- [go-redis](https://github.com/redis/go-redis) — клиент Redis
- [pterm](https://github.com/pterm/pterm) — форматированный вывод в терминал
- [cleanenv](https://github.com/ilyakaznacheev/cleanenv) — загрузка конфигурации
- [tint](https://github.com/lmittmann/tint) — цветное структурированное логирование
