package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	// When a calculation is started and received from the UI.
	eventStartCalculation = "start_calculation"

	// When a new calculation is sent back to the UI to be shown.
	eventNewCalculation = "new_calculation"

	// eventProgressCalculation = "progress_calculation"
	// eventFinishCalculation   = "finish_calculation"
)

const (
	defaultProgressSteps int = 5
)

var (
	errInvalidEvent = errors.New("the encountered event is invalid")
)

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
	r.attach(eventStartCalculation, r.handleStartCalculation)
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

/* Event handlers */

func (r *eventRouter) handleStartCalculation(ev *Event, c *client) error {
	log.Println("event router received event type:", ev.Type)

	var req StartCalculationRequest
	if err := json.Unmarshal(ev.Content, &req); err != nil {
		return fmt.Errorf("failed to unmarshal %+v: %w", req, err)
	}

	resp := StartCalculationResponse{
		Calculation: Calculation{
			ID:          uuid.NewString(),
			CreatedTime: time.Now().String(),
			Equation:    req.Equation,
			Progress: Progress{
				Current: 0,
				Outof:   defaultProgressSteps,
			},
			Result: "",
		},
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal %+v: %w", resp, err)
	}

	outEvent := Event{
		Type:    eventNewCalculation,
		Content: bs,
	}

	c.send(&outEvent)

	return nil
}
