package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

func (s *SongRepository) CreateLyrics(tx *sql.Tx, lyrics *musicmax.Lyrics) error {
	query := `INSERT INTO "lyrics"("text", "song_id") VALUES($1, $2) RETURNING id;`
	row := tx.QueryRow(query, lyrics.Text, lyrics.SongId)
	err := row.Scan(&lyrics.Id)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("SongRepository.CreateLyrics error scaning query %w", err)
		return err
	}

	return nil
}

func (s *SongRepository) GetLyrics(songId string) (*musicmax.Lyrics, error) {
	slog.Debug("SongRepository.GetLyrics id", slog.Any("id", songId))
	_, err := strconv.Atoi(songId)
	if err != nil {
		err := fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return nil, err
	}

	query := `SELECT COALESCE("l"."text", ''), "s"."song", "s"."group" FROM "lyrics" "l" INNER JOIN "songs" "s" ON "l"."song_id" = "s"."id" WHERE "s"."id" = $1`

	row := s.db.QueryRow(query, songId)

	lyrics := new(musicmax.Lyrics)
	err = row.Scan(&lyrics.Text, &lyrics.Song, &lyrics.Group)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("repository.GetAllSongs error scanning query result: %w", err)
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return nil, err
	}

	lyrics.SongId = songId
	return lyrics, nil
}

func (s *SongRepository) UpdateLyrics(tx *sql.Tx, songId string, lyrics *musicmax.Lyrics) error {
	if strings.TrimSpace(lyrics.Text) == "" {
		return fmt.Errorf("SongRepository.UpdateLyrics error:\n%w", musicmax.ErrBadRequest)
	}
	query := `UPDATE lyrics SET "text"=$1 WHERE "song_id"=$2`
	_, err := s.db.Exec(query, strings.TrimSpace(lyrics.Text), lyrics.SongId)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("SongRepository.UpdateLyrics error updating lyrics %w", err)
		return err
	}

	return nil
}
