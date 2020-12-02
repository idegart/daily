package api

import (
	"os"
)

type Config struct {
	BindAddr string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":" + os.Getenv("BIND_INTERNAL_ADDR"),
	}
}