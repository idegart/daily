package sqlx

import (
	"bot/internal/config"
	"fmt"
)

func NewConfig(driver, host, port, username, password, database string) *config.DB {
	return &config.DB{
		Driver: driver,
		Url: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, database,
		),
	}
}
