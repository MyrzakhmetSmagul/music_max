package service

import (
	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

type Song interface {
	GetSongs(filters map[string]string, page int, limit int) (*model.SongsResponse, error)
	CreateSong(songReq *model.SongRequest) error
	GetLyrics(id string, page int, limit int) (*model.LyricsResponse, error)
	DeleteSong(id string) error
	UpdateSong(id string, song *model.SongPatchRequest) error
}

type Lyrics interface {
	GetLyrics(songReq *model.SongRequest) (*model.LyricsAPIResponse, error)
}
