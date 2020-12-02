package api

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Router struct {
	server *Server
	router *mux.Router
}

func NewRouter(server *Server) *Router {
	router := &Router{
		server: server,
		router: mux.NewRouter(),
	}

	router.setupRoutes()

	return router
}

func (router *Router) setupRoutes() {
	router.router.HandleFunc("/", router.handleHello())
}

func (router *Router) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		router.server.logger.Info("Handle hello")
		io.WriteString(w, "Hello world")
	}
}
