package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohadayo/pulsehub/services/gateway/models"
)

func TestHealth(t *testing.T) {
	h := New("http://localhost:8001")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp models.HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
	if resp.Service != "PulseHub Gateway" {
		t.Errorf("expected service 'PulseHub Gateway', got %q", resp.Service)
	}
}

func TestReceiveEvent_Success(t *testing.T) {
	h := New("http://localhost:8001")

	body := map[string]interface{}{
		"event_type": "click",
		"source":     "test",
		"payload":    map[string]interface{}{"button": "submit"},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ReceiveEvent(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", w.Code)
	}

	var event models.Event
	if err := json.NewDecoder(w.Body).Decode(&event); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if event.ID == "" {
		t.Error("expected event ID to be set")
	}
	if event.EventType != "click" {
		t.Errorf("expected event_type 'click', got %q", event.EventType)
	}
	if event.Timestamp == nil {
		t.Error("expected timestamp to be set")
	}
}

func TestReceiveEvent_InvalidMethod(t *testing.T) {
	h := New("http://localhost:8001")
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	h.ReceiveEvent(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestReceiveEvent_InvalidJSON(t *testing.T) {
	h := New("http://localhost:8001")
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	h.ReceiveEvent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestReceiveEvent_MissingFields(t *testing.T) {
	h := New("http://localhost:8001")

	body := map[string]interface{}{"event_type": "click"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.ReceiveEvent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetStats(t *testing.T) {
	h := New("http://localhost:8001")
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	h.GetStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var stats models.GatewayStats
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if stats.EventsReceived != 0 {
		t.Errorf("expected 0 events received, got %d", stats.EventsReceived)
	}
}

func TestReceiveEvent_IncrementsCounter(t *testing.T) {
	h := New("http://localhost:8001")

	body := map[string]interface{}{
		"event_type": "click",
		"source":     "test",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ReceiveEvent(w, req)

	if h.eventsReceived.Load() != 1 {
		t.Errorf("expected 1 event received, got %d", h.eventsReceived.Load())
	}
}
