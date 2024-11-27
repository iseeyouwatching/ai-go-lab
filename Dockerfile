# Используем официальное изображение Go
FROM golang:1.23.3-alpine as builder

# Устанавливаем зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN go build -o bot .

# Минимизируем конечное изображение
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /app/bot .

# Переменная окружения для токена бота
ENV TELEGRAM_BOT_TOKEN="7444307382:AAEjAq-yEJDa5o44gMuu2uNC4Sy2CYeyJPs"

# Команда для запуска бота
CMD ["./bot"]