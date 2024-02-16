package repository

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DBName       string
	SSLMode      string
	MigrationURL string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	var dataSourseName string = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", dataSourseName)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationURL, cfg.DBName, driver)

	if err != nil {
		return nil, err
	}

	m.Up()

	return db, nil
}
