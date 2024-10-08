package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

// getLyrics return song lyrics by id.
//
// @Summary Получить текст песни
// @Description Получить текст песни по id
// @Tags Songs
// @Param id path string true "Song ID"
// @Produce  json
// @Success 200 {object} model.LyricsResponse "Текст песни с пагинацией по куплетам"
// @Failure 400 {object} model.Description "Bad Request"
// @Failure 500 {object} model.Description "Internal Server Error"
// @Router /songs/{id}/lyrics [get]
func (h *Handler) getLyrics(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET api/v1/songs/{id}/lycirs")
	id := r.PathValue("id")
	page, limit, err := musicmax.GetPageAndLimit(r)
	if err != nil {
		model.DefaultResponse(w, http.StatusBadRequest)
		slog.Error("getLyrics error during convertation page or limit value to int", slog.String("page", r.URL.Query().Get("page")), slog.Any("error", err))
		return
	}

	slog.Debug("getLyrics", slog.Any("page", page), slog.Any("limit", limit), slog.Any("id", id))

	lyrics, err := h.service.GetLyrics(id, page, limit)
	if err != nil && !errors.Is(err, model.ErrBadRequest) {
		model.DefaultResponse(w, http.StatusInternalServerError)
		slog.Error("error get lyrics", slog.Any("error", err))
		return
	} else if errors.Is(err, model.ErrBadRequest) {
		model.DefaultResponse(w, http.StatusBadRequest)
		slog.Warn("error get lyrics", slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(lyrics); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		slog.Error("failed to encode lyrics to JSON", slog.Any("error", err))
		return
	}
}
