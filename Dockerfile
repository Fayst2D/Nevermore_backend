FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/user-service

RUN go build -o nevermore

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/cmd/user-service .

COPY --from=builder /app/config/config.yaml /app/config/

RUN apk add --no-cache ca-certificates

CMD ["./user-service"]