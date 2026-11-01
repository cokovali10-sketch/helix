// Package api provides metrics reporter
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// APIMetricsReporter reports metrics
type APIMetricsReporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewAPIMetricsReporter creates reporter
func NewAPIMetricsReporter() *APIMetricsReporter {
	return &APIMetricsReporter{metrics: make(map[string]float64)}
}

// Report reports metric
func (r *APIMetricsReporter) Report(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = value
}

// Handler serves
func (r *APIMetricsReporter) Handler(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":   r.metrics,
		"timestamp": time.Now(),
	})
}

// Count returns count
func (r *APIMetricsReporter) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.metrics)
}

// String returns summary
func (r *APIMetricsReporter) String() string {
	return fmt.Sprintf("APIMetricsReporter: %d metrics", r.Count())
}
