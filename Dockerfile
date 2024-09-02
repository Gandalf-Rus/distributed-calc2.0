# Dockerfile для backend (сервер и агент)

# Используем базовый образ Golang для сборки и запуска
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum, устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходные файлы
COPY . ./

# Собираем сервер
RUN go build -o /server ./cmd/server

# Собираем программу (agent)
RUN go build -o /agent ./cmd/agent

# Используем тот же базовый образ для финального контейнера
FROM golang:1.21-alpine

WORKDIR /app

COPY --from=builder /server /app/server
COPY --from=builder /agent /app/agent

# Экспонируем порты для сервера
EXPOSE 8080
EXPOSE 5000

# CMD будет указываться в docker-compose.yml
