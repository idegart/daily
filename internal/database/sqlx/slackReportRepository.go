package sqlx

import (
	"bot/internal/model"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type SlackReportRepository struct {
	db *sqlx.DB
}

func (r *SlackReportRepository) Create(report *model.SlackReport) error {
	return r.db.QueryRow(
		"INSERT INTO daily_slack_reports (slack_channel_id, date, ts) VALUES ($1, $2, $3) RETURNING id",
		report.SlackChannelId,
		report.Date,
		report.Ts,
	).Scan(&report.Id)
}

func (r *SlackReportRepository) Update(report *model.SlackReport) error {
	_, err := r.db.NamedExec(
		"UPDATE daily_slack_reports SET ts=:ts, updated_at = now() WHERE id=:id",
		report,
	)

	return err
}

func (r *SlackReportRepository) GetAllByDate(time time.Time) ([]model.SlackReport, error) {
	var reports []model.SlackReport

	if err := r.db.Select(&reports, "SELECT * FROM daily_slack_reports WHERE date=$1", time); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *SlackReportRepository) FindBySlackChannelAndDate(slackChannelId string, time time.Time) (*model.SlackReport, error) {
	var slackReport model.SlackReport

	if err := r.db.Get(
		&slackReport,
		"SELECT * FROM daily_slack_reports WHERE slack_channel_id=$1 and date=$2",
		slackChannelId,
		time,
	); err != nil {
		return nil, err
	}

	return &slackReport, nil
}

func (r *SlackReportRepository) UpdateOrCreate(report *model.SlackReport) error {
	slackReportModel, err := r.FindBySlackChannelAndDate(report.SlackChannelId, report.Date)

	if errors.Is(err, sql.ErrNoRows) {
		return r.Create(report)
	}

	if err != nil {
		return err
	}

	report.Id = slackReportModel.Id

	return r.Update(report)
}
