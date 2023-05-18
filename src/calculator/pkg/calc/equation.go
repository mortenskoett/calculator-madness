package calc

import (
	"calculator/api/pb"
)

type Equation struct {
	Value string
}

func Solve(eq *pb.Equation) (float64, error) {
	return float64(len(eq.Value)), nil
}
