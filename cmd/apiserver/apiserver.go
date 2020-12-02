package main

import (
	"SlackBot/internal/database"
	"SlackBot/internal/logger"
	"SlackBot/internal/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	logger   *logrus.Logger
	database *database.Database
	server   *server.Server
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

	app.server = server.NewServer(server.NewConfig(), appLogger)

	app.setupRoutes()

	if err := app.server.Start(); err != nil {
		log.Fatal(err)
	}
}

func (app *App) setupRoutes() {
	app.server.Router.HandleFunc("/", app.handleHello())
}

func (app *App) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Handle hello")
		io.WriteString(w, "Hello world")
	}
}
