FROM golang:1.23.1

WORKDIR /app

# Копируем только main.go для ускорения сборки
COPY main.go .

# Инициализируем модуль и устанавливаем зависимости
RUN go mod init ari-app
RUN go get github.com/go-redis/redis/v8
RUN go get github.com/go-resty/resty/v2

# Копируем остальные файлы
COPY . .

# Сборка приложения
RUN go build -o ari-app .

# Запуск приложения
CMD ["./ari-app"]
