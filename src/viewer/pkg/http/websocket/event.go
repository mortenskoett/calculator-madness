package websocket

import (
	"encoding/json"
	"errors"
	"log"
)

// Event type.
const (
	eventNewCalculation = "new_calculation"
)

var (
	errInvalidEvent = errors.New("the encountered event is invalid")
)

type Event struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"` // Should be marshalled individually by handlers.
}

// handler handles each incoming event type.
type handler func(e *Event, c *client) error

type eventRouter struct {
	handlers map[string]handler
}

// newEventRouter returns a router instance to handle incoming events.
func newEventRouter() *eventRouter {
	r := eventRouter{
		handlers: map[string]handler{},
	}
	r.attach(eventNewCalculation, r.handleNewCalculation)
	return &r
}

func (r *eventRouter) route(ev *Event, c *client) error {
	handler, ok := r.handlers[ev.Type]
	if !ok {
		return errInvalidEvent
	}
	return handler(ev, c)
}

// attach affiliates an event type with a concrete handler implementation.
func (r *eventRouter) attach(eventType string, h handler) {
	r.handlers[eventType] = h
}

type NewCalculationEvent struct {
	ID       string `json:"id"`
	Equation string `json:"equation"`
}

func (r *eventRouter) handleNewCalculation(ev *Event, c *client) error {
	log.Println("event router called for event:", ev)
	return nil
}

