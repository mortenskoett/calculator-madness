package queue

import (
	"time"

	"github.com/google/uuid"
)

// Additional metadata of the message
type MessageMetadata struct {
	MessageID   string
	CreatedTime time.Time
}

/* Calculation started message */

// Specific message designated a calculation has started
type CalcStartedMessage struct {
	*MessageMetadata
	CalculationID uuid.UUID
}

func NewCalcStartedMessage() (*CalcStartedMessage, error) {
	mesg := CalcStartedMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:   CalcStartedID,
			CreatedTime: time.Now(),
		},
		CalculationID: uuid.New(),
	}
	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcStartedMessage) topic() string {
	return CalcStatusTopic
}

/* Calculation progress message */

// Specific message designated a calculation has started
type CalcProgressMessage struct {
	*MessageMetadata
	CalculationID uuid.UUID
}

func NewCalcProgressMessage(id uuid.UUID) (*CalcProgressMessage, error) {
	mesg := CalcProgressMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:   CalcProgressID,
			CreatedTime: time.Now(),
		},
		CalculationID: id,
	}
	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcProgressMessage) topic() string {
	return CalcStatusTopic
}

// Specific message designated a calculation has ended
type CalcEndedMessage struct {
	*MessageMetadata
	CalculationID uuid.UUID
}

func NewCalcEndedMessage(id uuid.UUID) (*CalcEndedMessage, error) {
	mesg := CalcEndedMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:   CalcEndedID,
			CreatedTime: time.Now(),
		},
		CalculationID: id,
	}
	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcEndedMessage) topic() string {
	return CalcStatusTopic
}
