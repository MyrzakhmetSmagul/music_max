package musicmax

import (
	"context"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           os.Getenv("HTTP_PORT"),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(GetIntEnv("HTTP_READ_TIMEOUT", 10)) * time.Second,
		WriteTimeout:   time.Duration(GetIntEnv("HTTP_WRITE_TIMEOUT", 10)) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
