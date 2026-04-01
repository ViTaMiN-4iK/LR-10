# LR-10
# Лабораторная работа №10: Веб-разработка — FastAPI (Python) vs Gin (Go)

## Цель работы
Сравнить подходы к созданию веб-сервисов на двух языках программирования (Python и Go), реализовать их взаимодействие через REST и gRPC, а также контейнеризировать оба сервиса с помощью Docker Compose.

## Выполненные задачи

### Средняя сложность
1. Создано API на Go (Gin) с 3 эндпоинтами:
   - `GET /health` — проверка работоспособности
   - `POST /items` — создание элемента (JSON: `{"name": "string", "price": float}`)
   - `GET /items/:id` — получение элемента по UUID

2. Создано API на Python (FastAPI) с аналогичными эндпоинтами

3. Реализовано взаимодействие сервисов через REST (Python прокси к Go)

4. Реализован graceful shutdown в обоих сервисах

### Повышенная сложность
1. Реализован gRPC-сервер на Go и gRPC-клиент на Python
   - `POST /grpc-items` — создание элемента через gRPC
   - `GET /grpc-items/{id}` — получение элемента через gRPC

2. Развернуты оба сервиса в Docker Compose с общей сетью

## Структура проекта
lab10/
├── .gitignore
├── docker-compose.yml
├── README.md
├── PROMPT_LOG.md
├── PLAN.md
├── proto/
│ └── items.proto
├── go-service/
│ ├── Dockerfile
│ ├── go.mod
│ ├── go.sum
│ ├── main.go
│ ├── handlers/
│ │ ├── item_handler.go
│ │ └── item_handler_test.go
│ ├── models/
│ │ └── item.go
│ ├── storage/
│ │ └── storage.go
│ ├── grpcserver/
│ │ └── server.go
│ └── grpc/
│ └── pb/
│ ├── items.pb.go
│ └── items_grpc.pb.go
└── python-service/
├── Dockerfile
├── requirements.txt
├── main.py
├── schemas.py
├── storage.py
├── grpc_client.py
├── test_main.py
├── items_pb2.py
└── items_pb2_grpc.py


## Запуск проекта

### Локальный запуск

**Go сервис:**
```bash
cd go-service
go mod tidy
go run main.go
# Сервер запустится на порту 8080 (REST) и 50051 (gRPC)

Python сервис:
cd python-service
pip install -r requirements.txt
python main.py
# Сервер запустится на порту 8000

Запуск через Docker Compose
docker compose up --build

Тестирование
REST эндпоинты
# Health check
curl http://localhost:8080/health
curl http://localhost:8000/health

# Создание элемента в Go
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":999.99}'

# Получение элемента через Python прокси
curl http://localhost:8000/proxy-items/{id}

gRPC эндпоинты (через Python API)
# Создание элемента через gRPC
curl -X POST http://localhost:8000/grpc-items \
  -H "Content-Type: application/json" \
  -d '{"name":"gRPC Item","price":777.77}'

# Получение элемента через gRPC
curl http://localhost:8000/grpc-items/{id}

Результаты тестирования
Эндпоинт	Метод	Статус	Описание
/health (Go)	GET	200 OK	Проверка работоспособности Go
/health (Python)	GET	200 OK	Проверка работоспособности Python
/items (Go)	POST	201 Created	Создание элемента в Go
/items/{id} (Go)	GET	200 OK	Получение элемента из Go
/grpc-items (Python)	POST	200 OK	Создание элемента через gRPC
/grpc-items/{id} (Python)	GET	200 OK	Получение элемента через gRPC
/proxy-items/{id} (Python)	GET	200 OK	Прокси к Go через REST

Сравнение подходов
Характеристика	FastAPI (Python)	Gin (Go)
Скорость разработки	Высокая (автодокументация, Pydantic)	Средняя
Производительность	Хорошая (асинхронная)	Отличная (горутины)
Валидация данных	Встроенная (Pydantic)	Требует дополнительных библиотек
gRPC поддержка	Через сторонние библиотеки	Нативная (protobuf)
Graceful shutdown	Lifespan context manager	signal.Notify + Shutdown
Размер образа Docker	~180 MB	~15 MB

Выводы
FastAPI предоставляет более высокую скорость разработки благодаря встроенной валидации, автодокументации и асинхронной поддержке.

Gin демонстрирует лучшую производительность и меньший размер Docker образа, что важно для микросервисной архитектуры.

gRPC обеспечивает более эффективную коммуникацию между сервисами по сравнению с REST (бинарный протокол, потоковая передача).

Docker Compose упрощает оркестрацию нескольких сервисов и их взаимодействие в изолированной сети.

Используемые технологии
Go 1.21 + Gin + gRPC
Python 3.11 + FastAPI + httpx + gRPC
Docker + Docker Compose
Protocol Buffers (protobuf)

Лабораторная работа выполнена в рамках курса по методологии Agentic Engineering.