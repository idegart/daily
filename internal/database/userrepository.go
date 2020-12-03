package database

import (
	"SlackBot/internal/models"
	"database/sql"
	"github.com/slack-go/slack"
)

type UserRepository struct {
	database *Database
}

func (u *UserRepository) Create(user *models.User) error {
	if err := u.database.db.QueryRow(
		"INSERT INTO users (slack_id, name) VALUES ($1, $2) RETURNING id",
		user.SlackId,
		user.Name,
	).Scan(&user.Id); err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) FindBySlackId(slackId string) (*models.User, error) {
	var user models.User

	if err := u.database.db.Get(
		&user,
		"SELECT * FROM users WHERE slack_id=$1",
		slackId,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindOrCreateBySlackUser(slackUser *slack.User) (*models.User, error) {
	user, err := u.FindBySlackId(slackUser.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			user = &models.User{
				SlackId: slackUser.ID,
				Name: slackUser.Name,
			}

			if err := u.Create(user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return user, nil
}