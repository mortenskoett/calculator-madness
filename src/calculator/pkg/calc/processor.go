package calc

import (
	"time"
)

const (
	progressIntervalSecs = 1
)

// Workload defines the content needed to track progress during the processing of an equation.
type workload struct {
	ticker   *time.Ticker
	equation *Equation
	progress *Progress
}

// Processor processes an equation and returns the result over a channel.
type Processor struct {
	results chan *EquationResult // Receive processed equation results on this channel.
	intake  chan *workload       // Equations waiting to be processed are placed here.
}

func NewProcessor(maxConcurrent int, maxChannelSize int) *Processor {
	h := &Processor{
		results: make(chan *EquationResult, maxChannelSize),
		intake:  make(chan *workload, maxChannelSize),
	}

	// Start workers.
	for i := 0; i < maxConcurrent; i++ {
		go h.createWorker(h.intake, h.results)
	}
	return h
}

// Add an equation to be processed.
func (h *Processor) Add(eq *Equation) {
	h.intake <- &workload{
		ticker:   time.NewTicker(progressIntervalSecs * time.Second),
		equation: eq,
		progress: &Progress{
			Current: 0,
			Outof:   len(eq.Expression),
		},
	}
}

func (h *Processor) Results() <-chan *EquationResult {
	return h.results
}

// Creates a worker to process Equations. Call as goroutine.
func (h *Processor) createWorker(in <-chan *workload, out chan<- *EquationResult) {
	for {
		select {
		case w := <-in:
			w.start(out)
		}
	}
}

// Process starts the interval outputted processing of an Equation.
func (w *workload) start(out chan<- *EquationResult) {
	for {
		select {
		case <-w.ticker.C:
			w.progress.Current++

			if w.progress.Current < w.progress.Outof {
				continue
			}

			finalResult := &EquationResult{
				Equation: w.equation,
				Result:   float64(len(w.equation.Expression)),
			}
			out <- finalResult
		}
	}
}
