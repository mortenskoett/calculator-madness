package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"shared/queue"
	"viewer/pkg/env"
	"viewer/pkg/http"
	"viewer/pkg/http/websocket"
)

const (
	serviceNameChannel string = "viewer-web"
)

var (
	nsqlookupAddr = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
	port          = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
)

func main() {
	flag.Parse()

	// Context used to synchronize shutdown of goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup websocket handling
	wsmanager := websocket.NewManager()

	// Start NSQ client
	nsqconsumer, err := queue.NewNSQConsumer(*nsqlookupAddr, queue.CalcStatusTopic, serviceNameChannel)
	if err != nil {
		log.Fatal(err)
	}
	defer nsqconsumer.Stop()
	nsqconsumer.AddCalcProgressHandler(wsmanager.NSQCalcProgressHandler)
	nsqconsumer.AddCalcEndedHandler(wsmanager.NSQCalcEndedHandler)
	go nsqconsumer.Start(ctx)

	// Start HTTP server
	config := http.Config{Port: *port}
	server := http.NewServer(&config, wsmanager)
	go server.ListenAndServe(ctx)

	// Listen for interrupt and stop server using contexts.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
	log.Println("main: received ctrl-c interrupt - shutting down")
	cancel()
}
