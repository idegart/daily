package sqlx

import (
	"bot/internal/model"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type DailyReportRepository struct {
	db *sqlx.DB
}

func (r *DailyReportRepository) Create(report *model.DailyReport) error {
	return r.db.QueryRow(
		"INSERT INTO daily_user_reports (user_id, date, done, will_do, blocker) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		report.UserId,
		report.Date,
		report.Done,
		report.WillDo,
		report.Blocker,
	).Scan(&report.Id)
}

func (r *DailyReportRepository) Update(report *model.DailyReport) error {
	_, err := r.db.NamedExec(
		"UPDATE daily_user_reports SET user_id=:user_id, date=:date, done=:done, will_do=:will_do, blocker=:blocker, updated_at = now() WHERE id=:id",
		report,
	)

	return err
}

func (r *DailyReportRepository) GetByDate(time time.Time) ([]model.DailyReport, error) {
	var reports []model.DailyReport

	query, args, err := sqlx.In("SELECT * FROM daily_user_reports where date = ?", time)

	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	if err := r.db.Select(&reports, query, args...); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *DailyReportRepository) FindByUserAndDate(userId int, time time.Time) (*model.DailyReport, error) {
	var dailyReport model.DailyReport

	if err := r.db.Get(
		&dailyReport,
		"SELECT * FROM daily_user_reports WHERE user_id=$1 and date=$2",
		userId,
		time,
	); err != nil {
		return nil, err
	}

	return &dailyReport, nil
}

func (r *DailyReportRepository) FindByUsersAndDate(usersId []int, time time.Time) ([]model.DailyReport, error) {
	var dailyReports []model.DailyReport

	query, args, err := sqlx.In("SELECT * FROM daily_user_reports WHERE user_id IN (?) AND date=?", usersId, time)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&dailyReports, r.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	return dailyReports, nil
}

func (r *DailyReportRepository) UpdateOrCreate(report *model.DailyReport) error {
	dailyReportModel, err := r.FindByUserAndDate(report.UserId, report.Date)

	if errors.Is(err, sql.ErrNoRows) {
		return r.Create(report)
	}

	if err != nil {
		return err
	}

	report.Id = dailyReportModel.Id

	return r.Update(report)
}

func (r *DailyReportRepository) GetLastUserReport(userId int) (*model.DailyReport, error) {
	var dailyReport model.DailyReport

	if err := r.db.Get(
		&dailyReport,
		"SELECT * FROM daily_user_reports WHERE user_id=$1 and date::date != now()::date order by date desc",
		userId,
	); err != nil {
		return nil, err
	}

	return &dailyReport, nil
}