package main

import (
	"calculator/pkg/api"
	"calculator/pkg/calc"
	"flag"
	"log"
	"os"
)

var (
	calcServerPort = flag.String("calc-server-port", getEnvVarOrDefault("SERVER_PORT", "8000"), "Port of calc grpc server")
	nsqClientAddr  = flag.String("nsq-client-addr", getEnvVarOrDefault("NSQ_CLIENT_ADDR", "127.0.0.1:4150"), "Address of nsq client server with port")
	help           = flag.Bool("help", false, "Show this help")
)

func main() {
	log.Println("starting calculator grpc protobuf service")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	queueProducer, err := api.NewQueueProducer(*nsqClientAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer queueProducer.Stop()

	equationProcessor := calc.NewDummyProcessor(2, 100)
	calculationService := calc.NewCalculatorService(equationProcessor, queueProducer)
	calcServer, err := api.NewGRPCServer(*calcServerPort, calculationService)
	if err != nil {
		log.Fatalf("failed to create calc server: %v", err)
	}

	log.Println("serving on port", *calcServerPort)
	if err := calcServer.Serve(); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}

func getEnvVarOrDefault(envName string, def string) string {
	envvar := os.Getenv(envName)
	if len(envvar) == 0 {
		return def
	}
	return envvar
}
