package websocket

import (
	"log"
	"shared/queue"
)

func (m *Manager) NSQCalcProgressHandler(msg *queue.CalcProgressMessage, err error) error {
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

func (m *Manager) NSQCalcEndedHandler(msg *queue.CalcEndedMessage, err error) error {
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
