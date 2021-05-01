package main

import (
	"bot/internal/apps/daily"
	"bot/internal/database/sqlx"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	baseLogger "bot/internal/logger"
	"bot/internal/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	logger, err := baseLogger.NewLogger(baseLogger.NewConfig(os.Getenv("LOG_LEVEL")))

	if err != nil {
		log.Fatal(err)
	}

	if err := configureDaily(logger).Serve(); err != nil {
		logger.Fatal(err)
	}
}

func configureDaily(logger *logrus.Logger) *daily.Daily {
	database := sqlx.New(
		sqlx.NewConfig(
			"postgres",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
		),
		logger,
	)

	if err := database.Open(); err != nil {
		logger.Fatal(err)
	}

	a := airtable.NewAirtable(
		airtable.NewConfig(
			os.Getenv("AIRTABLE_API_KEY"),
		),
		logger,
	)

	a.SetupTeam(os.Getenv("AIRTABLE_TEAM_APP_ID"))
	a.SetupProjects(os.Getenv("AIRTABLE_PROJECTS_APP_ID"))

	s := slack.NewSlack(slack.NewConfig(
		os.Getenv("SLACK_API_TOKEN"),
		os.Getenv("SLACK_VERIFICATION_TOKEN"),
		os.Getenv("SLACK_SIGNING_SECRET"),
	), logger)

	go s.StartSendingMessages()

	serv := server.NewServer(server.NewConfig(os.Getenv("DAILY_PORT")), logger)

	return daily.NewDaily(
		logger,
		serv,
		database,
		a,
		s,
	)
}