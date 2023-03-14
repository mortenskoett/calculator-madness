// This package contains topics and messages to be used when producing and consuming from the
// calculatoer status queue.
package queue

// Topic type passed to the queue.
const (
	CalcStatusTopic string = "calc_status"
)

// Calculation message IDs
const (
	CalcStartedID  string = "calc_status_started"
	CalcProgressID string = "calc_status_progress"
	CalcEndedID    string = "calc_status_ended"
)

