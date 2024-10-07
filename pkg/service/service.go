package service

import (
	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

type Song interface {
	GetSongs(filters map[string]string, page int, limit int) (*musicmax.SongsResponse, error)
	CreateSong(songReq *musicmax.SongRequest) error
	GetLyrics(id string, page int, limit int) (*musicmax.LyricsResponse, error)
	DeleteSong(id string) error
	UpdateSong(id string, song *musicmax.SongPatchRequest) error
}

type Lyrics interface {
	GetLyrics(songReq *musicmax.SongRequest) (*musicmax.LyricsAPIResponse, error)
}
