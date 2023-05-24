package calc

type ClientInfo struct {
	ClientID      string
	CalculationID string
}

type Equation struct {
	*ClientInfo
	Expression string
}

type EquationResult struct {
	Equation *Equation
	Result   float64
}

type Progress struct {
	Current int
	Outof   int
}

type ProgressEvent struct {
	*ClientInfo
	*Progress
}

type EndedEvent struct {
	*ClientInfo
	Result float64
}

type EquationProcessor interface {
	Add(*Equation)
	Results() <-chan *EquationResult
}

type ResultNotifier interface {
	Progress(*ProgressEvent) error
	Ended(*EndedEvent) error
}

type calculatorService struct {
	notifier ResultNotifier
}

func NewCalculatorService(notifier ResultNotifier) *calculatorService {
	return &calculatorService{
		notifier: notifier,
	}
}

// Solve enqueues an equation for solving. The result is returned through the ResultNotifier.
func (c *calculatorService) Solve(eq *Equation) error {
	res := float64(len(eq.Expression))

	// proc := NewEquationProcessor(20, 100)

	endEvent := EndedEvent{
		ClientInfo: &ClientInfo{
			ClientID:      eq.ClientID,
			CalculationID: eq.CalculationID,
		},
		Result: res,
	}

	err := c.notifier.Ended(&endEvent)
	if err != nil {
		return err
	}
	return nil
}
