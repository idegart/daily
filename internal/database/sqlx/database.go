package sqlx

import (
	"bot/internal/config"
	"bot/internal/database"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Database struct {
	config                *config.DB
	db                    *sqlx.DB
	logger                *logrus.Logger
	userRepository        *UserRepository
	dailyReportRepository *DailyReportRepository
	slackReportRepository *SlackReportRepository
}

func New(config *config.DB, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Open() error {
	d.logger.Info("Open DB connection")

	db, err := sqlx.Connect(d.config.Driver, d.config.Url)

	if err != nil {
		return err
	}

	d.db = db

	d.logger.Info("Check for DB connection")

	if err := db.Ping(); err != nil {
		return err
	}

	d.logger.Info("Start migrations")

	if err := d.migrate(); err != nil {
		return err
	}

	return nil
}

func (d *Database) migrate() error {
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			d.logger.Info("Nothing to migrate")
			return nil
		}

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
