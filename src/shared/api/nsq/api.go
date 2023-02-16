// This package contains NSQ specific topics and messages to be used when producing and consuming
// from the queue
package api

import "time"

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

// Additional metadata of the message
type MessageMetadata struct {
	Time string
}

// Specific message designated a calculation has started
type CalcStartedMessage struct {
	MessageMetadata
	MessageID string
}

func newMessageMetadata() *MessageMetadata {
	return &MessageMetadata{
		Time: time.Now().String(),
	}
}

func NewCalcStartedMessage() (*CalcStartedMessage, error) {
	mesg := CalcStartedMessage{
		MessageMetadata: MessageMetadata{},
		MessageID:       CalcStartedMsg,
	}

	return &mesg, nil
}
