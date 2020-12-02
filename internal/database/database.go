package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type Database struct {
	config *Config
	db     *sql.DB
}

func NewDatabase(config *Config) *Database {
	return &Database{
		config: config,
	}
}

func (d Database) Open() error {
	db, err := sql.Open("postgres", d.config.DatabaseUrl)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func (d Database) Close() error {

	d.db.Close()

	return nil
}