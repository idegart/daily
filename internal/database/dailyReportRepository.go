package database

import (
	"SlackBot/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type DailyReportRepository struct {
	db *sqlx.DB
}

func (r *DailyReportRepository) Create(d *models.DailyReport) error {
	if err := r.db.QueryRow(
		"INSERT INTO daily_reports (user_id, date, done, will_do, blocker) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		d.UserId,
		d.Date,
		d.Done,
		d.WillDo,
		d.Blocker,
	).Scan(&d.Id); err != nil {
		return err
	}

	return nil
}

func (r *DailyReportRepository) Update(d *models.DailyReport) error {
	if _, err := r.db.NamedExec(
		"UPDATE daily_reports SET done=:done, will_do=:will_do, blocker=:blocker, updated_at = now() WHERE id=:id",
		d,
	); err != nil {
		return err
	}

	return nil
}

func (r *DailyReportRepository) FindByDateAndUser(date time.Time, userId int) (*models.DailyReport, error) {
	var report models.DailyReport

	if err := r.db.Get(
		&report,
		"SELECT * FROM daily_reports WHERE date=$1 and user_id=$2",
		date,
		userId,
	); err != nil {
		return nil, err
	}

	return &report, nil
}
