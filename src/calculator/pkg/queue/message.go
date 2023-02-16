package queue

import (
	nsqapi "calculator/shared/nsq"
	"encoding/json"
	"fmt"
	"log"
	"time"
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
	// Anonymous struct containing fields that is send to the queue
	tmp := struct {
		Time      string
		MessageID string
	}{
		Time:      time.Now().String(),
		MessageID: nsqapi.CalcStartedMsg,
	}

	bytes, err := toByteSlice(tmp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message to bytes: %v", err)
	}

	log.Println(bytes)

	mesg := calcStartedMessage{
		time:      time.Now(),
		messageID: nsqapi.CalcStartedMsg,
		payload:   bytes,
	}

	return &mesg, nil
}

func (m *calcStartedMessage) Topic() string {
	return nsqapi.CalcStatusTopic
}

func (m *calcStartedMessage) Message() []byte {
	return m.payload
}

func toByteSlice(v any) ([]byte, error) {
	res, err := json.Marshal(v)
	log.Println("to bytes", res, string(res))
	return res, err
}
