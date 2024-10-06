package handler

import (
	"net/http"

	"github.com/MyrzakhmetSmagul/music_max/pkg/service"
)

type Handler struct {
	service service.Song
}

func NewHandler(service service.Song) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/songs", h.getSongs)
	mux.HandleFunc("GET /api/v1/songs/{id}/lyrics", h.getLyrics)
	mux.HandleFunc("POST /api/v1/songs", h.createSong)
	return mux
}
