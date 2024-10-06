package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
	_ "github.com/lib/pq"
)

type SongRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SongRepository{
		db: db,
	}
}

func (s *SongRepository) GetSongs(filters map[string]string, page int, limit int) ([]musicmax.Song, error) {
	query := "SELECT \"id\", \"song\", \"group\", \"release_date\" FROM \"songs\""
	args := []interface{}{}
	conditions := []string{}
	counter := 1
	if group, ok := filters["group"]; ok && group != "" {
		conditions = append(conditions, fmt.Sprintf("\"group\" ILIKE $%d", counter))
		args = append(args, group)
		counter++
	}

	if song, ok := filters["song"]; ok && song != "" {
		conditions = append(conditions, fmt.Sprintf("\"song\" ILIKE $%d", counter))
		args = append(args, song)
		counter++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY \"id\" LIMIT $%d OFFSET $%d", counter, counter+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		err = fmt.Errorf("repository.GetAllSongs error during executing query: %w", err)
		return nil, err
	}
	defer rows.Close()

	songs := make([]musicmax.Song, 0)
	for rows.Next() {
		song := new(musicmax.Song)
		err = rows.Scan(&song.Id, &song.Name, &song.Group, &song.Release)
		if err != nil {
			err = fmt.Errorf("repository.GetAllSongs error during scanning row: %w", err)
			return nil, err
		}

		songs = append(songs, *song)
	}

	slog.Debug("songs selected")

	return songs, nil
}

func (s *SongRepository) CreateSong(song *musicmax.Song) error {
	query := `insert into "songs"("song", "group", "release_date") VALUES ($1, $2, $3) returning "id";`
	res := s.db.QueryRow(query, song.Name, song.Group, song.Release)
	err := res.Scan(&song.Id)
	if err != nil {
		err = fmt.Errorf("songs.CreateSong postgres exec error: %w", err)
		return err
	}
	slog.Debug("db create song", slog.Any("song id", song.Id))

	return nil
}
