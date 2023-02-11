// Contains the grpc server and endpoint of the calculation server
package calc

import (
	"calculator/api/pb"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CalculationServer interface {
	Serve() error
	Run(context context.Context, calculationRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error)
}

type CalcServerConfig struct {
	Address string
}

// calculationServer implements CalculationServiceServer interface
type calculationServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
	Server                                   *grpc.Server
	Listener                                 net.Listener
}

func NewGRPCServer(config CalcServerConfig) (CalculationServer, error) {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()

	calcServer := &calculationServer{
		UnimplementedCalculationServiceServer: pb.UnimplementedCalculationServiceServer{},
		Server:                                server,
		Listener:                              listener,
	}

	pb.RegisterCalculationServiceServer(server, calcServer)

	// Setup reflection to be able to work with grpcurl
	reflection.Register(server)

	return calcServer, nil
}

func (s *calculationServer) Serve() error {
	return s.Server.Serve(s.Listener)
}

/* GRPC Protobuf end points */

func (s *calculationServer) Run(context context.Context, calculationRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	log.Println("Run called with:", calculationRequest.Equation)
	eq := Equation{Value: calculationRequest.Equation}
	result, err := Solve(eq)
	if err != nil {
		return nil, fmt.Errorf("failed to solve equation: %w", err)
	}
	return &pb.RunCalculationResponse{
		Result: result,
	}, nil
}
