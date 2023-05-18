package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager manages incoming clients that have created a websocket to the server. It facilitates
// concurrent R/W access for the clients and maintains the state of open connections.
type Manager struct {
	upgrader websocket.Upgrader
	clients  map[string]*client
	router   *EventRouter
	sync.RWMutex
}

// NewManager returns a websocket manager to concurrently handle socket connections and events
func NewManager(router *EventRouter) *Manager {
	return &Manager{
		upgrader: websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		clients:  map[string]*client{},
		router:   router,
		RWMutex:  sync.RWMutex{},
	}
}

func (m *Manager) HandleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := m.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("websocket upgrade failed: ", err)
			return
		}
		log.Print("websocket upgrade successfully made")

		client := newClient(conn, m.router, m.remove)
		m.add(client)

		go client.readMessages()
		go client.writeMessages()
	}
}

// add marks the client as active in the manager
func (m *Manager) add(c *client) {
	m.Lock()
	defer m.Unlock()
	m.clients[c.id] = c
}

// remove makes sure a client does not linger around when it is done
func (m *Manager) remove(c *client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[c.id]; ok {
		c.connection.Close()
		delete(m.clients, c.id)
	}
}
