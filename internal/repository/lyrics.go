package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/MyrzakhmetSmagul/music_max/internal/model"
)

func (s *SongRepository) CreateLyrics(tx *sql.Tx, lyrics *model.Lyrics) error {
	query := `INSERT INTO "lyrics"("text", "song_id") VALUES($1, $2) RETURNING id;`
	row := tx.QueryRow(query, lyrics.Text, lyrics.SongId)
	err := row.Scan(&lyrics.Id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("songs.CreateLyrics: failed to rollback transaction after error: %v, original error: %w", rbErr, err)
			return err
		}

		err = fmt.Errorf("SongRepository.CreateLyrics error scaning query %w", err)
		return err
	}

	return nil
}

func (s *SongRepository) GetLyrics(songId string) (*model.Lyrics, error) {
	slog.Debug("SongRepository.GetLyrics id", slog.Any("id", songId))
	_, err := strconv.Atoi(songId)
	if err != nil {
		err := fmt.Errorf("%w:\n%w", model.ErrBadRequest, err)
		return nil, err
	}

	query := `SELECT COALESCE("l"."text", ''), "s"."song", "s"."group" FROM "lyrics" "l" INNER JOIN "songs" "s" ON "l"."song_id" = "s"."id" WHERE "s"."id" = $1`

	row := s.db.QueryRow(query, songId)

	lyrics := new(model.Lyrics)
	err = row.Scan(&lyrics.Text, &lyrics.Song, &lyrics.Group)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("repository.GetAllSongs error scanning query result: %w", err)
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("%w:\n%w", model.ErrBadRequest, err)
		return nil, err
	}

	lyrics.SongId = songId
	return lyrics, nil
}

func (s *SongRepository) UpdateLyrics(tx *sql.Tx, songId string, lyrics *model.Lyrics) error {
	if strings.TrimSpace(lyrics.Text) == "" {
		return fmt.Errorf("SongRepository.UpdateLyrics error:\n%w", model.ErrBadRequest)
	}
	query := `UPDATE lyrics SET "text"=$1 WHERE "song_id"=$2`
	_, err := s.db.Exec(query, strings.TrimSpace(lyrics.Text), lyrics.SongId)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("songs.UpdateLyrics: failed to rollback transaction after error: %v, original error: %w", rbErr, err)
			return err
		}

		err = fmt.Errorf("SongRepository.UpdateLyrics error updating lyrics %w", err)
		return err
	}

	return nil
}
