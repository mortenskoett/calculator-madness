package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"shared/queue"
	"viewer/api/pb"
	"viewer/pkg/env"
	"viewer/pkg/http"
	"viewer/pkg/http/websocket"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	nsqClientAddr  = flag.String("nsq-client-addr", env.GetEnvVarOrDefault("NSQ_CLIENT_ADDR", "127.0.0.1:4150"), "Address of nsq client endpoint")
	nsqHttpAddr    = flag.String("nsq-http-addr", env.GetEnvVarOrDefault("NSQ_HTTP_ADDR", "127.0.0.1:4151"), "Address of nsq HTTP endpoint")
	calcServerAddr = flag.String("calculator-addr", env.GetEnvVarOrDefault("CALCULATOR_ADDR", "127.0.0.1:8000"), "Port of Calculator server")
	port           = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
)

func main() {
	flag.Parse()

	// Unique topic used by this service
	calcResultTopic := "web-viewer" + "-" + uuid.NewString()

	// Initiate topic on queue
	queueMaintainer := queue.NewQueueMaintainer(*nsqHttpAddr)
	queueMaintainer.CreateTopic(calcResultTopic)

	// Context used to synchronize shutdown of goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Calculator service grpc client
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(*calcServerAddr, opts...)
	if err != nil {
		log.Fatal("failed to create grpc conn:", err)
	}
	defer conn.Close()
	calcClient := pb.NewCalculationServiceClient(conn)

	// Websocket handling
	wsrouter := websocket.NewEventRouter(calcClient, calcResultTopic)
	wsmanager := websocket.NewManager(wsrouter)

	// Queue clients
	queueConsumer, err := queue.NewNSQConsumer[queue.Enqueable](*nsqClientAddr, calcResultTopic)

	defer func() {
		queueConsumer.Stop()
		// Clean up topic
		queueMaintainer.DeleteTopic(calcResultTopic)
	}()
	go queueConsumer.Start(ctx)

	queueConsumer.SetHandler(func(msg queue.Enqueable) error {
		switch m := msg.(type) {
		case queue.CalcProgressMessage:
			if err := wsmanager.CalcProgressHandler(&m); err != nil {
				return err
			}
		case queue.CalcEndedMessage:
			if err := wsmanager.CalcEndedHandler(&m); err != nil {
				return err
			}
		}
		return nil
	})

	// HTTP server
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
