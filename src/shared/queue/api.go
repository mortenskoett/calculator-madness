// This package contains topics and messages to be used when producing and consuming from the
// calculatoer status queue.
package queue

// Topic type passed to the queue.
const (
	CalcStatusTopic string = "calc_status"
)

// Message type passed to the queue.
const (
	CalcStartedMsg  string = "calc_status_started"
	CalcProgressMsg string = "calc_status_progress"
	CalcEndedMsg    string = "calc_status_ended"
)

