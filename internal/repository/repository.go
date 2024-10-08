package repository

import (
	"database/sql"

	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

type Repository interface {
	GetSongs(filters map[string]string, page int, limit int) ([]model.Song, error)
	CreateSong(tx *sql.Tx, song *model.Song) error
	CreateLyrics(tx *sql.Tx, lyrics *model.Lyrics) error
	DeleteSong(id string) error
	GetLyrics(songId string) (*model.Lyrics, error)
	UpdateSong(tx *sql.Tx, song *model.Song) error
	UpdateLyrics(tx *sql.Tx, songId string, lyrics *model.Lyrics) error
	BeginTx() (*sql.Tx, error)
}
