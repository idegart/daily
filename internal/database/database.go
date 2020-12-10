package database

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Database struct {
	config *Config
	logger *logrus.Logger
	db     *sqlx.DB

	userRepository *UserRepository
	dailyReportRepository *DailyReportRepository
}

func NewDatabase(config *Config, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Open() error {
	d.logger.Info("Open DB connection")

	db, err := sqlx.Connect("postgres", d.config.DatabaseUrl)

	if err != nil {
		return err
	}

	d.db = db

	d.logger.Info("Check for DB connection")

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() error {
	d.logger.Info("Close DB connection")

	if err := d.db.Close(); err != nil {
		return err
	}

	return nil
}

func (d *Database) Query() *sqlx.DB {
	return d.db
}

func (d *Database) UserRepository() *UserRepository {
	if d.userRepository == nil {
		d.userRepository = &UserRepository{
			db: d.db,
		}
	}

	return d.userRepository
}

func (d *Database) DailyReportRepository() *DailyReportRepository {
	if d.dailyReportRepository == nil {
		d.dailyReportRepository = &DailyReportRepository{
			db: d.db,
		}
	}

	return d.dailyReportRepository
}

