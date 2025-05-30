
# Marketplace Project

Микросервисный backend-проект для симуляции маркетплейса. Реализована базовая бизнес-логика с инфраструктурой, включая мониторинг и обмен событиями между сервисами.

## 🧩 Архитектура

- **NGINX** — API Gateway
- **Orders Service** — оформление заказов
- **Customers Service** — управление клиентами
- **PostgreSQL** — основное хранилище
- **Kafka** — асинхронное взаимодействие между сервисами
- **Prometheus + Grafana** — мониторинг и визуализация
- **Docker + Docker Compose** — сборка и оркестрация

## 📁 Структура

![Снимок экрана от 2025-05-28 11-13-20](https://github.com/user-attachments/assets/77656e98-4e81-4fa0-9d0f-0f8a8372c686)


```
/marketplace-project
│
├── services/
│   ├── orders/
│   │   ├── cmd/            # main.go и инициализация сервиса
│   │   ├── internal/       # бизнес-логика (DDD-модули)
│   │   └── Dockerfile
│   │
│   ├── customers/
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── Dockerfile
│
├── api-gateway/            # NGINX 
│   ├── nginx.conf
│   └── Dockerfile
│
├── monitoring/            # Мониторинг с Prometheus
│   ├── prometheus/
│
├── docker-compose.yaml     # локальный запуск 
└── README.md
```

## ⚙️ Функциональность

### Orders Service
- Создание заказов.
- Публикация события `order_placed` в Kafka.

### Customers Service
- CRUD клиентов.
- Обработка Kafka-события `order_placed` с повышением `activity_score` у клиента.

### Kafka Events

```json
{
  "event": "order_placed",
  "timestamp": "2025-05-28T08:00:00Z",
  "order_id": 123,
  "customer_id": 456,
  "status": "created"
}
````

### Мониторинг

* Метрики собираются Prometheus.
* Дашборд на Grafana ([http://localhost:3000](http://localhost:3000)).

![Снимок экрана от 2025-05-27 15-46-31](https://github.com/user-attachments/assets/ab20c19e-fe5c-45d2-87d2-094c63c1fa62)


## 🚀 Запуск проекта

```bash
docker-compose up --build // Запуск контейнеров

cd services/orders/       // Выполнить миграции
go run scripts/migrate.go
cd services/customers/
go run scripts/migrate.go
```

Ожидаемые сервисы:

* `http://localhost:8080` — NGINX (прокси к сервисам)
* `http://localhost:9090` — Prometheus
* `http://localhost:3000` — Grafana
* `http://localhost:9092` — Kafka

## 📦 Используемые технологии

* Golang (без ORM, только `pgxpool`)
* PostgreSQL
* Apache Kafka
* NGINX
* Docker
* Prometheus & Grafana
* Structured Logging via `slog`

## 📌 ToDo

* Добавить кэширование запросов через Redis, обращаться с Gateway сначала к нему
* Сделать Auth Service, с которым будет общаться Gateway
* Развернуть приложение в Kubernetes

