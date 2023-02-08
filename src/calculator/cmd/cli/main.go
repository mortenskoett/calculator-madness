package main

import (
	"calculator/api/pb"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverPort = flag.String("port", "8000", "Port of server")
	serverAddr = flag.String("addr", "localhost", "Address of server")
	equation   = flag.String("eq", "1+1", "Equation to send")
	help       = flag.Bool("help", false, "Show this help")
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

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

	result, err := client.Run(context.Background(), &pb.RunCalculationRequest{Equation: *equation})
	if err != nil {
		log.Println("an error occurred while solving the equation", err)
	}

	log.Println(*equation, "=", result.Result)
}
