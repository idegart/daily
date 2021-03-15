package sqlx

import (
	"bot/internal/external/airtable"
	"bot/internal/model"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, name, airtable_id, slack_id) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email,
		user.Name,
		user.AirtableId,
		user.SlackId,
	).Scan(&user.Id)
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.db.NamedExec(
		"UPDATE users SET email=:email, name=:name, airtable_id=:airtable_id, slack_id=:slack_id, updated_at = now() WHERE id=:id",
		user,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	if err := r.db.Get(
		&user,
		"SELECT * FROM users WHERE email=$1",
		email,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateOrCreate(user *model.User) error {
	userModel, err := r.FindByEmail(user.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return r.Create(user)
	}

	if err != nil {
		return err
	}

	user.Id = userModel.Id
	user.Emoji = userModel.Emoji

	return r.Update(user)
}

func (r *UserRepository) GenerateFromAirtable(airUsers []airtable.User) ([]model.User, error) {
	var users []model.User

	for _, airtableUser := range airUsers{
		user := &model.User{
			Name: airtableUser.Fields.Name,
			Email: airtableUser.Fields.Email,
			AirtableId: &airtableUser.Fields.ID,
			SlackId: airtableUser.Fields.SlackUserID,
		}

		if err := r.UpdateOrCreate(user); err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

