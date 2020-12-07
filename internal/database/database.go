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

	userRepository             *UserRepository
	dailySessionRepository     *DailySessionRepository
	userDailySessionRepository *UserDailySessionRepository
}

func NewDatabase(config *Config, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Open() error {
	d.logger.Info("Open DB connection: ", d.config.DatabaseUrl)

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

func (d *Database) DailySession() *DailySessionRepository {
	if d.dailySessionRepository == nil {
		d.dailySessionRepository = &DailySessionRepository{
			database: d,
		}
	}

	return d.dailySessionRepository
}

func (d *Database) User() *UserRepository {
	if d.userRepository == nil {
		d.userRepository = &UserRepository{
			database: d,
		}
	}

	return d.userRepository
}

func (d *Database) UserDailySession() *UserDailySessionRepository {
	if d.userDailySessionRepository == nil {
		d.userDailySessionRepository = &UserDailySessionRepository{
			database: d,
		}
	}

	return d.userDailySessionRepository
}
