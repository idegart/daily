package main

import (
	"bot/internal/apps/staff"
	"bot/internal/database"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	baseLogger "bot/internal/logger"
	"bot/internal/model"
	"bot/internal/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type App struct {
	logger   *logrus.Logger
	server   *server.Server
	database database.Database
	airtable *airtable.Airtable
	slack    *slack.Slack

	staff *staff.Staff
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

	if err := app.server.Start(); err != nil {
		app.logger.Fatal(err)
	}
}

func (a *App) close() {
	if err := app.database.Close(); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) SendInitialMessages() ([]model.User, error) {
	users, err := a.staff.GetUsers(true)

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == "a.degtyarev@proscom.ru" {
			if err := a.SendSlackInitialMessageToUser(user); err != nil {
				a.logger.Error(err)
			}
		}

	}

	return users, nil
}

//func (a *App) FindUserBySlackId(slackId string) *model.User {
//	for i := range a.users {
//		if a.users[i].SlackId == slackId {
//			return &a.users[i]
//		}
//	}
//
//	return nil
//}
//
//func (a App) GetUsersBySlackUsers(slackUsers []slackGo.User) []model.User {
//	var users []model.User
//
//	for i := range a.users {
//		for j := range slackUsers {
//			if a.users[i].SlackId == slackUsers[j].ID {
//				users = append(users, a.users[i])
//			}
//		}
//	}
//
//	return users
//}
