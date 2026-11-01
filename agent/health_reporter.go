// Package agent provides health reporter
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthReporter reports health
type HealthReporter struct {
	status string
	checks map[string]string
}

// NewHealthReporter creates reporter
func NewHealthReporter() *HealthReporter {
	return &HealthReporter{status: "healthy", checks: make(map[string]string)}
}

// SetStatus sets status
func (hr *HealthReporter) SetStatus(status string) {
	hr.status = status
}

// SetCheck sets check
func (hr *HealthReporter) SetCheck(name, status string) {
	hr.checks[name] = status
}

// Status returns status
func (hr *HealthReporter) Status() string {
	return hr.status
}

// Checks returns checks
func (hr *HealthReporter) Checks() map[string]string {
	out := make(map[string]string, len(hr.checks))
	for k, v := range hr.checks {
		out[k] = v
	}
	return out
}

// Handler serves
func (hr *HealthReporter) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": hr.status,
		"checks": hr.checks,
		"time":   time.Now(),
	})
}

// IsHealthy returns health
func (hr *HealthReporter) IsHealthy() bool {
	return hr.status == "healthy"
}

// String returns summary
func (hr *HealthReporter) String() string {
	return fmt.Sprintf("HealthReporter: %s (%d checks)", hr.status, len(hr.checks))
}
