package models

import "time"

type SlackReport struct {
	Id             int       `json:"id" db:"id"`
	Date           time.Time `json:"date" db:"date"`
	SlackChannelId string    `json:"slack_channel_id" db:"slack_channel_id"`
	Ts             string    `json:"ts" db:"ts"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
