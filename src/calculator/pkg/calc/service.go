package calc

type ClientInfo struct {
	ClientID      string
	CalculationID string
}

type Equation struct {
	*ClientInfo
	Expression string
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

// Enqueue enqueues an equation for solving. The result is returned through the ResultNotifier.
func (c *calculatorService) Enqueue(eq *Equation) error {
	res := float64(len(eq.Expression))

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
