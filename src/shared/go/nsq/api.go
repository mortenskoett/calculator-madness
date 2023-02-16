// This package contains NSQ specific topics and messages to be used when producing and consuming
// from the queue
package api

// Topic types passed to the queue.
const (
	CalcStatusTopic string = "calc_status"
)

// Message type passed to the queue.
const (
	CalcStartedMsg  string = "calc_status_started"
	CalcProgressMsg string = "calc_status_progress"
	CalcEndedMsg    string = "calc_status_ended"
)
