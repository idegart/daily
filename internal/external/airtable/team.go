package airtable

import (
	"github.com/brianloveswords/airtable"
	"time"
)

type User struct {
	airtable.Record
	Fields struct {
		ID          int
		Name        string
		Email       string
		Phone       string
		Status      string
		SlackUserID string `json:"Slack User ID"`
	}
}

type AbsentUser struct {
	airtable.Record
	Fields struct {
		Date        string `json:"Date"`
		Email       []string  `json:"Email (from Me In Team)"`
		SlackUserID []string  `json:"Slack User ID (from Me In Team)"`
	}
}

const (
	activeTeamTable = "tblKBZennm0q7jeX5"
	activeTeamView  = "viwV8yNHIA0FGfgSJ"

	absentTeamTable = "tblKM7RONuyOy32Yj"
	absentTeamView  = "viwd79j0syxJCFVgv"
)

func (a *Airtable) GetActiveUsers() ([]User, error) {
	a.logger.Info("Load active airtable users")

	var users []User

	a.client.BaseID = a.team.appID

	usersTable := a.client.Table(activeTeamTable)

	if err := usersTable.List(&users, &airtable.Options{
		View: activeTeamView,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.logger.Info("Total active airtable users: ", len(users))

	return users, nil
}

func (a *Airtable) GetAbsentUsersForDate(t time.Time) ([]AbsentUser, error) {
	users, err := a.GetAbsentUsers()

	if err != nil {
		return nil, err
	}

	var todayAbsent []AbsentUser

	for _, user := range users {
		if user.Date().Equal(t) {
			todayAbsent = append(todayAbsent, user)
		}
	}

	return todayAbsent, nil
}

func (a *Airtable) GetAbsentUsers() ([]AbsentUser, error) {
	a.logger.Info("Load absent airtable users")

	var users []AbsentUser

	a.client.BaseID = a.team.appID

	usersTable := a.client.Table(absentTeamTable)

	if err := usersTable.List(&users, &airtable.Options{
		View: absentTeamView,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.logger.Info("Total absent airtable users: ", len(users))

	return users, nil
}

func (u AbsentUser) Email() string {
	return u.Fields.Email[0]
}

func (u AbsentUser) Date() *time.Time {
	t, err := time.Parse("2006-01-02", u.Fields.Date)

	if err != nil {
		return nil
	}

	return &t
}
