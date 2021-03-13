package model

import "time"

type Project struct {
	Id             int       `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	AirtableId     string    `json:"airtable_id" db:"airtable_id"`
	SlackId        string    `json:"slack_id" db:"slack_id"`
	IsInfographics bool      `json:"is_infographics" db:"is_infographics"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	Users []User
}
