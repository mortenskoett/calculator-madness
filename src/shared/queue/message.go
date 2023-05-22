package queue

import (
	"time"

	"github.com/google/uuid"
)

// Additional metadata of the message
type Progress struct {
	Current int
	Outof   int
}

type Status struct {
	Progress Progress
}

type MessageMetadata struct {
	MessageID   string
	CreatedTime time.Time
}

// Specific message designated a calculation has started
type CalcProgressMessage struct {
	*MessageMetadata
	ClientID      string
	CalculationID string
	Status        *Status
}

func NewCalcProgressMessage(clientID string, calcID string, status *Status) (*CalcProgressMessage, error) {
	mesg := CalcProgressMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:   uuid.NewString(),
			CreatedTime: time.Now(),
		},
		ClientID:      clientID,
		CalculationID: calcID,
		Status:        status,
	}
	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcProgressMessage) topic() string {
	return CalculationStatusTopic
}

// Specific message designated a calculation has ended
type CalcEndedMessage struct {
	*MessageMetadata
	ClientID      string
	CalculationID string
	Result        float64
}

func NewCalcEndedMessage(clientID string, calcID string, result float64) (*CalcEndedMessage, error) {
	mesg := CalcEndedMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:   uuid.NewString(),
			CreatedTime: time.Now(),
		},
		ClientID:      clientID,
		CalculationID: calcID,
		Result:        result,
	}
	return &mesg, nil
}

// topic implements the Enquable interface
func (m CalcEndedMessage) topic() string {
	return CalculationStatusTopic
}
