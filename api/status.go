// Package api provides status endpoint
package api

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// StatusResponse represents status
type StatusResponse struct {
	Status    string `json:"status"`
	GoVersion string `json:"go_version"`
	Uptime    string `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

// StatusHandler handles status
type StatusHandler struct {
	start time.Time
}

// NewStatusHandler creates handler
func NewStatusHandler() *StatusHandler {
	return &StatusHandler{start: time.Now()}
}

// ServeHTTP serves
func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatusResponse{
		Status:    "running",
		GoVersion: runtime.Version(),
		Uptime:    time.Since(h.start).String(),
		Timestamp: time.Now(),
	})
}
