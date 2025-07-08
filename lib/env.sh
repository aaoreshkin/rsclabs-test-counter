#!/bin/sh

# Включение режима отладки
# В продакшене должно быть установлено в 0 или не задано
export DEBUG=1

# Порт для HTTP сервера
# Используется в cmd/main.go для запуска сервера
export SERVICE_PORT=3000

# Настройка подключения к базе данных в зависимости от режима
if [ "$DEBUG" = 1 ]; then
    # Локальная PostgreSQL для разработки
    # sslmode=disable используется только для локальной разработки
    # В продакшене ОБЯЗАТЕЛЬНО должен быть включен SSL
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
fi

# Генерация случайного секретного ключа при каждом запуске
# В продакшене должен быть статичным и храниться в безопасном месте
export SECRET_KEY="$(openssl rand -base64 32)"
