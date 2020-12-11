package database

import (
	"SlackBot/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) Create(u *models.User) error {
	if err := r.db.QueryRow(
		"INSERT INTO users (email, name, airtable_id, slack_id) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.Name,
		u.AirtableId,
		u.SlackId,
	).Scan(&u.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Get(
		&user,
		"SELECT * FROM users WHERE email=$1",
		email,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmailOrCreate(email string, name string, airtableId int, slackId string) (*models.User, error) {
	user, err := r.FindByEmail(email)

	if err == sql.ErrNoRows {
		err = nil

		user := &models.User{
			Email:      email,
			Name:       name,
			SlackId:    slackId,
			AirtableId: airtableId,
		}

		err = r.Create(user)

		return user, err
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
