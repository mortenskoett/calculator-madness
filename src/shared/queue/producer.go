package queue

import (
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"
)

// Interface for types that can be enqueued.
type Enqueable interface {
	topic() string // Topic returns the topic to which the message should be published.
}

type nsqProducer struct {
	producer *nsq.Producer
}

// Creates a new NSQ producer client
func NewNSQProducer(serverAddr string) (*nsqProducer, error) {
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(serverAddr, nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create nsq producer: %w", err)
	}
	return &nsqProducer{producer: producer}, nil
}

func (n *nsqProducer) Publish(msg Enqueable) error {
	bytes, err := toByteSlice(msg)
	if err != nil {
		return err
	}

	err = n.producer.Publish(msg.topic(), bytes)
	if err != nil {
		return err
	}

	return nil
}

func (n *nsqProducer) Stop() {
	n.producer.Stop()
}

func toByteSlice(v any) ([]byte, error) {
	return json.Marshal(v)
}
