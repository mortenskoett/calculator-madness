package websocket

import (
	"encoding/json"
	"log"
	"shared/queue"
)

// Callback to handle incoming calculation ended events on NSQ.
// Warning: When returning error here, the message will be resent with a backoff.
func (m *Manager) NSQCalcProgressHandler(msg *queue.CalcProgressMessage) error {
	if msg == nil {
		return nil
	}

	log.Printf("received progress calculation for client %v for calculation %v", msg.ClientID, msg.CalculationID)

	c, ok := m.clients[msg.ClientID]
	if !ok {
		log.Println("failed to handle progress message: client did not exist")
		return nil
	}

	resp := ProgressCalculationResponse{
		ID: msg.CalculationID,
		Progress: Progress{
			Current: msg.Status.Progress.Current,
			Outof:   msg.Status.Progress.Outof,
		},
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal progress response: %+v: %v", resp, err)
		return nil
	}

	ev := Event{
		Type:     eventProgressCalculation,
		Contents: bs,
	}

	c.send(&ev)

	return nil
}

// Callback to handle incoming calculation ended events on NSQ.
// Warning: When returning error here, the message will be resent with a backoff.
func (m *Manager) NSQCalcEndedHandler(msg *queue.CalcEndedMessage) error {
	if msg == nil {
		return nil
	}

	log.Printf("received ended calculation for client %v for calculation %v", msg.ClientID, msg.CalculationID)

	c, ok := m.clients[msg.ClientID]
	if !ok {
		log.Println("failed to handle ended calculation: client did not exist")
		return nil
	}

	resp := EndCalculationResponse{
		ID:     msg.CalculationID,
		Result: msg.Result,
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal ended calc response %+v: %v", resp, err)
		return nil
	}

	ev := Event{
		Type:     eventEndedCalculation,
		Contents: bs,
	}

	c.send(&ev)

	return nil
}
