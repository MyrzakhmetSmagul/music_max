package model

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Description struct {
	Description string `json:"description"`
}

var ErrBadRequest = errors.New("bad request")
var ErrInternalServerError = errors.New("internal server error")

func DefaultResponse(w http.ResponseWriter, statusCode int) {
	resp := Description{
		Description: http.StatusText(statusCode),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		slog.Error("failed to encode resp to JSON", slog.Any("error", err))
		return
	}
}
