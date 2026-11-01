// Package api provides health endpoint
package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthStatus represents health
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// HealthHandler handles health
type HealthHandler struct {
	start time.Time
}

// NewHealthHandler creates handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{start: time.Now()}
}

// ServeHTTP serves
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.start).String(),
	})
}

// Liveness returns liveness
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Readiness returns readiness
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
