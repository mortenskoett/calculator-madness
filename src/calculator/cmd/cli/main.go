package main

import (
	"calculator/api/pb"
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverPort = flag.String("port", "8000", "Port of server")
	serverAddr = flag.String("addr", "localhost", "Address of server")
	equation   = flag.String("equation", "1+1", "Equation to send")
)

func main() {
	flag.Parse()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	host := *serverAddr + ":" + *serverPort

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatalln("failed to create grpc connection", err)
	}

	defer conn.Close()

	client := pb.NewCalculationServiceClient(conn)

	result, err := client.Run(context.Background(), &pb.RunCalculationRequest{Equation: "1+1"})
	if err != nil {
		log.Println("an error occurred while solving the equation", err)
	}

	fmt.Println(*equation, "=", result.Result)
}
