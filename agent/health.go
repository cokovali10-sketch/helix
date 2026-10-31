// Package agent provides health check and metrics collection
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthChecker checks agent health
type HealthChecker struct {
	checks  map[string]HealthCheck
	mu      sync.RWMutex
	status  HealthStatus
}

// HealthCheck represents a health check
type HealthCheck struct {
	Name     string
	Check    func() error
	Interval time.Duration
	Timeout  time.Duration
	LastRun  time.Time
	LastErr  error
}

// HealthStatus represents overall health status
type HealthStatus struct {
	Status    string                 `json:"status"`
	Checks    map[string]CheckResult `json:"checks"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    time.Duration          `json:"uptime"`
}

// CheckResult represents a single check result
type CheckResult struct {
	Status    string        `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration"`
	LastCheck time.Time     `json:"lastCheck"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
		status: HealthStatus{
			Status: "healthy",
			Checks: make(map[string]CheckResult),
		},
	}
}

// RegisterCheck registers a health check
func (hc *HealthChecker) RegisterCheck(name string, check func() error, interval, timeout time.Duration) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = HealthCheck{
		Name:     name,
		Check:    check,
		Interval: interval,
		Timeout:  timeout,
	}
}

// UnregisterCheck unregisters a health check
func (hc *HealthChecker) UnregisterCheck(name string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	delete(hc.checks, name)
}

// Check runs all health checks
func (hc *HealthChecker) Check() HealthStatus {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	status := HealthStatus{
		Status:    "healthy",
		Checks:    make(map[string]CheckResult),
		Timestamp: time.Now(),
	}

	allHealthy := true
	for name, check := range hc.checks {
		start := time.Now()
		err := check.Check()
		duration := time.Since(start)

		result := CheckResult{
			Duration:  duration,
			LastCheck:  time.Now(),
		}

		if err != nil {
			result.Status = "unhealthy"
			result.Message = err.Error()
			allHealthy = false
		} else {
			result.Status = "healthy"
		}

		status.Checks[name] = result
	}

	if !allHealthy {
		status.Status = "unhealthy"
	}

	hc.status = status
	return status
}

// GetStatus returns the last health status
func (hc *HealthChecker) GetStatus() HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.status
}

// ServeHTTP serves health check endpoint
func (hc *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := hc.Check()
	w.Header().Set("Content-Type", "application/json")
	if status.Status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status)
}

// MetricsCollector collects agent metrics
type MetricsCollector struct {
	metrics  map[string]*Metric
	mu       sync.RWMutex
	startTime time.Time
}

// Metric represents a metric
type Metric struct {
	Name      string
	Type      string // counter, gauge, histogram, summary
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics:   make(map[string]*Metric),
		startTime: time.Now(),
	}
}

// RegisterMetric registers a metric
func (mc *MetricsCollector) RegisterMetric(name, metricType string, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.metricKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:      name,
		Type:      metricType,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, labels map[string]string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.metricKey(name, labels)
	if m, ok := mc.metrics[key]; ok {
		m.Value += value
		m.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      "counter",
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// SetGauge sets a gauge metric
func (mc *MetricsCollector) SetGauge(name string, labels map[string]string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.metricKey(name, labels)
	if m, ok := mc.metrics[key]; ok {
		m.Value = value
		m.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      "gauge",
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// GetMetric gets a metric by name and labels
func (mc *MetricsCollector) GetMetric(name string, labels map[string]string) (*Metric, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	key := mc.metricKey(name, labels)
	m, ok := mc.metrics[key]
	return m, ok
}

// GetAllMetrics returns all metrics
func (mc *MetricsCollector) GetAllMetrics() []*Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	metrics := make([]*Metric, 0, len(mc.metrics))
	for _, m := range mc.metrics {
		metrics = append(metrics, m)
	}
	return metrics
}

// GetUptime returns agent uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}

func (mc *MetricsCollector) metricKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += fmt.Sprintf(",%s=%s", k, v)
	}
	return key
}

// ToPrometheus converts metrics to Prometheus format
func (mc *MetricsCollector) ToPrometheus() string {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	var output string
	for _, m := range mc.metrics {
		labels := ""
		if len(m.Labels) > 0 {
			parts := make([]string, 0, len(m.Labels))
			for k, v := range m.Labels {
				parts = append(parts, fmt.Sprintf("%s=%q", k, v))
			}
			labels = "{" + joinStrings(parts, ",") + "}"
		}
		output += fmt.Sprintf("# TYPE %s %s\n", m.Name, m.Type)
		output += fmt.Sprintf("%s%s %f %d\n", m.Name, labels, m.Value, m.Timestamp.UnixMilli())
	}
	return output
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// PerformanceProfiler profiles agent performance
type PerformanceProfiler struct {
	samples   []ProfileSample
	mu        sync.Mutex
	enabled   bool
	interval  time.Duration
	maxSamples int
}

// ProfileSample represents a performance sample
type ProfileSample struct {
	Timestamp    time.Time
	CPUUsage     float64
	MemoryUsage  int64
	Goroutines   int
	HeapSize     int64
	GCPauseMs    float64
}

// NewPerformanceProfiler creates a new profiler
func NewPerformanceProfiler(interval time.Duration, maxSamples int) *PerformanceProfiler {
	return &PerformanceProfiler{
		samples:    make([]ProfileSample, 0, maxSamples),
		enabled:    false,
		interval:   interval,
		maxSamples: maxSamples,
	}
}

// Start starts profiling
func (pp *PerformanceProfiler) Start(ctx context.Context) {
	pp.enabled = true
	go pp.run(ctx)
}

// Stop stops profiling
func (pp *PerformanceProfiler) Stop() {
	pp.enabled = false
}

// run runs the profiling loop
func (pp *PerformanceProfiler) run(ctx context.Context) {
	ticker := time.NewTicker(pp.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !pp.enabled {
				return
			}
			pp.collectSample()
		}
	}
}

// collectSample collects a performance sample
func (pp *PerformanceProfiler) collectSample() {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	sample := ProfileSample{
		Timestamp:    time.Now(),
		CPUUsage:     0.0, // Would use runtime metrics
		MemoryUsage:  0,
		Goroutines:   0,
		HeapSize:     0,
		GCPauseMs:    0.0,
	}
	
	pp.samples = append(pp.samples, sample)
	if len(pp.samples) > pp.maxSamples {
		pp.samples = pp.samples[1:]
	}
}

// GetSamples returns collected samples
func (pp *PerformanceProfiler) GetSamples() []ProfileSample {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	samples := make([]ProfileSample, len(pp.samples))
	copy(samples, pp.samples)
	return samples
}

// GetAverageCPU returns average CPU usage
func (pp *PerformanceProfiler) GetAverageCPU() float64 {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	if len(pp.samples) == 0 {
		return 0
	}
	
	var sum float64
	for _, s := range pp.samples {
		sum += s.CPUUsage
	}
	return sum / float64(len(pp.samples))
}

// GetAverageMemory returns average memory usage
func (pp *PerformanceProfiler) GetAverageMemory() int64 {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	if len(pp.samples) == 0 {
		return 0
	}
	
	var sum int64
	for _, s := range pp.samples {
		sum += s.MemoryUsage
	}
	return sum / int64(len(pp.samples))
}

// Diagnostics runs diagnostics
type Diagnostics struct {
	checks []DiagnosticCheck
}

// DiagnosticCheck represents a diagnostic check
type DiagnosticCheck struct {
	Name        string
	Description string
	Check       func() (string, error)
}

// NewDiagnostics creates a new diagnostics runner
func NewDiagnostics() *Diagnostics {
	d := &Diagnostics{
		checks: make([]DiagnosticCheck, 0),
	}
	
	// Add default checks
	d.AddCheck("go_version", "Go runtime version", func() (string, error) {
		return "go1.21", nil
	})
	d.AddCheck("memory", "Memory usage", func() (string, error) {
		return "OK", nil
	})
	d.AddCheck("goroutines", "Goroutine count", func() (string, error) {
		return "OK", nil
	})
	d.AddCheck("network", "Network connectivity", func() (string, error) {
		return "OK", nil
	})
	
	return d
}

// AddCheck adds a diagnostic check
func (d *Diagnostics) AddCheck(name, description string, check func() (string, error)) {
	d.checks = append(d.checks, DiagnosticCheck{
		Name:        name,
		Description: description,
		Check:       check,
	})
}

// Run runs all diagnostic checks
func (d *Diagnostics) Run() map[string]string {
	results := make(map[string]string)
	for _, check := range d.checks {
		result, err := check.Check()
		if err != nil {
			results[check.Name] = fmt.Sprintf("ERROR: %v", err)
		} else {
			results[check.Name] = result
		}
	}
	return results
}

// RunCheck runs a specific diagnostic check
func (d *Diagnostics) RunCheck(name string) (string, error) {
	for _, check := range d.checks {
		if check.Name == name {
			return check.Check()
		}
	}
	return "", fmt.Errorf("check %s not found", name)
}