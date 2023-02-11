package queue

import (
	"encoding/json"
	"fmt"
	"time"
)

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

type Message interface {
	Topic() string
	Message() []byte
}

type calcStartedMessage struct {
	time      time.Time
	messageID string
	payload   []byte
}

func NewCalcStartedMessage() (Message, error) {
	tmp := struct {
		time      time.Time
		messageID string
	}{
		time:      time.Now(),
		messageID: CalcStartedMsg,
	}

	bytes, err := toByteSlice(tmp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message to bytes: %v", err)
	}

	return calcStartedMessage{
		time:      time.Now(),
		messageID: CalcStartedMsg,
		payload:   bytes,
	}, nil
}

func (m calcStartedMessage) Topic() string {
	return CalcStatusTopic
}

func (m calcStartedMessage) Message() []byte {
	return m.payload
}

func toByteSlice(v any) ([]byte, error) {
	return json.Marshal(v)
}
