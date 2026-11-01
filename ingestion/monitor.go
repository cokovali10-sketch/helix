// Package ingestion provides additional ingestion features
// including advanced monitoring, health checks, and metrics
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// IngestionMonitor monitors ingestion health and performance
type IngestionMonitor struct {
	gateway     *IngestionGateway
	metrics     *IngestionMetrics
	healthChecks []HealthCheck
	alerts      []MonitorAlert
	mu          sync.RWMutex
	running     bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// IngestionMetrics tracks ingestion metrics
type IngestionMetrics struct {
	EventsReceived   int64
	EventsProcessed  int64
	EventsDropped    int64
	EventsForwarded  int64
	BytesReceived    int64
	BytesForwarded   int64
	ProcessingTimeMs int64
	ForwardTimeMs     int64
	ErrorCount       int64
	ActiveTenants    int
	QueueDepth       int
	LastEventTime    time.Time
	mu               sync.RWMutex
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

// MonitorAlert represents a monitoring alert
type MonitorAlert struct {
	ID        string
	Type      string // warning, critical
	Message   string
	Timestamp time.Time
	Resolved  bool
}

// NewIngestionMonitor creates a new ingestion monitor
func NewIngestionMonitor(gateway *IngestionGateway) *IngestionMonitor {
	return &IngestionMonitor{
		gateway:     gateway,
		metrics:     &IngestionMetrics{},
		healthChecks: make([]HealthCheck, 0),
		alerts:      make([]MonitorAlert, 0),
	}
}

// Start starts the monitor
func (im *IngestionMonitor) Start(ctx context.Context) {
	im.ctx, im.cancel = context.WithCancel(ctx)
	im.running = true

	// Register default health checks
	im.RegisterHealthCheck("gateway", func() error {
		if !im.gateway.IsRunning() {
			return fmt.Errorf("gateway not running")
		}
		return nil
	}, 10*time.Second, 5*time.Second)

	im.RegisterHealthCheck("nats", func() error {
		// Check NATS connection
		return nil
	}, 30*time.Second, 10*time.Second)

	im.RegisterHealthCheck("kafka", func() error {
		// Check Kafka connection
		return nil
	}, 30*time.Second, 10*time.Second)

	go im.run()
}

// Stop stops the monitor
func (im *IngestionMonitor) Stop() {
	im.running = false
	if im.cancel != nil {
		im.cancel()
	}
}

// run runs the monitoring loop
func (im *IngestionMonitor) run() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-im.ctx.Done():
			return
		case <-ticker.C:
			im.collectMetrics()
			im.runHealthChecks()
			im.checkAlerts()
		}
	}
}

// collectMetrics collects metrics from gateway
func (im *IngestionMonitor) collectMetrics() {
	gwMetrics := im.gateway.GetMetrics()
	im.metrics.mu.Lock()
	defer im.metrics.mu.Unlock()

	im.metrics.EventsReceived = gwMetrics.EventsReceived
	im.metrics.EventsProcessed = gwMetrics.EventsProcessed
	im.metrics.EventsDropped = gwMetrics.EventsDropped
	im.metrics.BytesReceived = gwMetrics.BytesReceived
	im.metrics.ActiveTenants = len(im.gateway.ListTenants())
	im.metrics.LastEventTime = time.Now()
}

// runHealthChecks runs all health checks
func (im *IngestionMonitor) runHealthChecks() {
	im.mu.RLock()
	checks := make([]HealthCheck, len(im.healthChecks))
	copy(checks, im.healthChecks)
	im.mu.RUnlock()

	for _, check := range checks {
		start := time.Now()
		err := check.Check()
		duration := time.Since(start)

		check.LastRun = time.Now()
		check.LastErr = err

		if err != nil {
			im.addAlert("critical", fmt.Sprintf("Health check %s failed: %v", check.Name, err))
		}

		_ = duration
	}
}

// checkAlerts checks for alert conditions
func (im *IngestionMonitor) checkAlerts() {
	im.metrics.mu.RLock()
	dropRate := float64(im.metrics.EventsDropped) / float64(im.metrics.EventsReceived+1)
	queueDepth := im.metrics.QueueDepth
	im.metrics.mu.RUnlock()

	if dropRate > 0.1 {
		im.addAlert("warning", fmt.Sprintf("High drop rate: %.2f%%", dropRate*100))
	}

	if queueDepth > 10000 {
		im.addAlert("critical", fmt.Sprintf("Queue depth too high: %d", queueDepth))
	}
}

// addAlert adds a monitoring alert
func (im *IngestionMonitor) addAlert(alertType, message string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	alert := MonitorAlert{
		ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
		Type:      alertType,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  false,
	}
	im.alerts = append(im.alerts, alert)
}

// RegisterHealthCheck registers a health check
func (im *IngestionMonitor) RegisterHealthCheck(name string, check func() error, interval, timeout time.Duration) {
	im.mu.Lock()
	defer im.mu.Unlock()

	hc := HealthCheck{
		Name:     name,
		Check:    check,
		Interval: interval,
		Timeout:  timeout,
	}
	im.healthChecks = append(im.healthChecks, hc)
}

// GetMetrics returns ingestion metrics
func (im *IngestionMonitor) GetMetrics() *IngestionMetrics {
	im.metrics.mu.RLock()
	defer im.metrics.mu.RUnlock()

	metrics := *im.metrics
	return &metrics
}

// GetAlerts returns monitoring alerts
func (im *IngestionMonitor) GetAlerts() []MonitorAlert {
	im.mu.RLock()
	defer im.mu.RUnlock()

	alerts := make([]MonitorAlert, len(im.alerts))
	copy(alerts, im.alerts)
	return alerts
}

// GetUnresolvedAlerts returns unresolved alerts
func (im *IngestionMonitor) GetUnresolvedAlerts() []MonitorAlert {
	im.mu.RLock()
	defer im.mu.RUnlock()

	var result []MonitorAlert
	for _, alert := range im.alerts {
		if !alert.Resolved {
			result = append(result, alert)
		}
	}
	return result
}

// ResolveAlert resolves an alert
func (im *IngestionMonitor) ResolveAlert(id string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	for i, alert := range im.alerts {
		if alert.ID == id {
			im.alerts[i].Resolved = true
			break
		}
	}
}

// ServeHTTP serves metrics endpoint
func (im *IngestionMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics := im.GetMetrics()
	alerts := im.GetUnresolvedAlerts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"alerts":  alerts,
		"status":  "healthy",
	})
}

// IngestionDashboard provides a dashboard view
type IngestionDashboard struct {
	monitor *IngestionMonitor
}

// NewIngestionDashboard creates a new dashboard
func NewIngestionDashboard(monitor *IngestionMonitor) *IngestionDashboard {
	return &IngestionDashboard{monitor: monitor}
}

// GetSummary returns a summary of ingestion status
func (id *IngestionDashboard) GetSummary() map[string]interface{} {
	metrics := id.monitor.GetMetrics()
	alerts := id.monitor.GetUnresolvedAlerts()

	dropRate := float64(0)
	if metrics.EventsReceived > 0 {
		dropRate = float64(metrics.EventsDropped) / float64(metrics.EventsReceived)
	}

	return map[string]interface{}{
		"events_received":   metrics.EventsReceived,
		"events_processed":  metrics.EventsProcessed,
		"events_dropped":    metrics.EventsDropped,
		"drop_rate":         dropRate,
		"active_tenants":    metrics.ActiveTenants,
		"bytes_received":    metrics.BytesReceived,
		"queue_depth":       metrics.QueueDepth,
		"unresolved_alerts": len(alerts),
		"last_event_time":   metrics.LastEventTime,
		"status":            "healthy",
	}
}

// GetTenantStats returns per-tenant statistics
func (id *IngestionDashboard) GetTenantStats() []map[string]interface{} {
	tenants := id.monitor.gateway.ListTenants()
	stats := make([]map[string]interface{}, 0, len(tenants))

	for _, tenant := range tenants {
		stats = append(stats, map[string]interface{}{
			"tenant_id":   tenant.ID,
			"tenant_name": tenant.Name,
			"rate_limit":  tenant.RateLimit,
			"quota":       tenant.Quota,
			"active":      tenant.Active,
		})
	}

	return stats
}

// GetThroughput returns throughput metrics
func (id *IngestionDashboard) GetThroughput() map[string]interface{} {
	metrics := id.monitor.GetMetrics()

	eventsPerSec := float64(0)
	if !metrics.LastEventTime.IsZero() {
		elapsed := time.Since(metrics.LastEventTime).Seconds()
		if elapsed > 0 {
			eventsPerSec = float64(metrics.EventsReceived) / elapsed
		}
	}

	return map[string]interface{}{
		"events_per_second": eventsPerSec,
		"bytes_per_second":  float64(metrics.BytesReceived) / 1.0,
		"queue_depth":       metrics.QueueDepth,
	}
}

// GetErrorRate returns error rate metrics
func (id *IngestionDashboard) GetErrorRate() map[string]interface{} {
	metrics := id.monitor.GetMetrics()

	errorRate := float64(0)
	if metrics.EventsReceived > 0 {
		errorRate = float64(metrics.ErrorCount) / float64(metrics.EventsReceived)
	}

	return map[string]interface{}{
		"error_count": metrics.ErrorCount,
		"error_rate":  errorRate,
		"total":       metrics.EventsReceived,
	}
}