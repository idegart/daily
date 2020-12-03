package main

import (
	"SlackBot/internal/models"
	"github.com/slack-go/slack"
	"sync"
)

func (app *App) initDailyNotifier() {
	var wg sync.WaitGroup

	users, err := app.bot.GetActiveUsers()

	app.users = make(map[string]*models.User, len(users))
	app.sessions = make(map[int]*models.UserDailySession, len(users))

	if err != nil {
		app.logger.Error(err)
	}

	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go app.sendInitialMessage(&wg, &users[i])
	}
}

func (app *App) sendInitialMessage(wg *sync.WaitGroup, user *slack.User) {

	realUser, err := app.database.User().FindOrCreateBySlackUser(user)

	app.users[user.ID] = realUser

	app.logger.Info("Test: ", app.users)

	if err != nil {
		app.logger.Error(err)
		return
	}

	session, err := app.database.UserDailySession().FindOrCreateByDailySessionAndUser(app.session, realUser)

	app.sessions[realUser.Id] = session

	if err != nil {
		app.logger.Error(err)
		return
	}

	app.logger.Info("Session: ", session)

	_, _, err = app.bot.Api.PostMessage(
		user.ID,
		slack.MsgOptionText("Привет, настало время чтобы рассказать чем ты занимался вчера и чем планируешь заняться сегодня", false),
	)

	if err != nil {
		app.logger.Errorf("Error when send daily initial message to %s: %s", user.ID, err)
	}

	wg.Done()
}
