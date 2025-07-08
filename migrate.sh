#!/bin/bash

# Загрузка переменных окружения
. ./lib/env.sh

# Обработка аргументов командной строки
while [[ "$#" -gt 0 ]]; do
  case "$1" in
  -up)
    # Применение всех неприменённых миграций
    # Безопасная операция - применяет только новые миграции
    migrate -path migrations -database "${DATABASE_URL}" up
    shift
    ;;
  -down)
    # Откат ВСЕХ миграций (ОПАСНО!)
    # Полностью очищает схему БД
    migrate -path migrations -database "${DATABASE_URL}" down
    shift
    ;;
  -drop)
    # Удаление всех таблиц и данных (НЕВЕРОЯТНО ОПАСНО!)
    # Полностью уничтожает всю схему и данные
    migrate -path migrations -database "${DATABASE_URL}" drop
    shift
    ;;
  -create)
    # Создание новой миграции с указанным именем
    # Создает пару файлов: up и down миграции
    # Требует обязательный параметр - имя миграции
    if [ -n "$2" ]; then
      migrate create -ext sql -dir migrations "$2"
      shift 2
    else
      echo "Error: Missing migration name." >&2
      echo "Usage: ./migrate.sh -create <migration_name>" >&2
      exit 1
    fi
    ;;
  -goto)
    # Переход к конкретной версии миграции
    # Может применять или откатывать миграции для достижения целевой версии
    # Требует номер версии в качестве параметра
    if [ -n "$2" ]; then
      migrate -path migrations -database "${DATABASE_URL}" goto "$2"
      shift 2
    else
      echo "Error: Missing version number." >&2
      echo "Usage: ./migrate.sh -goto <version_number>" >&2
      exit 1
    fi
    ;;
  -fix)
    # Принудительная установка версии миграции (для исправления ошибок)
    # Используется когда миграция завершилась с ошибкой и нужно исправить состояние
    # ОСТОРОЖНО: не выполняет саму миграцию, только обновляет версию в schema_migrations
    if [ -n "$2" ]; then
      migrate -path migrations -database "${DATABASE_URL}" force "$2"
      shift 2
    else
      echo "Error: Missing version number for fix." >&2
      echo "Usage: ./migrate.sh -fix <version_number>" >&2
      exit 1
    fi
    ;;
  *)
    # Обработка неизвестных опций с выводом справки
    echo "Invalid option: $1" >&2
    echo "" >&2
    echo "Available options:" >&2
    echo "  -up                    Apply all pending migrations" >&2
    echo "  -down                  Rollback ALL migrations (DESTRUCTIVE)" >&2
    echo "  -drop                  Drop all database objects (DESTRUCTIVE)" >&2
    echo "  -create <name>         Create new migration files" >&2
    echo "  -goto <version>        Migrate to specific version" >&2
    echo "  -fix <version>         Force set migration version (for error recovery)" >&2
    echo "" >&2
    echo "Examples:" >&2
    echo "  ./migrate.sh -up" >&2
    echo "  ./migrate.sh -create add_users_table" >&2
    echo "  ./migrate.sh -goto 1" >&2
    exit 1
    ;;
  esac
done
