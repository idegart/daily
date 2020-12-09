package main

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/daily"
	"SlackBot/internal/database"
	"SlackBot/internal/env"
	"SlackBot/internal/logger"
	"SlackBot/internal/server"
	"SlackBot/internal/slackbot"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
)

type App struct {
	logger   *logrus.Logger
	database *database.Database
	server   *server.Server
	slackBot *slackbot.SlackBot
	dailyBot *daily.Bot
	airtable *airtable.Airtable
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var app = &App{}

	app.configureLogger()

	app.logger.Info("Start app")

	app.configureDatabase()

	app.configureServer()

	app.configureSlackBot()

	app.configureAirtable()

	app.configureDailyBot()

	defer app.gracefullyStop()

	if err := app.server.Start(); err != nil {
		app.logger.Error(err)
	}
}

func (a *App) configureLogger() {
	botLogger, err := logger.NewLogger(logger.NewConfig())

	if err != nil {
		log.Fatal(err)
	}

	a.logger = botLogger
}

func (a *App) configureDatabase() {
	a.database = database.NewDatabase(database.NewConfig(), a.logger)

	if err := a.database.Open(); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) configureServer() {
	a.server = server.NewServer(
		server.NewConfig(env.Get("SERVER_BIND_ADDR", "")),
		a.logger,
	)

	a.configureRouter()
}

func (a *App) configureSlackBot() {
	bot, err := slackbot.NewSlackBot(slackbot.NewConfig(), a.logger)

	if err != nil {
		a.logger.Fatal(err)
	}

	a.slackBot = bot
}

func (a *App) configureAirtable() {
	a.airtable = airtable.NewAirtable(airtable.NewConfig())
}

func (a *App) configureDailyBot() {
	a.dailyBot = daily.NewDailyBot(a.logger, a.database, a.slackBot, a.airtable)
}

func (a *App) gracefullyStop() {
	a.logger.Info("Stopping services")

	if err := a.database.Close(); err != nil {
		a.logger.Error(err)
	}
}
