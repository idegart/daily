package main

import (
	"SlackBot/internal/database"
	"SlackBot/internal/env"
	"SlackBot/internal/logger"
	"SlackBot/internal/models"
	"SlackBot/internal/server"
	"SlackBot/internal/slackbot"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

type App struct {
	logger   *logrus.Logger
	database *database.Database
	server   *server.Server
	bot      *slackbot.SlackBot
	session  *models.DailySession
	users    map[string]*models.User
	sessions map[int]*models.UserDailySession
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := App{}

	appLogger, err := logger.NewLogger(logger.NewConfig())

	if err != nil {
		log.Fatal(err)
	}

	app.logger = appLogger

	app.database = database.NewDatabase(database.NewConfig(), app.logger)

	if err := app.database.Open(); err != nil {
		log.Fatal(err)
	}

	defer app.database.Close()

	session, err := app.database.DailySession().FindOrCreateByDate(time.Now())

	if err != nil {
		app.logger.Fatal(err)
	}

	app.session = session

	app.logger.Info("Current session is: ", app.session)

	app.server = server.NewServer(
		server.NewConfig(env.Get("DAILY_BOT_SERVER_BIND_INTERNAL_ADDR", "")),
		appLogger,
	)

	appBot, err := slackbot.NewSlackBot(slackbot.NewConfig(), appLogger)

	if err != nil {
		log.Fatal(err)
	}

	app.bot = appBot

	app.setupRoutes()

	go app.initDailyNotifier()

	if err := app.server.Start(); err != nil {
		log.Fatal(err)
	}
}

func (app *App) setupRoutes() {
	app.server.Router.HandleFunc("/", app.handleHello())
	app.server.Router.HandleFunc("/events-endpoint", app.handleEvents())
}

func (app *App) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Handle hello from daily bot")
		io.WriteString(w, "Hello world from daily bot")
	}
}

//
//func sendDailyInitialMessage(bot *slackbot.SlackBot, user slack.User) {
//	a, b, err := bot.Api.PostMessage(
//		user.ID,
//		slack.MsgOptionText("Hello world 123!", false),
//		//slack.MsgOptionBlocks(
//		//
//		//	slack.NewActionBlock(
//		//		"hello_there1",
//		//		slack.NewButtonBlockElement("test1", "test", slack.NewTextBlockObject(slack.PlainTextType, "One", false, false)),
//		//		slack.NewButtonBlockElement("test2", "test", slack.NewTextBlockObject(slack.PlainTextType, "Two", false, false)),
//		//	),
//		//	slack.NewActionBlock(
//		//		"hello_there2",
//		//		slack.NewButtonBlockElement("test3", "test", slack.NewTextBlockObject(slack.PlainTextType, "Three", false, false)),
//		//		slack.NewButtonBlockElement("test4", "test", slack.NewTextBlockObject(slack.PlainTextType, "Four", false, false)),
//		//	),
//		//),
//	)
//
//	bot.Logger.Infof("%s | %s", a, b)
//
//	if err != nil {
//		bot.Logger.Error(err)
//		return
//	}
//
//	bot.Logger.Infof("%s | %s", a, b)
//
//}
