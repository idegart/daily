package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int           `json:"id" db:"id"`
	Email      string        `json:"email" db:"email"`
	Name       string        `json:"name" db:"name"`
	AirtableId sql.NullInt64 `json:"airtable_id" db:"airtable_id"`
	SlackId    string        `json:"slack_id" db:"slack_id"`
	Emoji      string        `json:"emoji" db:"emoji"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at"`
}
