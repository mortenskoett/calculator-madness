package calc

type Equation struct {
	Value string
}

func Solve(eq Equation) (float64, error) {
	return float64(len(eq.Value)), nil
}
