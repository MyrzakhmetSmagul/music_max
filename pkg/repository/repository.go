package repository

import (
	"database/sql"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

type Repository interface {
	GetSongs(filters map[string]string, page int, limit int) ([]musicmax.Song, error)
	CreateSong(tx *sql.Tx, song *musicmax.Song) error
	CreateLyrics(tx *sql.Tx, lyrics *musicmax.Lyrics) error
	DeleteSong(id string) error
	GetLyrics(songId string) (*musicmax.Lyrics, error)
	UpdateSong(tx *sql.Tx, song *musicmax.Song) error
	UpdateLyrics(tx *sql.Tx, songId string, lyrics *musicmax.Lyrics) error
	BeginTx() (*sql.Tx, error)
}
