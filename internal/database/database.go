package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)


type Database struct {
	config *Config
	logger *logrus.Logger
	db     *sql.DB
}

func NewDatabase(config *Config, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d Database) Open() error {
	d.logger.Info("Open DB connection")

	db, err := sql.Open("postgres", d.config.DatabaseUrl)

	if err != nil {
		return err
	}

	d.logger.Info("Check for DB connection")

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func (d Database) Close() error {

	d.logger.Info("Close DB connection")

	d.db.Close()

	return nil
}