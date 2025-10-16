FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/hackathon

RUN go build -o hackathon

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/cmd/hackathon/hackathon .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Копируем миграции
COPY --from=builder /app/migration ./migration

RUN apk add --no-cache ca-certificates

# Создаем скрипт запуска
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Waiting for database..."' >> /app/start.sh && \
    echo 'sleep 5' >> /app/start.sh && \
    echo 'echo "Running migrations..."' >> /app/start.sh && \
    echo 'goose -dir ./migration postgres "user=postgres password=1 dbname=indev sslmode=disable host=db" up' >> /app/start.sh && \
    echo 'echo "Starting application..."' >> /app/start.sh && \
    echo 'exec ./hackathon' >> /app/start.sh && \
    chmod +x /app/start.sh

CMD ["/app/start.sh"]
