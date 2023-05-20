package websocket

import (
	"encoding/json"
	"log"
	"shared/queue"
)

// Callback to handle incoming calculation ended events on NSQ.
// Warning: When returning error here, the message will be resent with a backoff.
// func (m *Manager) NSQCalcProgressHandler(msg *queue.CalcProgressMessage, err error) error {
// 	m.clients[msg.ClientID].send()

// 	log.Printf("Calc progress: calcID: %+v, msgID: %+v, time: %+v\n",
// 		msg.CalculationID,
// 		msg.MessageID,
// 		msg.CreatedTime,
// 	)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Callback to handle incoming calculation ended events on NSQ.
// Warning: When returning error here, the message will be resent with a backoff.
func (m *Manager) NSQCalcEndedHandler(msg *queue.CalcEndedMessage, err error) error {
	log.Printf("received ended calculation for client %v for calculation %v", msg.ClientID, msg.CalculationID)

	c, ok := m.clients[msg.ClientID]
	if !ok {
		log.Println("failed to handle ended calculation: client did not exist: %w", err)
		return nil
	}

	resp := EndCalculationResponse{
		ID:     msg.CalculationID,
		Result: msg.Result,
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal %+v: %v", resp, err)
		return nil
	}

	ev := Event{
		Type:     eventEndedCalculation,
		Contents: bs,
	}

	c.send(&ev)

	return nil
}
