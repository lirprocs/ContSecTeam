package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ContSecTeam/internal/model"
	"ContSecTeam/internal/service"
)

type Handler struct {
	Srv *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Srv: s}
}

func (h *Handler) Enqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if t.ID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.Srv.Enqueue(&t); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "enqueued task %s\n", t.ID)
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
