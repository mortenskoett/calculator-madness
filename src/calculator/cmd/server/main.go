package main

import (
	"calculator/pkg/calc"
	"calculator/pkg/queue"
	"flag"
	"log"
	"os"
)

var (
	calcServerPort = flag.String("calc-server-port", getEnvVarOrDefault("SERVER_PORT", ":8000"), "Port of calc grpc server")
	nsqAddr        = flag.String("nsq-addr", getEnvVarOrDefault("NSQ_ADDR", "127.0.0.1:4151"), "Address of nsq server with port")
	help           = flag.Bool("help", false, "Show this help")
)

func main() {
	log.Println("starting calculator GRPC protobuf service")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	// Create queue handler
	producer,err := queue.NewNSQProducer(*nsqAddr)
	if err != nil {
		log.Fatal("failed to create NSQ producer: ", err)
	}
	defer producer.Stop()

	log.Println("nsq producer client created OK")

	// messageBody := []byte("hello world")
	// topicName := "morten-topic"

	// err = producer.Publish(queue.NewCalculationStatusMessage{})
	// if err != nil {
		// log.Fatal("failed to publish to nsq: ", err)
	// }

	// Create GRPC endpoint
	serverConfig := calc.CalcServerConfig{Port: *calcServerPort}
	calcServer, err := calc.NewGRPCServer(serverConfig, producer)
	// TODO: Create server. Then give to calc service ish. Move/register grpc request handles on the
	// server here
	if err != nil {
		log.Fatalf("failed to create calc server: %v", err)
	}

	if err := calcServer.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnvVarOrDefault(envName string, def string) string {
	envvar := os.Getenv(envName)
	if len(envvar) == 0 {
		return def
	}
	return envvar
}
