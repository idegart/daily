package sqlx

import (
	"bot/internal/external/airtable"
	"bot/internal/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type AbsentUserRepository struct {
	db *sqlx.DB
}

func (r *AbsentUserRepository) Create(absentUser *model.AbsentUser) error {
	return r.db.QueryRow(
		"INSERT INTO absent_users (user_id, date) VALUES ($1, $2) RETURNING id",
		absentUser.UserId,
		absentUser.Date,
	).Scan(&absentUser.Id)
}

func (r *AbsentUserRepository) GetAllByDate(time time.Time) ([]model.AbsentUser, error) {
	var absentUsers []model.AbsentUser

	query, args, err := sqlx.In("SELECT * FROM absent_users where date = ?", time)

	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	if err := r.db.Select(&absentUsers, query, args...); err != nil {
		return nil, err
	}

	return absentUsers, nil
}

func (r *AbsentUserRepository) GenerateFromAirtableForDate(airUsers []airtable.AbsentUser, users []model.User, date time.Time) ([]model.AbsentUser, error) {
	absentUsers, err := r.GetAllByDate(date)

	if err != nil {
		return nil, err
	}

	LOOP:
	for _, airUser := range airUsers {
		for _, user := range users {
			if airUser.Email() != user.Email {
				continue
			}

			for _, absentUser := range absentUsers {
				if absentUser.UserId == user.Id {
					continue LOOP
				}
			}

			absentUser := &model.AbsentUser{
				UserId: user.Id,
				Date: date,
			}
			r.Create(absentUser)

			absentUsers = append(absentUsers, *absentUser)
		}
	}

	return absentUsers, nil
}
