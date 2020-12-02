package database

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseUrl string
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_INTERNAL_HOST"),
			os.Getenv("DB_INTERNAL_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
		),
	}
}