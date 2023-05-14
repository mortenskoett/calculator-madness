package http

/* Setup of the HTTP server */

import (
	"context"
	"log"
	"net/http"
	"viewer/pkg/http/websocket"
)

type Config struct {
	Port string
}

type server struct {
	mux       *http.ServeMux
	config    *Config
	wsmanager *websocket.Manager
}

func NewServer(config *Config, wsmanager *websocket.Manager) *server {
	log.Println("starting calculator viewer http server")
	s := &server{
		mux:       http.NewServeMux(),
		config:    config,
		wsmanager: wsmanager,
	}
	s.attachRoutes()
	return s
}

// Blocking call.
func (s *server) ListenAndServe(ctx context.Context) {
		log.Println("http server listening at port", s.config.Port)
	go func() {
		log.Fatalln(http.ListenAndServe(":"+s.config.Port, s.mux))
	}()

	// Handle shutdown using context
	<-ctx.Done()
	log.Println("stopping http server: cancelled by context.")
	return
}
