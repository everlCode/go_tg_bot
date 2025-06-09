FROM golang:1.24-bullseye AS builder
# или
# FROM golang:1.24-bookworm AS builder

RUN apt-get update && apt-get install -y gcc musl-dev sqlite3 libsqlite3-dev

WORKDIR /app

# Копируем только файлы модулей для кеша зависимостей
COPY app/go.mod app/go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY app/. .

# Собираем бинарник
RUN go build -o bot .

# --- Финишный образ ---
FROM debian:bullseye-slim

ENV TZ=Europe/Moscow

RUN apt-get update && apt-get install -y sqlite3 ca-certificates tzdata && \
    ln -snf /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo "Europe/Moscow" > /etc/timezone

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/bot .

CMD ["./bot"]
