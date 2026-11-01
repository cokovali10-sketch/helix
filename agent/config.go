// Package agent provides additional agent features
// including configuration management, update mechanism, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// ConfigManager manages agent configuration
type ConfigManager struct {
	configPath string
	config     *AgentConfig
	mu         sync.RWMutex
	watchers   []func(*AgentConfig)
}

// NewConfigManager creates a new config manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
		watchers:   make([]func(*AgentConfig), 0),
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() (*AgentConfig, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, err
	}

	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	cm.config = &config
	return &config, nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save(config *AgentConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(cm.configPath)
	os.MkdirAll(dir, 0755)

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return err
	}

	cm.config = config
	cm.notifyWatchers(config)
	return nil
}

// Get returns current configuration
func (cm *ConfigManager) Get() *AgentConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Watch adds a configuration watcher
func (cm *ConfigManager) Watch(watcher func(*AgentConfig)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

func (cm *ConfigManager) notifyWatchers(config *AgentConfig) {
	for _, watcher := range cm.watchers {
		go watcher(config)
	}
}

// Updater manages agent updates
type Updater struct {
	currentVersion string
	updateURL      string
	checkInterval  time.Duration
	autoUpdate     bool
	mu             sync.RWMutex
}

// NewUpdater creates a new updater
func NewUpdater(currentVersion, updateURL string) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		updateURL:      updateURL,
		checkInterval:  24 * time.Hour,
		autoUpdate:     false,
	}
}

// CheckForUpdates checks for available updates
func (u *Updater) CheckForUpdates() (string, bool, error) {
	// In production, would check update server
	return u.currentVersion, false, nil
}

// DownloadUpdate downloads an update
func (u *Updater) DownloadUpdate(version string) (string, error) {
	// In production, would download update
	return "", nil
}

// ApplyUpdate applies an update
func (u *Updater) ApplyUpdate(path string) error {
	// In production, would apply update
	return nil
}

// SetAutoUpdate sets auto-update
func (u *Updater) SetAutoUpdate(enabled bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.autoUpdate = enabled
}

// GetAutoUpdate returns auto-update status
func (u *Updater) GetAutoUpdate() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.autoUpdate
}

// DiagnosticRunner runs diagnostics
type DiagnosticRunner struct {
	checks []Diagnostic
	mu     sync.RWMutex
}

// Diagnostic represents a diagnostic check
type Diagnostic struct {
	Name        string
	Description string
	Category    string
	Check       func() (string, error)
}

// NewDiagnosticRunner creates a new diagnostic runner
func NewDiagnosticRunner() *DiagnosticRunner {
	dr := &DiagnosticRunner{
		checks: make([]Diagnostic, 0),
	}

	// Add default diagnostics
	dr.AddDiagnostic(Diagnostic{
		Name:        "go_version",
		Description: "Go runtime version",
		Category:    "runtime",
		Check: func() (string, error) {
			return runtime.Version(), nil
		},
	})

	dr.AddDiagnostic(Diagnostic{
		Name:        "os",
		Description: "Operating system",
		Category:    "system",
		Check: func() (string, error) {
			return runtime.GOOS, nil
		},
	})

	dr.AddDiagnostic(Diagnostic{
		Name:        "arch",
		Description: "Architecture",
		Category:    "system",
		Check: func() (string, error) {
			return runtime.GOARCH, nil
		},
	})

	dr.AddDiagnostic(Diagnostic{
		Name:        "goroutines",
		Description: "Number of goroutines",
		Category:    "runtime",
		Check: func() (string, error) {
			return fmt.Sprintf("%d", runtime.NumGoroutine()), nil
		},
	})

	dr.AddDiagnostic(Diagnostic{
		Name:        "memory",
		Description: "Memory usage",
		Category:    "runtime",
		Check: func() (string, error) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return fmt.Sprintf("Alloc: %d MB, Sys: %d MB", m.Alloc/1024/1024, m.Sys/1024/1024), nil
		},
	})

	dr.AddDiagnostic(Diagnostic{
		Name:        "cpu",
		Description: "Number of CPUs",
		Category:    "system",
		Check: func() (string, error) {
			return fmt.Sprintf("%d", runtime.NumCPU()), nil
		},
	})

	return dr
}

// AddDiagnostic adds a diagnostic check
func (dr *DiagnosticRunner) AddDiagnostic(d Diagnostic) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	dr.checks = append(dr.checks, d)
}

// Run runs all diagnostics
func (dr *DiagnosticRunner) Run() map[string]string {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	results := make(map[string]string)
	for _, check := range dr.checks {
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
func (dr *DiagnosticRunner) RunCategory(category string) map[string]string {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	results := make(map[string]string)
	for _, check := range dr.checks {
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
func (dr *DiagnosticRunner) ListChecks() []Diagnostic {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	checks := make([]Diagnostic, len(dr.checks))
	copy(checks, dr.checks)
	return checks
}

// AgentProfiler profiles agent performance
type AgentProfiler struct {
	samples   []ProfileSample
	mu        sync.Mutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	interval  time.Duration
	maxSamples int
}

// ProfileSample represents a performance sample
type ProfileSample struct {
	Timestamp    time.Time
	CPUUsage     float64
	MemoryUsage  uint64
	Goroutines   int
	HeapSize     uint64
	GCPauseMs    float64
	EventsPerSec float64
}

// NewAgentProfiler creates a new profiler
func NewAgentProfiler(interval time.Duration, maxSamples int) *AgentProfiler {
	return &AgentProfiler{
		samples:    make([]ProfileSample, 0, maxSamples),
		interval:   interval,
		maxSamples: maxSamples,
	}
}

// Start starts profiling
func (ap *AgentProfiler) Start(ctx context.Context) {
	ap.ctx, ap.cancel = context.WithCancel(ctx)
	ap.running = true
	go ap.run()
}

// Stop stops profiling
func (ap *AgentProfiler) Stop() {
	ap.running = false
	if ap.cancel != nil {
		ap.cancel()
	}
}

// run runs the profiling loop
func (ap *AgentProfiler) run() {
	ticker := time.NewTicker(ap.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ap.ctx.Done():
			return
		case <-ticker.C:
			if !ap.running {
				return
			}
			ap.collectSample()
		}
	}
}

// collectSample collects a performance sample
func (ap *AgentProfiler) collectSample() {
	ap.mu.Lock()
	defer ap.mu.Unlock()

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

	ap.samples = append(ap.samples, sample)
	if len(ap.samples) > ap.maxSamples {
		ap.samples = ap.samples[1:]
	}
}

// GetSamples returns collected samples
func (ap *AgentProfiler) GetSamples() []ProfileSample {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	samples := make([]ProfileSample, len(ap.samples))
	copy(samples, ap.samples)
	return samples
}

// GetAverageMemory returns average memory usage
func (ap *AgentProfiler) GetAverageMemory() uint64 {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if len(ap.samples) == 0 {
		return 0
	}

	var sum uint64
	for _, s := range ap.samples {
		sum += s.MemoryUsage
	}
	return sum / uint64(len(ap.samples))
}

// GetAverageGoroutines returns average goroutine count
func (ap *AgentProfiler) GetAverageGoroutines() float64 {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if len(ap.samples) == 0 {
		return 0
	}

	var sum int
	for _, s := range ap.samples {
		sum += s.Goroutines
	}
	return float64(sum) / float64(len(ap.samples))
}

// GetMaxMemory returns maximum memory usage
func (ap *AgentProfiler) GetMaxMemory() uint64 {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	var max uint64
	for _, s := range ap.samples {
		if s.MemoryUsage > max {
			max = s.MemoryUsage
		}
	}
	return max
}

// Export exports profile data
func (ap *AgentProfiler) Export(format string) ([]byte, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	switch format {
	case "json":
		return json.Marshal(ap.samples)
	case "csv":
		var result string
		result += "timestamp,memory,goroutines,heap,gc_pause\n"
		for _, s := range ap.samples {
			result += fmt.Sprintf("%d,%d,%d,%d,%.2f\n",
				s.Timestamp.Unix(), s.MemoryUsage, s.Goroutines, s.HeapSize, s.GCPauseMs)
		}
		return []byte(result), nil
	default:
		return json.Marshal(ap.samples)
	}
}

// AgentMonitor monitors agent health
type AgentMonitor struct {
	agent    *Agent
	profiler *AgentProfiler
	diag    *DiagnosticRunner
	metrics  *MetricsCollector
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewAgentMonitor creates a new agent monitor
func NewAgentMonitor(agent *Agent) *AgentMonitor {
	return &AgentMonitor{
		agent:    agent,
		profiler: NewAgentProfiler(10*time.Second, 1000),
		diag:     NewDiagnosticRunner(),
		metrics:  NewMetricsCollector(),
	}
}

// Start starts monitoring
func (am *AgentMonitor) Start(ctx context.Context) {
	am.ctx, am.cancel = context.WithCancel(ctx)
	am.running = true

	am.profiler.Start(am.ctx)

	go am.run()
}

// Stop stops monitoring
func (am *AgentMonitor) Stop() {
	am.running = false
	if am.cancel != nil {
		am.cancel()
	}
	am.profiler.Stop()
}

// run runs the monitoring loop
func (am *AgentMonitor) run() {
	ticker := time.NewTicker(15 * time.Second)
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

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// GetStatus returns agent status
func (am *AgentMonitor) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"running":    am.agent.IsRunning(),
		"metrics":    am.metrics.GetAllMetrics(),
		"diagnostics": am.diag.Run(),
		"profile": map[string]interface{}{
			"avg_memory":   am.profiler.GetAverageMemory(),
			"avg_goroutines": am.profiler.GetAverageGoroutines(),
			"max_memory":   am.profiler.GetMaxMemory(),
		},
	}
}

// GetPrometheusMetrics returns Prometheus-formatted metrics
func (am *AgentMonitor) GetPrometheusMetrics() string {
	return am.metrics.ToPrometheus()
}

// AgentManager manages multiple agents
type AgentManager struct {
	agents   map[string]*Agent
	monitors map[string]*AgentMonitor
	mu       sync.RWMutex
}

// NewAgentManager creates a new agent manager
func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents:   make(map[string]*Agent),
		monitors: make(map[string]*AgentMonitor),
	}
}

// Register registers an agent
func (am *AgentManager) Register(id string, agent *Agent) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	am.agents[id] = agent
	am.monitors[id] = NewAgentMonitor(agent)
	return nil
}

// Unregister unregisters an agent
func (am *AgentManager) Unregister(id string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	if monitor, ok := am.monitors[id]; ok {
		monitor.Stop()
	}
	delete(am.agents, id)
	delete(am.monitors, id)
	return nil
}

// Get gets an agent by ID
func (am *AgentManager) Get(id string) (*Agent, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	agent, ok := am.agents[id]
	return agent, ok
}

// GetMonitor gets a monitor by agent ID
func (am *AgentManager) GetMonitor(id string) (*AgentMonitor, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	monitor, ok := am.monitors[id]
	return monitor, ok
}

// List lists all agents
func (am *AgentManager) List() []*Agent {
	am.mu.RLock()
	defer am.mu.RUnlock()

	agents := make([]*Agent, 0, len(am.agents))
	for _, agent := range am.agents {
		agents = append(agents, agent)
	}
	return agents
}

// StartAll starts all agents
func (am *AgentManager) StartAll(ctx context.Context) error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for id, agent := range am.agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
		if monitor, ok := am.monitors[id]; ok {
			monitor.Start(ctx)
		}
	}
	return nil
}

// StopAll stops all agents
func (am *AgentManager) StopAll() error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for id, agent := range am.agents {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
		if monitor, ok := am.monitors[id]; ok {
			monitor.Stop()
		}
	}
	return nil
}

// GetStatus returns status of all agents
func (am *AgentManager) GetStatus() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	status := make(map[string]interface{})
	for id, monitor := range am.monitors {
		status[id] = monitor.GetStatus()
	}
	return status
}