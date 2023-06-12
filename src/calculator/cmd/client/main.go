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

	resp, err := client.Run(context.Background(), &pb.RunCalculationRequest{
		ClientId: "calculator-cli-client-id",
		Equation: &pb.Equation{
			Id:    "calculator-cli-equation-id",
			Value: *equation,
		},
	})
	if err != nil {
		log.Fatal("failed to send calculation request:", err)
	}

	if resp.Error != nil {
		log.Fatalf("an error occurred while solving the equation: %+v", resp.Error)
	}

	log.Println("equation successfully sent to calculation backend")
}
