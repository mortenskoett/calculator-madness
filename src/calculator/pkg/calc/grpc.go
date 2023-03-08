// Contains the grpc server and endpoint of the calculation server
package calc

import (
	"calculator/api/pb"
	"context"
	"fmt"
	"log"
	"net"
	"shared/queue"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: Either use or throw away
// type CalculationServer interface {
// 	Serve() error
// Run(context context.Context, calculationRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error)
// }

type CalcServerConfig struct {
	Port string
}

type QueueProducer interface {
	Publish(queue.Enqueable) error
	Stop()
}

// CalculationGRPCServer implements CalculationServiceServer interface
type CalculationGRPCServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
	server                                   *grpc.Server
	listener                                 net.Listener
	queueProducer                            QueueProducer
}

func NewGRPCServer(config CalcServerConfig, producer QueueProducer) (*CalculationGRPCServer, error) {
	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()

	calcServer := &CalculationGRPCServer{
		UnimplementedCalculationServiceServer: pb.UnimplementedCalculationServiceServer{},
		server:                                server,
		listener:                              listener,
		queueProducer:                         producer,
	}

	// Register endpoint
	pb.RegisterCalculationServiceServer(server, calcServer)

	// Setup reflection to be able to work with grpcurl
	reflection.Register(server)

	return calcServer, nil
}

func (s *CalculationGRPCServer) Serve() error {
	return s.server.Serve(s.listener)
}

/* GRPC Protobuf end points */

func (s *CalculationGRPCServer) Run(context context.Context, calculationRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	log.Println("Request received to Run equation", calculationRequest.Equation)

	msg, err := queue.NewCalcStartedMessage()
	if err != nil {
		return nil, err
	}

	err = s.queueProducer.Publish(msg)
	if err != nil {
		return nil, err
	}

	eq := Equation{Value: calculationRequest.Equation}

	result, err := Solve(eq)
	if err != nil {
		return nil, fmt.Errorf("failed to solve equation: %w", err)
	}
	return &pb.RunCalculationResponse{
		Result: result,
	}, nil
}
