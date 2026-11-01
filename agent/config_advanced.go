// Package agent provides additional agent features
// including configuration, monitoring, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// AgentDashboard provides a dashboard for agent monitoring
type AgentDashboard struct {
	agent    *Agent
	metrics  *MetricsCollector
	health   *HealthChecker
	profiler *Profiler
	mu       sync.RWMutex
}

// NewAgentDashboard creates a new agent dashboard
func NewAgentDashboard(agent *Agent) *AgentDashboard {
	return &AgentDashboard{
		agent:    agent,
		metrics:  NewMetricsCollector(),
		health:   NewHealthChecker(),
		profiler: NewProfiler(10*time.Second, 1000),
	}
}

// Start starts the dashboard
func (ad *AgentDashboard) Start(ctx context.Context) {
	ad.profiler.Start(ctx)
}

// Stop stops the dashboard
func (ad *AgentDashboard) Stop() {
	ad.profiler.Stop()
}

// GetStatus returns agent status
func (ad *AgentDashboard) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"running":    ad.agent.IsRunning(),
		"metrics":    ad.metrics.GetAllMetrics(),
		"health":     ad.health.GetStatus(),
		"profile": map[string]interface{}{
			"avg_memory":     ad.profiler.GetAverageMemory(),
			"avg_goroutines": ad.profiler.GetAverageGoroutines(),
			"max_memory":     ad.profiler.GetMaxMemory(),
		},
		"diagnostics": NewDiagnostics().Run(),
	}
}

// GetPrometheusMetrics returns Prometheus metrics
func (ad *AgentDashboard) GetPrometheusMetrics() string {
	return ad.metrics.ToPrometheus()
}

// ServeHTTP serves the dashboard
func (ad *AgentDashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad.GetStatus())
}

// AgentConfig provides agent configuration
type AgentConfig struct {
	AgentID          string            `json:"agentId"`
	TenantID         string            `json:"tenantId"`
	MaxEventsPerSec  int               `json:"maxEventsPerSec"`
	BufferSize       int               `json:"bufferSize"`
	BatchSize        int               `json:"batchSize"`
	FlushInterval    time.Duration     `json:"flushInterval"`
	EnableTLS        bool              `json:"enableTLS"`
	CertFile         string            `json:"certFile"`
	KeyFile          string            `json:"keyFile"`
	CAFile           string            `json:"caFile"`
	IngestionEndpoint string           `json:"ingestionEndpoint"`
	Plugins          []string          `json:"plugins"`
	Sources          []SourceConfig    `json:"sources"`
	Exporters        []ExporterConfig  `json:"exporters"`
	Metadata         map[string]string `json:"metadata"`
}

// SourceConfig configures a data source
type SourceConfig struct {
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Interval time.Duration          `json:"interval"`
	Config   map[string]interface{} `json:"config"`
}

// ExporterConfig configures an exporter
type ExporterConfig struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// DefaultAgentConfig returns default configuration
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		AgentID:           "helix-agent-1",
		TenantID:          "default",
		MaxEventsPerSec:   10000,
		BufferSize:        10000,
		BatchSize:         100,
		FlushInterval:     10 * time.Second,
		EnableTLS:         false,
		IngestionEndpoint: "http://localhost:8080",
		Plugins:           []string{},
		Sources:           []SourceConfig{},
		Exporters:         []ExporterConfig{},
		Metadata:          make(map[string]string),
	}
}

// Validate validates the configuration
func (c *AgentConfig) Validate() []string {
	var errors []string

	if c.AgentID == "" {
		errors = append(errors, "agentId is required")
	}
	if c.TenantID == "" {
		errors = append(errors, "tenantId is required")
	}
	if c.MaxEventsPerSec <= 0 {
		errors = append(errors, "maxEventsPerSec must be positive")
	}
	if c.BufferSize <= 0 {
		errors = append(errors, "bufferSize must be positive")
	}
	if c.BatchSize <= 0 {
		errors = append(errors, "batchSize must be positive")
	}
	if c.IngestionEndpoint == "" {
		errors = append(errors, "ingestionEndpoint is required")
	}

	return errors
}

// IsValid returns true if config is valid
func (c *AgentConfig) IsValid() bool {
	return len(c.Validate()) == 0
}

// AgentConfigValidator validates agent configuration
type AgentConfigValidator struct {
	rules []ConfigRule
}

// ConfigRule represents a validation rule
type ConfigRule struct {
	Field    string
	Required bool
	Type     string
	Min      interface{}
	Max      interface{}
	Pattern  string
}

// NewAgentConfigValidator creates a new validator
func NewAgentConfigValidator() *AgentConfigValidator {
	return &AgentConfigValidator{
		rules: []ConfigRule{
			{Field: "agentId", Required: true, Type: "string"},
			{Field: "tenantId", Required: true, Type: "string"},
			{Field: "maxEventsPerSec", Required: true, Type: "int", Min: 1, Max: 1000000},
			{Field: "bufferSize", Required: true, Type: "int", Min: 1, Max: 1000000},
			{Field: "batchSize", Required: true, Type: "int", Min: 1, Max: 10000},
			{Field: "ingestionEndpoint", Required: true, Type: "string"},
		},
	}
}

// Validate validates a config
func (v *AgentConfigValidator) Validate(config *AgentConfig) []string {
	var errors []string

	for _, rule := range v.rules {
		switch rule.Field {
		case "agentId":
			if rule.Required && config.AgentID == "" {
				errors = append(errors, "agentId is required")
			}
		case "tenantId":
			if rule.Required && config.TenantID == "" {
				errors = append(errors, "tenantId is required")
			}
		case "maxEventsPerSec":
			if config.MaxEventsPerSec <= 0 {
				errors = append(errors, "maxEventsPerSec must be positive")
			}
		case "bufferSize":
			if config.BufferSize <= 0 {
				errors = append(errors, "bufferSize must be positive")
			}
		case "batchSize":
			if config.BatchSize <= 0 {
				errors = append(errors, "batchSize must be positive")
			}
		case "ingestionEndpoint":
			if rule.Required && config.IngestionEndpoint == "" {
				errors = append(errors, "ingestionEndpoint is required")
			}
		}
	}

	return errors
}

// AgentConfigWatcher watches for config changes
type AgentConfigWatcher struct {
	configPath string
	onChange   func(*AgentConfig)
	interval   time.Duration
	mu         sync.Mutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewAgentConfigWatcher creates a new config watcher
func NewAgentConfigWatcher(configPath string, onChange func(*AgentConfig)) *AgentConfigWatcher {
	return &AgentConfigWatcher{
		configPath: configPath,
		onChange:   onChange,
		interval:   5 * time.Second,
	}
}

// Start starts watching
func (w *AgentConfigWatcher) Start(ctx context.Context) {
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.running = true
	go w.watch()
}

// Stop stops watching
func (w *AgentConfigWatcher) Stop() {
	w.running = false
	if w.cancel != nil {
		w.cancel()
	}
}

// watch watches for changes
func (w *AgentConfigWatcher) watch() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	var lastMod time.Time
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			info, err := os.Stat(w.configPath)
			if err != nil {
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				data, err := os.ReadFile(w.configPath)
				if err != nil {
					continue
				}
				var config AgentConfig
				if err := json.Unmarshal(data, &config); err != nil {
					continue
				}
				if w.onChange != nil {
					w.onChange(&config)
				}
			}
		}
	}
}

// AgentConfigExporter exports agent config
type AgentConfigExporter struct {
	config *AgentConfig
}

// NewAgentConfigExporter creates a new config exporter
func NewAgentConfigExporter(config *AgentConfig) *AgentConfigExporter {
	return &AgentConfigExporter{config: config}
}

// ExportJSON exports config as JSON
func (ace *AgentConfigExporter) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(ace.config, "", "  ")
}

// ExportYAML exports config as YAML
func (ace *AgentConfigExporter) ExportYAML() ([]byte, error) {
	// In production, would use yaml library
	return json.MarshalIndent(ace.config, "", "  ")
}

// ExportTOML exports config as TOML
func (ace *AgentConfigExporter) ExportTOML() ([]byte, error) {
	// In production, would use toml library
	return json.MarshalIndent(ace.config, "", "  ")
}

// AgentConfigImporter imports agent config
type AgentConfigImporter struct{}

// NewAgentConfigImporter creates a new config importer
func NewAgentConfigImporter() *AgentConfigImporter {
	return &AgentConfigImporter{}
}

// ImportJSON imports config from JSON
func (aci *AgentConfigImporter) ImportJSON(data []byte) (*AgentConfig, error) {
	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ImportYAML imports config from YAML
func (aci *AgentConfigImporter) ImportYAML(data []byte) (*AgentConfig, error) {
	// In production, would use yaml library
	return aci.ImportJSON(data)
}

// ImportTOML imports config from TOML
func (aci *AgentConfigImporter) ImportTOML(data []byte) (*AgentConfig, error) {
	// In production, would use toml library
	return aci.ImportJSON(data)
}

// AgentConfigValidator2 validates agent config
type AgentConfigValidator2 struct {
	rules []ConfigValidationRule
}

// ConfigValidationRule represents a validation rule
type ConfigValidationRule struct {
	Field    string
	Required bool
	Type     string
	Min      interface{}
	Max      interface{}
	Pattern  string
	Custom   func(interface{}) error
}

// NewAgentConfigValidator2 creates a new validator
func NewAgentConfigValidator2() *AgentConfigValidator2 {
	return &AgentConfigValidator2{
		rules: []ConfigValidationRule{
			{Field: "agentId", Required: true, Type: "string"},
			{Field: "tenantId", Required: true, Type: "string"},
			{Field: "maxEventsPerSec", Required: true, Type: "int", Min: 1},
			{Field: "bufferSize", Required: true, Type: "int", Min: 1},
			{Field: "batchSize", Required: true, Type: "int", Min: 1},
			{Field: "ingestionEndpoint", Required: true, Type: "string"},
		},
	}
}

// Validate validates a config
func (v *AgentConfigValidator2) Validate(config *AgentConfig) []string {
	var errors []string

	if config.AgentID == "" {
		errors = append(errors, "agentId is required")
	}
	if config.TenantID == "" {
		errors = append(errors, "tenantId is required")
	}
	if config.MaxEventsPerSec <= 0 {
		errors = append(errors, "maxEventsPerSec must be positive")
	}
	if config.BufferSize <= 0 {
		errors = append(errors, "bufferSize must be positive")
	}
	if config.BatchSize <= 0 {
		errors = append(errors, "batchSize must be positive")
	}
	if config.IngestionEndpoint == "" {
		errors = append(errors, "ingestionEndpoint is required")
	}

	return errors
}

// IsValid returns true if config is valid
func (v *AgentConfigValidator2) IsValid(config *AgentConfig) bool {
	return len(v.Validate(config)) == 0
}