package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

type removeFunc func(*client)

type client struct {
	connection *websocket.Conn
	rmFunc     removeFunc
	outbox     chan []byte
}

// newClient instantiates an incoming websocket connection client. It needs a function to remove
// itself from the manager when it is done working.
func newClient(conn *websocket.Conn, rm removeFunc) *client {
	return &client{
		connection: conn,
		rmFunc:     rm,
		outbox:     make(chan []byte),
	}
}

// readMessages blocks and continues to read messages for this client.
func (c *client) readMessages() {
	defer c.rmFunc(c)
	for {
		mtype, p, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
			websocket.CloseGoingAway, 
			websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		// TODO: Handle messages
		log.Println("MessageType: ", mtype, "Payload: ", string(p))

		// TODO: Should be removed when done testing
		c.outbox <- p
	}
}

func (c *client) writeMessages() {
	defer c.rmFunc(c)
	for {
		select {
		case msg, ok := <-c.outbox:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("failed to write channel closed message: ", err)
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("failed to send message %s from client: %v", msg, err)
			}
		}
	}
}
