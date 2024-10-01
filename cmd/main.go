package main

import (
	"log/slog"
	"os"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	"github.com/MyrzakhmetSmagul/music_max/pkg/handler"
)

func main() {
	handler := new(handler.Handler)
	srv := new(musicmax.Server)
	slog.Info("http://localhost:3003 starting")
	if err := srv.Run(handler.InitRoutes()); err != nil {
		slog.Error("error occured while running http server:", slog.Any("error", err))
		os.Exit(1)
	}
}
