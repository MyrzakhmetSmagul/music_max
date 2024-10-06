package service

import (
	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

type Song interface {
	GetSongs(filters map[string]string, page int, limit int) (*musicmax.SongsResponse, error)
	CreateSong(songReq *musicmax.SongRequest) error
	GetLyrics(song *musicmax.Song) error
}

type Lyrics interface {
	GetLyrics(songReq *musicmax.SongRequest) (*musicmax.Song, error)
}
