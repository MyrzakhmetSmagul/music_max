package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

func (s *SongRepository) CreateLyrics(lyrics *musicmax.Lyrics) error {
	query := `INSERT INTO "lyrics"("text", "song_id") VALUES($1, $2) RETURNING id;`
	row := s.db.QueryRow(query, lyrics.Text, lyrics.SongId)
	err := row.Scan(&lyrics.Id)
	if err != nil {
		err = fmt.Errorf("SongRepository.CreateLyrics error scaning query %w", err)
		return err
	}

	return nil
}

func (s *SongRepository) GetLyrics(id string) (*musicmax.Lyrics, error) {
	slog.Debug("SongRepository.GetLyrics id", slog.Any("id", id))
	_, err := strconv.Atoi(id)
	if err != nil {
		err := fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return nil, err
	}

	query := `SELECT COALESCE("l"."text", ''), "s"."song", "s"."group" FROM "lyrics" "l" INNER JOIN "songs" "s" ON "l"."song_id" = "s"."id" WHERE "s"."id" = $1`

	row := s.db.QueryRow(query, id)

	lyrics := new(musicmax.Lyrics)
	err = row.Scan(&lyrics.Text, &lyrics.Song, &lyrics.Group)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("repository.GetAllSongs error scanning query result: %w", err)
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return nil, err
	}

	lyrics.SongId = id
	return lyrics, nil
}
