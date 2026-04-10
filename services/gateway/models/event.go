package models

import "time"

type Event struct {
	ID        string            `json:"id,omitempty"`
	EventType string            `json:"event_type"`
	Source    string            `json:"source"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type GatewayStats struct {
	EventsReceived int64    `json:"events_received"`
	EventsForwarded int64   `json:"events_forwarded"`
	EventsFailed   int64    `json:"events_failed"`
	Uptime         string   `json:"uptime"`
}
