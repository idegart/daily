package daily

import (
	"bot/internal/database"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	"bot/internal/model"
	"bot/internal/server"
	"github.com/sirupsen/logrus"
)

const (
	SIDailyReportCallbackStart  = "daily_report_start"
	SIDailyReportCallbackFinish = "daily_report_finish"

	SIDailyReportDone    = "Done"
	SIDailyReportWillDo  = "WillDo"
	SIDailyReportBlocker = "Blocker"
)

type Daily struct {
	logger   *logrus.Logger
	server   *server.Server
	database database.Database
	airtable *airtable.Airtable
	slack    *slack.Slack

	projects     []model.Project
	users        []model.User
	absentUsers  []model.AbsentUser

	projectsToReport chan model.Project
	usersToInitiate  chan model.User
}

func NewDaily(
	logger *logrus.Logger,
	server *server.Server,
	database database.Database,
	airtable *airtable.Airtable,
	slack *slack.Slack,
) *Daily {
	return &Daily{
		logger:           logger,
		server:           server,
		database:         database,
		airtable:         airtable,
		slack:            slack,
		projectsToReport: make(chan model.Project),
		usersToInitiate:  make(chan model.User),
	}
}

func (d *Daily) Serve() error {
	go d.slack.StartSendingMessages()

	go d.startSendingInitiations()
	go d.startSendingReports()

	d.configureServer()

	if err := d.Init(); err != nil {
		return err
	}

	if err := d.server.Start(); err != nil {
		return err
	}

	return nil
}

func (d *Daily) GetUserBySlackId(slackID string) *model.User {
	for i := range d.users {
		if d.users[i].SlackId == slackID {
			return &d.users[i]
		}
	}

	return nil
}

