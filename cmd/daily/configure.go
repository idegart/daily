package main

import (
	"bot/internal/airtable"
	"bot/internal/database/sqlx"
	"bot/internal/slack"
	"os"
)

func (a *App) configure() error {
	if err := a.configureDatabase(); err != nil {
		return err
	}

	a.configureServer()

	a.configureAirtable()

	a.configureSlack()

	return nil
}

func (a *App) configureDatabase() error {
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

func (a *App) configureAirtable() {
	a.airtable = airtable.NewAirtable(
		airtable.NewConfig(
			os.Getenv("AIRTABLE_API_KEY"),
			os.Getenv("AIRTABLE_BASE_ID"),
			os.Getenv("AIRTABLE_USERS_TABLE"),
		),
		a.logger,
	)
}

func (a *App) configureSlack() {
	a.slack = slack.NewSlack(slack.NewConfig(
		os.Getenv("SLACK_API_TOKEN"),
		os.Getenv("SLACK_VERIFICATION_TOKEN"),
		os.Getenv("SLACK_SIGNING_SECRET"),
	), a.logger)
}
