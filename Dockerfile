FROM golang:1.24-alpine AS builder

#RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest
# Копируем только файлы модулей для кеша зависимостей
COPY app/go.mod app/go.sum ./
RUN go mod download

# Копируем весь исходный код проекта из go_tg_bot
COPY app/. .

RUN go build -o bot .

FROM alpine:3.18

#RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/bot .

CMD CompileDaemon --build="go build -o bot ." --command="./bot"
