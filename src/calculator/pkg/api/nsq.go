package api

import (
	"calculator/pkg/calc"
	"fmt"
	"log"
	"shared/queue"
)

type nsqProducer struct {
	producer *queue.NsqQueueProducer
}

func NewNSQProducer(addr string) (*nsqProducer, error) {
	producer, err := queue.NewNSQProducer(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}
	return &nsqProducer{
		producer: producer,
	}, err
}

// Implements ResultNotifier.
func (n *nsqProducer) Progress(ev *calc.ProgressEvent) error {
	progMsg, err := queue.NewCalcProgressMessage(ev.ClientID, ev.CalculationID, &queue.Status{
		Progress: queue.Progress{
			Current: ev.Current,
			Outof:   ev.Outof,
		},
	})
	if err != nil {
		return err
	}

	err = n.producer.Publish(progMsg)
	if err != nil {
		return err
	}
	log.Println("progress calculation message sent to queue: clientId:", ev.ClientID)
	return nil
}

// Implements ResultNotifier.
func (n *nsqProducer) Ended(ev *calc.EndedEvent) error {
	endMsg, err := queue.NewCalcEndedMessage(ev.ClientID, ev.CalculationID, ev.Result)
	if err != nil {
		return err
	}

	err = n.producer.Publish(endMsg)
	if err != nil {
		return err
	}
	log.Println("end calculation message sent to queue: clientId:", ev.ClientID)
	return nil
}

// Stop service gracefully.
func (n *nsqProducer) Stop() {
	n.producer.Stop()
}
