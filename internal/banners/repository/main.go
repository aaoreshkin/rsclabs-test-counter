package repository

import (
	"context"
	"time"

	"github.com/aaoreshkin/click-counter/internal/banners/model"
	"github.com/aaoreshkin/click-counter/internal/provider/database"
	"github.com/jackc/pgx/v5"
)

type (
	// Repository предоставляет доступ к данным счетчиков баннеров в базе данных.
	Repository struct {
		connection *database.Connection
	}
)

// Новый экземпляр Repository с переданным подключением к БД.
func New(connection *database.Connection) *Repository {

	return &Repository{
		connection,
	}
}

// Сохраняет батч данных из шардов в БД.
// Xранения агрегированных данных как в ТЗ по минутам.
func (r *Repository) BatchData(ctx context.Context, data map[int]int64) error {
	if len(data) == 0 {
		return nil
	}
	const query = `
	INSERT INTO banners_counter (
			banner_id,
			ts,
			v
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (banner_id, ts)
		DO UPDATE SET v = banners_counter.v + EXCLUDED.v
	`

	batch := &pgx.Batch{}

	for id, v := range data {
		batch.Queue(query, id, time.Now().Truncate(time.Minute), v)
	}

	br := r.connection.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(data); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

// Возвращает статистику по баннеру за указанный период времени.
// Данные возвращаются отсортированными по времени.
func (r *Repository) GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]model.Counter, error) {
	const query = `
		SELECT banner_id, ts, v
		FROM banners_counter
		WHERE banner_id = $1 AND ts >= $2 AND ts <= $3
		ORDER BY ts
	`

	rows, err := r.connection.Query(ctx, query, bannerID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.Counter
	for rows.Next() {
		var counter model.Counter
		if err := rows.Scan(&counter.ID, &counter.TS, &counter.V); err != nil {
			return nil, err
		}
		stats = append(stats, counter)
	}

	return stats, rows.Err()
}
