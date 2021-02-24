package main

import (
	"bot/internal/database/sqlx"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	"bot/internal/server"
	"os"
)

func (a *App) configure() error {
	if err := configureDatabase(a); err != nil {
		return err
	}

	configureAirtable(a)

	configureSlack(a)

	configureServer(a)

	return nil
}

func configureDatabase(a *App) error {
	a.database = sqlx.New(
		sqlx.NewConfig(
			"postgres",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
		),
		app.logger,
	)

	if err := app.database.Open(); err != nil {
		return err
	}

	return nil
}

func configureAirtable(a *App) {
	a.airtable = airtable.NewAirtable(
		airtable.NewConfig(
			os.Getenv("AIRTABLE_API_KEY"),
			os.Getenv("AIRTABLE_APP_ID"),
		),
		a.logger,
	)

	a.airtable.SetupActive(
		airtable.NewActiveConfig(
			os.Getenv("AIRTABLE_ACTIVE_PROJECTS_TABLE"),
			os.Getenv("AIRTABLE_ACTIVE_PROJECTS_VIEW"),
			os.Getenv("AIRTABLE_ACTIVE_TEAM_TABLE"),
			os.Getenv("AIRTABLE_ACTIVE_TEAM_VIEW"),
		),
	)

	a.airtable.SetupInfographics(
		airtable.NewInfographicsConfig(
			os.Getenv("AIRTABLE_ACTIVE_INFOGRAPHICS_TEAM_TABLE"),
			os.Getenv("AIRTABLE_ACTIVE_INFOGRAPHICS_TEAM_VIEW"),
		),
	)
}

func configureSlack(a *App) {
	a.slack = slack.NewSlack(slack.NewConfig(
		os.Getenv("SLACK_API_TOKEN"),
		os.Getenv("SLACK_VERIFICATION_TOKEN"),
		os.Getenv("SLACK_SIGNING_SECRET"),
	), a.logger)
}

func configureServer(a *App) {
	a.server = server.NewServer(server.NewConfig(os.Getenv("DAILY_PORT")), a.logger)

	a.server.Router().HandleFunc("/health", handleHealth(a))
	a.server.Router().HandleFunc("/callback/interactive", handleSlackInteractiveCallback(a))

	a.server.Router().HandleFunc("/start-daily", handleStartDaily(a))
	a.server.Router().HandleFunc("/send-reports", handleSendReports(a))
}