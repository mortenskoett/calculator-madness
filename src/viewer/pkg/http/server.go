package http

/* Setup of the HTTP server */

import (
	"log"
	"net/http"
	"viewer/pkg/http/socket"
)

type Config struct {
	Port string
}

type server struct {
	mux       *http.ServeMux
	config    *Config
	wsmanager *socket.Manager
}

func NewServer(config *Config, wsmanager *socket.Manager) *server {
	log.Println("starting calculator viewer web service")
	s := &server{
		mux:       http.NewServeMux(),
		config:    config,
		wsmanager: wsmanager,
	}
	s.attachRoutes()
	return s
}

func (s *server) ListenAndServe() {
	log.Println("http server listening at port", s.config.Port)
	log.Fatalln(http.ListenAndServe(":"+s.config.Port, s.mux))
}