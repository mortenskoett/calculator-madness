package websocket

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type     string          `json:"type"`
	Contents json.RawMessage `json:"contents"` // Should be marshalled individually by handlers.
}

type Progress struct {
	Current int `json:"current"`
	Outof   int `json:"outof"`
}

// Request from UI.
type StartCalculationRequest struct {
	Equation string `json:"equation"`
}

// Response to UI.
type StartCalculationResponse struct {
	ID          string    `json:"id"`
	Equation    string    `json:"equation"`
	CreatedTime time.Time `json:"created_time"`
	Progress    Progress  `json:"progress"`
}

// Response to UI.
type EndCalculationResponse struct {
	ID        string    `json:"id"`
	Result    float64   `json:"result"`
}
