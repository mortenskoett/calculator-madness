package queue

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Additional metadata of the message
type Progress struct {
	Current int `json:"current"`
	Outof   int `json:"outof"`
}

type Status struct {
	Progress Progress `json:"progress"`
}

type MessageMetadata struct {
	MessageID    string    `json:"message_id"`
	MessageTopic string    `json:"message_topic"` // Topic that the producer should return the result to.
	CreatedTime  time.Time `json:"created_time"`
}

// Specific message designated a calculation has started
type CalcProgressMessage struct {
	*MessageMetadata `json:"message_metadata"`
	ClientID         string  `json:"client_id"`
	CalculationID    string  `json:"calculation_id"`
	Status           *Status `json:"status"`
}

// Specific message designated a calculation has ended
type CalcEndedMessage struct {
	*MessageMetadata
	ClientID      string  `json:"client_id"`
	CalculationID string  `json:"calculation_id"`
	Result        float64 `json:"result"`
}

func NewCalcProgressMessage(clientID, calcID, topic string, status *Status) (*CalcProgressMessage, error) {
	mesg := CalcProgressMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:    uuid.NewString(),
			MessageTopic: topic,
			CreatedTime:  time.Now(),
		},
		ClientID:      clientID,
		CalculationID: calcID,
		Status:        status,
	}
	return &mesg, nil
}

func NewCalcEndedMessage(clientID, calcID, topic string, result float64) (*CalcEndedMessage, error) {
	mesg := CalcEndedMessage{
		MessageMetadata: &MessageMetadata{
			MessageID:    uuid.NewString(),
			MessageTopic: topic,
			CreatedTime:  time.Now(),
		},
		ClientID:      clientID,
		CalculationID: calcID,
		Result:        result,
	}
	return &mesg, nil
}

// Topic implements the Enquable interface
func (m CalcProgressMessage) topic() string {
	return m.MessageTopic
}

// Topic implements the Enquable interface
func (m CalcEndedMessage) topic() string {
	return m.MessageTopic
}

// Unpacker assists in unmarshalling interface types into structs.
type Unpacker struct {
	Data interface{}
}

func NewUnpacker() *Unpacker {
	return &Unpacker{
		Data: nil,
	}
}

func (u *Unpacker) UnmarshalJSON(b []byte) error {
	prog := CalcProgressMessage{}
	err := json.Unmarshal(b, &prog)

	// Validate progress message
	if err == nil && prog.Status != nil {
		u.Data = prog
		return nil
	}

	// Exit if other than type error
	if _, ok := err.(*json.UnmarshalTypeError); err != nil && !ok {
		return err
	}

	end := CalcEndedMessage{}
	err = json.Unmarshal(b, &end)
	if err == nil {
		u.Data = end
		return nil
	}

	return nil
}

func (u *Unpacker) Get() Enqueable {
	switch d := u.Data.(type) {
	case CalcProgressMessage:
		return CalcProgressMessage(d)
	case CalcEndedMessage:
		return CalcEndedMessage(d)
	}
	return nil
}
