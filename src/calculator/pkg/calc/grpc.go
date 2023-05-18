// Contains the grpc server and endpoint of the calculation server
package calc

import (
	"calculator/api/pb"
	"context"
	"fmt"
	"log"
	"net"
	"shared/queue"
	"time"

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
	listener, err := net.Listen("tcp", ":"+config.Port)
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

func (s *CalculationGRPCServer) Run(context context.Context, calcRequest *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	log.Printf("request received to Run equation: %+v", calcRequest.Equation)

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

	result, err := Solve(calcRequest.Equation)
	if err != nil {
		return nil, fmt.Errorf("failed to solve equation: %w", err)
	}

	// progMsg, err := queue.NewCalcProgressMessage(startMsg.CalculationID)
	// if err != nil {
	// 	return nil, err
	// }

	// err = s.queueProducer.Publish(progMsg)
	// if err != nil {
	// return nil, err
	// }

	// TODO: Artifical work: Do someting to emulate long processing time
	log.Println("sleeping...")
	time.Sleep(time.Duration(result) * time.Second)

	endMsg, err := queue.NewCalcEndedMessage(calcRequest.GetClientId(), calcRequest.GetEquation().Id)
	if err != nil {
		return nil, err
	}
	log.Println("end calculation message sent to queue: clientId:", calcRequest.ClientId)

	err = s.queueProducer.Publish(endMsg)
	if err != nil {
		return nil, err
	}

	// TODO: Make the above stuff run in a goroutine so that we don't wait long for this return
	return &pb.RunCalculationResponse{
		Error: nil,
	}, nil
}
