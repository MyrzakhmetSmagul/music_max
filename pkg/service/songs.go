package service

import (
	"fmt"
	"log/slog"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	"github.com/MyrzakhmetSmagul/music_max/pkg/repository"
)

type SongService struct {
	repo   repository.Repository
	lyrics Lyrics
}

func NewSongService(repo repository.Repository) Song {
	return &SongService{
		repo: repo,
	}
}

func (s *SongService) GetSongs(filters map[string]string, page int, limit int) (*musicmax.SongsResponse, error) {
	temp, err := s.repo.GetSongs(filters, page, limit)
	if err != nil {
		err = fmt.Errorf("service.GetSongs error during get songs from repository: %w", err)
		return nil, err
	}

	slog.Debug("songs from db", slog.Any("songs", temp))
	slog.Debug("length of temp", slog.Any("len", len(temp)))
	songs := make([]musicmax.SongResponse, len(temp))
	for i, v := range temp {
		song := new(musicmax.SongResponse)
		song.Id = v.Id
		song.Group = v.Group
		song.Name = v.Name
		song.Release = v.Release
		songs[i] = *song
	}
	resp := new(musicmax.SongsResponse)
	resp.Songs = songs
	resp.Page = page
	resp.Limit = limit
	slog.Debug("musicmax.SongsResponse response", slog.Any("response", resp))
	return resp, nil
}

func (s *SongService) CreateSong(songReq *musicmax.SongRequest) error {
	song := &musicmax.Song{
		Group: songReq.Group,
		Name:  songReq.Name,
	}

	err := s.repo.CreateSong(song)
	if err != nil {
		err = fmt.Errorf("service.CreateSong error during get songs from repository: %w", err)
		return err
	}

	return nil
}

func (s *SongService) GetLyrics(song *musicmax.Song) error {
	return nil
}

func (s *SongService) DeleteSong(id string) error {
	err := s.repo.DeleteSong(id)
	if err != nil {
		err = fmt.Errorf("SongService.DeleteSong error deleting song from repo:\n%w", err)
		return err
	}
	return nil
}
