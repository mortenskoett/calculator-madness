package api

import (
	"calculator/pkg/calc"
	"fmt"
	"shared/queue"
)

type QueueProducer interface {
	Publish(queue.Enqueable) error
	Stop()
}

type queueProducer struct {
	producer QueueProducer
}

func NewQueueProducer(addr string) (*queueProducer, error) {
	producer, err := queue.NewNSQProducer(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}
	return &queueProducer{
		producer: producer,
	}, err
}

// Implements ResultNotifier.
func (n *queueProducer) Progress(ev *calc.ProgressEvent) error {
	progMsg, err := queue.NewCalcProgressMessage(ev.ClientID, ev.CalculationID, ev.ResultTopic, &queue.Status{
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
	return nil
}

// Implements ResultNotifier.
func (n *queueProducer) Ended(ev *calc.EndedEvent) error {
	endMsg, err := queue.NewCalcEndedMessage(ev.ClientID, ev.CalculationID, ev.ResultTopic, ev.Result)

	if err != nil {
		return err
	}

	err = n.producer.Publish(endMsg)

	if err != nil {
		return err
	}
	return nil
}

// Stop service gracefully.
func (n *queueProducer) Stop() {
	n.producer.Stop()
}
