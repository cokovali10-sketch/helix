// Package ingestion provides metrics reporter
package ingestion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// IngestMetricsReporter reports metrics
type IngestMetricsReporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewIngestMetricsReporter creates reporter
func NewIngestMetricsReporter() *IngestMetricsReporter {
	return &IngestMetricsReporter{metrics: make(map[string]float64)}
}

// Report reports metric
func (r *IngestMetricsReporter) Report(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = value
}

// Handler serves
func (r *IngestMetricsReporter) Handler(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":   r.metrics,
		"timestamp": time.Now(),
	})
}

// Count returns count
func (r *IngestMetricsReporter) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.metrics)
}

// String returns summary
func (r *IngestMetricsReporter) String() string {
	return fmt.Sprintf("IngestMetricsReporter: %d metrics", r.Count())
}
