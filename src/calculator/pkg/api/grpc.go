// Contains the grpc server and endpoint of the calculation server
package api

import (
	"calculator/api/pb"
	"calculator/pkg/calc"
	"context"
	"fmt"
	"log"
	"net"
	"shared/queue"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CalculationService interface {
	Enqueue(eq *calc.Equation) error
}

type QueueProducer interface {
	Publish(queue.Enqueable) error
	Stop()
}

type CalcServerConfig struct {
	Port string
}

// CalculationGRPCServer implements CalculationServiceServer interface
type CalculationGRPCServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
	server                                   *grpc.Server
	listener                                 net.Listener
	calcService                              CalculationService
}

func NewGRPCServer(port string, calcService CalculationService) (*CalculationGRPCServer, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()

	calcServer := &CalculationGRPCServer{
		UnimplementedCalculationServiceServer: pb.UnimplementedCalculationServiceServer{},
		server:                                server,
		listener:                              listener,
		calcService:                           calcService,
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

func (s *CalculationGRPCServer) Run(context context.Context, req *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	log.Printf("request received to run equation: %+v", req.Equation.Id)

	// // Fix this
	// startMsg, err := queue.NewCalcStartedMessage(calcRequest.GetClientId(), calcRequest.GetEquation().Id)
	// if err != nil {
	// 	return nil, err
	// }

	// err = s.queueProducer.Publish(startMsg)
	// if err != nil {
	// return nil, err
	// }
	// log.Println("new calculation message sent to queue: clientId:", calcRequest.ClientId)

	equation := &calc.Equation{Value: req.Equation.Value}
	err := s.calcService.Enqueue(equation)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue equation: %w", err)
	}

	// progMsg, err := queue.NewCalcProgressMessage(startMsg.CalculationID)
	// if err != nil {
	// 	return nil, err
	// }

	// err = s.queueProducer.Publish(progMsg)
	// if err != nil {
	// return nil, err
	// }

	// TODO: Make the above stuff run in a goroutine so that we don't wait long for this return
	return &pb.RunCalculationResponse{}, nil
}
