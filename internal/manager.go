package internal

import (
	"context"
	"runtime"
	"time"

	"github.com/aaoreshkin/click-counter/internal/banners"
	"github.com/aaoreshkin/click-counter/internal/provider/database"
	"github.com/aaoreshkin/click-counter/internal/provider/inmemory"
)

const (
	// Количество фоновых воркеров для периодического сброса кэша в БД.
	// Оптимальное значение зависит от нагрузки и производительности БД.
	workers = 4

	// Количество шардов кэша.
	shards = 128

	// Интервал между сбросами кэша в БД.
	// Меньший интервал = меньше потерь при сбоях, но больше нагрузка на БД.
	interval = time.Second
)

type (
	// Корневой менеджер приложения, координирующий все модули.
	// Инкапсулирует создание общих зависимостей (БД, кэш) и инициализацию модулей.
	// Служит точкой входа для настройки всей архитектуры приложения.
	Manager struct {
		connection *database.Connection
		cache      *inmemory.Cache

		// Менеджер модуля баннеров, предоставляющий доступ к его функциональности.
		Banners *banners.Manager
	}
)

// Новый экземпляр корневого Manager с инициализацией всех модулей.
// Количество шардов кэша на основе количества CPU ядер (NumCPU * 2) как временное решение
// для оптимального баланса между производительностью и потреблением памяти.
// Инициализирует модуль баннеров с предустановленными параметрами воркеров и интервала сброса.
func New(ctx context.Context, connection *database.Connection) (*Manager, error) {

	cache := inmemory.New(runtime.NumCPU() * 2)

	return &Manager{
		Banners: banners.New(ctx, connection, cache, workers, interval),
	}, nil
}
