package calc

import "log"

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
	Process(*Equation)
	GetResults() <-chan *EndedEvent
	GetProgress() <-chan *ProgressEvent
}

type ResultNotifier interface {
	Progress(*ProgressEvent) error
	Ended(*EndedEvent) error
}

type calculatorService struct {
	processor EquationProcessor
	notifier  ResultNotifier
}

func NewCalculatorService(processor EquationProcessor, notifier ResultNotifier) *calculatorService {
	c := &calculatorService{
		processor: processor,
		notifier:  notifier,
	}
	go c.handleEvents()
	return c
}

// Solve enqueues an equation for solving. The result is returned through the ResultNotifier.
func (c *calculatorService) Solve(eq *Equation) error {
	c.processor.Process(eq)
	return nil
}

func (c *calculatorService) handleEvents() error {
	log.Println("listening for equation events")
	for {
		select {
		case p := <-c.processor.GetProgress():
			err := c.notifier.Progress(p)
			if err != nil {
				log.Println("failed to send progress message calculationID:", p.CalculationID)
			}
		case r := <-c.processor.GetResults():
			err := c.notifier.Ended(r)
			if err != nil {
				log.Println("failed to send end message for calculationID:", r.CalculationID)
			}
		}
	}
}
