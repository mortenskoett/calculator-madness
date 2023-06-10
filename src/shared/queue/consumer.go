package queue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

const (
	channelSize int = 1000
)

type NSQUniqueConsumer[T Enqueable] struct {
	topic         string
	channel       string
	nsqlookupAddr string
	consumer      *nsq.Consumer
	results       chan T
	stopSignal    chan struct{}
}

func NewNSQUniqueConsumer[T Enqueable](nsqlookupdAddr string, uniqueTopic string) (*NSQUniqueConsumer[T], error) {
	log.Println("creating new nsq consumer")

	uniqueChannel := uniqueTopic + "-channel"
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(uniqueTopic, uniqueChannel, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nsq consumer")
	}

	cons := NSQUniqueConsumer[T]{
		topic:         uniqueTopic,
		channel:       uniqueChannel,
		nsqlookupAddr: nsqlookupdAddr,
		consumer:      consumer,
		results:       make(chan T, channelSize),
		stopSignal:    make(chan struct{}, 1),
	}

	cons.setNSQEventHandler()

	return &cons, nil
}

func (c *NSQUniqueConsumer[T]) Start(ctx context.Context) {
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
func (c *NSQUniqueConsumer[T]) Stop() {
	// Gracefully stop the consumer.
	log.Println("stopping nsq consumer")
	c.consumer.Stop()

	log.Println("stop signal sent")
	c.stopSignal <- struct{}{}
}

// Handler used to process consumed messages implementing interface type T.
type MsgHandler[T any] func(msg T) error

// Set a handler for incoming messages. The handler should handle multiple concrete implementations
// of the T interface. The callback is run as a goroutine and will be called everytime a message is
// consumed.
func (c *NSQUniqueConsumer[T]) SetHandler(callbackFn MsgHandler[T]) {
	go func() {
		log.Println("starting handler loop")
		for {
			select {
			case msg := <-c.results:
				callbackFn(msg)
			case <-c.stopSignal:
				log.Println("stopping handler loop")
				return
			}
		}
	}()
}

// The default handler unmarshals all incoming NSQ event message types.
func (c *NSQUniqueConsumer[T]) setNSQEventHandler() {
	c.consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			// Returning nil will send a FIN command to NSQ marking the message as processed.
			return nil
		}

		var unpacker = NewUnpacker()
		err := json.Unmarshal(m.Body, &unpacker)
		if err != nil {
			log.Printf("failed to unmarshal received message into unpacker: %v\n", err)
			return nil
		}

		// Necessary to assert type of unpacked data before sending to channel.
		c.results <- unpacker.Get().(T)

		return nil
	}))
}
