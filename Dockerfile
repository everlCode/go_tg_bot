FROM golang:1.22-alpine

WORKDIR /app

# Копируем только main.go — чтобы инициализировать модуль
COPY main.go .

# Инициализируем go.mod внутри контейнера
RUN go mod init telegram-bot && \
    go mod tidy

# Затем копируем всё остальное (если будут другие файлы)
COPY . .

# Собираем бинарник
RUN go build -o bot

# Запускаем
CMD ["./bot"]
