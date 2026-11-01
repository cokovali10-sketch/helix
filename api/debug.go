// Package api provides additional API handlers
// including tracing, profiling, and debug handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// TracingHandler handles tracing endpoints
type TracingHandler struct {
	spans  []*Span
	mu     sync.RWMutex
}

// Span represents a trace span
type Span struct {
	TraceID   string            `json:"traceId"`
	SpanID    string            `json:"spanId"`
	ParentID  string            `json:"parentId,omitempty"`
	Operation string            `json:"operation"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  time.Duration     `json:"duration"`
	Tags      map[string]string `json:"tags"`
	Logs      []SpanLog         `json:"logs"`
	Status    string            `json:"status"`
}

// SpanLog represents a span log entry
type SpanLog struct {
	Timestamp time.Time         `json:"timestamp"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
}

// NewTracingHandler creates a new tracing handler
func NewTracingHandler() *TracingHandler {
	return &TracingHandler{
		spans: make([]*Span, 0),
	}
}

// AddSpan adds a span
func (th *TracingHandler) AddSpan(span *Span) {
	th.mu.Lock()
	defer th.mu.Unlock()
	th.spans = append(th.spans, span)
}

// GetSpans returns all spans
func (th *TracingHandler) GetSpans(traceID string) []*Span {
	th.mu.RLock()
	defer th.mu.RUnlock()

	if traceID == "" {
		return th.spans
	}

	var result []*Span
	for _, span := range th.spans {
		if span.TraceID == traceID {
			result = append(result, span)
		}
	}
	return result
}

// GetSpan gets a span by ID
func (th *TracingHandler) GetSpan(spanID string) (*Span, bool) {
	th.mu.RLock()
	defer th.mu.RUnlock()

	for _, span := range th.spans {
		if span.SpanID == spanID {
			return span, true
		}
	}
	return nil, false
}

// ServeHTTP serves tracing endpoint
func (th *TracingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	traceID := r.URL.Query().Get("traceId")

	if r.Method == http.MethodGet {
		spans := th.GetSpans(traceID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"spans": spans,
			"total": len(spans),
		})
	} else if r.Method == http.MethodPost {
		var span Span
		if err := json.NewDecoder(r.Body).Decode(&span); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		th.AddSpan(&span)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(span)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ProfilingHandler handles profiling endpoints
type ProfilingHandler struct {
	profiles map[string][]byte
	mu       sync.RWMutex
}

// NewProfilingHandler creates a new profiling handler
func NewProfilingHandler() *ProfilingHandler {
	return &ProfilingHandler{
		profiles: make(map[string][]byte),
	}
}

// CaptureProfile captures a profile
func (ph *ProfilingHandler) CaptureProfile(profileType string) ([]byte, error) {
	var buf []byte
	var err error

	switch profileType {
	case "cpu":
		// In production, would start CPU profiling
		buf = []byte("cpu profile")
	case "heap":
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		data, _ := json.Marshal(m)
		buf = data
	case "goroutine":
	buf = []byte("pprof profile")
	case "block":
	buf = []byte("pprof profile")
	case "mutex":
	buf = []byte("pprof profile")
	default:
		return nil, fmt.Errorf("unknown profile type: %s", profileType)
	}

	ph.mu.Lock()
	ph.profiles[profileType] = buf
	ph.mu.Unlock()

	return buf, err
}

// ServeHTTP serves profiling endpoint
func (ph *ProfilingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	profileType := r.URL.Query().Get("type")
	if profileType == "" {
		profileType = "heap"
	}

	data, err := ph.CaptureProfile(profileType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.prof", profileType))
	w.Write(data)
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

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"heap_alloc":  m.HeapAlloc,
			"heap_sys":    m.HeapSys,
			"heap_idle":   m.HeapIdle,
			"heap_inuse":  m.HeapInuse,
			"heap_released": m.HeapReleased,
			"heap_objects": m.HeapObjects,
		},
		"gc": map[string]interface{}{
			"num_gc":        m.NumGC,
			"gc_cpu_fraction": m.GCCPUFraction,
			"last_gc":       time.Unix(0, int64(m.LastGC)).Format(time.RFC3339),
		},
		"runtime": map[string]interface{}{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"num_cpu":    runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
		},
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
	name := r.URL.Path[len("/debug/pprof/"):]
	if name == "" {
		// Index page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>pprof</h1></body></html>"))
		return
	}

	profile := pprof.Lookup(name)
	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	profile.WriteTo(w, 0)
}

// StatusHandler handles status endpoints
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
		"goVersion": runtime.Version(),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"numCPU":    runtime.NumCPU(),
	})
}

// ReadinessHandler handles readiness endpoint
type ReadinessHandler struct {
	checks map[string]func() error
	mu     sync.RWMutex
	ready  bool
}

// NewReadinessHandler creates a new readiness handler
func NewReadinessHandler() *ReadinessHandler {
	return &ReadinessHandler{
		checks: make(map[string]func() error),
		ready:  false,
	}
}

// RegisterCheck registers a readiness check
func (rh *ReadinessHandler) RegisterCheck(name string, check func() error) {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	rh.checks[name] = check
}

// SetReady sets readiness state
func (rh *ReadinessHandler) SetReady(ready bool) {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	rh.ready = ready
}

// ServeHTTP serves readiness endpoint
func (rh *ReadinessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.mu.RLock()
	defer rh.mu.RUnlock()

	if !rh.ready {
		http.Error(w, "Not ready", http.StatusServiceUnavailable)
		return
	}

	allReady := true
	checks := make(map[string]string)
	for name, check := range rh.checks {
		if err := check(); err != nil {
			allReady = false
			checks[name] = err.Error()
		} else {
			checks[name] = "ok"
		}
	}

	if !allReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ready":  allReady,
		"checks": checks,
	})
}

// LivenessHandler handles liveness endpoint
type LivenessHandler struct{}

// NewLivenessHandler creates a new liveness handler
func NewLivenessHandler() *LivenessHandler {
	return &LivenessHandler{}
}

// ServeHTTP serves liveness endpoint
func (lh *LivenessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// MetricsExporter exports metrics
type MetricsExporter struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

// NewMetricsExporter creates a new metrics exporter
func NewMetricsExporter() *MetricsExporter {
	return &MetricsExporter{
		metrics: make(map[string]float64),
	}
}

// Set sets a metric
func (me *MetricsExporter) Set(name string, value float64) {
	me.mu.Lock()
	defer me.mu.Unlock()
	me.metrics[name] = value
}

// Increment increments a metric
func (me *MetricsExporter) Increment(name string, value float64) {
	me.mu.Lock()
	defer me.mu.Unlock()
	me.metrics[name] += value
}

// Get gets a metric
func (me *MetricsExporter) Get(name string) (float64, bool) {
	me.mu.RLock()
	defer me.mu.RUnlock()
	v, ok := me.metrics[name]
	return v, ok
}

// GetAll gets all metrics
func (me *MetricsExporter) GetAll() map[string]float64 {
	me.mu.RLock()
	defer me.mu.RUnlock()
	result := make(map[string]float64)
	for k, v := range me.metrics {
		result[k] = v
	}
	return result
}

// ServeHTTP serves metrics endpoint
func (me *MetricsExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	w.Header().Set("Content-Type", "text/plain")
	for name, value := range me.metrics {
		fmt.Fprintf(w, "%s %f\n", name, value)
	}
}
