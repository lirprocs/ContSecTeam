# ContSecTeam Queue Service
## Требования
- Go 1.24.0
- Без внешних зависимостей

## Конфигурация (env)
- `WORKERS` — число воркеров, по умолчанию 4
- `QUEUE_SIZE` — размер буферизированной очереди, по умолчанию 64
- `PORT` — порт HTTP сервера, по умолчанию 8080

## HTTP API
- `POST /enqueue`
  - Тело JSON:
    ```json
    {"id":"<string>","payload":"<string>","max_retries":<int>}
    ```
  - Ответы:
    - `202 Accepted` — задача добавлена в очередь
    - `400 Bad Request` — невалидный JSON или отсутствует `id`
    - `503 Service Unavailable` — очередь переполнена
  - Авторизация не требуется

- `GET /healthz` → `200 OK`

## Обработка задач
- Задачи попадают в буферизированный канал-очередь (`QUEUE_SIZE`).
- Пул воркеров (`WORKERS`) читает задачи и обрабатывает их.
- «Обработка» симулируется задержкой 100–500 мс.
- 20% задач «падают» (см. `pkg.ShouldFail`). Выполняются повторы с экспоненциальным бэкоффом и джиттером до `max_retries`.
- Статусы задач сохраняются в `sync.Map` и обновляются: `queued | running | done | failed`.

## Graceful Shutdown
- По `SIGINT/SIGTERM` сервер перестаёт принимать соединения и корректно завершается.
- Очередь закрывается, сервис дожидается завершения текущих задач (см. `srv.Stop()`).

## Файлы
- `cmd/main.go` — запуск сервера, маршруты, graceful shutdown
- `config/config.go` — чтение переменных окружения
- `config/defaults.go` — переменные окружения по умолчанию
- `internal/handler/handler.go` — HTTP-обработчики
- `internal/service/service.go` — бизнес-логика, очередь, ретраи
- `internal/model/model.go` — модель задачи и статусы
- `pkg/utils.go` — случайные задержки, вероятность отказа, бэкофф
- `pkg/worker/worker.go` — переиспользуемый пул воркеров 

## Запуск
### Linux/Mac:
```bash
WORKERS=4 QUEUE_SIZE=64 PORT=8080 go run ./cmd
```
### Windows (Powershell):
```bash
$env:WORKERS=4; $env:QUEUE_SIZE=64; $env:PORT=8080; go run ./cmd
```
### Windows (CMD)
```bash
set WORKERS=4 && set QUEUE_SIZE=64 && set PORT=8080 && go run ./cmd
```

Примеры запросов:
```bash
Invoke-RestMethod -Uri "http://localhost:8080/healthz" -Method GET
```
```bash
curl -X POST http://localhost:8080/enqueue \
  -H 'Content-Type: application/json' \
  -d '{"id":"t1","payload":"data","max_retries":3}'
```

## Тестирование
### Код автоматически тестируется при пуше на GitHub
### Ручное тестирование:
```bash
go test ./...
```

