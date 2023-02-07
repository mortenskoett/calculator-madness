package calc

type Equation struct {
	Value string
}

func Solve(eq Equation) (float64, error) {
	// TODO: Dummy return value
	return float64(len(eq.Value)), nil
}
