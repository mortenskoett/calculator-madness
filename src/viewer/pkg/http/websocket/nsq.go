package websocket

import (
	"log"
	"shared/queue"
)

// Handles incoming calculation progress events on NSQ.
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

// Handles incoming calculation ended events on NSQ.
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
