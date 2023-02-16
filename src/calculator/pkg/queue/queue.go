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

type MMessage struct {
	Name      string
	Content   string
	Timestamp string
}

func (n *nsqProducer) Publish(msg Message) error {

	//Init topic name and message
	// shit := MMessage{
	// 	Name:      "Message Name Example",
	// 	Content:   "Message Content Example",
	// 	Timestamp: time.Now().String(),
	// }

	//Convert message as []byte
	// payload, err := json.Marshal(shit)
	// if err != nil {
	// 	log.Println(err)
	// }

	log.Println("MSK1", msg.Message())

	err := n.producer.Publish(msg.Topic(), msg.Message())
	if err != nil {
		return err
	}
	log.Println("Sent message", msg.Message(), " to queue on topic", msg.Topic())
	return nil
}

func (n *nsqProducer) Stop() {
	n.producer.Stop()
}
