package main

import (
	"SlackBot/internal/api"
	"SlackBot/internal/database"
	"SlackBot/internal/logger"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	apiLogger, err := logger.NewLogger(logger.NewConfig())

	if err != nil {
		log.Fatal(err)
	}

	db := database.NewDatabase(database.NewConfig())

	server := api.NewServer(api.NewConfig(), apiLogger, db)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
