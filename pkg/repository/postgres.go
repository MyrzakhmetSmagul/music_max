package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		err = fmt.Errorf("NewPostgresDB can not open connection\nerror: %w", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("NewPostgresDB can not ping db\nerror: %w", err)
		return nil, err
	}

	return db, nil
}
