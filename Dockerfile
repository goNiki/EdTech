# Этап сборки (Builder)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Отключаем CGO для создания полностью статического бинарника
ENV CGO_ENABLED=0 GOOS=linux

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник (на выходе получаем один файл /bin/edtech)
RUN go build -a -installsuffix cgo -o /bin/edtech ./cmd/edtech/main.go


# Финальный этап (Минимальный образ)
FROM alpine:latest

# Добавляем сертификаты и таймзоны (часто нужны для работы с сетью и временем)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем только готовый бинарник из первого этапа
COPY --from=builder /bin/edtech /app/edtech

# Экспортируем порт (информативно)
EXPOSE 8082

# Запускаем скомпилированное приложение
CMD ["/app/edtech"]
