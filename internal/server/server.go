package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func NewServer(config *Config, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.logger.WithField("port", s.config.BindAddr).Info("Starting server")
	return http.ListenAndServe(":" + s.config.BindAddr, s.router)
}

func (s *Server) Router() *mux.Router {
	return s.router
}