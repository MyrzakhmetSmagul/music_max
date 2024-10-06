package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

func (h *Handler) getSongs(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET api/v1/songs")

	page, limit, err := musicmax.GetPageAndLimit(r)

	if err != nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		slog.Error("error during convertation page or limit value to int",
			slog.String("page", r.URL.Query().Get("page")),
			slog.String("limit", r.URL.Query().Get("limit")),
			slog.Any("error", err))
		return
	}

	slog.Debug("getSongs", slog.Any("page", page), slog.Any("limit", limit))
	slog.Debug("query params", slog.Any("queryParams", r.URL.Query().Encode()))

	filters, err := musicmax.GetQueryParams(r)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		slog.Error("error during get filters", slog.Any("error", err))
		return
	}

	songs, err := h.service.GetSongs(filters, page, limit)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		slog.Error("error during convertation \"page\" value to int", slog.String("page", r.URL.Query().Get("page")), slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(songs)
}

func (h *Handler) createSong(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST api/v1/songs")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("handler.createSong error during read from body", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	songReq := new(musicmax.SongRequest)
	json.Unmarshal(body, songReq)
	slog.Debug("song request", slog.Any("songReq", songReq))
	if strings.TrimSpace(songReq.Group) == "" || strings.TrimSpace(songReq.Name) == "" {
		slog.Error("bad request")
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	err = h.service.CreateSong(songReq)
	if err != nil {
		slog.Error("handler.createSong error during creating song", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		return
	}

	musicmax.DefaultResponse(w, http.StatusOK)
}
