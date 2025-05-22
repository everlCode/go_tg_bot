FROM golang:1.24

WORKDIR /app

ENV GOFLAGS="-buildvcs=false"

RUN apt-get update && apt-get install -y gcc libsqlite3-dev && apt install -y sqlite3

RUN go install github.com/air-verse/air@latest

COPY . .
# Скачиваем зависимости и собираем приложение
RUN go mod tidy && go build -o main .

RUN mkdir -p /app/tmp

CMD ["air"]