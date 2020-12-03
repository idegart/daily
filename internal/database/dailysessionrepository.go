package database

import (
	"SlackBot/internal/models"
	"database/sql"
	"time"
)

type DailySessionRepository struct {
	database *Database
}

func (r *DailySessionRepository) Create(session *models.DailySession) error {
	if err := r.database.db.QueryRow(
		"INSERT INTO daily_session (date) VALUES ($1) RETURNING id",
		session.Date,
	).Scan(&session.Id); err != nil {
		return err
	}

	return nil
}

func (r *DailySessionRepository) FindByDate(date time.Time) (*models.DailySession, error) {
	var dailySession models.DailySession

	if err := r.database.db.Get(
		&dailySession,
		"SELECT * FROM daily_session WHERE date=$1",
		date,
	); err != nil {
		return nil, err
	}

	return &dailySession, nil
}

func (r *DailySessionRepository) FindOrCreateByDate(date time.Time) (*models.DailySession, error) {
	session, err := r.FindByDate(date)

	if err != nil {
		if err == sql.ErrNoRows {
			session = &models.DailySession{
				Date: date,
			}
			if err := r.Create(session); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return session, nil
}
