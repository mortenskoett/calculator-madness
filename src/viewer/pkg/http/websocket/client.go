package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Notes:
// PING and PONG messages are described in the RFC. In summary, peers (including the browser)
// automatically respond to a PING message with a PONG message.

// The best practice for detecting a dead client is to read with a deadline. If the client
// application does not send messages frequently enough for the deadline you want, then send PING
// messages to induce the client to send a PONG. Update the deadline in the pong handler and after
// reading a message.

var (
	readDeadline = 20 * time.Second
	pingInterval = (readDeadline * 9) / 10 // Calculate 90% without decimals.
)

type cleanupFn func(*client)

type client struct {
	connection *websocket.Conn
	router     *eventRouter
	cleanupFn  cleanupFn
	outbox     chan Event
}

// newClient instantiates an incoming websocket connection client. It needs a function to remove
// itself from the manager when it is done working.
func newClient(conn *websocket.Conn, router *eventRouter, cleanFn cleanupFn) *client {
	c := &client{
		connection: conn,
		router:     router,
		cleanupFn:  cleanFn,
		outbox:     make(chan Event),
	}
	c.connection.SetPongHandler(c.pongHandler)
	c.connection.SetReadLimit(512)
	return c
}

func (c *client) pongHandler(s string) error {
	log.Println("pong received")
	c.setReadDeadline(readDeadline)
	return nil
}

func (c *client) setReadDeadline(dur time.Duration) error {
	if err := c.connection.SetReadDeadline(time.Now().Add(dur)); err != nil {
		return fmt.Errorf("failed to set read deadline: %w", err)
	}
	return nil
}

// readMessages blocks and continues to read messages for this client. Should be called as
// a goroutine.
func (c *client) readMessages() {
	defer c.cleanupFn(c)

	// To detect a dead client.
	if err := c.setReadDeadline(readDeadline); err != nil {
		log.Println(err)
		return
	}

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

	}
}

// writeMessages empties the outbox continously and sends the found messages to the client. Should
// be called as goroutine.
func (c *client) writeMessages() {
	defer c.cleanupFn(c)
	pingTicker := time.NewTicker(pingInterval)

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
				log.Printf("failed to send message %s to client: %v", string(bs), err)
			}

		case <-pingTicker.C:
			log.Println("ping sent")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("failed to send ping message: ", err)
				return
			}
		}
	}
}
