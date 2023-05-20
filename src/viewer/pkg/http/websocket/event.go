package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"viewer/api/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	// When a calculation is started and received from the UI.
	eventStartCalculation = "start_calculation"

	// When a new calculation is sent back to the UI to be shown.
	eventNewCalculation = "new_calculation"

	// eventProgressCalculation = "progress_calculation"

	// When a calculation ended is sent back to UI.
	eventEndedCalculation = "ended_calculation"
)

const (
	initialStepCount int = 5
)

var (
	errInvalidEvent = errors.New("the encountered event is invalid")
)

type Calculator interface {
	Run(context.Context, *pb.RunCalculationRequest, ...grpc.CallOption) (*pb.RunCalculationResponse, error)
}

// handler handles each incoming event type.
type handler func(e *Event, c *client) error

// EventRouter routes incoming websocket events to the right handlers.
type EventRouter struct {
	handlers   map[string]handler
	calculator Calculator
}

// NewEventRouter returns a router instance to handle incoming events.
func NewEventRouter(calc Calculator) *EventRouter {
	r := EventRouter{
		handlers:   map[string]handler{},
		calculator: calc,
	}
	r.attach(eventStartCalculation, r.handleStartCalculation)
	return &r
}

func (r *EventRouter) route(ev *Event, c *client) error {
	handler, ok := r.handlers[ev.Type]
	if !ok {
		return errInvalidEvent
	}
	return handler(ev, c)
}

// attach affiliates an event type with a concrete handler implementation.
func (r *EventRouter) attach(eventType string, h handler) {
	r.handlers[eventType] = h
}

/* Event handlers */

func (r *EventRouter) handleStartCalculation(ev *Event, c *client) error {
	log.Println("router received event:", ev.Type)

	var req StartCalculationRequest
	if err := json.Unmarshal(ev.Contents, &req); err != nil {
		return fmt.Errorf("failed to unmarshal %+v: %w", req, err)
	}

	resp := StartCalculationResponse{
		ID:          uuid.NewString(),
		CreatedTime: time.Now(),
		Equation:    req.Equation,
		Progress: Progress{
			Current: 0,
			Outof:   initialStepCount,
		},
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal %+v: %w", resp, err)
	}

	outEvent := Event{
		Type:     eventNewCalculation,
		Contents: bs,
	}

	// Send event back to UI using websocket.
	c.send(&outEvent)

	// Send equation to calculation backend.
	calcresp, err := r.calculator.Run(context.TODO(), &pb.RunCalculationRequest{
		ClientId: c.id,
		Equation: &pb.Equation{
			Id:    resp.ID,
			Value: req.Equation,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send calculation request: %w", err)
	}

	e := calcresp.Error
	if e != nil {
		return fmt.Errorf("failed to process calculation request: id: %v: message: %v", e.GetCode(), e.GetMessage())
	}
	log.Println("router successfully sent calculation request to backend")

	return nil
}
