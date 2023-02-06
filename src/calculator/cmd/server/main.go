package main

import (
	"calculator/api/pb"
	"calculator/pkg/calc"
	"context"
	"fmt"
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
	fmt.Println("vim-go")
}
