#!/bin/sh

# Загрузка переменных окружения
. ./lib/env.sh

echo "Starting application in development mode..."
echo "Server will be available at http://localhost:$SERVICE_PORT"
echo "Press Ctrl+C to stop"
go run cmd/*.go
