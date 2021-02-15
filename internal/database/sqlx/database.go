package sqlx

import (
	"bot/internal/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Database struct {
	config                *Config
	db                    *sqlx.DB
	logger                *logrus.Logger
	userRepository        *UserRepository
	dailyReportRepository *DailyReportRepository
	slackReportRepository *SlackReportRepository
}

func New(config *Config, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Open() error {
	d.logger.Info("Open DB connection")

	db, err := sqlx.Connect(d.config.driver, d.config.url)

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

func (d *Database) User() database.UserRepository {
	if d.userRepository == nil {
		d.userRepository = &UserRepository{
			db: d.db,
		}
	}

	return d.userRepository
}

func (d *Database) DailyReport() database.DailyReportRepository {
	if d.dailyReportRepository == nil {
		d.dailyReportRepository = &DailyReportRepository{
			db: d.db,
		}
	}

	return d.dailyReportRepository
}

func (d *Database) SlackReport() database.SlackReportRepository {
	if d.slackReportRepository == nil {
		d.slackReportRepository = &SlackReportRepository{
			db: d.db,
		}
	}

	return d.slackReportRepository
}
