package server

import "bot/internal/config"

func NewConfig(addr string) *config.Server {
	return &config.Server{
		BindAddr: addr,
	}
}
