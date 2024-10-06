package repository

import (
	"fmt"

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
