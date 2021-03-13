package daily

import (
	"bot/internal/model"
	"time"
)

func (d *Daily) Init() error {
	if err := d.initiateUsers(); err != nil {
		return err
	}

	if err := d.initiateAbsentUsers(); err != nil {
		return err
	}

	if err := d.initiateProjects(); err != nil {
		return err
	}

	if err := d.initiateProjectUsers(); err != nil {
		return err
	}

	return nil
}

func (d *Daily) initiateUsers() error {
	d.logger.Info("Initiate users")

	airUsers, err := d.airtable.GetActiveUsers()

	if err != nil {
		return err
	}

	dbUsers, err := d.database.User().GenerateFromAirtable(airUsers)

	if err != nil {
		return err
	}

	d.users = dbUsers

	return nil
}

func (d *Daily) initiateAbsentUsers() error {
	d.logger.Info("Initiate absent users")

	airAbsentUsers, err := d.airtable.GetAbsentUsersForDate(time.Now())

	if err != nil {
		return err
	}

	dbAbsentUsers, err := d.database.AbsentUser().GenerateFromAirtableForDate(airAbsentUsers, d.users, time.Now())

	if err != nil {
		return err
	}

	for i := range d.users {
		for j := range dbAbsentUsers {
			if dbAbsentUsers[j].UserId == d.users[i].Id {
				dbAbsentUsers[j].User = &d.users[i]
			}
		}
	}

	d.absentUsers = dbAbsentUsers

	return nil
}

func (d *Daily) initiateProjects() error {
	d.logger.Info("Initiate projects")

	airProjects, err := d.airtable.GetActiveProjects()

	if err != nil {
		return err
	}

	dbProjects, err := d.database.Project().GenerateFromAirtable(airProjects)

	if err != nil {
		return err
	}

	d.projects = dbProjects

	for _, airProject := range airProjects {
		for i := range dbProjects {
			if airProject.Fields.ID == dbProjects[i].AirtableId {
				for _, slackID := range airProject.GetSlackIds() {
					user := d.getUserBySlackId(slackID)
					if user != nil {
						dbProjects[i].Users = append(dbProjects[i].Users, *user)
					}
				}
			}
		}
	}

	return nil
}

func (d *Daily) initiateProjectUsers() error {
	d.logger.Info("Initiate project users")

	for i := range d.projects {
		if err := d.database.Project().SyncUsers(d.projects[i], d.projects[i].Users); err != nil {
			return err
		}
	}

	return nil
}

func (d *Daily) getUserBySlackId(slackID string) *model.User  {
	for _, user := range d.users {
		if user.SlackId == slackID {
			return &user
		}
	}

	return nil
}