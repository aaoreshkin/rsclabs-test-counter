package model

import (
	"context"
	"time"
)

type (
	// Представляет агрегированные данные счетчика баннера за определенный период времени.
	// ID не включается в JSON ответ, TS - временная метка, V - значение счетчика.
	Counter struct {
		ID int       `json:"-"`
		TS time.Time `json:"ts"`
		V  int       `json:"v"`
	}

	// Представляет запрос на получение статистики с временными границами.
	// Время передается в строковом формате для удобства JSON сериализации.
	Stats struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	// Представляет ответ с массивом статистических данных.
	StatsResponse struct {
		Stats []Counter `json:"stats"`
	}

	// Определяет интерфейс бизнес-логики для работы со счетчиками баннеров.
	// Абстрагирует controller от конкретной реализации usecase слоя.
	Usecase interface {
		Increment(int)
		GetStats(context.Context, int, time.Time, time.Time) ([]Counter, error)
	}

	// Repository определяет интерфейс доступа к данным счетчиков.
	// Абстрагирует usecase от конкретной реализации хранилища данных.
	Repository interface {
		BatchData(context.Context, map[int]int64) error
		GetStats(context.Context, int, time.Time, time.Time) ([]Counter, error)
	}
)
