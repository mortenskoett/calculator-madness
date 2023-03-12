package main

import (
	"flag"
	"log"
	"os"
	"shared/queue"
)

const (
ServiceNameChannel string = "viewer"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", getEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
)

func calcStartedHandler(msg *queue.CalcStartedMessage, err error) error {
	log.Printf("Hello from the viewer: time: %+v, msgID: %+v\n", msg.CreatedTime, msg.MessageID)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("starting calculator viewer cli client")
	flag.Parse()

	consumer, err := queue.NewNSQConsumer(*nsqlookupdAddr, queue.CalcStatusTopic, ServiceNameChannel)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddCalcStartedHandler(calcStartedHandler)
	consumer.Start()
	defer consumer.Stop()
}

func getEnvVarOrDefault(envName string, def string) string {
	envvar := os.Getenv(envName)
	if len(envvar) == 0 {
		return def
	}
	return envvar
}

