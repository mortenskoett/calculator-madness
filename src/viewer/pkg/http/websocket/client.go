package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Notes: PING and PONG messages are described in the RFC. In summary, peers (including the browser)
// automatically respond to a PING message with a PONG message.

// The best practice for detecting a dead client is to read with a deadline. If the client
// application does not send messages frequently enough for the deadline you want, then send PING
// messages to induce the client to send a PONG. Update the deadline in the pong handler and after
// reading a message.

var (
	readDeadline = 30 * time.Second        // If no ping happens before deadline the socket is closed.
	pingInterval = (readDeadline * 8) / 10 // Calculate percentage of readDeadline.
)

type cleanupFn func(*client)

type client struct {
	id         string
	connection *websocket.Conn
	router     *EventRouter
	cleanupFn  cleanupFn
	outbox     chan Event
}

// newClient instantiates an incoming websocket connection client. It needs a function to remove
// itself from the manager when it is done working.
func newClient(conn *websocket.Conn, router *EventRouter, cleanFn cleanupFn) *client {
	c := &client{
		id:         uuid.NewString(),
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

// Send puts the given Event in the outbox of the client to be sent to the UI.
func (c *client) send(ev *Event) {
	c.outbox <- *ev
}

// ReadMessages blocks and continues to read messages for this client. Should be called as
// a goroutine.
func (c *client) readMessages() {
	defer c.cleanupFn(c)

	for {
		_, bs, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("error reading incoming ws message: %v", err)
			}
			continue
		}

		var req Event

		if err := json.Unmarshal(bs, &req); err != nil {
			log.Println("failed to unmarshal incoming ws event:", string(bs))
			continue
		}

		// Route events in separate goroutines.
		go func() {
			if err := c.router.route(&req, c); err != nil {
				log.Println("failed to route incoming ws event:", err)
			}
		}()
	}
}

// WriteMessages empties the outbox continously and sends the found messages to the client. Should
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
				log.Printf("failed to send ping message to client %v: %v", c.id, err)
				return
			}
		}
	}
}
