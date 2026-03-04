# Сборка
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/sprin1 ./cmd/app

# Финальный образ
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/sprin1 .

EXPOSE 8080

# Переменные окружения задаются при запуске (docker run -e или compose env)
ENTRYPOINT ["./sprin1"]
