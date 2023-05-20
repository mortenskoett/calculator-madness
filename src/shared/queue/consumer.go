package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

type NsqQueueConsumer struct {
	topic         string
	channel       string
	nsqlookupAddr string
	consumer      *nsq.Consumer
}

// Create a new NSQ consumer that subcribes to the given service channel and consumes messages from
// one or more topics. Handlers to consume topics must be added to it after instantiation.
func NewNSQConsumer(nsqlookupdAddr, topic, channel string) (*NsqQueueConsumer, error) {
	log.Println("creating new nsq consumer")
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nsq consumer")
	}

	return &NsqQueueConsumer{
		topic:         topic,
		channel:       channel,
		nsqlookupAddr: nsqlookupdAddr,
		consumer:      consumer,
	}, nil
}

// Start listening for incoming messages on the set topic. Blocking call.
func (c *NsqQueueConsumer) Start(ctx context.Context) {
	log.Println("connecting to nsqlookupd")
	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err := c.consumer.ConnectToNSQLookupd(c.nsqlookupAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Handle shutdown using context
	<-ctx.Done()
	log.Println("stopping nsq consumer: cancelled by context.")
	return
}

// Should be a deferred call to stop the consumer gracefully.
func (c *NsqQueueConsumer) Stop() {
	// Gracefully stop the consumer.
	c.consumer.Stop()
	log.Println("nsq consumer stopped")
}

// Callback types
type callback[T Enqueable] func(*T, error) error
type CalcProgressCallback func(*CalcProgressMessage, error) error
type CalcEndedCallback func(*CalcEndedMessage, error) error

/* Add handlers of a specific message types. Panics if called after Start(). */

func (c *NsqQueueConsumer) AddCalcProgressHandler(fn CalcProgressCallback) {
	addCallback(fn, c)
}

func (c *NsqQueueConsumer) AddCalcEndedHandler(fn CalcEndedCallback) {
	addCallback(fn, c)
}

/* Code to handle each incoming message */

func addCallback[T Enqueable](fn func(*T, error) error, c *NsqQueueConsumer) {
	c.consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		return attachCallback(m, fn)
	}))
}

func attachCallback[T Enqueable](m *nsq.Message, callback callback[T]) error {
	msg, err := unmarshalMessage[T](m)
	if err != nil {
		return callback(nil, err)
	}
	return callback(msg, nil)
}

// Unmarshal a received nsq message.
func unmarshalMessage[T Enqueable](m *nsq.Message) (*T, error) {
	if len(m.Body) == 0 {
		// Returning nil will send a FIN command to NSQ marking the message as processed.
		return nil, nil
	}

	var msg T
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the message body: %w", err)
	}

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return &msg, nil
}
