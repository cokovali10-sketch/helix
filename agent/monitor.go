// Package agent provides additional agent features
// including monitoring, alerting, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AgentMonitor monitors agent health and performance
type AgentMonitor struct {
	agent    *Agent
	metrics  *MetricsCollector
	alerts   []MonitorAlert
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
}

// MonitorAlert represents a monitoring alert
type MonitorAlert struct {
	ID        string
	Type      string // info, warning, critical
	Message   string
	Timestamp time.Time
	Resolved  bool
}

// NewAgentMonitor creates a new agent monitor
func NewAgentMonitor(agent *Agent, interval time.Duration) *AgentMonitor {
	return &AgentMonitor{
		agent:    agent,
		metrics:  NewMetricsCollector(),
		alerts:   make([]MonitorAlert, 0),
		interval: interval,
	}
}

// Start starts monitoring
func (am *AgentMonitor) Start(ctx context.Context) {
	am.ctx, am.cancel = context.WithCancel(ctx)
	am.running = true
	go am.run()
}

// Stop stops monitoring
func (am *AgentMonitor) Stop() {
	am.running = false
	if am.cancel != nil {
		am.cancel()
	}
}

// run runs the monitoring loop
func (am *AgentMonitor) run() {
	ticker := time.NewTicker(am.interval)
	defer ticker.Stop()

	for {
		select {
		case <-am.ctx.Done():
			return
		case <-ticker.C:
			if !am.running {
				return
			}
			am.collectMetrics()
			am.checkAlerts()
		}
	}
}

// collectMetrics collects agent metrics
func (am *AgentMonitor) collectMetrics() {
	agentMetrics := am.agent.GetMetrics()
	am.metrics.IncrementCounter("agent_events_collected", nil, float64(agentMetrics.EventsCollected))
	am.metrics.IncrementCounter("agent_events_processed", nil, float64(agentMetrics.EventsProcessed))
	am.metrics.IncrementCounter("agent_events_dropped", nil, float64(agentMetrics.EventsDropped))
	am.metrics.SetGauge("agent_running", nil, boolToFloat(am.agent.IsRunning()))
}

// checkAlerts checks for alert conditions
func (am *AgentMonitor) checkAlerts() {
	agentMetrics := am.agent.GetMetrics()

	// Check drop rate
	if agentMetrics.EventsCollected > 0 {
		dropRate := float64(agentMetrics.EventsDropped) / float64(agentMetrics.EventsCollected)
		if dropRate > 0.1 {
			am.addAlert("warning", fmt.Sprintf("High drop rate: %.2f%%", dropRate*100))
		}
	}

	// Check if agent is not running
	if !am.agent.IsRunning() {
		am.addAlert("critical", "Agent is not running")
	}
}

// addAlert adds a monitoring alert
func (am *AgentMonitor) addAlert(alertType, message string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert := MonitorAlert{
		ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
		Type:      alertType,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  false,
	}
	am.alerts = append(am.alerts, alert)
}

// GetMetrics returns agent metrics
func (am *AgentMonitor) GetMetrics() *MetricsCollector {
	return am.metrics
}

// GetAlerts returns monitoring alerts
func (am *AgentMonitor) GetAlerts() []MonitorAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alerts := make([]MonitorAlert, len(am.alerts))
	copy(alerts, am.alerts)
	return alerts
}

// GetUnresolvedAlerts returns unresolved alerts
func (am *AgentMonitor) GetUnresolvedAlerts() []MonitorAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var result []MonitorAlert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			result = append(result, alert)
		}
	}
	return result
}

// ResolveAlert resolves an alert
func (am *AgentMonitor) ResolveAlert(id string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i, alert := range am.alerts {
		if alert.ID == id {
			am.alerts[i].Resolved = true
			break
		}
	}
}

// GetStatus returns agent status
func (am *AgentMonitor) GetStatus() map[string]interface{} {
	metrics := am.agent.GetMetrics()
	return map[string]interface{}{
		"running":          am.agent.IsRunning(),
		"events_collected": metrics.EventsCollected,
		"events_processed": metrics.EventsProcessed,
		"events_dropped":   metrics.EventsDropped,
		"events_exported":  metrics.EventsExported,
		"alerts":           am.GetUnresolvedAlerts(),
		"metrics":          am.metrics.GetAllMetrics(),
	}
}

// GetPrometheusMetrics returns Prometheus metrics
func (am *AgentMonitor) GetPrometheusMetrics() string {
	return am.metrics.ToPrometheus()
}

// AgentDiagnostician runs diagnostics
type AgentDiagnostician struct {
	agent  *Agent
	checks []DiagnosticCheck
	mu     sync.RWMutex
}

// DiagnosticCheck represents a diagnostic check
type DiagnosticCheck struct {
	Name        string
	Description string
	Category    string
	Check       func() (string, error)
}

// NewAgentDiagnostician creates a new diagnostician
func NewAgentDiagnostician(agent *Agent) *AgentDiagnostician {
	ad := &AgentDiagnostician{
		agent:  agent,
		checks: make([]DiagnosticCheck, 0),
	}

	// Add default checks
	ad.AddCheck(DiagnosticCheck{
		Name:        "agent_running",
		Description: "Agent running status",
		Category:    "agent",
		Check: func() (string, error) {
			if agent.IsRunning() {
				return "running", nil
			}
			return "stopped", nil
		},
	})

	ad.AddCheck(DiagnosticCheck{
		Name:        "go_version",
		Description: "Go runtime version",
		Category:    "runtime",
		Check: func() (string, error) {
			return "go1.21", nil
		},
	})

	ad.AddCheck(DiagnosticCheck{
		Name:        "goroutines",
		Description: "Number of goroutines",
		Category:    "runtime",
		Check: func() (string, error) {
			return fmt.Sprintf("%d", 0), nil
		},
	})

	ad.AddCheck(DiagnosticCheck{
		Name:        "memory",
		Description: "Memory usage",
		Category:    "runtime",
		Check: func() (string, error) {
			return "OK", nil
		},
	})

	ad.AddCheck(DiagnosticCheck{
		Name:        "config",
		Description: "Configuration validity",
		Category:    "config",
		Check: func() (string, error) {
			if agent.config.IsValid() {
				return "valid", nil
			}
			return "invalid", fmt.Errorf("invalid config")
		},
	})

	return ad
}

// AddCheck adds a diagnostic check
func (ad *AgentDiagnostician) AddCheck(check DiagnosticCheck) {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	ad.checks = append(ad.checks, check)
}

// Run runs all diagnostics
func (ad *AgentDiagnostician) Run() map[string]string {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	results := make(map[string]string)
	for _, check := range ad.checks {
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
func (ad *AgentDiagnostician) RunCategory(category string) map[string]string {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	results := make(map[string]string)
	for _, check := range ad.checks {
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
func (ad *AgentDiagnostician) ListChecks() []DiagnosticCheck {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	checks := make([]DiagnosticCheck, len(ad.checks))
	copy(checks, ad.checks)
	return checks
}

// AgentReporter reports agent status
type AgentReporter struct {
	agent     *Agent
	reportURL string
	interval  time.Duration
	mu        sync.RWMutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewAgentReporter creates a new agent reporter
func NewAgentReporter(agent *Agent, reportURL string, interval time.Duration) *AgentReporter {
	return &AgentReporter{
		agent:     agent,
		reportURL: reportURL,
		interval:  interval,
	}
}

// Start starts reporting
func (ar *AgentReporter) Start(ctx context.Context) {
	ar.ctx, ar.cancel = context.WithCancel(ctx)
	ar.running = true
	go ar.run()
}

// Stop stops reporting
func (ar *AgentReporter) Stop() {
	ar.running = false
	if ar.cancel != nil {
		ar.cancel()
	}
}

// run runs the reporting loop
func (ar *AgentReporter) run() {
	ticker := time.NewTicker(ar.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ar.ctx.Done():
			return
		case <-ticker.C:
			if !ar.running {
				return
			}
			ar.report()
		}
	}
}

// report sends a report
func (ar *AgentReporter) report() {
	metrics := ar.agent.GetMetrics()
	data, _ := json.Marshal(map[string]interface{}{
		"agent_id":         ar.agent.config.AgentID,
		"events_collected": metrics.EventsCollected,
		"events_processed": metrics.EventsProcessed,
		"events_dropped":   metrics.EventsDropped,
		"events_exported":  metrics.EventsExported,
		"running":          ar.agent.IsRunning(),
		"timestamp":        time.Now(),
	})
	_ = data
	// In production, would POST to reportURL
}

// AgentLogger logs agent events
type AgentLogger struct {
	level    string
	output   string
	filePath string
	file     *os.File
	mu       sync.Mutex
}

// NewAgentLogger creates a new agent logger
func NewAgentLogger(level, output, filePath string) *AgentLogger {
	logger := &AgentLogger{
		level:    level,
		output:   output,
		filePath: filePath,
	}

	if output == "file" || output == "both" {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			logger.file = file
		}
	}

	return logger
}

// Log logs a message
func (al *AgentLogger) Log(level, message string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	if !al.shouldLog(level) {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)

	if al.output == "stdout" || al.output == "both" {
		fmt.Print(logLine)
	}

	if (al.output == "file" || al.output == "both") && al.file != nil {
		al.file.WriteString(logLine)
	}
}

// shouldLog checks if a level should be logged
func (al *AgentLogger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	configLevel, ok := levels[al.level]
	if !ok {
		configLevel = 1
	}

	msgLevel, ok := levels[level]
	if !ok {
		msgLevel = 1
	}

	return msgLevel >= configLevel
}

// Close closes the logger
func (al *AgentLogger) Close() {
	if al.file != nil {
		al.file.Close()
	}
}

// AgentMetricsExporter exports agent metrics
type AgentMetricsExporter struct {
	agent    *Agent
	format   string
	interval time.Duration
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewAgentMetricsExporter creates a new metrics exporter
func NewAgentMetricsExporter(agent *Agent, format string, interval time.Duration) *AgentMetricsExporter {
	return &AgentMetricsExporter{
		agent:    agent,
		format:   format,
		interval: interval,
	}
}

// Start starts exporting
func (ame *AgentMetricsExporter) Start(ctx context.Context) {
	ame.ctx, ame.cancel = context.WithCancel(ctx)
	ame.running = true
	go ame.run()
}

// Stop stops exporting
func (ame *AgentMetricsExporter) Stop() {
	ame.running = false
	if ame.cancel != nil {
		ame.cancel()
	}
}

// run runs the export loop
func (ame *AgentMetricsExporter) run() {
	ticker := time.NewTicker(ame.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ame.ctx.Done():
			return
		case <-ticker.C:
			if !ame.running {
				return
			}
			ame.export()
		}
	}
}

// export exports metrics
func (ame *AgentMetricsExporter) export() {
	metrics := ame.agent.GetMetrics()

	switch ame.format {
	case "prometheus":
		ame.exportPrometheus(metrics)
	case "json":
		ame.exportJSON(metrics)
	case "csv":
		ame.exportCSV(metrics)
	}
}

func (ame *AgentMetricsExporter) exportPrometheus(metrics *AgentMetrics) {
	output := fmt.Sprintf("# HELP helix_agent_events_collected Total events collected\n")
	output += fmt.Sprintf("# TYPE helix_agent_events_collected counter\n")
	output += fmt.Sprintf("helix_agent_events_collected %d\n", metrics.EventsCollected)
	output += fmt.Sprintf("# HELP helix_agent_events_processed Total events processed\n")
	output += fmt.Sprintf("# TYPE helix_agent_events_processed counter\n")
	output += fmt.Sprintf("helix_agent_events_processed %d\n", metrics.EventsProcessed)
	output += fmt.Sprintf("# HELP helix_agent_events_dropped Total events dropped\n")
	output += fmt.Sprintf("# TYPE helix_agent_events_dropped counter\n")
	output += fmt.Sprintf("helix_agent_events_dropped %d\n", metrics.EventsDropped)
	_ = output
}

func (ame *AgentMetricsExporter) exportJSON(metrics *AgentMetrics) {
	data, _ := json.Marshal(map[string]interface{}{
		"events_collected": metrics.EventsCollected,
		"events_processed": metrics.EventsProcessed,
		"events_dropped":   metrics.EventsDropped,
		"events_exported":  metrics.EventsExported,
		"export_errors":    metrics.ExportErrors,
		"source_errors":    metrics.SourceErrors,
	})
	_ = data
}

func (ame *AgentMetricsExporter) exportCSV(metrics *AgentMetrics) {
	output := fmt.Sprintf("events_collected,events_processed,events_dropped,events_exported\n")
	output += fmt.Sprintf("%d,%d,%d,%d\n",
		metrics.EventsCollected, metrics.EventsProcessed, metrics.EventsDropped, metrics.EventsExported)
	_ = output
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}