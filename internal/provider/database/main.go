package database

import (
	"context"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	// Обертка над пулом соединений PostgreSQL.
	// Инкапсулирует pgxpool.Pool с оптимальными настройками.
	// Предоставляет высокопроизводительные соединения с автоматическим управлением жизненным циклом.
	Connection struct {
		*pgxpool.Pool
	}
)

var (
	// Пакетные переменные для инициализации пула соединений.
	dbpoolConfig *pgxpool.Config
	dbpool       *pgxpool.Pool

	err error
)

// Новое соединение с PostgreSQL с оптимальными настройками пула.
//
// Конфигурация пула соединений:
// - MaxConns: NumCPU * 5 - максимальное количество соединений на основе количества CPU ядер
// - MinConns: 5 - минимальное количество постоянно открытых соединений
// - HealthCheckPeriod: 1 минута - интервал проверки здоровья соединений
// - MaxConnLifetime: 12 часов - максимальное время жизни соединения
// - MaxConnIdleTime: 15 минут - максимальное время простоя соединения
// - ConnectTimeout: 10 секунд - таймаут установки соединения
//
// Настройки TCP соединения:
// - KeepAlive: синхронизирован с HealthCheckPeriod для оптимального обнаружения обрывов связи
// - Timeout: соответствует ConnectTimeout для консистентности
//
// Параметры оптимизированы для:
// - Высоконагруженных приложений с частыми короткими запросами
// - Минимизации латентности за счет поддержания минимального количества соединений
// - Обнаружения и восстановления от сетевых сбоев
//
// Ожидает DATABASE_URL в переменных окружения в формате PostgreSQL connection string.
func New(ctx context.Context) (*Connection, error) {
	// Парсинг из переменной окружения
	if dbpoolConfig, err = pgxpool.ParseConfig(os.Getenv("DATABASE_URL")); err != nil {
		return nil, err
	}

	// Максимальное количество соединений: 5 на каждое CPU ядро
	// Оптимально для I/O-интенсивных приложений
	dbpoolConfig.MaxConns = int32(runtime.NumCPU() * 5)

	// Минимальное количество постоянно открытых соединений
	// Обеспечивает мгновенный отклик на запросы
	dbpoolConfig.MinConns = 5

	// Период проверки здоровья соединений
	// Позволяет быстро обнаруживать и заменять нерабочие соединения
	dbpoolConfig.HealthCheckPeriod = 1 * time.Minute

	// Максимальное время жизни соединения
	// Предотвращает накопление долгоживущих соединений с потенциальными проблемами
	dbpoolConfig.MaxConnLifetime = 12 * time.Hour

	// Максимальное время простоя соединения
	// Освобождает неиспользуемые соединения, оставляя минимальное количество
	dbpoolConfig.MaxConnIdleTime = 15 * time.Minute

	// Таймаут установки нового соединения
	// Предотвращает зависание при проблемах с сетью или БД
	dbpoolConfig.ConnConfig.ConnectTimeout = 10 * time.Second

	// Настройка TCP соединения с оптимизацией для обнаружения обрывов связи
	dbpoolConfig.ConnConfig.DialFunc = (&net.Dialer{
		// KeepAlive синхронизирован с HealthCheckPeriod
		// Обеспечивает согласованность между TCP и приложением в обнаружении проблем
		KeepAlive: dbpoolConfig.HealthCheckPeriod,

		// Timeout соответствует ConnectTimeout
		// Обеспечивает консистентность в обработке таймаутов
		Timeout: dbpoolConfig.ConnConfig.ConnectTimeout,
	}).DialContext

	// Создание пула с конфигом
	if dbpool, err = pgxpool.NewWithConfig(ctx, dbpoolConfig); err != nil {
		return nil, err
	}

	return &Connection{
		dbpool,
	}, nil
}
