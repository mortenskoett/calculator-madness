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
	log.Printf("Calc started: calcID: %+v, msgID: %+v, time: %+v\n",
		msg.CalculationID,
		msg.MessageID,
		msg.CreatedTime,
	)

	if err != nil {
		return err
	}
	return nil
}

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

	consumer, err := queue.NewNSQConsumer(*nsqlookupdAddr, queue.CalcStatusTopic, ServiceNameChannel)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddCalcStartedHandler(calcStartedHandler)
	consumer.AddCalcProgressHandler(calcProgressHandler)
	consumer.AddCalcEndedHandler(calcEndedHandler)
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
