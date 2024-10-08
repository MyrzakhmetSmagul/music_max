package service

import (
	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

type LyricsMockService struct {
}

var counter int

func NewLyricsMockService() Lyrics {
	return &LyricsMockService{}
}

func (l *LyricsMockService) GetLyrics(songReq *model.SongRequest) (*model.LyricsAPIResponse, error) {
	return l.mock()
}

func (l *LyricsMockService) mock() (*model.LyricsAPIResponse, error) {
	defer func() {
		counter++
	}()

	switch counter % 3 {
	case 0:
		date := "16.07.2006"
		resp := &model.LyricsAPIResponse{
			Release: &date,
			Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}
		return resp, nil
	case 1:
		return nil, model.ErrBadRequest
	default:
		return nil, model.ErrInternalServerError
	}
}
