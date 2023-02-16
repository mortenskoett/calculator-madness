package queue

import (
	nsqapi "shared/api/nsq"
)

// Interface for enqueable messages
type Enqueable interface {
	Topic() string
}

func (m *nsqapi.CalcStartedMessage) Topic() string {
	return nsqapi.CalcStatusTopic
}

