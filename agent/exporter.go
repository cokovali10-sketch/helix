// Package agent provides metrics exporter
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Exporter exports metrics
type Exporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewExporter creates exporter
func NewExporter() *Exporter {
	return &Exporter{metrics: make(map[string]float64)}
}

// Set sets metric
func (e *Exporter) Set(name string, value float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics[name] = value
}

// Get gets metric
func (e *Exporter) Get(name string) (float64, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.metrics[name]
	return v, ok
}

// Increment increments metric
func (e *Exporter) Increment(name string, delta float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics[name] += delta
}

// Handler serves metrics
func (e *Exporter) Handler(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e.metrics)
}

// Prometheus formats metrics
func (e *Exporter) Prometheus() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := ""
	for k, v := range e.metrics {
		out += fmt.Sprintf("%s %f\n", k, v)
	}
	return out
}

// Count returns count
func (e *Exporter) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.metrics)
}

// Reset resets metrics
func (e *Exporter) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metrics = make(map[string]float64)
}

// Timestamp returns current
func (e *Exporter) Timestamp() time.Time {
	return time.Now()
}
