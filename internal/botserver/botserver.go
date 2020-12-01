package botserver

import (
	"SlackBot/internal/slack"
	"SlackBot/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BotServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
	slack  *slack.Slack
}

func New(config *Config) *BotServer {
	return &BotServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (bot BotServer) Start() error {
	if err := bot.configureLogger(); err != nil {
		return err
	}

	if err := bot.configureStore(); err != nil {
		return err
	}

	if err := bot.configureSlack(); err != nil {
		return err
	}

	bot.configureRouter()

	bot.logger.Info("Starting bot server")

	return http.ListenAndServe(bot.config.BindAddr, bot.router)
}

func (bot *BotServer) configureLogger() error {
	level, err := logrus.ParseLevel(bot.config.LogLevel)

	if err != nil {
		return err
	}

	bot.logger.SetLevel(level)

	return nil
}

func (bot *BotServer) configureStore() error {
	store := store.New(bot.config.Store)

	if err := store.Open(); err != nil {
		return err
	}

	bot.store = store

	return nil
}

func (bot *BotServer) configureSlack() error {
	slack := slack.New(bot.config.Slack)

	if err := slack.SetApi(); err != nil {
		return err
	}

	bot.slack = slack

	return nil
}

func (bot *BotServer) configureRouter() {
	bot.router.HandleFunc("/", bot.handleHello())
}

func (bot *BotServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world")
	}
}
