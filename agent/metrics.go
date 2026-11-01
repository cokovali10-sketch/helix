// Package agent provides additional agent features
// including metrics collection, tracing, and profiling
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MetricsCollector collects and aggregates metrics
type MetricsCollector struct {
	metrics   map[string]*Metric
	histograms map[string]*Histogram
	summaries  map[string]*Summary
	mu        sync.RWMutex
	startTime time.Time
}

// Metric represents a metric
type Metric struct {
	Name      string
	Type      string
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
	Help      string
}

// Histogram represents a histogram metric
type Histogram struct {
	Name    string
	Buckets []float64
	Counts  []uint64
	Sum     float64
	Count   uint64
	Labels  map[string]string
	mu      sync.Mutex
}

// Summary represents a summary metric
type Summary struct {
	Name       string
	Quantiles  []float64
	Values     []float64
	Sum        float64
	Count      uint64
	Labels     map[string]string
	mu         sync.Mutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics:    make(map[string]*Metric),
		histograms: make(map[string]*Histogram),
		summaries:  make(map[string]*Summary),
		startTime:  time.Now(),
	}
}

// RegisterMetric registers a metric
func (mc *MetricsCollector) RegisterMetric(name, metricType, help string, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metricKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:      name,
		Type:      metricType,
		Labels:    labels,
		Timestamp: time.Now(),
		Help:      help,
	}
}

// IncrementCounter increments a counter
func (mc *MetricsCollector) IncrementCounter(name string, labels map[string]string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metricKey(name, labels)
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

// SetGauge sets a gauge value
func (mc *MetricsCollector) SetGauge(name string, labels map[string]string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metricKey(name, labels)
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

// ObserveHistogram observes a histogram value
func (mc *MetricsCollector) ObserveHistogram(name string, labels map[string]string, value float64, buckets []float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metricKey(name, labels)
	h, ok := mc.histograms[key]
	if !ok {
		h = &Histogram{
			Name:    name,
			Buckets: buckets,
			Counts:  make([]uint64, len(buckets)+1),
			Labels:  labels,
		}
		mc.histograms[key] = h
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.Count++
	h.Sum += value

	for i, bucket := range buckets {
		if value <= bucket {
			h.Counts[i]++
			return
		}
	}
	h.Counts[len(buckets)]++
}

// ObserveSummary observes a summary value
func (mc *MetricsCollector) ObserveSummary(name string, labels map[string]string, value float64, quantiles []float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metricKey(name, labels)
	s, ok := mc.summaries[key]
	if !ok {
		s = &Summary{
			Name:      name,
			Quantiles: quantiles,
			Values:    make([]float64, 0),
			Labels:    labels,
		}
		mc.summaries[key] = s
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Count++
	s.Sum += value
	s.Values = append(s.Values, value)

	// Keep only last 1000 values
	if len(s.Values) > 1000 {
		s.Values = s.Values[len(s.Values)-1000:]
	}
}

// GetMetric gets a metric
func (mc *MetricsCollector) GetMetric(name string, labels map[string]string) (*Metric, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	key := metricKey(name, labels)
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

// GetAllHistograms returns all histograms
func (mc *MetricsCollector) GetAllHistograms() []*Histogram {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	histograms := make([]*Histogram, 0, len(mc.histograms))
	for _, h := range mc.histograms {
		histograms = append(histograms, h)
	}
	return histograms
}

// GetAllSummaries returns all summaries
func (mc *MetricsCollector) GetAllSummaries() []*Summary {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	summaries := make([]*Summary, 0, len(mc.summaries))
	for _, s := range mc.summaries {
		summaries = append(summaries, s)
	}
	return summaries
}

// GetUptime returns uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}

func metricKey(name string, labels map[string]string) string {
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

	// Counters and gauges
	for _, m := range mc.metrics {
		labels := formatLabels(m.Labels)
		output += fmt.Sprintf("# TYPE %s %s\n", m.Name, m.Type)
		output += fmt.Sprintf("%s%s %f %d\n", m.Name, labels, m.Value, m.Timestamp.UnixMilli())
	}

	// Histograms
	for _, h := range mc.histograms {
		h.mu.Lock()
		labels := formatLabels(h.Labels)
		output += fmt.Sprintf("# TYPE %s histogram\n", h.Name)
		cumulative := uint64(0)
		for i, bucket := range h.Buckets {
			cumulative += h.Counts[i]
			output += fmt.Sprintf("%s_bucket%s%.10g} %d\n", h.Name, labels[:len(labels)-1], bucket, cumulative)
		}
		output += fmt.Sprintf("%s_bucket%s+Inf} %d\n", h.Name, labels[:len(labels)-1], h.Count)
		output += fmt.Sprintf("%s_sum%s %f\n", h.Name, labels, h.Sum)
		output += fmt.Sprintf("%s_count%s %d\n", h.Name, labels, h.Count)
		h.mu.Unlock()
	}

	// Summaries
	for _, s := range mc.summaries {
		s.mu.Lock()
		labels := formatLabels(s.Labels)
		output += fmt.Sprintf("# TYPE %s summary\n", s.Name)
		for _, q := range s.Quantiles {
			value := calculateQuantile(s.Values, q)
			output += fmt.Sprintf("%s%s%.2f} %f\n", s.Name, labels[:len(labels)-1], q, value)
		}
		output += fmt.Sprintf("%s_sum%s %f\n", s.Name, labels, s.Sum)
		output += fmt.Sprintf("%s_count%s %d\n", s.Name, labels, s.Count)
		s.mu.Unlock()
	}

	return output
}

func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	parts := make([]string, 0, len(labels))
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=%q", k, v))
	}
	return "{" + joinStrings(parts, ",") + "}"
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

func calculateQuantile(values []float64, q float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	// Simple sort
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	idx := int(q * float64(len(sorted)-1))
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// Tracer provides distributed tracing
type Tracer struct {
	serviceName string
	spans      []*Span
	mu         sync.Mutex
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

// NewTracer creates a new tracer
func NewTracer(serviceName string) *Tracer {
	return &Tracer{
		serviceName: serviceName,
		spans:       make([]*Span, 0),
	}
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(operation string, parentID string) *Span {
	span := &Span{
		TraceID:   fmt.Sprintf("trace-%d", time.Now().UnixNano()),
		SpanID:    fmt.Sprintf("span-%d", time.Now().UnixNano()),
		ParentID:  parentID,
		Operation: operation,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Logs:      make([]SpanLog, 0),
		Status:    "ok",
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = append(t.spans, span)
	return span
}

// FinishSpan finishes a span
func (t *Tracer) FinishSpan(span *Span) {
	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime)
}

// AddTag adds a tag to a span
func (t *Tracer) AddTag(span *Span, key, value string) {
	span.Tags[key] = value
}

// AddLog adds a log to a span
func (t *Tracer) AddLog(span *Span, message string, fields map[string]string) {
	span.Logs = append(span.Logs, SpanLog{
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	})
}

// GetSpans returns all spans
func (t *Tracer) GetSpans() []*Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	spans := make([]*Span, len(t.spans))
	copy(spans, t.spans)
	return spans
}

// GetSpansByTrace returns spans for a trace
func (t *Tracer) GetSpansByTrace(traceID string) []*Span {
	t.mu.Lock()
	defer t.mu.Unlock()

	var result []*Span
	for _, span := range t.spans {
		if span.TraceID == traceID {
			result = append(result, span)
		}
	}
	return result
}

// Profiler profiles agent performance
type Profiler struct {
	samples   []ProfileSample
	mu        sync.Mutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	interval  time.Duration
	maxSamples int
}

// ProfileSample represents a profile sample
type ProfileSample struct {
	Timestamp    time.Time
	CPUUsage     float64
	MemoryUsage  uint64
	Goroutines   int
	HeapSize     uint64
	GCPauseMs    float64
	EventsPerSec float64
}

// NewProfiler creates a new profiler
func NewProfiler(interval time.Duration, maxSamples int) *Profiler {
	return &Profiler{
		samples:    make([]ProfileSample, 0, maxSamples),
		interval:   interval,
		maxSamples: maxSamples,
	}
}

// Start starts profiling
func (p *Profiler) Start(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.running = true
	go p.run()
}

// Stop stops profiling
func (p *Profiler) Stop() {
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
}

// run runs the profiling loop
func (p *Profiler) run() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			if !p.running {
				return
			}
			p.collectSample()
		}
	}
}

// collectSample collects a profile sample
func (p *Profiler) collectSample() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sample := ProfileSample{
		Timestamp:    time.Now(),
		CPUUsage:     0.0,
		MemoryUsage:  m.Alloc,
		Goroutines:   runtime.NumGoroutine(),
		HeapSize:     m.HeapAlloc,
		GCPauseMs:    float64(m.PauseNs[(m.NumGC+255)%256]) / 1e6,
		EventsPerSec: 0,
	}

	p.samples = append(p.samples, sample)
	if len(p.samples) > p.maxSamples {
		p.samples = p.samples[1:]
	}
}

// GetSamples returns collected samples
func (p *Profiler) GetSamples() []ProfileSample {
	p.mu.Lock()
	defer p.mu.Unlock()

	samples := make([]ProfileSample, len(p.samples))
	copy(samples, p.samples)
	return samples
}

// GetAverageMemory returns average memory usage
func (p *Profiler) GetAverageMemory() uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.samples) == 0 {
		return 0
	}

	var sum uint64
	for _, s := range p.samples {
		sum += s.MemoryUsage
	}
	return sum / uint64(len(p.samples))
}

// GetAverageGoroutines returns average goroutine count
func (p *Profiler) GetAverageGoroutines() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.samples) == 0 {
		return 0
	}

	var sum int
	for _, s := range p.samples {
		sum += s.Goroutines
	}
	return float64(sum) / float64(len(p.samples))
}

// GetMaxMemory returns maximum memory usage
func (p *Profiler) GetMaxMemory() uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	var max uint64
	for _, s := range p.samples {
		if s.MemoryUsage > max {
			max = s.MemoryUsage
		}
	}
	return max
}

// Export exports profile data
func (p *Profiler) Export(format string) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch format {
	case "json":
		return json.Marshal(p.samples)
	case "csv":
		var result string
		result += "timestamp,memory,goroutines,heap,gc_pause\n"
		for _, s := range p.samples {
			result += fmt.Sprintf("%d,%d,%d,%d,%.2f\n",
				s.Timestamp.Unix(), s.MemoryUsage, s.Goroutines, s.HeapSize, s.GCPauseMs)
		}
		return []byte(result), nil
	default:
		return json.Marshal(p.samples)
	}
}

// Diagnostics runs diagnostics
type Diagnostics struct {
	checks []Diagnostic
}

// Diagnostic represents a diagnostic check
type Diagnostic struct {
	Name        string
	Description string
	Category    string
	Check       func() (string, error)
}

// NewDiagnostics creates a new diagnostics runner
func NewDiagnostics() *Diagnostics {
	d := &Diagnostics{
		checks: make([]Diagnostic, 0),
	}

	// Add default diagnostics
	d.AddDiagnostic(Diagnostic{
		Name:        "go_version",
		Description: "Go runtime version",
		Category:    "runtime",
		Check: func() (string, error) {
			return runtime.Version(), nil
		},
	})

	d.AddDiagnostic(Diagnostic{
		Name:        "os",
		Description: "Operating system",
		Category:    "system",
		Check: func() (string, error) {
			return runtime.GOOS, nil
		},
	})

	d.AddDiagnostic(Diagnostic{
		Name:        "arch",
		Description: "Architecture",
		Category:    "system",
		Check: func() (string, error) {
			return runtime.GOARCH, nil
		},
	})

	d.AddDiagnostic(Diagnostic{
		Name:        "goroutines",
		Description: "Number of goroutines",
		Category:    "runtime",
		Check: func() (string, error) {
			return fmt.Sprintf("%d", runtime.NumGoroutine()), nil
		},
	})

	d.AddDiagnostic(Diagnostic{
		Name:        "memory",
		Description: "Memory usage",
		Category:    "runtime",
		Check: func() (string, error) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return fmt.Sprintf("Alloc: %d MB, Sys: %d MB", m.Alloc/1024/1024, m.Sys/1024/1024), nil
		},
	})

	d.AddDiagnostic(Diagnostic{
		Name:        "cpu",
		Description: "Number of CPUs",
		Category:    "system",
		Check: func() (string, error) {
			return fmt.Sprintf("%d", runtime.NumCPU()), nil
		},
	})

	return d
}

// AddDiagnostic adds a diagnostic check
func (d *Diagnostics) AddDiagnostic(diag Diagnostic) {
	d.checks = append(d.checks, diag)
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

// RunCategory runs diagnostics for a category
func (d *Diagnostics) RunCategory(category string) map[string]string {
	results := make(map[string]string)
	for _, check := range d.checks {
		if check.Category == category {
			result, err := check.Check()
			if err != nil {
				results[check.Name] = fmt.Sprintf("ERROR: %v", err)
			} else {
				results[check.Name] = result
			}
		}
	}
	return results
}

// ListChecks lists all diagnostic checks
func (d *Diagnostics) ListChecks() []Diagnostic {
	checks := make([]Diagnostic, len(d.checks))
	copy(checks, d.checks)
	return checks
}