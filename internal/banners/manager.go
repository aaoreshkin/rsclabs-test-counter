package banners

import (
	"context"
	"time"

	"github.com/aaoreshkin/click-counter/internal/banners/controller"
	"github.com/aaoreshkin/click-counter/internal/banners/repository"
	"github.com/aaoreshkin/click-counter/internal/banners/usecase"
	"github.com/aaoreshkin/click-counter/internal/provider/database"
	"github.com/aaoreshkin/click-counter/internal/provider/inmemory"
)

type (
	// Manager управляет жизненным циклом всех компонентов модуля баннеров.
	// Инкапсулирует создание зависимостей, запуск фоновых воркеров для сброса кэша
	// и предоставляет доступ к HTTP контроллеру для роутера.
	Manager struct {
		repository *repository.Repository
		usecase    *usecase.Usecase
		controller *controller.Controller
	}
)

// Новый экземпляр Manager с полной инициализацией всех компонентов.
// Запускает указанное количество воркеров для периодического сброса кэша в БД.
// Воркеры автоматически останавливаются при отмене контекста.
func New(ctx context.Context, connection *database.Connection, cache *inmemory.Cache, workers int, interval time.Duration) *Manager {

	repository := repository.New(connection)
	usecase := usecase.New(repository, cache)
	controller := controller.New(usecase)

	for range workers {
		go func() {
			ticker := time.NewTicker(interval)

			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					usecase.FlushToDB(ctx)
				}
			}
		}()
	}

	return &Manager{
		repository,
		usecase,
		controller,
	}
}

// Возвращает HTTP контроллер для регистрации роутов.
// Используется роутером для настройки эндпоинтов модуля баннеров.
func (m *Manager) Controller() *controller.Controller {

	return m.controller
}
