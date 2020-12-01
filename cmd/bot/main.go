package main

import (
	"SlackBot/internal/botserver"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := botserver.NewConfig()

	bot := botserver.New(config)

	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
