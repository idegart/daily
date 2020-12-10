package database

import (
	"SlackBot/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) Create(u *models.User) error {
	if err := r.db.QueryRow(
		"INSERT INTO users (email, airtable_id, slack_id) VALUES ($1, $2, $3) RETURNING id",
		u.Email,
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

