# Укажите базовый образ Go
FROM golang:1.21.1-alpine AS builder

# Установим необходимые пакеты
RUN apk update && apk add --no-cache git

# Создадим рабочую директорию
WORKDIR /app

# Скопируем исходный код
COPY . .

# Соберем Go-приложение
RUN go build -o main ./cmd/app

# Базовый образ для финального контейнера
FROM alpine:latest

# Установим необходимые пакеты
RUN apk --no-cache add ca-certificates

# Создадим рабочую директорию
WORKDIR /root/

# Скопируем бинарный файл
COPY --from=builder /app/main .

COPY .env .

# Запуск приложения
CMD ./main
