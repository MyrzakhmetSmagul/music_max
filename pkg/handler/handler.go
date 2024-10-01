package handler

import "net/http"

type Handler struct {
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/songs", h.getSongs)
	return mux
}
