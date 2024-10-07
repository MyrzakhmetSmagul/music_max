package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	"github.com/MyrzakhmetSmagul/music_max/config"
	"github.com/MyrzakhmetSmagul/music_max/pkg/handler"
	"github.com/MyrzakhmetSmagul/music_max/pkg/repository"
	"github.com/MyrzakhmetSmagul/music_max/pkg/service"
)

// @title Music API
// @version 1.0
// @description This is a sample server for music management.

// @contact.name Developer
// @contact.email smagulmyrzakhmet@gmail.com

// @host localhost:3030
// @BasePath /api/v1

func main() {
	mode := os.Getenv("MODE")
	options := new(slog.HandlerOptions)

	switch mode {
	case "PROD":
		options.Level = slog.LevelWarn
	case "DEV":
		options.Level = slog.LevelInfo
	default:
		options.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, options))
	slog.SetDefault(logger)

	err := config.InitConfig("./config/.env")
	if err != nil {
		slog.Error("error occured during initialization of configurations:", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		slog.Error("error occured during opening new connection with DB:", slog.Any("error", err))
		os.Exit(1)
	}

	repo := repository.NewRepository(db)
	lyricsService := service.NewLyricsService()
	service := service.NewSongService(repo, lyricsService)
	handler := handler.NewHandler(service)
	srv := new(musicmax.Server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info(fmt.Sprintf("Starting server on %s", os.Getenv("HTTP_PORT")))
		if err := srv.Run(handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			slog.Error("error occured during running http server:", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("Shutting down the server..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error occured when server shutting down", slog.Any("error", err))
		os.Exit(1)
	}

	if err := db.Close(); err != nil {
		slog.Error("error occured when server closing db connection", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
