package database

import (
	"SlackBot/internal/models"
	"database/sql"
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

func (r *DailyReportRepository) FindAllByDateAndUsers(date time.Time, users []models.User) ([]models.DailyReport, error) {
	var reports []models.DailyReport

	var ids []int

	for _, u := range users {
		ids = append(ids, u.Id)
	}

	query, args, err := sqlx.In("SELECT * FROM daily_reports where date = ? and user_id in (?)", date,  ids)

	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	if err := r.db.Select(&reports, query, args...); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *DailyReportRepository) CreateOrUpdateByDateAndUser(date time.Time, userId int, done string, willDo string, blocker string) (*models.DailyReport, error) {
	report, err := r.FindByDateAndUser(date, userId)

	if err == sql.ErrNoRows {
		report = &models.DailyReport{
			UserId:  userId,
			Date:    time.Now(),
			Done:    done,
			WillDo:  willDo,
			Blocker: blocker,
		}

		err = r.Create(report)

		if err != nil {
			return nil, err
		}

		return report, err
	}

	if err != nil {
		return nil, err
	}

	report.Done = done
	report.WillDo = willDo
	report.Blocker = blocker

	err = r.Update(report)

	if err != nil {
		return nil, err
	}

	return report, nil
}
