// Package catalog provides metrics reporter
package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CatalogMetricsReporter reports metrics
type CatalogMetricsReporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewCatalogMetricsReporter creates reporter
func NewCatalogMetricsReporter() *CatalogMetricsReporter {
	return &CatalogMetricsReporter{metrics: make(map[string]float64)}
}

// Report reports metric
func (r *CatalogMetricsReporter) Report(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = value
}

// Handler serves
func (r *CatalogMetricsReporter) Handler(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":   r.metrics,
		"timestamp": time.Now(),
	})
}

// Count returns count
func (r *CatalogMetricsReporter) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.metrics)
}

// String returns summary
func (r *CatalogMetricsReporter) String() string {
	return fmt.Sprintf("CatalogMetricsReporter: %d metrics", r.Count())
}
