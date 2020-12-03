package models

import "time"

type User struct {
	Id        int       `json:"id" db:"id"`
	SlackId   string    `json:"slack_id" db:"slack_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
