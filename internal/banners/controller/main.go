package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaoreshkin/click-counter/common"
	"github.com/aaoreshkin/click-counter/internal/banners/model"
)

type (
	// Controller обрабатывает HTTP запросы для работы со счетчиками баннеров.
	Controller struct {
		usecase model.Usecase
	}
)

// Новый экземпляр Controller с переданным usecase.
func New(usecase model.Usecase) *Controller {

	return &Controller{
		usecase,
	}
}

// Обрабатывает клики по баннеру, увеличивая счетчик на 1.
// Ожидает bannerID в параметрах запроса. Возвращает 204 No Content.
func (c *Controller) HandleClick(w http.ResponseWriter, r *http.Request) {
	bannerID, err := common.IntParam(r, "bannerID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.usecase.Increment(bannerID)
	w.WriteHeader(http.StatusNoContent)
}

// Возвращает статистику по баннеру за указанный период.
// Ожидает bannerID в параметрах и JSON с полями from/to в теле запроса.
// Время должно быть в формате RFC3339. Пример запроса в Readme.
func (c *Controller) HandleStats(w http.ResponseWriter, r *http.Request) {

	// Обертка в директории common
	bannerID, err := common.IntParam(r, "bannerID")
	if err != nil {
		http.Error(w, "invalid bannerID", http.StatusBadRequest)
		return
	}

	data, err := common.DecodeJSON[model.Stats](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, data.From)
	if err != nil {
		http.Error(w, "invalid from time format", http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, data.To)
	if err != nil {
		http.Error(w, "invalid to time format", http.StatusBadRequest)
		return
	}

	stats, err := c.usecase.GetStats(r.Context(), bannerID, from, to)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.StatsResponse{Stats: stats})
}
