package calc

type Progress struct {
	Current int
	Outof   int
}

type ProgressEvent struct {
	ClientID      string
	CalculationID string
	*Progress
}

type EndedEvent struct {
	ClientID      string
	CalculationID string
	Result        float64
}

type ResultNotifier interface {
	Progress(*ProgressEvent) error
	Ended(*EndedEvent) error
}

type Equation struct {
	Value string
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
	// TODO: Make sure the notifier is used to return result
	// return float64(len(eq.Value)),
	return nil
}

// // TODO: Artifical work: Do someting to emulate long processing time
// log.Println("sleeping...")
// time.Sleep(time.Duration(result) * time.Second)
