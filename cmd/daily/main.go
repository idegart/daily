package main

import (
	"bot/internal/database"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	baseLogger "bot/internal/logger"
	"bot/internal/model"
	"bot/internal/server"
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

// https://daily-bot.proscom.tech/callback/interactive
type App struct {
	logger   *logrus.Logger
	server   *server.Server
	database database.Database
	airtable *airtable.Airtable
	slack    *slack.Slack

	teamUsers         []airtable.User
	teamProjects      []airtable.Project
	users             []model.User
	infographicsUsers []airtable.User
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

	app.SendReports()

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
	if err := a.PrepareTeam(); err != nil {
		return nil, err
	}

	for _, user := range a.users {
		go func(user model.User) {
			if err := a.SendSlackInitialMessageToUser(user); err != nil {
				a.logger.Error(err)
			}
		}(user)
	}

	return a.users, nil
}

func (a *App) GetUserBySlackId(slackId string) *model.User {
	for i := range a.users {
		if a.users[i].SlackId == slackId {
			return &a.users[i]
		}
	}

	return nil
}

func (a *App) GetUsersBySlackUsersId(slackUsersId []string) []model.User {
	var users []model.User

	for _, slackUserId := range slackUsersId {
		user := a.GetUserBySlackId(slackUserId)
		if user != nil {
			users = append(users, *user)
		}
	}

	return users
}

func (a *App) SendReports() error {
	if err := a.PrepareTeam(); err != nil {
		return err
	}

	_, err := a.database.DailyReport().GetByDate(time.Now())

	if err != nil {
		return err
	}

	for _, project := range a.teamProjects {
		go a.SendReportToProject(project)
	}

	return nil
}

func (a *App) SendReportToProject(project airtable.Project) {
	if _, _, _, err := a.slack.Client().JoinConversation(project.Fields.SlackID); err != nil {
		a.logger.Error(err)
	}

	projectUsers := a.GetUsersBySlackUsersId(project.GetSlackIds())

	var usersId []int
	var badUsers []model.User

	for _, user := range projectUsers {
		usersId = append(usersId, user.Id)
	}

	reports, err := a.database.DailyReport().FindByUsersAndDate(usersId, time.Now())

	if err != nil {
		a.logger.Error(err)
	}

LOOP:
	for _, user := range projectUsers {
		for _, report := range reports {
			if user.Id == report.UserId {
				continue LOOP
			}
		}
		badUsers = append(badUsers, user)
	}

	slackReport, err := a.database.SlackReport().FindBySlackChannelAndDate(project.Fields.SlackID, time.Now())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.logger.Error(err)
	}

	var ts string

	if slackReport != nil {
		ts = slackReport.Ts
	}

	_, ts, err = a.SendSlackReportToChannel(
		project.Fields.SlackID,
		projectUsers,
		badUsers,
		reports,
		ts,
	)

	if err != nil {
		a.logger.Error(err)
	}

	if err := a.database.SlackReport().UpdateOrCreate(&model.SlackReport{
		Date:           time.Now(),
		SlackChannelId: project.Fields.SlackID,
		Ts:             ts,
	}); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) ResendReportsByUser(user *model.User) {
	for _, project := range a.teamProjects {
		for _, slackId := range project.GetSlackIds() {
			if user.SlackId == slackId {
				report, err := a.database.SlackReport().FindBySlackChannelAndDate(project.Fields.SlackID, time.Now())

				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					a.logger.Error(err)
				}

				if report != nil {
					go a.SendReportToProject(project)
				}
			}
		}
	}
}
