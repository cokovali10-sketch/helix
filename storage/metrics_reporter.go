// Package storage provides metrics reporter
package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// StorageMetricsReporter reports metrics
type StorageMetricsReporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewStorageMetricsReporter creates reporter
func NewStorageMetricsReporter() *StorageMetricsReporter {
	return &StorageMetricsReporter{metrics: make(map[string]float64)}
}

// Report reports metric
func (r *StorageMetricsReporter) Report(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = value
}

// Handler serves
func (r *StorageMetricsReporter) Handler(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":   r.metrics,
		"timestamp": time.Now(),
	})
}

// Count returns count
func (r *StorageMetricsReporter) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.metrics)
}

// String returns summary
func (r *StorageMetricsReporter) String() string {
	return fmt.Sprintf("StorageMetricsReporter: %d metrics", r.Count())
}
