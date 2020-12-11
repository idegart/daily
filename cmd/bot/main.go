package main

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/app/slackBot"
	"SlackBot/internal/database"
	"SlackBot/internal/env"
	"SlackBot/internal/logger"
	"SlackBot/internal/server"
	"SlackBot/internal/slack"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	bot := configureSlackBot()

	if err := bot.Serve(); err != nil {
		log.Fatal(err)
	}
}

func configureSlackBot() *slackBot.SlackBot {
	l, err := logger.NewLogger(logger.NewConfig())

	if err != nil {
		log.Fatal(err)
	}

	db := database.NewDatabase(database.NewConfig(), l)

	if err := db.Open(); err != nil {
		log.Fatal(err)
	}

	serv := server.NewServer(server.NewConfig(env.Get("SERVER_BIND_ADDR", "")), l)

	sl := slack.NewSlack(slack.NewConfig(), l)

	air := airtable.NewAirtable(airtable.NewConfig(), l)

	return slackBot.NewSlackBot(l, db, serv, sl, air)
}