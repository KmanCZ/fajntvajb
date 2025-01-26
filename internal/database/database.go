package database

import (
	"database/sql"
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"fmt"
	_ "github.com/lib/pq"
	"os"

	"fajntvajb/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type DB struct {
	Users *repository.Users
	Vajbs *repository.Vajbs
}

func migrate(db *sql.DB) error {
	log := logger.Get()
	goose.SetBaseFS(files.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error().Err(err).Msg("Failed to set dialect")
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Error().Err(err).Msg("Failed to run migrations")
	}

	return nil
}

func getConnectionString() (string, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if user == "" {
		return "", fmt.Errorf("DB_USER is not set")
	}
	if password == "" {
		return "", fmt.Errorf("DB_PASSWORD is not set")
	}
	if dbname == "" {
		return "", fmt.Errorf("DB_NAME is not set")
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode), nil
}

func New() (*DB, error) {
	log := logger.Get()

	connectionString, err := getConnectionString()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get connection string")
		return nil, err
	}

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	if err := migrate(db.DB); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
		return nil, err
	}

	res := DB{
		Users: repository.NewUsers(db),
		Vajbs: repository.NewVajbs(db),
	}

	return &res, nil
}
