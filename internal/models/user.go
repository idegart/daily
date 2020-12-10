package models

import "time"

type User struct {
	Id         int       `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	AirtableId int       `json:"airtable_id" db:"airtable_id"`
	SlackId    string    `json:"slack_id" db:"slack_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
