FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev
# Копируем только файлы модулей для кеша зависимостей
COPY app/go.mod app/go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY app/. .

# Собираем бинарник
RUN go build -o bot .

FROM alpine:3.18

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/bot .

# Запускаем без live reload
CMD ["./bot"]
