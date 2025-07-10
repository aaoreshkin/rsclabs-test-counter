package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aaoreshkin/click-counter/internal"
	"github.com/aaoreshkin/click-counter/internal/provider/database"
	"github.com/aaoreshkin/click-counter/internal/router"
)

var (
	connection *database.Connection
	mux        *router.Mux

	err error
)

// Точка входа.
// Инициализирует контекст и запускает основную логику приложения.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Printf("Application error: %v", err)
	}
}

// Выполняет основную логику приложения:
// - подключается к базе данных
// - инициализирует корневой менеджер (контролит других менеджеров отвечающих за модуль)
// - настраивает HTTP роутер
// - запускает HTTP сервер на порту из переменной окружения SERVICE_PORT
func run(ctx context.Context) error {
	if connection, err = database.New(ctx); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}
	defer connection.Close()

	manager, err := internal.New(ctx, connection)
	if err != nil {
		return err
	}

	if mux, err = router.New(ctx, manager); err != nil {
		log.Println(err)
	}

	server := &http.Server{
		Addr:           ":" + os.Getenv("SERVICE_PORT"),
		Handler:        mux,
		ReadTimeout:    0, // Без таймаутов
		WriteTimeout:   0,
		IdleTimeout:    0,
		MaxHeaderBytes: 1 << 10, // 1KB - минимум для заголовков
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP server error: %v\n", err)
	}

	return nil
}
