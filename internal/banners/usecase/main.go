package usecase

import (
	"context"
	"maps"
	"time"

	"github.com/aaoreshkin/click-counter/internal/banners/model"
	"github.com/aaoreshkin/click-counter/internal/provider/inmemory"
)

type (
	// Слой бизнес-логики для работы со счетчиками баннеров.
	Usecase struct {
		repository model.Repository
		cache      *inmemory.Cache
		buffer     map[int]int64 // переиспользуемый буфер
	}
)

// Новый экземпляр Usecase.
func New(repository model.Repository, cache *inmemory.Cache) *Usecase {
	return &Usecase{
		repository: repository,
		cache:      cache,
		buffer:     make(map[int]int64),
	}
}

// Увеличивает счетчик баннера на 1.
// Операция потокобезопасная благодаря шардированному кэшу с мьютексами.
func (u *Usecase) Increment(id int) {
	sh := u.cache.GetShard(id)

	sh.Mu.Lock()
	defer sh.Mu.Unlock()

	sh.Data[id]++
}

// Сбрасывает все накопленные в кэше данные в базу данных батчами.
// Операция атомарна для каждого шарда.
func (u *Usecase) FlushToDB(ctx context.Context) {
	// Каждый шард обрабатывается независимо
	for _, sh := range u.cache.Shards {
		sh.Mu.Lock()

		// Очищаем переиспользуемый буфер
		for k := range u.buffer {
			delete(u.buffer, k)
		}

		maps.Copy(u.buffer, sh.Data)

		sh.Data = make(map[int]int64) // После сброса кэш очищается.
		sh.Mu.Unlock()

		if len(u.buffer) == 0 {
			continue
		}
		// Ошибка в одном батче не останавливает другие
		u.repository.BatchData(ctx, u.buffer)
	}
}

// Возвращает статистику по баннеру за указанный период времени.
// Пример запроса в Readme.
func (u *Usecase) GetStats(ctx context.Context, bannerID int, from, to time.Time) ([]model.Counter, error) {
	return u.repository.GetStats(ctx, bannerID, from, to)
}
