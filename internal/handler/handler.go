package handler

import (
	"net/http"

	_ "github.com/MyrzakhmetSmagul/music_max/docs"
	"github.com/MyrzakhmetSmagul/music_max/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
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
	mux.HandleFunc("DELETE /api/v1/songs/{id}", h.deleteSong)
	mux.HandleFunc("PATCH /api/v1/songs/{id}", h.updateSong)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	return mux
}
