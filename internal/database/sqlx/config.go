package sqlx

import (
	"fmt"
)

type Config struct {
	driver string
	url    string
}

func NewConfig(driver, host, port, username, password, database string) *Config {
	return &Config{
		driver: driver,
		url: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, database,
		),
	}
}
