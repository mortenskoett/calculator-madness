package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"viewer/api/pb"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type Calculator interface {
	Run(context.Context, *pb.RunCalculationRequest, ...grpc.CallOption) (*pb.RunCalculationResponse, error)
}

// Manager manages incoming clients that have created a websocket to the server. It facilitates
// concurrent R/W access for the clients and maintains the state of open connections.
type Manager struct {
	upgrader   websocket.Upgrader
	clients    map[*client]bool
	router     *eventRouter
	calculator Calculator
	sync.RWMutex
}

func NewManager(calc Calculator) *Manager {
	return &Manager{
		upgrader:   websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		clients:    map[*client]bool{},
		router:     newEventRouter(),
		RWMutex:    sync.RWMutex{},
		calculator: calc,
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
	m.clients[c] = true
}

// remove makes sure a client does not linger around when it is done
func (m *Manager) remove(c *client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}
