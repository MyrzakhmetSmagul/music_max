package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/MyrzakhmetSmagul/music_max/internal/model"
	"github.com/MyrzakhmetSmagul/music_max/internal/repository"
)

type SongService struct {
	repo   repository.Repository
	lyrics Lyrics
}

func NewSongService(repo repository.Repository, lyrics Lyrics) Song {
	return &SongService{
		repo:   repo,
		lyrics: lyrics,
	}
}

func (s *SongService) GetSongs(filters map[string]string, page int, limit int) (*model.SongsResponse, error) {
	slog.Debug("SongService.GetSong")
	temp, err := s.repo.GetSongs(filters, page, limit)
	if err != nil {
		err = fmt.Errorf("service.GetSongs error during get songs from repository: %w", err)
		return nil, err
	}

	slog.Debug("songs from db", slog.Any("songs", temp))
	slog.Debug("length of temp", slog.Any("len", len(temp)))

	songs := make([]model.SongResponse, len(temp))
	for i, v := range temp {
		release := ""
		if v.Release != nil {
			release = v.Release.Format("02.01.2006")
		}
		song := new(model.SongResponse)
		song.Id = v.Id
		song.Group = v.Group
		song.Name = v.Name
		song.Release = release
		song.Link = v.Link
		songs[i] = *song
	}
	resp := new(model.SongsResponse)
	resp.Songs = songs
	resp.Page = page
	resp.Limit = limit
	slog.Debug("model.SongsResponse response", slog.Any("response", resp))
	return resp, nil
}

func (s *SongService) CreateSong(songReq *model.SongRequest) error {
	slog.Debug("SongService.CreateSong")
	lyricsResp, err := s.lyrics.GetLyrics(songReq)
	if err != nil {
		err = fmt.Errorf("SongService.CreateSong error get lyrics:\n%w", err)
		return err
	}

	release, err := time.Parse("02.01.2006", *lyricsResp.Release)
	if err != nil {
		err = fmt.Errorf("SongService.CreateSong parse time error:\n%w\n%w", err, model.ErrBadRequest)
		return err
	}
	song := &model.Song{
		Group:   songReq.Group,
		Name:    songReq.Name,
		Release: &release,
		Link:    lyricsResp.Link,
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		err = fmt.Errorf("SongService.CreateSong error begin tx:\n%w", err)
		return err
	}

	err = s.repo.CreateSong(tx, song)
	if err != nil {
		err = fmt.Errorf("service.CreateSong error creating songs in repository: %w", err)
		return err
	}

	lyrics := new(model.Lyrics)
	lyrics.Text = lyricsResp.Text
	lyrics.SongId = song.Id
	err = s.repo.CreateLyrics(tx, lyrics)
	if err != nil {
		err = fmt.Errorf("service.CreateSong error creating lyrics in repository: %w", err)
		return err
	}

	return tx.Commit()
}

func (s *SongService) GetLyrics(id string, page, limit int) (*model.LyricsResponse, error) {
	lyrics, err := s.repo.GetLyrics(id)
	if err != nil {
		err = fmt.Errorf("SongService.DeleteSong error deleting song from repo:\n%w", err)
		return nil, err
	}

	temp := strings.Split(lyrics.Text, "\n\n")
	lyricsResp := new(model.LyricsResponse)
	lyricsResp.Group = lyrics.Group
	lyricsResp.Song = lyrics.Song
	lyricsResp.Page = page
	lyricsResp.Limit = limit
	verses := make([]string, 0)
	for i := (page - 1) * limit; i < len(temp) && i-(page-1)*limit < limit; i++ {
		verses = append(verses, temp[i])
	}
	lyricsResp.Text = verses
	return lyricsResp, nil
}

func (s *SongService) DeleteSong(id string) error {
	err := s.repo.DeleteSong(id)
	if err != nil {
		err = fmt.Errorf("SongService.DeleteSong error deleting song from repo:\n%w", err)
		return err
	}

	return nil
}

func (s *SongService) UpdateSong(id string, songReq *model.SongPatchRequest) error {
	song := new(model.Song)
	song.Id = id
	updateSong := false
	updateText := false
	if songReq.Name != nil && *songReq.Name != "" {
		song.Name = *songReq.Name
		updateSong = true
	}

	if songReq.Group != nil && *songReq.Group != "" {
		song.Group = *songReq.Group
		updateSong = true
	}

	if songReq.Link != nil && *songReq.Link != "" {
		song.Link = *songReq.Link
		updateSong = true
	}

	if songReq.Release != nil {
		release, parseErr := time.Parse("02.01.2006", *songReq.Release)
		if parseErr != nil {
			parseErr = fmt.Errorf("SongService.UpdateSong parse time error:\n%w\n%w", parseErr, model.ErrBadRequest)
			return parseErr
		}

		song.Release = &release
		updateSong = true
	}

	if songReq.Text != nil {
		song.Lyrics = new(model.Lyrics)
		song.Lyrics.SongId = id
		song.Lyrics.Text = *songReq.Text
		updateText = true
	}

	if !updateSong && !updateText {
		return fmt.Errorf("SongService.UpdateSong %w", model.ErrBadRequest)
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		err = fmt.Errorf("SongService.UpdateSong error begin tx:\n%w", err)
		return err
	}

	if updateSong {
		err = s.repo.UpdateSong(tx, song)
		if err != nil {
			err = fmt.Errorf("SongService.UpdateSong error updating song from repo:\n%w", err)
			return err
		}
	}

	if updateText {
		err = s.repo.UpdateLyrics(tx, song.Id, song.Lyrics)
		if err != nil {
			err = fmt.Errorf("SongService.UpdateSong error updating song from repo:\n%w", err)
			return err
		}
	}

	return tx.Commit()
}
