package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) getSongs(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status  string
		Message string
	}{
		Status:  "success",
		Message: "Data processed successfully",
	}
	slog.Debug("api/v1/songs")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
