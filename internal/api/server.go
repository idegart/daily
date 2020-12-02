package api

import (
	"SlackBot/internal/database"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config   *Config
	logger   *logrus.Logger
	router   *Router
	database *database.Database
}

func NewServer(config *Config, logger *logrus.Logger, database *database.Database) *Server {
	server := &Server{
		config: config,
		logger: logger,
		database: database,
	}

	server.router = NewRouter(server)

	return server
}

func (s *Server) Start() error {
	if err := s.database.Open(); err != nil {
		return err
	}

	defer s.database.Close()

	return http.ListenAndServe(s.config.BindAddr, s.router.router)
}
