package main

import (
	"calculator/api/pb"
	"calculator/pkg/calc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// calculationServer implements CalculationServiceServer interface
type calculationServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
}

func (s *calculationServer) Run(context context.Context, calculationRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	eq := calc.Equation{Value: calculationRequest.Equation}
	result, err := calc.Solve(eq)
	if err != nil {
		return nil, fmt.Errorf("Failed to solve equation: %w", err)
	}
	return &pb.RunCalculationResponse{
		Result: result,
	}, nil
}

func main() {
	log.Println("Starting calculator GRPC protobuf service")
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCalculationServiceServer(server, &calculationServer{})

	// setup reflection to be able to work with grpcurl
	reflection.Register(server)

	log.Println("Serving at :8000")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
