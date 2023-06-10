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
	nsqlookupAddr  = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
	port           = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
	calcServerAddr = flag.String("calculator-addr", env.GetEnvVarOrDefault("CALCULATOR_ADDR", "127.0.0.1:8000"), "Port of Calculator server")
)

func main() {
	flag.Parse()

	// Unique topic used by this service
	calcResultTopic := "web-viewer" + "-" + uuid.NewString()

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

	// NSQ clients
	nsqConsumer, err := queue.NewNSQUniqueConsumer[queue.Enqueable](*nsqlookupAddr, calcResultTopic)
	defer nsqConsumer.Stop()
	go nsqConsumer.Start(ctx)

	nsqConsumer.SetHandler(func(msg queue.Enqueable) error {
		switch m := msg.(type) {
		case queue.CalcProgressMessage:
			if err := wsmanager.NSQCalcProgressHandler(&m); err != nil {
				return err
			}
		case queue.CalcEndedMessage:
			if err := wsmanager.NSQCalcEndedHandler(&m); err != nil {
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
