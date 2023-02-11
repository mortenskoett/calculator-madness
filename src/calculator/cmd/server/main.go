package main

import (
	"calculator/pkg/calc"
	"log"
)

func main() {
	log.Println("starting calculator GRPC protobuf service")

	serverConfig := calc.CalcServerConfig{Address: ":8000"}
	calcServer, err := calc.NewGRPCServer(serverConfig)
	if err != nil {
		log.Fatalf("failed to create calc server: %v", err)
	}

	if err := calcServer.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
