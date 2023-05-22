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
	nsqAddr        = flag.String("nsq-addr", getEnvVarOrDefault("NSQ_ADDR", "127.0.0.1:4151"), "Address of nsq server with port")
	help           = flag.Bool("help", false, "Show this help")
)

func main() {
	log.Println("starting calculator grpc protobuf service")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	nsqproducer, err := api.NewNSQProducer(*nsqAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nsqproducer.Stop()

	calculationService := calc.NewCalculatorService(nsqproducer)
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
