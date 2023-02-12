package queue

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

type QueueProducer interface {
	Publish(Message) error
	Stop()
}

type nsqProducer struct {
	producer *nsq.Producer
}

func NewNSQProducer(serverAddr string) (QueueProducer, error) {
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(serverAddr, nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create nsq producer: %w", err)
	}

	return &nsqProducer{producer: producer}, nil
}

func (n *nsqProducer) Publish(msg Message) error {
	err := n.producer.Publish(msg.Topic(), msg.Message())
	if err != nil {
		return err
	}
	log.Println("Sent message to queue on topic", msg.Topic())
	return nil
}

func (n *nsqProducer) Stop() {
	n.producer.Stop()
}
