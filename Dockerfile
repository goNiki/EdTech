# Build stage
FROM golang:1.27-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарники с оптимизацией размера (-w -s убирают debug-информацию)
RUN go build -ldflags="-w -s" -o /bin/edtech ./cmd/edtech/main.go
RUN go build -ldflags="-w -s" -o /bin/migration ./cmd/migration/migration.go

# Production stage (минимальный образ)
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Копируем скомпилированные бинарники
COPY --from=builder /bin/edtech /app/edtech
COPY --from=builder /bin/migration /app/migration

# Копируем SQL-файлы миграций
COPY --from=builder /app/migrators /app/migrators

# Создаем папку под uploads с правами для appuser
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app

USER appuser
EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8082/api/v1/categories || exit 1

CMD ["/app/edtech"]
