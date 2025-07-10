package common

import (
	"fmt"
	"net/http"

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
func IntParam(r *http.Request, param string) (int, error) {
	s := chi.URLParam(r, param)
	if s == "" {
		return 0, fmt.Errorf("parameter %s not found", param)
	}

	result := 0

	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, fmt.Errorf("invalid integer: %s", s)
		}
		result = result*10 + int(s[i]-'0')
	}
	return result, nil
}
