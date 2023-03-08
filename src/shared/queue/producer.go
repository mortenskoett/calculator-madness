package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

// Interface for types that can be enqueued.
type Enqueable interface {
	topic() string
}

type NsqQueueProducer struct {
	producer *nsq.Producer
}

// Creates a new NSQ producer client
func NewNSQProducer(serverAddr string) (*NsqQueueProducer, error) {
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(serverAddr, nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create nsq producer: %w", err)
	}
	return &NsqQueueProducer{producer: producer}, nil
}

func (n *NsqQueueProducer) Publish(msg Enqueable) error {
	bytes, err := toByteSlice(msg)
	if err != nil {
		return err
	}

	err = n.producer.Publish(msg.topic(), bytes)
	if err != nil {
		return err
	}

	log.Println("sent message to queue on topic:", msg.topic())

	return nil
}

func (n *NsqQueueProducer) Stop() {
	n.producer.Stop()
}

func toByteSlice(v any) ([]byte, error) {
	return json.Marshal(v)
}
