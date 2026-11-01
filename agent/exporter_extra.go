// Package agent provides additional exporter features
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ExporterExtra extends exporter
type ExporterExtra struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewExporterExtra creates exporter
func NewExporterExtra() *ExporterExtra {
	return &ExporterExtra{metrics: make(map[string]float64)}
}

// Set sets metric
func (e *ExporterExtra) Set(name string, value float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics[name] = value
}

// Get gets metric
func (e *ExporterExtra) Get(name string) (float64, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.metrics[name]
	return v, ok
}

// Increment increments
func (e *ExporterExtra) Increment(name string, delta float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics[name] += delta
}

// Handler serves
func (e *ExporterExtra) Handler(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e.metrics)
}

// Prometheus formats
func (e *ExporterExtra) Prometheus() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := ""
	for k, v := range e.metrics {
		out += fmt.Sprintf("%s %f\n", k, v)
	}
	return out
}

// Count returns count
func (e *ExporterExtra) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.metrics)
}

// Reset resets
func (e *ExporterExtra) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics = make(map[string]float64)
}

// Timestamp returns now
func (e *ExporterExtra) Timestamp() time.Time {
	return time.Now()
}
