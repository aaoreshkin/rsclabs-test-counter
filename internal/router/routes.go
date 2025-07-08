package router

import (
	"github.com/go-chi/chi/v5"
)

// Регистрирует эндпоинты модуля Banners:
func (mux *Mux) routeBanners() chi.Router {
	router := chi.NewRouter()

	controller := mux.manager.Banners.Controller()

	// - GET /counter/{bannerID} - инкремент счетчика баннера
	router.Get("/counter/{bannerID}", controller.HandleClick)

	// - POST /stats/{bannerID} - получение статистики по баннеру за период
	router.Post("/stats/{bannerID}", controller.HandleStats)

	return router
}
