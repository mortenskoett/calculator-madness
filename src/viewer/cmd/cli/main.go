package main

import (
	"context"
	"flag"
	"log"
	"shared/queue"
	"viewer/pkg/env"
)

const (
	ServiceNameChannel string = "viewer-cli"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
)

func calcProgressHandler(msg *queue.CalcProgressMessage, err error) error {
	log.Printf("Calc progress: calcID: %+v, msgID: %+v, time: %+v\n",
		msg.CalculationID,
		msg.MessageID,
		msg.CreatedTime,
	)

	if err != nil {
		return err
	}
	return nil
}

func calcEndedHandler(msg *queue.CalcEndedMessage, err error) error {
	log.Printf("Calc ended: calcID: %+v, msgID: %+v, time: %+v\n",
		msg.CalculationID,
		msg.MessageID,
		msg.CreatedTime,
	)

	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("starting calculator viewer cli client")
	flag.Parse()

	consumer, err := queue.NewNSQConsumer(*nsqlookupdAddr, queue.CalculationStatusTopic, ServiceNameChannel)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Stop()

	consumer.AddCalcProgressHandler(calcProgressHandler)
	consumer.AddCalcEndedHandler(calcEndedHandler)
	consumer.Start(context.Background())
}
