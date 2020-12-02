package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config *Config
	logger *logrus.Logger
	Router *mux.Router
}

func NewServer(config *Config, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
		Router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.logger.Infof("Staring api server on port %s", s.config.BindAddr)
	return http.ListenAndServe(s.config.BindAddr, s.Router)
}