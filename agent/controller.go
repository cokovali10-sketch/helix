// Package agent provides additional agent features
// including configuration, monitoring, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AgentController controls agent operations
type AgentController struct {
	agent    *Agent
	config   *AgentConfig
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewAgentController creates a new agent controller
func NewAgentController(agent *Agent, config *AgentConfig) *AgentController {
	return &AgentController{
		agent:  agent,
		config: config,
	}
}

// Start starts the agent
func (ac *AgentController) Start(ctx context.Context) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if ac.running {
		return fmt.Errorf("agent already running")
	}

	ac.ctx, ac.cancel = context.WithCancel(ctx)
	ac.running = true

	return ac.agent.Start(ac.ctx)
}

// Stop stops the agent
func (ac *AgentController) Stop() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if !ac.running {
		return fmt.Errorf("agent not running")
	}

	ac.running = false
	if ac.cancel != nil {
		ac.cancel()
	}

	return ac.agent.Stop()
}

// Restart restarts the agent
func (ac *AgentController) Restart(ctx context.Context) error {
	if err := ac.Stop(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return ac.Start(ctx)
}

// IsRunning returns if agent is running
func (ac *AgentController) IsRunning() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.running
}

// GetConfig returns the agent config
func (ac *AgentController) GetConfig() *AgentConfig {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.config
}

// UpdateConfig updates the agent config
func (ac *AgentController) UpdateConfig(config *AgentConfig) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if ac.running {
		return fmt.Errorf("cannot update config while running")
	}

	ac.config = config
	return nil
}

// AgentSupervisor supervises agents
type AgentSupervisor struct {
	agents     map[string]*AgentController
	mu         sync.RWMutex
	healthCheck time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	running    bool
}

// NewAgentSupervisor creates a new supervisor
func NewAgentSupervisor(healthCheck time.Duration) *AgentSupervisor {
	return &AgentSupervisor{
		agents:      make(map[string]*AgentController),
		healthCheck: healthCheck,
	}
}

// Register registers an agent
func (as *AgentSupervisor) Register(id string, controller *AgentController) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	as.agents[id] = controller
	return nil
}

// Unregister unregisters an agent
func (as *AgentSupervisor) Unregister(id string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	delete(as.agents, id)
	return nil
}

// Start starts all agents
func (as *AgentSupervisor) Start(ctx context.Context) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.ctx, as.cancel = context.WithCancel(ctx)
	as.running = true

	for id, controller := range as.agents {
		if err := controller.Start(as.ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
	}

	go as.supervise()
	return nil
}

// Stop stops all agents
func (as *AgentSupervisor) Stop() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.running = false
	if as.cancel != nil {
		as.cancel()
	}

	for id, controller := range as.agents {
		if err := controller.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
	}

	return nil
}

// supervise supervises agents
func (as *AgentSupervisor) supervise() {
	ticker := time.NewTicker(as.healthCheck)
	defer ticker.Stop()

	for {
		select {
		case <-as.ctx.Done():
			return
		case <-ticker.C:
			if !as.running {
				return
			}
			as.checkAgents()
		}
	}
}

// checkAgents checks all agents
func (as *AgentSupervisor) checkAgents() {
	as.mu.RLock()
	agents := make(map[string]*AgentController)
	for id, controller := range as.agents {
		agents[id] = controller
	}
	as.mu.RUnlock()

	for id, controller := range agents {
		if !controller.IsRunning() {
			// Restart agent
			ctx := context.Background()
			if err := controller.Start(ctx); err != nil {
				// Log error
				_ = err
			} else {
				_ = id
			}
		}
	}
}

// GetStatus returns supervisor status
func (as *AgentSupervisor) GetStatus() map[string]interface{} {
	as.mu.RLock()
	defer as.mu.RUnlock()

	status := make(map[string]interface{})
	status["running"] = as.running
	status["agent_count"] = len(as.agents)

	agents := make(map[string]bool)
	for id, controller := range as.agents {
		agents[id] = controller.IsRunning()
	}
	status["agents"] = agents

	return status
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
		"agent_id":          ar.agent.config.AgentID,
		"events_collected":  metrics.EventsCollected,
		"events_processed":  metrics.EventsProcessed,
		"events_dropped":    metrics.EventsDropped,
		"events_exported":   metrics.EventsExported,
		"running":           ar.agent.IsRunning(),
		"timestamp":         time.Now(),
	})
	_ = data
	// In production, would POST to reportURL
}

// AgentLogger logs agent events
type AgentLogger struct {
	level    string
	output   string // stdout, file, both
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
		configLevel = 1 // default to info
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
	format   string // prometheus, json, csv
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
	// In production, would write to file or endpoint
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