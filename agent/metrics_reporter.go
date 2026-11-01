// Package agent provides metrics reporter
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// MetricsReporter reports metrics
type MetricsReporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewMetricsReporter creates reporter
func NewMetricsReporter() *MetricsReporter {
	return &MetricsReporter{metrics: make(map[string]float64)}
}

// Report reports metric
func (mr *MetricsReporter) Report(name string, value float64) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.metrics[name] = value
}

// Handler serves
func (mr *MetricsReporter) Handler(w http.ResponseWriter, r *http.Request) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":   mr.metrics,
		"timestamp": time.Now(),
	})
}

// Count returns count
func (mr *MetricsReporter) Count() int {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	return len(mr.metrics)
}

// String returns summary
func (mr *MetricsReporter) String() string {
	return fmt.Sprintf("MetricsReporter: %d metrics", mr.Count())
}
