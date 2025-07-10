package router

import (
	"context"
	"net/http"

	"github.com/aaoreshkin/click-counter/internal"
	"github.com/go-chi/chi/v5"
)

type (
	// Mux представляет основной HTTP роутер приложения.
	Mux struct {
		*chi.Mux

		manager *internal.Manager // Ссылка на корневой менеджер для доступа к модулям.
	}
)

// Новый экземпляр Mux с полной настройкой роутинга.
// Настраивает:
// - Мидлвар для автоматической установки Content-Type: application/json
// - Версионированные роуты под префиксом /v1
// - Монтирование модуля баннеров по пути /v1/banners
// - Эндпоинт проверки состояния /v1/healthcheck, но по большей части для прогрева TCP для тестов
func New(ctx context.Context, manager *internal.Manager) (*Mux, error) {

	router := &Mux{chi.NewRouter(), manager}

	router.Route("/v1", func(r chi.Router) {

		r.Mount("/banners", router.routeBanners())

		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return router, nil
}
