// Package api provides additional API handlers
// including metrics, health, and monitoring handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// MetricsHandler handles metrics endpoints
type MetricsHandler struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
}

// Metric represents an API metric
type Metric struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
	Help      string            `json:"help"`
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		metrics: make(map[string]*Metric),
	}
}

// RegisterMetric registers a metric
func (mh *MetricsHandler) RegisterMetric(name, metricType, help string, labels map[string]string) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	key := metricKey(name, labels)
	mh.metrics[key] = &Metric{
		Name:      name,
		Type:      metricType,
		Labels:    labels,
		Timestamp: time.Now(),
		Help:      help,
	}
}

// IncrementCounter increments a counter metric
func (mh *MetricsHandler) IncrementCounter(name string, labels map[string]string, value float64) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	key := metricKey(name, labels)
	if m, ok := mh.metrics[key]; ok {
		m.Value += value
		m.Timestamp = time.Now()
	} else {
		mh.metrics[key] = &Metric{
			Name:      name,
			Type:      "counter",
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// SetGauge sets a gauge metric
func (mh *MetricsHandler) SetGauge(name string, labels map[string]string, value float64) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	key := metricKey(name, labels)
	if m, ok := mh.metrics[key]; ok {
		m.Value = value
		m.Timestamp = time.Now()
	} else {
		mh.metrics[key] = &Metric{
			Name:      name,
			Type:      "gauge",
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// ObserveHistogram observes a histogram value
func (mh *MetricsHandler) ObserveHistogram(name string, labels map[string]string, value float64, buckets []float64) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	key := metricKey(name, labels)
	if _, ok := mh.metrics[key]; !ok {
		mh.metrics[key] = &Metric{
			Name:      name,
			Type:      "histogram",
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// GetMetric gets a metric
func (mh *MetricsHandler) GetMetric(name string, labels map[string]string) (*Metric, bool) {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	key := metricKey(name, labels)
	m, ok := mh.metrics[key]
	return m, ok
}

// GetAllMetrics returns all metrics
func (mh *MetricsHandler) GetAllMetrics() []*Metric {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	metrics := make([]*Metric, 0, len(mh.metrics))
	for _, m := range mh.metrics {
		metrics = append(metrics, m)
	}
	return metrics
}

// ServeHTTP serves metrics endpoint
func (mh *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := mh.GetAllMetrics()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"total":   len(metrics),
	})
}

func metricKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += fmt.Sprintf(",%s=%s", k, v)
	}
	return key
}

// HealthHandler handles health endpoints
type HealthHandler struct {
	checks map[string]func() error
	mu     sync.RWMutex
	ready  bool
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		checks: make(map[string]func() error),
		ready:  false,
	}
}

// RegisterCheck registers a health check
func (hh *HealthHandler) RegisterCheck(name string, check func() error) {
	hh.mu.Lock()
	defer hh.mu.Unlock()
	hh.checks[name] = check
}

// SetReady sets readiness state
func (hh *HealthHandler) SetReady(ready bool) {
	hh.mu.Lock()
	defer hh.mu.Unlock()
	hh.ready = ready
}

// ServeHTTP serves health endpoint
func (hh *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.mu.RLock()
	defer hh.mu.RUnlock()

	status := "healthy"
	checks := make(map[string]string)

	for name, check := range hh.checks {
		if err := check(); err != nil {
			status = "unhealthy"
			checks[name] = err.Error()
		} else {
			checks[name] = "ok"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"checks":  checks,
		"version": "0.1.0",
	})
}

// ReadyHandler handles readiness endpoint
type ReadyHandler struct {
	handler *HealthHandler
}

// NewReadyHandler creates a new ready handler
func NewReadyHandler(handler *HealthHandler) *ReadyHandler {
	return &ReadyHandler{handler: handler}
}

// ServeHTTP serves readiness endpoint
func (rh *ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.handler.mu.RLock()
	defer rh.handler.mu.RUnlock()

	if !rh.handler.ready {
		http.Error(w, "Not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// LiveHandler handles liveness endpoint
type LiveHandler struct{}

// NewLiveHandler creates a new live handler
func NewLiveHandler() *LiveHandler {
	return &LiveHandler{}
}

// ServeHTTP serves liveness endpoint
func (lh *LiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// VersionHandler handles version endpoint
type VersionHandler struct {
	version   string
	buildTime string
	gitCommit string
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(version, buildTime, gitCommit string) *VersionHandler {
	return &VersionHandler{
		version:   version,
		buildTime: buildTime,
		gitCommit: gitCommit,
	}
}

// ServeHTTP serves version endpoint
func (vh *VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version":   vh.version,
		"buildTime": vh.buildTime,
		"gitCommit": vh.gitCommit,
		"goVersion": "go1.21",
	})
}

// StatusHandler handles status endpoint
type StatusHandler struct {
	startTime time.Time
	version   string
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(version string) *StatusHandler {
	return &StatusHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// ServeHTTP serves status endpoint
func (sh *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"version":   sh.version,
		"uptime":    time.Since(sh.startTime).String(),
		"startTime": sh.startTime,
		"goVersion": "go1.21",
		"os":        "linux",
		"arch":      "amd64",
		"numCPU":    4,
	})
}