package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Декодирует JSON из тела HTTP запроса в указанный тип T.
// Использует дженерики для типобезопасности и возвращает указатель на декодированную структуру.
func DecodeJSON[T any](r *http.Request) (*T, error) {
	var payload T
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}
	return &payload, nil
}

// Извлекает целочисленный параметр из URL пути HTTP запроса.
// Использует chi для получения параметра и конвертирует его в int.
func IntParam(r *http.Request, param string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, param))
}
