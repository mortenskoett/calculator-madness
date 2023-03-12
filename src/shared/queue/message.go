package queue

import (
	"time"
)

// Additional metadata of the message
type MessageMetadata struct {
	CreatedTime time.Time
}

// Specific message designated a calculation has started
type CalcStartedMessage struct {
	*MessageMetadata
	MessageID string
}

func newMessageMetadata() *MessageMetadata {
	return &MessageMetadata{
		CreatedTime: time.Now(),
	}
}

func NewCalcStartedMessage() (*CalcStartedMessage, error) {
	mesg := CalcStartedMessage{
		MessageMetadata: newMessageMetadata(),
		MessageID:       CalcStartedMsg,
	}

	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcStartedMessage) topic() string {
	return CalcStatusTopic
}
