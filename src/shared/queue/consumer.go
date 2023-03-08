package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

type NsqQueueConsumer struct {
	Topic         string
	Channel       string
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
		Topic:         topic,
		Channel:       channel,
		nsqlookupAddr: nsqlookupdAddr,
		consumer:      consumer,
	}, nil
}

// Start listening for incoming messages on the set topic.
func (c *NsqQueueConsumer) Start() {
	log.Println("connecting to nsqlookupd")
	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err := c.consumer.ConnectToNSQLookupd(c.nsqlookupAddr)
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

// Should be a deferred call to stop the consumer gracefully.
func (c *NsqQueueConsumer) Stop() {
	// Gracefully stop the consumer.
	c.consumer.Stop()
}

// Contain
type callbackHolder[T Enqueable] struct {
	callback func(*T, error) error
}

type CalcStartedCallback func(*CalcStartedMessage, error) error

// Add handler callback of a specific message type. Panics if called after Start().
func (c *NsqQueueConsumer) AddCalcStartedHandler(fn CalcStartedCallback) {
	c.consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		msgHandler := callbackHolder[CalcStartedMessage]{
			callback: fn,
		}
		return handleCallback(m, msgHandler)
	}))
}

func handleCallback[T Enqueable](m *nsq.Message, handler callbackHolder[T]) error {
	msg, err := unmarshalMessage[T](m)
	if err != nil {
		return handler.callback(nil, err)
	}
	return handler.callback(msg, nil)
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
