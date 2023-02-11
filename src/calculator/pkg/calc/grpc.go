// Contains the grpc server and endpoint of the calculation server
package calc

import (
	"calculator/api/pb"
	"calculator/pkg/queue"
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
	Port string
}

// calculationServer implements CalculationServiceServer interface
type calculationServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
	Server                                   *grpc.Server
	Listener                                 net.Listener
	QueueProducer                            queue.QueueProducer
}

func NewGRPCServer(config CalcServerConfig, producer queue.QueueProducer) (CalculationServer, error) {
	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()

	calcServer := &calculationServer{
		UnimplementedCalculationServiceServer: pb.UnimplementedCalculationServiceServer{},
		Server:                                server,
		Listener:                              listener,
		QueueProducer:                         producer,
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

	s.QueueProducer.Publish(queue.NewCalcStartedMessage())

	eq := Equation{Value: calculationRequest.Equation}
	result, err := Solve(eq)
	if err != nil {
		return nil, fmt.Errorf("failed to solve equation: %w", err)
	}
	return &pb.RunCalculationResponse{
		Result: result,
	}, nil
}
