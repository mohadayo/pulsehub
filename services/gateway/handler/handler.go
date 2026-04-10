package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mohadayo/pulsehub/services/gateway/models"
)

type Handler struct {
	analyzerURL     string
	eventsReceived  atomic.Int64
	eventsForwarded atomic.Int64
	eventsFailed    atomic.Int64
	startTime       time.Time
}

func New(analyzerURL string) *Handler {
	return &Handler{
		analyzerURL: analyzerURL,
		startTime:   time.Now(),
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := models.HealthResponse{
		Status:  "ok",
		Service: "PulseHub Gateway",
		Version: "1.0.0",
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse{
			Error:   "method_not_allowed",
			Message: "Only POST is allowed",
		})
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("[ERROR] Failed to decode event: %v", err)
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_payload",
			Message: "Could not parse event JSON",
		})
		return
	}

	if event.EventType == "" || event.Source == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "event_type and source are required",
		})
		return
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	event.Timestamp = &now

	h.eventsReceived.Add(1)
	log.Printf("[INFO] Received event id=%s type=%s source=%s", event.ID, event.EventType, event.Source)

	writeJSON(w, http.StatusAccepted, event)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Round(time.Second).String()
	stats := models.GatewayStats{
		EventsReceived:  h.eventsReceived.Load(),
		EventsForwarded: h.eventsForwarded.Load(),
		EventsFailed:    h.eventsFailed.Load(),
		Uptime:          uptime,
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR] Failed to write response: %v", err)
	}
}
