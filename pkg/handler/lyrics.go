package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

func (h *Handler) getLyrics(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET api/v1/songs/{id}/lycirs")
	id := r.PathValue("id")
	page, limit, err := musicmax.GetPageAndLimit(r)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		slog.Error("getLyrics error during convertation page or limit value to int", slog.String("page", r.URL.Query().Get("page")), slog.Any("error", err))
		return
	}

	slog.Debug("getLyrics", slog.Any("page", page), slog.Any("limit", limit), slog.Any("id", id))

	songs, err := h.service.GetSongs(nil, page, limit)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		slog.Error("error during convertation \"page\" value to int", slog.String("page", r.URL.Query().Get("page")), slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(songs)
}
