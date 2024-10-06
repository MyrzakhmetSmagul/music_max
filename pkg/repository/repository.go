package repository

import (
	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

type Repository interface {
	GetSongs(filters map[string]string, page int, limit int) ([]musicmax.Song, error)
	CreateSong(*musicmax.Song) error
	CreateLyrics()
	DeleteSong(id string) error
}
