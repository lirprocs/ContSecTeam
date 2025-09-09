package handler

import (
	"ContSecTeam/internal/model"
	"ContSecTeam/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T, queueSize, workers int) *service.Service {
	t.Helper()
	s := service.NewService(queueSize)
	s.Start(context.Background(), workers)
	t.Cleanup(func() { s.Stop() })
	return s
}

func TestHealthz(t *testing.T) {
	s := newTestServer(t, 1, 1)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	h.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestEnqueue_Success(t *testing.T) {
	s := newTestServer(t, 2, 1)
	h := NewHandler(s)

	payload := model.Task{ID: "t1", Payload: "data", MaxRetries: 0}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Enqueue(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 Accepted, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestEnqueue_MethodNotAllowed(t *testing.T) {
	s := newTestServer(t, 1, 1)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/enqueue", nil)
	rec := httptest.NewRecorder()

	h.Enqueue(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestEnqueue_BadRequest_NoID(t *testing.T) {
	s := newTestServer(t, 1, 1)
	h := NewHandler(s)

	payload := map[string]any{"payload": "data"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Enqueue(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEnqueue_QueueOverflow(t *testing.T) {
	s := newTestServer(t, 1, 1)
	h := NewHandler(s)

	payload1 := model.Task{ID: "t1", MaxRetries: 0}
	b1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	h.Enqueue(rec1, req1)
	if rec1.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec1.Code)
	}

	time.Sleep(50 * time.Millisecond)

	payload2 := model.Task{ID: "t2", MaxRetries: 0}
	b2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	h.Enqueue(rec2, req2)
	if rec2.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec2.Code)
	}

	payload3 := model.Task{ID: "t3", MaxRetries: 0}
	b3, _ := json.Marshal(payload3)
	req3 := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(b3))
	req3.Header.Set("Content-Type", "application/json")
	rec3 := httptest.NewRecorder()
	h.Enqueue(rec3, req3)
	if rec3.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec3.Code)
	}
}
