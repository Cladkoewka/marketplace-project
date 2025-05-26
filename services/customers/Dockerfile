# 1. Этап сборки
FROM golang:latest AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .  

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o customers-service ./cmd

# 2. Минимальный рантайм образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/customers-service .

ENTRYPOINT ["./customers-service"]
