package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

type LyricsService struct {
	address string
	client  *http.Client
}

func NewLyricsService() Lyrics {
	return &LyricsService{
		address: os.Getenv("API_ADDR"),
		client:  new(http.Client),
	}
}

func (l *LyricsService) GetLyrics(songReq *model.SongRequest) (*model.LyricsAPIResponse, error) {
	slog.Debug("LyricsService", slog.Any("address", l.address))

	reqURL, err := url.Parse(l.address)
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error parsing URL:\n%w", err)
		return nil, err
	}
	params := url.Values{}
	params.Add("group", songReq.Group)
	params.Add("song", songReq.Name)

	reqURL.RawQuery = params.Encode()

	resp, err := l.client.Get(reqURL.String())
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error occured during get request:\n%w", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, model.ErrBadRequest
		default:
			err := fmt.Errorf("status code: %d\nerror: %w", resp.StatusCode, model.ErrInternalServerError)
			return nil, err
		}
	}

	response := new(model.LyricsAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error decoding response:\n%w", err)
		return nil, err
	}

	slog.Debug("api response", slog.Any("response", response))
	return response, nil
}
