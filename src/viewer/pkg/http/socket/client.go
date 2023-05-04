package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

type removeFunc func(*client)

type client struct {
	connection *websocket.Conn
	rmFunc     removeFunc
}

// newClient instantiates an incoming websocket connection client. It needs a function to remove
// itself from the manager when it is done working.
func newClient(conn *websocket.Conn, rm removeFunc) *client {
	return &client{
		connection: conn,
		rmFunc:     rm,
	}
}

// readMessages blocks and continues to read messages from this client.
func (c *client) readMessages() {
	defer c.rmFunc(c)
	for {
		mtype, p, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		log.Println("MessageType: ", mtype)
		log.Println("Payload: ", string(p))
	}
}
