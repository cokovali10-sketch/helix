// Package api provides additional API handlers and middleware
package api

import (
	"encoding/json"
	"fmt"
	"io"
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
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		metrics: make(map[string]*Metric),
	}
}

// RegisterMetric registers a metric
func (mh *MetricsHandler) RegisterMetric(name, metricType string, labels map[string]string) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	
	key := metricKey(name, labels)
	mh.metrics[key] = &Metric{
		Name:      name,
		Type:      metricType,
		Labels:    labels,
		Timestamp: time.Now(),
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
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// HealthCheck represents a health check
type HealthCheck struct {
	Name     string
	Check    func() error
	Interval time.Duration
	Timeout  time.Duration
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		checks: make(map[string]HealthCheck),
	}
}

// RegisterCheck registers a health check
func (hh *HealthHandler) RegisterCheck(name string, check func() error, interval, timeout time.Duration) {
	hh.mu.Lock()
	defer hh.mu.Unlock()
	hh.checks[name] = HealthCheck{
		Name:     name,
		Check:    check,
		Interval: interval,
		Timeout:  timeout,
	}
}

// ServeHTTP serves health endpoint
func (hh *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.mu.RLock()
	defer hh.mu.RUnlock()
	
	status := "healthy"
	checks := make(map[string]string)
	
	for name, check := range hh.checks {
		if err := check.Check(); err != nil {
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
	ready bool
	mu    sync.RWMutex
}

// NewReadyHandler creates a new ready handler
func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{ready: false}
}

// SetReady sets readiness state
func (rh *ReadyHandler) SetReady(ready bool) {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	rh.ready = ready
}

// ServeHTTP serves readiness endpoint
func (rh *ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.mu.RLock()
	defer rh.mu.RUnlock()
	
	if !rh.ready {
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
		"version":    vh.version,
		"buildTime":  vh.buildTime,
		"gitCommit":  vh.gitCommit,
		"goVersion":  "go1.21",
	})
}

// ConfigHandler handles configuration endpoints
type ConfigHandler struct {
	config map[string]interface{}
	mu     sync.RWMutex
}

// NewConfigHandler creates a new config handler
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		config: make(map[string]interface{}),
	}
}

// Set sets a config value
func (ch *ConfigHandler) Set(key string, value interface{}) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.config[key] = value
}

// Get gets a config value
func (ch *ConfigHandler) Get(key string) (interface{}, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	v, ok := ch.config[key]
	return v, ok
}

// ServeHTTP serves config endpoint
func (ch *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ch.config)
}

// DebugHandler handles debug endpoints
type DebugHandler struct{}

// NewDebugHandler creates a new debug handler
func NewDebugHandler() *DebugHandler {
	return &DebugHandler{}
}

// ServeHTTP serves debug endpoint
func (dh *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"goroutines": 0,
		"memory":     map[string]interface{}{},
		"gc":         map[string]interface{}{},
	})
}

// PprofHandler handles pprof endpoints
type PprofHandler struct{}

// NewPprofHandler creates a new pprof handler
func NewPprofHandler() *PprofHandler {
	return &PprofHandler{}
}

// ServeHTTP serves pprof endpoint
func (ph *PprofHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pprof endpoint"))
}

// WebhookHandler handles webhook endpoints
type WebhookHandler struct {
	handlers map[string]func([]byte) error
	mu       sync.RWMutex
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		handlers: make(map[string]func([]byte) error),
	}
}

// RegisterHandler registers a webhook handler
func (wh *WebhookHandler) RegisterHandler(path string, handler func([]byte) error) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	wh.handlers[path] = handler
}

// ServeHTTP serves webhook endpoint
func (wh *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	wh.mu.RLock()
	handler, ok := wh.handlers[path]
	wh.mu.RUnlock()

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}