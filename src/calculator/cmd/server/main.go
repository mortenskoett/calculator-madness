package main

import (
	"calculator/api/pb"
	"context"
	"fmt"
)

// calculationServer implements CalculationServiceServer interface
type calculationServer struct {
	pb.UnimplementedCalculationServiceServer // for forward compat
}

func (s *calculationServer) Run(context.Context, *pb.RunCalculationRequest) (*pb.RunCalculationResponse, error) {
	return nil,nil
}

func main() {
	fmt.Println("vim-go")
}
