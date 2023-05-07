package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type removeFunc func(*client)

type client struct {
	connection *websocket.Conn
	rmFunc     removeFunc
	outbox     chan Event
	router     *eventRouter
}

// newClient instantiates an incoming websocket connection client. It needs a function to remove
// itself from the manager when it is done working.
func newClient(conn *websocket.Conn, router *eventRouter, rm removeFunc) *client {
	return &client{
		connection: conn,
		rmFunc:     rm,
		outbox:     make(chan Event),
		router:     router,
	}
}

// readMessages blocks and continues to read messages for this client.
func (c *client) readMessages() {
	defer c.rmFunc(c)
	for {
		_, bs, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var req Event

		if err := json.Unmarshal(bs, &req); err != nil {
			log.Println("failed unmarshal incoming event:", string(bs))
			break
		}

		if err := c.router.route(&req, c); err != nil {
			log.Println("failed route incoming event:", err)
			break
		}

		// // TODO: Handle messages
		// log.Println("MessageType: ", mtype, "Content: ", string(p))
		// // TODO: Should be removed when done testing
		c.outbox <- req
	}
}

func (c *client) writeMessages() {
	defer c.rmFunc(c)
	for {
		select {
		case event, ok := <-c.outbox:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("failed to write channel closed message: ", err)
				}
				return
			}

			bs, err := json.Marshal(event)
			if err != nil {
				log.Println("failed to marshal event before sending:", event)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, bs); err != nil {
				log.Printf("failed to send message %s from client: %v", bs, err)
			}
		}
	}
}
