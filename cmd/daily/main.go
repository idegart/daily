package main

import (
	"bot/internal/airtable"
	"bot/internal/database"
	baseLogger "bot/internal/logger"
	"bot/internal/model"
	"bot/internal/server"
	"bot/internal/slack"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	slackGo "github.com/slack-go/slack"
	"log"
	"os"
)

type App struct {
	logger   *logrus.Logger
	server   *server.Server
	database database.Database
	airtable *airtable.Airtable
	slack    *slack.Slack

	cron *cron.Cron

	airtableUsers []airtable.User
	slackUsers    []slackGo.User
	users         []model.User
	slackProjects []slackGo.Channel
}

var app App

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	logger, err := baseLogger.NewLogger(baseLogger.NewConfig(os.Getenv("LOG_LEVEL")))

	if err != nil {
		log.Fatal(err)
	}

	app.logger = logger
}

func main() {
	if err := app.configure(); err != nil {
		app.logger.Fatal(err)
	}

	defer app.close()

	app.cron = cron.New()

	if _, err := app.cron.AddFunc("0 8 * * *", func() {
		go app.sendInitialMessages()
	}); err != nil {
		app.logger.Fatal(err)
	}

	if _, err := app.cron.AddFunc("0 10 * * *", func() {
		go app.sendReports()
	}); err != nil {
		app.logger.Fatal(err)
	}

	app.cron.Start()

	if err := app.server.Start(); err != nil {
		app.logger.Fatal(err)
	}
}

func (a *App) close() {
	if err := app.database.Close(); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) FindUserBySlackId(slackId string) *model.User  {
	for i := range a.users {
		if a.users[i].SlackId == slackId {
			return &a.users[i]
		}
	}

	return nil
}

func (a App) GetUsersBySlackUsers(slackUsers []slackGo.User) []model.User {
	var users []model.User

	for i := range a.users {
		for j := range slackUsers {
			if a.users[i].SlackId == slackUsers[j].ID {
				users = append(users, a.users[i])
			}
		}
	}

	return users
}

