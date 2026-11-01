// Package agent provides metrics collection
package agent

import (
	"sync"
	"time"
)

// MetricsCollector collects metrics
type MetricsCollector struct {
	gauges   map[string]float64
	counters map[string]float64
	mu       sync.RWMutex
}

// NewMetricsCollector creates collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{gauges: make(map[string]float64), counters: make(map[string]float64)}
}

// SetGauge sets gauge
func (mc *MetricsCollector) SetGauge(name string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.gauges[name] = value
}

// IncCounter increments counter
func (mc *MetricsCollector) IncCounter(name string, delta float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.counters[name] += delta
}

// GetGauge gets gauge
func (mc *MetricsCollector) GetGauge(name string) (float64, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	v, ok := mc.gauges[name]
	return v, ok
}

// GetCounter gets counter
func (mc *MetricsCollector) GetCounter(name string) float64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.counters[name]
}

// Snapshot returns snapshot
func (mc *MetricsCollector) Snapshot() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return map[string]interface{}{
		"gauges":   mc.gauges,
		"counters": mc.counters,
		"time":     time.Now(),
	}
}

// Reset resets
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.gauges = make(map[string]float64)
	mc.counters = make(map[string]float64)
}
