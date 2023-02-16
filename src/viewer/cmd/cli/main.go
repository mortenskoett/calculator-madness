package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

const (
	CalcStatusQueue    string = "calc_status"
	ServiceNameChannel string = "viewer"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", getEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
)

type CalcStartedMessage struct {
	Time      string
	MessageID string
}

type MMessage struct {
 Name      string
 Content   string
 Timestamp string
}

type myMessageHandler struct{}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	log.Println("recieved message", m.Body, "on topic", CalcStatusQueue, "using channel", ServiceNameChannel)

	// do whatever actual message processing is desired
	var msg CalcStartedMessage
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		log.Println("failed to unmarshal the message body:", err)
		return err
	}

	log.Println("unmarshalled message contents: ", msg.MessageID, msg.Time)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func main() {
	log.Println("starting calculator viewer CLI client")

	flag.Parse()

	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(CalcStatusQueue, ServiceNameChannel, config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(&myMessageHandler{})

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd(*nsqlookupdAddr)
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}

func getEnvVarOrDefault(envName string, def string) string {
	envvar := os.Getenv(envName)
	if len(envvar) == 0 {
		return def
	}
	return envvar
}
