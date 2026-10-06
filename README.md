# Лабораторная работа №10

**Студент:** Сурков Всеволод Сергеевич
**Группа:** 220032-11
**Вариант:** №2

**Тема:** Веб-разработка: FastAPI (Python) vs Gin (Go)

По варианту №2 выполняются задания средней сложности №2, №4 и №6.

---

## Содержание

1. [Задание 2 — Middleware для логирования в Go](#задание-2--middleware-для-логирования-в-go)
2. [Задание 4 — FastAPI-сервис, вызывающий Go-сервис по HTTP](#задание-4--fastapi-сервис-вызывающий-go-сервис-по-http)
3. [Задание 6 — Сравнение скорости FastAPI и Gin под нагрузкой](#задание-6--сравнение-скорости-fastapi-и-gin-под-нагрузкой)

---

## Структура проекта

```
.
├── go.mod                       один Go-модуль на все задания
├── go.sum
├── requirements.txt             зависимости Python-сервисов
├── task2/
│   ├── main.go                  Gin + middleware логирования
│   ├── result.txt               фактический вывод
│   └── server.log               лог сервера
├── task4/
│   ├── go-service/main.go       Go-сервис, принимающий JSON
│   ├── fastapi-service/main.py  Python-сервис, вызывающий Go по HTTP
│   └── result.txt               фактический вывод
├── task6/
│   ├── go-service/main.go       Gin для замера
│   ├── fastapi-service/main.py  FastAPI для замера
│   ├── comparison.txt           разбор результатов и выводы
│   ├── result.txt               фактический вывод
│   ├── test_go.txt              полный лог hey для Gin
│   └── test_fastapi.txt         полный лог hey для FastAPI
└── PROMPT_LOG.md                лог промптов
```

---

## Подготовка окружения

```bash
# Go-зависимости
go mod download

# Python-зависимости
pip install -r requirements.txt
```

Проверка, что всё собирается:

```bash
go build ./...
go vet ./...
```

---

## Задание 2 — Middleware для логирования в Go

Собственный middleware `LoggerMiddleware` пишет в лог метод, путь, код ответа
и время обработки каждого запроса. Роутер создаётся через `gin.New()`, чтобы
встроенный логгер Gin не дублировал записи.

### Запуск

```bash
go run ./task2
```

### Проверка эндпоинтов

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/hello/World
```

**Фактический вывод:**

```json
{"message":"pong"}
{"hello":"World"}
```

**Лог сервера:**

```
2026/10/06 17:55:06 Gin-сервис слушает http://localhost:8080
2026/10/06 17:55:08 [GIN] GET /ping -> 200 (0s)
2026/10/06 17:55:08 [GIN] GET /hello/World -> 200 (0s)
```

Полный вывод: [task2/result.txt](task2/result.txt), [task2/server.log](task2/server.log)

---

## Задание 4 — FastAPI-сервис, вызывающий Go-сервис по HTTP

Go-сервис принимает структуру `UserRequest` и возвращает `UserResponse`.
Python-сервис валидирует входные данные через Pydantic и передаёт их в Go
через `httpx.AsyncClient`. Адрес Go-сервиса настраивается переменной
окружения `GO_SERVICE_URL` (по умолчанию `http://localhost:8080`).

### Запуск

Два терминала:

```bash
# Терминал 1 — Go-сервис
go run ./task4/go-service
```

```bash
# Терминал 2 — FastAPI-сервис
cd task4/fastapi-service
uvicorn main:app --host 127.0.0.1 --port 8000
```

### Проверка

```bash
# Проверка доступности Go-сервиса
curl http://127.0.0.1:8000/ping-go
```

```json
{"go_status":"pong"}
```

```bash
# Успешная передача данных
curl -X POST http://127.0.0.1:8000/process-user \
     -H "Content-Type: application/json" \
     -d '{"user_id":1,"name":"Alice","age":30}'
```

```json
{"status":"ok","message":"Processed user Alice"}
```

```bash
# Невалидные данные — 422 от валидации Pydantic
curl -X POST http://127.0.0.1:8000/process-user \
     -H "Content-Type: application/json" \
     -d '{"user_id":1,"name":"","age":300}'
```

```json
{"detail":[{"type":"string_too_short","loc":["body","name"], ... },
           {"type":"less_than_equal","loc":["body","age"], ... }]}
HTTP 422
```

**Цепочка запроса:** `curl -> FastAPI :8000 -> валидация Pydantic -> Go :8080 -> ответ`

Полный вывод: [task4/result.txt](task4/result.txt)

---

## Задание 6 — Сравнение скорости FastAPI и Gin под нагрузкой

Оба сервиса отдают одинаковый ответ `{"message":"pong"}` (18 байт).
Логирование и access-логи отключены, чтобы замерять чистую обработку запроса.
Замер сделан утилитой `hey` — кроссплатформенной заменой `wrk`/`ab`.

### Запуск

```bash
# Терминал 1
go run ./task6/go-service
```

```bash
# Терминал 2
cd task6/fastapi-service
uvicorn main:app --host 127.0.0.1 --port 8000 --no-access-log
```

### Замер

```bash
hey -n 10000 -c 100 http://127.0.0.1:8080/ping  > test_go.txt
hey -n 10000 -c 100 http://127.0.0.1:8000/ping  > test_fastapi.txt
```

**Фактический вывод:**

```
Gin:      Total: 0.1297 secs  Average: 0.0012 secs  Requests/sec: 77088.9367
FastAPI:  Total: 3.7543 secs  Average: 0.0354 secs  Requests/sec: 2663.6238
```

| Метрика              | Gin (Go)    | FastAPI (Python) |
|----------------------|-------------|------------------|
| Requests/sec         | 77088.94    | 2663.62          |
| Среднее время ответа | 0.0012 сек | 0.0354 сек       |
| Полное время 10 000  | 0.1297 сек  | 3.7543 сек       |
| Самое медленное      | 0.0136 сек  | 0.6894 сек       |

**Вывод:** Gin примерно в 29 раз быстрее по пропускной способности, среднее
время ответа меньше в 30 раз. Причина — нативный код Go против интерпретатора
Python. Полный разбор: [task6/comparison.txt](task6/comparison.txt),
логи: [test_go.txt](task6/test_go.txt), [test_fastapi.txt](task6/test_fastapi.txt).

---

## Лог промптов

Все промпты, по которым писался код, зафиксированы в [PROMPT_LOG.md](PROMPT_LOG.md).