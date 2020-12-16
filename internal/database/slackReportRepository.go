package database

import (
	"SlackBot/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/slack-go/slack"
	"time"
)

type SlackReportRepository struct {
	db *sqlx.DB
}

func (r *SlackReportRepository) Create(d *models.SlackReport) error {
	if err := r.db.QueryRow(
		"INSERT INTO slack_reports (slack_channel_id, date, ts) VALUES ($1, $2, $3) RETURNING id",
		d.SlackChannelId,
		d.Date,
		d.Ts,
	).Scan(&d.Id); err != nil {
		return err
	}

	return nil
}

func (r *SlackReportRepository) Update(d *models.SlackReport) error {
	if _, err := r.db.NamedExec(
		"UPDATE daily_reports SET updated_at = now() WHERE id=:id",
		d,
	); err != nil {
		return err
	}

	return nil
}

func (r *SlackReportRepository) FindByDateAndSlackChannel(date time.Time, slackChannelId string) (*models.SlackReport, error) {
	var report models.SlackReport

	if err := r.db.Get(
		&report,
		"SELECT * FROM slack_reports WHERE date=$1 and slack_channel_id=$2",
		date,
		slackChannelId,
	); err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *SlackReportRepository) FindOrCreateByDateAndSlackChannel(slackReport *models.SlackReport) error {
	storedReport, err := r.FindByDateAndSlackChannel(slackReport.Date, slackReport.SlackChannelId)

	if err == nil {
		slackReport.Id = storedReport.Id
		return nil
	}

	if err == sql.ErrNoRows {
		return r.Create(slackReport)
	}

	return err
}

func (r *SlackReportRepository) FindAllByDateAndChannels(date time.Time, users []slack.Channel) ([]models.SlackReport, error) {
	var reports []models.SlackReport

	var ids []string

	for _, u := range users {
		ids = append(ids, u.ID)
	}

	query, args, err := sqlx.In("SELECT * FROM slack_reports where date = ? and slack_channel_id in (?)", date,  ids)

	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	if err := r.db.Select(&reports, query, args...); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *SlackReportRepository) DateReportsExists(date time.Time) (bool, error) {
	var result models.SlackReport

	if err := r.db.Get(&result, "SELECT * FROM slack_reports WHERE date=$1 LIMIT 1", date); err != nil {
		return false, err
	}

	return true, nil
}