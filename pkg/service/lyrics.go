package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
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

func (l *LyricsService) GetLyrics(songReq *musicmax.SongRequest) (*musicmax.LyricsAPIResponse, error) {
	slog.Debug("LyricsService", slog.Any("address", l.address))
	jsonData, err := json.Marshal(songReq)
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error marshalling JSON:\n%w", err)
		return nil, err
	}
	slog.Debug("request data", slog.Any("data", songReq))
	req, err := http.NewRequest("GET", l.address, bytes.NewBuffer(jsonData))
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error making GET request:\n%w", err)
		return nil, err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error occured during get request:\n%w", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, musicmax.ErrBadRequest
		default:
			err := fmt.Errorf("status code: %d\nerror: %w", resp.StatusCode, musicmax.ErrInternalServerError)
			return nil, err
		}
	}

	response := new(musicmax.LyricsAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		err := fmt.Errorf("LyricsService.GetLyrics error decoding response:\n%w", err)
		return nil, err
	}

	slog.Debug("api response", slog.Any("response", response))
	return response, nil
}
