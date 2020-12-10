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
	slackBot := configureSlackBot()

	if err := slackBot.Serve(); err != nil {
		log.Fatal(err)
	}
}

func configureSlackBot() *slackBot.SlackBot {
	logger, err := logger.NewLogger(logger.NewConfig())

	if err != nil {
		log.Fatal(err)
	}

	database := database.NewDatabase(database.NewConfig(), logger)

	if err := database.Open(); err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(server.NewConfig(env.Get("SERVER_BIND_ADDR", "")), logger)

	slack := slack.NewSlack(slack.NewConfig(), logger)

	airtable := airtable.NewAirtable(airtable.NewConfig(), logger)

	return slackBot.NewSlackBot(logger, database, server, slack, airtable)
}