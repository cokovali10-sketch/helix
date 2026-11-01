// Package agent provides health checking
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthChecker monitors agent health
type HealthChecker struct {
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// HealthCheck represents a check
type HealthCheck struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	LastChecked time.Time `json:"last_checked"`
	Duration    time.Duration `json:"duration"`
}

// NewHealthChecker creates checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{checks: make(map[string]HealthCheck)}
}

// Register registers a check
func (hc *HealthChecker) Register(name string, fn func() (bool, string)) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	start := time.Now()
	ok, msg := fn()
	status := "healthy"
	if !ok {
		status = "unhealthy"
	}
	hc.checks[name] = HealthCheck{
		Name:        name,
		Status:      status,
		Message:     msg,
		LastChecked: time.Now(),
		Duration:    time.Since(start),
	}
}

// Run runs all checks
func (hc *HealthChecker) Run() map[string]HealthCheck {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	out := make(map[string]HealthCheck, len(hc.checks))
	for k, v := range hc.checks {
		out[k] = v
	}
	return out
}

// IsHealthy returns overall health
func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	for _, c := range hc.checks {
		if c.Status != "healthy" {
			return false
		}
	}
	return true
}

// Handler returns http handler
func (hc *HealthChecker) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !hc.IsHealthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(hc.Run())
}

// String summary
func (hc *HealthChecker) String() string {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return fmt.Sprintf("HealthChecker: %d checks", len(hc.checks))
}
