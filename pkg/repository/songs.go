package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

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

func (s *SongRepository) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *SongRepository) GetSongs(filters map[string]string, page int, limit int) ([]musicmax.Song, error) {
	query := `SELECT "id", "song", "group", "release_date", COALESCE("link",'') FROM "songs"`
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

	if start, ok := filters["startDate"]; ok && start != "" {
		startDate, err := time.Parse("02.01.2006", start)
		if err != nil {
			err = fmt.Errorf("SongRepository.GetSongs parse time error:\n%w\n%w", err, musicmax.ErrBadRequest)
			return nil, err
		}
		between := false
		var endDate time.Time
		if end, ok := filters["endDate"]; ok && end != "" {
			between = true
			endDate, err = time.Parse("02.01.2006", end)
			if err != nil {
				err = fmt.Errorf("SongRepository.GetSongs parse time error:\n%w\n%w", err, musicmax.ErrBadRequest)
				return nil, err
			}
		}
		if between {
			conditions = append(conditions, fmt.Sprintf("\"release_date\" BETWEEN $%d AND $%d", counter, counter+1))
			args = append(args, startDate, endDate)
			counter += 2
		} else {
			conditions = append(conditions, fmt.Sprintf("\"release_date\" >= %d", counter))
			args = append(args, startDate)
			counter++
		}

	} else if end, ok := filters["endDate"]; ok && end != "" {
		endDate, err := time.Parse("02.01.2006", end)
		if err != nil {
			err = fmt.Errorf("SongRepository.GetSongs parse time error:\n%w\n%w", err, musicmax.ErrBadRequest)
			return nil, err
		}

		conditions = append(conditions, fmt.Sprintf("\"release_date\" <= $%d", counter))
		args = append(args, endDate)
		counter++
	}

	if link, ok := filters["link"]; ok && link != "" {
		conditions = append(conditions, fmt.Sprintf("\"link\" = $%d", counter))
		args = append(args, link)
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
		err = rows.Scan(&song.Id, &song.Name, &song.Group, &song.Release, &song.Link)
		if err != nil {
			err = fmt.Errorf("repository.GetAllSongs error during scanning row: %w", err)
			return nil, err
		}

		songs = append(songs, *song)
	}

	slog.Debug("songs selected")

	return songs, nil
}

func (s *SongRepository) CreateSong(tx *sql.Tx, song *musicmax.Song) error {
	query := `INSERT INTO "songs"("song", "group", "release_date", "link") VALUES ($1, $2, $3, $4) RETURNING "id";`
	res := tx.QueryRow(query, song.Name, song.Group, song.Release, song.Link)
	err := res.Scan(&song.Id)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("songs.CreateSong postgres exec error: %w", err)
		return err
	}
	slog.Debug("db create song", slog.Any("song id", song.Id))

	return nil
}

func (s *SongRepository) DeleteSong(id string) error {
	slog.Debug("SongRepository.DeleteSong id", slog.Any("id", id))
	_, err := strconv.Atoi(id)
	if err != nil {
		err := fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return err
	}

	query := `DELETE FROM songs WHERE "id" = $1`
	_, err = s.db.Exec(query, id)
	if err != nil {
		err = fmt.Errorf("SongRepository.DeleteSong error executing query: %w", err)
		return err
	}

	return nil
}

func (s *SongRepository) UpdateSong(tx *sql.Tx, song *musicmax.Song) error {
	slog.Debug("SongRepository.UpdateSong id", slog.Any("id", song.Id))
	_, err := strconv.Atoi(song.Id)
	if err != nil {
		err := fmt.Errorf("%w:\n%w", musicmax.ErrBadRequest, err)
		return err
	}

	if strings.TrimSpace(song.Name) == "" &&
		strings.TrimSpace(song.Group) == "" &&
		strings.TrimSpace(song.Link) == "" &&
		song.Release == nil {
		return fmt.Errorf("SongRepository.UpdateSong empty input: %w", musicmax.ErrBadRequest)
	}

	args := []interface{}{}
	conditions := []string{}
	counter := 1
	query := `UPDATE "songs" SET `
	if strings.TrimSpace(song.Name) != "" {
		conditions = append(conditions, fmt.Sprintf("song= $%d", counter))
		args = append(args, strings.TrimSpace(song.Name))
		counter++
	}

	if strings.TrimSpace(song.Group) != "" {
		conditions = append(conditions, fmt.Sprintf("group= $%d", counter))
		args = append(args, strings.TrimSpace(song.Group))
		counter++
	}

	if strings.TrimSpace(song.Link) != "" {
		conditions = append(conditions, fmt.Sprintf("link= $%d", counter))
		args = append(args, strings.TrimSpace(song.Link))
		counter++
	}
	if song.Release != nil {
		conditions = append(conditions, fmt.Sprintf("release_date= $%d", counter))
		args = append(args, song.Release)
		counter++
	}

	query += strings.Join(conditions, ", ")
	query = fmt.Sprintf("%s WHERE \"id\"=$%d", query, counter)
	args = append(args, song.Id)

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("SongRepository.UpdateSong error executing query: %w", err)
		return err
	}

	return nil
}
