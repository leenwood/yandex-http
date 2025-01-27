# Используем базовый образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Собираем приложение
RUN go build -o app ./cmd/client/main.go

# Открываем порт приложения
EXPOSE 9000

# Запуск приложения
CMD ["./app"]