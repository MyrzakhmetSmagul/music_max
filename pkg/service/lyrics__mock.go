package service

import (
	"fmt"
	"time"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

type LyricsMockService struct {
}

var counter int

func NewLyricsMockService() Lyrics {
	return &LyricsMockService{}
}

func (l *LyricsMockService) GetLyrics(songReq *musicmax.SongRequest) (*musicmax.LyricsAPIResponse, error) {
	return l.mock()
}

func (l *LyricsMockService) mock() (*musicmax.LyricsAPIResponse, error) {
	defer func() {
		counter++
	}()

	day, err := time.Parse("02.01.2006", "16.07.2006")
	if err != nil {
		err = fmt.Errorf("%w\nMock error parse time:%w", musicmax.ErrInternalServerError, err)
		return nil, err
	}

	switch counter % 3 {
	case 0:
		resp := &musicmax.LyricsAPIResponse{
			Release: &day,
			Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}
		return resp, nil
	case 1:
		return nil, musicmax.ErrBadRequest
	default:
		return nil, musicmax.ErrInternalServerError
	}
}
