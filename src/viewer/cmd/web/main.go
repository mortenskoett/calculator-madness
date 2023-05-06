package main

import (
	"flag"
	"viewer/pkg/env"
	"viewer/pkg/http"
	"viewer/pkg/http/websocket"
)

const (
	ServiceNameChannel string = "viewer"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
	port           = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
)

func main() {
	flag.Parse()
	wsmanager := websocket.NewManager()
	config := http.Config{Port: *port}
	server := http.NewServer(&config, wsmanager)
	server.ListenAndServe()
}
