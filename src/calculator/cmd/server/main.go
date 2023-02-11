package main

import (
	"calculator/pkg/calc"
	"flag"
	"log"
	"os"

	"github.com/nsqio/go-nsq"
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

	// Instantiate the NSQ producer.
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(*nsqAddr, nsqConfig)
	if err != nil {
		log.Fatal("failed to create nsq producer: ", err)
	}

	log.Println("nsq producer client created successfully")

	messageBody := []byte("hello world")
	topicName := "morten-topic"

	// Synchronously publish a single message to the specified topic. Messages can also be sent asynchronously and/or in batches.
	err = producer.Publish(topicName, messageBody)
	if err != nil {
		log.Fatal("failed to publish to nsq: ", err)
	}
	defer producer.Stop()

	// // Start serving GRPC endpoint
	serverConfig := calc.CalcServerConfig{Port: *calcServerPort}
	calcServer, err := calc.NewGRPCServer(serverConfig)
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
