package calc

import (
	"time"
)

const (
	progressIntervalSecs int = 2
)

// Workload defines the content needed to track progress during the processing of an equation.
type workload struct {
	ticker   *time.Ticker
	equation *Equation
	progress *Progress
}

type dummyProcessor struct {
	results  chan *EndedEvent    // Receive processed equation results on this channel.
	progress chan *ProgressEvent // Receive progress on equations on this channel.
	intake   chan *workload      // Equations waiting to be processed are placed here.
}

// NewDummyProcessor creates a new processor to process equations concurrently. It returns results
// and progress messages over separate channels. This particular processor is quite stupid. It
// counts the chars N of the given equation and returns N as result. Furthermore it spends
// N * progressIntervalSecs seconds processing the equation while sending a progress message every
// progressIntervalSecs second.
func NewDummyProcessor(maxConcurrent int, maxChannelSize int) *dummyProcessor {
	h := &dummyProcessor{
		results:  make(chan *EndedEvent, maxChannelSize),
		progress: make(chan *ProgressEvent, maxChannelSize),
		intake:   make(chan *workload, maxChannelSize),
	}

	// Start workers.
	for i := 0; i < maxConcurrent; i++ {
		go h.createWorker(h.intake, h.results, h.progress)
	}
	return h
}

// Process enqueues an equation for processing. If the intake channel is full the function will
// block.
func (h *dummyProcessor) Process(eq *Equation) {
	h.intake <- &workload{
		ticker:   time.NewTicker(time.Duration(progressIntervalSecs) * time.Second),
		equation: eq,
		progress: &Progress{
			Current: 0,
			Outof:   len(eq.Expression),
		},
	}
}

// GetResults returns the channel on which equation result messages are posted when they are done.
func (h *dummyProcessor) GetResults() <-chan *EndedEvent {
	return h.results
}

// GetProgress returns the channel on which equation progress messages are posted when they are done.
func (h *dummyProcessor) GetProgress() <-chan *ProgressEvent {
	return h.progress
}

// Creates a worker to process Equations. Call as goroutine.
func (h *dummyProcessor) createWorker(in <-chan *workload, resultOut chan<- *EndedEvent, progressOut chan<- *ProgressEvent) {
	for {
		select {
		case w := <-in:
			w.start(resultOut, progressOut)
		}
	}
}

// Process starts the interval outputted processing of an Equation.
func (w *workload) start(resultOut chan<- *EndedEvent, progressOut chan<- *ProgressEvent) {
	for {
		select {
		case <-w.ticker.C:
			w.progress.Current++

			// Send Progress message back.
			if w.progress.Current < w.progress.Outof {
				progress := &ProgressEvent{
					ClientInfo: &ClientInfo{
						ClientID:      w.equation.ClientID,
						CalculationID: w.equation.CalculationID,
					},
					Progress: w.progress,
				}
				progressOut <- progress
				break
			}

			// Send Result message back.
			result := &EndedEvent{
				ClientInfo: &ClientInfo{
					ClientID:      w.equation.ClientID,
					CalculationID: w.equation.CalculationID,
				},
				Result: float64(len(w.equation.Expression)),
			}
			resultOut <- result
			return
		}
	}
}
