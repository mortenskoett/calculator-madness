package main

import (
	"calculator/pkg/calc"
	"flag"
	"log"
	"os"
	"shared/queue"
)

var (
	calcServerPort = flag.String("calc-server-port", getEnvVarOrDefault("SERVER_PORT", "8000"), "Port of calc grpc server")
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

	// Create GRPC endpoint
	serverConfig := calc.CalcServerConfig{Port: *calcServerPort}

	// TODO: Separate creation of server from biz logic. E.g. give server to calc service ish.
	// Move/register grpc request handles on the server here
	calcServer, err := calc.NewGRPCServer(serverConfig, producer)
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
