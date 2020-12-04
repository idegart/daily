package main

import (
	"SlackBot/internal/daily"
	"SlackBot/internal/database"
	"SlackBot/internal/env"
	"SlackBot/internal/logger"
	"SlackBot/internal/server"
	"SlackBot/internal/slackbot"
	"github.com/sirupsen/logrus"
	"log"
)

type App struct {
	logger   *logrus.Logger
	database *database.Database
	server   *server.Server
	slackBot      *slackbot.SlackBot
	dailyBot *daily.Bot
}

func main() {
	var app = &App{}

	app.configureLogger()

	app.logger.Info("Start app")

	app.configureDatabase()
	defer app.database.Close()

	app.configureServer()

	app.configureSlackBot()

	app.configureDailyBot()

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
		server.NewConfig(env.Get("SERVER_INTERNAL_BIND_ADDR", "")),
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

func (a *App) configureDailyBot() {
	a.dailyBot = daily.NewDailyBot(a.logger, a.database, a.slackBot)
}
