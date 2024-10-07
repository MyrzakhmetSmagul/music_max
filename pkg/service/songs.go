package service

import (
	"fmt"
	"log/slog"
	"strings"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	"github.com/MyrzakhmetSmagul/music_max/pkg/repository"
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

func (s *SongService) GetSongs(filters map[string]string, page int, limit int) (*musicmax.SongsResponse, error) {
	slog.Debug("SongService.GetSong")
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
		song.Link = v.Link
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
	slog.Debug("SongService.CreateSong")
	lyricsResp, err := s.lyrics.GetLyrics(songReq)
	if err != nil {
		err = fmt.Errorf("SongService.CreateSong error get lyrics:\n%w", err)
		return err
	}

	song := &musicmax.Song{
		Group:   songReq.Group,
		Name:    songReq.Name,
		Release: lyricsResp.Release,
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

	lyrics := new(musicmax.Lyrics)
	lyrics.Text = lyricsResp.Text
	lyrics.SongId = song.Id
	err = s.repo.CreateLyrics(tx, lyrics)
	if err != nil {
		err = fmt.Errorf("service.CreateSong error creating lyrics in repository: %w", err)
		return err
	}

	return tx.Commit()
}

func (s *SongService) GetLyrics(id string, page, limit int) (*musicmax.LyricsResponse, error) {
	lyrics, err := s.repo.GetLyrics(id)
	if err != nil {
		err = fmt.Errorf("SongService.DeleteSong error deleting song from repo:\n%w", err)
		return nil, err
	}

	temp := strings.Split(lyrics.Text, "\n\n")
	lyricsResp := new(musicmax.LyricsResponse)
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

func (s *SongService) UpdateSong(id string, songReq *musicmax.SongPatchRequest) error {
	song := new(musicmax.Song)
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
		song.Release = songReq.Release
		updateSong = true
	}

	if songReq.Text != nil {
		song.Lyrics = new(musicmax.Lyrics)
		song.Lyrics.SongId = id
		song.Lyrics.Text = *songReq.Text
		updateText = true
	}

	if !updateSong && !updateText {
		return fmt.Errorf("SongService.UpdateSong %w", musicmax.ErrBadRequest)
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
