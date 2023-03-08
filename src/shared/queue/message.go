package queue

import (
	"time"
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

// topic implements the Enquable interface
func (m CalcStartedMessage) topic() string {
	return CalcStatusTopic
}
