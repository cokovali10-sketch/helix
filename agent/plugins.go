// Package agent provides additional agent features
// including advanced plugins, sources, and exporters
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// AdvancedPlugin is an advanced plugin interface
type AdvancedPlugin interface {
	Processor
	Initialize(config map[string]interface{}) error
	Shutdown() error
	Health() error
	Metrics() map[string]interface{}
}

// MetricsCollectorPlugin collects metrics from various sources
type MetricsCollectorPlugin struct {
	sources    []MetricSource
	interval   time.Duration
	metrics    map[string]float64
	mu         sync.RWMutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

// MetricSource represents a metric source
type MetricSource struct {
	Name     string
	Type     string // prometheus, statsd, jmx, custom
	Endpoint string
	Interval time.Duration
	Labels   map[string]string
}

// NewMetricsCollectorPlugin creates a new metrics collector plugin
func NewMetricsCollectorPlugin(interval time.Duration) *MetricsCollectorPlugin {
	return &MetricsCollectorPlugin{
		sources:  make([]MetricSource, 0),
		interval: interval,
		metrics:  make(map[string]float64),
	}
}

func (m *MetricsCollectorPlugin) Name() string { return "metrics_collector" }
func (m *MetricsCollectorPlugin) Priority() int { return 5 }
func (m *MetricsCollectorPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	// Enrich event with collected metrics
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if event.Metrics == nil {
		event.Metrics = make(map[string]float64)
	}
	
	for key, value := range m.metrics {
		event.Metrics[key] = value
	}
	
	return event, nil
}

func (m *MetricsCollectorPlugin) Initialize(config map[string]interface{}) error {
	if sources, ok := config["sources"].([]interface{}); ok {
		for _, s := range sources {
			if src, ok := s.(map[string]interface{}); ok {
				metricSource := MetricSource{
					Name:     getString(src, "name"),
					Type:     getString(src, "type"),
					Endpoint: getString(src, "endpoint"),
					Labels:   getMap(src, "labels"),
				}
				m.sources = append(m.sources, metricSource)
			}
		}
	}
	return nil
}

func (m *MetricsCollectorPlugin) Shutdown() error {
	m.running = false
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

func (m *MetricsCollectorPlugin) Health() error {
	if !m.running {
		return fmt.Errorf("metrics collector not running")
	}
	return nil
}

func (m *MetricsCollectorPlugin) Metrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range m.metrics {
		result[k] = v
	}
	return result
}

func (m *MetricsCollectorPlugin) Start(ctx context.Context) {
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.running = true
	
	go m.collectLoop()
}

func (m *MetricsCollectorPlugin) collectLoop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectMetrics()
		}
	}
}

func (m *MetricsCollectorPlugin) collectMetrics() {
	for _, source := range m.sources {
		switch source.Type {
		case "prometheus":
			m.collectPrometheus(source)
		case "statsd":
			m.collectStatsD(source)
		case "jmx":
			m.collectJMX(source)
		case "custom":
			m.collectCustom(source)
		}
	}
}

func (m *MetricsCollectorPlugin) collectPrometheus(source MetricSource) {
	// In production, would scrape Prometheus endpoint
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[fmt.Sprintf("%s_cpu", source.Name)] = 50.0
	m.metrics[fmt.Sprintf("%s_memory", source.Name)] = 75.0
}

func (m *MetricsCollectorPlugin) collectStatsD(source MetricSource) {
	// In production, would receive StatsD metrics
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[fmt.Sprintf("%s_requests", source.Name)] = 100.0
}

func (m *MetricsCollectorPlugin) collectJMX(source MetricSource) {
	// In production, would query JMX
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[fmt.Sprintf("%s_heap", source.Name)] = 512.0
}

func (m *MetricsCollectorPlugin) collectCustom(source MetricSource) {
	// Custom collection logic
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[fmt.Sprintf("%s_custom", source.Name)] = 1.0
}

// LogCollectorPlugin collects logs
type LogCollectorPlugin struct {
	paths    []string
	patterns []string
	events   chan Event
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewLogCollectorPlugin(paths []string) *LogCollectorPlugin {
	return &LogCollectorPlugin{
		paths:  paths,
		events: make(chan Event, 1000),
	}
}

func (l *LogCollectorPlugin) Name() string { return "log_collector" }
func (l *LogCollectorPlugin) Priority() int { return 15 }
func (l *LogCollectorPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	return event, nil
}

func (l *LogCollectorPlugin) Initialize(config map[string]interface{}) error {
	if paths, ok := config["paths"].([]interface{}); ok {
		for _, p := range paths {
			if path, ok := p.(string); ok {
				l.paths = append(l.paths, path)
			}
		}
	}
	if patterns, ok := config["patterns"].([]interface{}); ok {
		for _, p := range patterns {
			if pattern, ok := p.(string); ok {
				l.patterns = append(l.patterns, pattern)
			}
		}
	}
	return nil
}

func (l *LogCollectorPlugin) Shutdown() error {
	l.running = false
	if l.cancel != nil {
		l.cancel()
	}
	close(l.events)
	return nil
}

func (l *LogCollectorPlugin) Health() error {
	if !l.running {
		return fmt.Errorf("log collector not running")
	}
	return nil
}

func (l *LogCollectorPlugin) Metrics() map[string]interface{} {
	return map[string]interface{}{
		"paths":    len(l.paths),
		"patterns": len(l.patterns),
	}
}

func (l *LogCollectorPlugin) Start(ctx context.Context) {
	l.ctx, l.cancel = context.WithCancel(ctx)
	l.running = true
	go l.collectLoop()
}

func (l *LogCollectorPlugin) collectLoop() {
	for _, path := range l.paths {
		go l.tailFile(path)
	}
	<-l.ctx.Done()
}

func (l *LogCollectorPlugin) tailFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Seek to end
	file.Seek(0, 2)

	reader := make([]byte, 4096)
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			n, err := file.Read(reader)
			if err != nil && err != io.EOF {
				return
			}
			if n > 0 {
				// Process log lines
				lines := string(reader[:n])
				for _, line := range splitLines(lines) {
					event := Event{
						Timestamp: time.Now(),
						Source:    "log",
						Metrics:   map[string]float64{},
						Tags:      map[string]string{"file": path},
					}
					// Parse log line
					event.Tags["message"] = line
					select {
					case l.events <- event:
					default:
					}
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (l *LogCollectorPlugin) Events() <-chan Event {
	return l.events
}

// NetworkCollectorPlugin collects network metrics
type NetworkCollectorPlugin struct {
	interfaces []string
	metrics    map[string]float64
	mu         sync.RWMutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewNetworkCollectorPlugin() *NetworkCollectorPlugin {
	return &NetworkCollectorPlugin{
		interfaces: make([]string, 0),
		metrics:    make(map[string]float64),
	}
}

func (n *NetworkCollectorPlugin) Name() string { return "network_collector" }
func (n *NetworkCollectorPlugin) Priority() int { return 25 }
func (n *NetworkCollectorPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	if event.Metrics == nil {
		event.Metrics = make(map[string]float64)
	}
	
	for key, value := range n.metrics {
		event.Metrics[key] = value
	}
	
	return event, nil
}

func (n *NetworkCollectorPlugin) Initialize(config map[string]interface{}) error {
	if ifaces, ok := config["interfaces"].([]interface{}); ok {
		for _, iface := range ifaces {
			if name, ok := iface.(string); ok {
				n.interfaces = append(n.interfaces, name)
			}
		}
	}
	return nil
}

func (n *NetworkCollectorPlugin) Shutdown() error {
	n.running = false
	if n.cancel != nil {
		n.cancel()
	}
	return nil
}

func (n *NetworkCollectorPlugin) Health() error {
	if !n.running {
		return fmt.Errorf("network collector not running")
	}
	return nil
}

func (n *NetworkCollectorPlugin) Metrics() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range n.metrics {
		result[k] = v
	}
	return result
}

func (n *NetworkCollectorPlugin) Start(ctx context.Context) {
	n.ctx, n.cancel = context.WithCancel(ctx)
	n.running = true
	go n.collectLoop()
}

func (n *NetworkCollectorPlugin) collectLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-n.ctx.Done():
			return
		case <-ticker.C:
			n.collectNetworkMetrics()
		}
	}
}

func (n *NetworkCollectorPlugin) collectNetworkMetrics() {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	for _, iface := range n.interfaces {
		// In production, would read /proc/net/dev
		n.metrics[fmt.Sprintf("net_%s_rx_bytes", iface)] = 1000000.0
		n.metrics[fmt.Sprintf("net_%s_tx_bytes", iface)] = 500000.0
		n.metrics[fmt.Sprintf("net_%s_rx_packets", iface)] = 1000.0
		n.metrics[fmt.Sprintf("net_%s_tx_packets", iface)] = 500.0
		n.metrics[fmt.Sprintf("net_%s_errors", iface)] = 0.0
	}
}

// ProcessCollectorPlugin collects process metrics
type ProcessCollectorPlugin struct {
	processes []string
	metrics   map[string]float64
	mu        sync.RWMutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewProcessCollectorPlugin() *ProcessCollectorPlugin {
	return &ProcessCollectorPlugin{
		processes: make([]string, 0),
		metrics:   make(map[string]float64),
	}
}

func (p *ProcessCollectorPlugin) Name() string { return "process_collector" }
func (p *ProcessCollectorPlugin) Priority() int { return 30 }
func (p *ProcessCollectorPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if event.Metrics == nil {
		event.Metrics = make(map[string]float64)
	}
	
	for key, value := range p.metrics {
		event.Metrics[key] = value
	}
	
	return event, nil
}

func (p *ProcessCollectorPlugin) Initialize(config map[string]interface{}) error {
	if procs, ok := config["processes"].([]interface{}); ok {
		for _, proc := range procs {
			if name, ok := proc.(string); ok {
				p.processes = append(p.processes, name)
			}
		}
	}
	return nil
}

func (p *ProcessCollectorPlugin) Shutdown() error {
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func (p *ProcessCollectorPlugin) Health() error {
	if !p.running {
		return fmt.Errorf("process collector not running")
	}
	return nil
}

func (p *ProcessCollectorPlugin) Metrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range p.metrics {
		result[k] = v
	}
	return result
}

func (p *ProcessCollectorPlugin) Start(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.running = true
	go p.collectLoop()
}

func (p *ProcessCollectorPlugin) collectLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.collectProcessMetrics()
		}
	}
}

func (p *ProcessCollectorPlugin) collectProcessMetrics() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for _, proc := range p.processes {
		// In production, would read /proc
		p.metrics[fmt.Sprintf("proc_%s_cpu", proc)] = 10.0
		p.metrics[fmt.Sprintf("proc_%s_memory", proc)] = 100.0
		p.metrics[fmt.Sprintf("proc_%s_threads", proc)] = 5.0
		p.metrics[fmt.Sprintf("proc_%s_fds", proc)] = 20.0
	}
}

// ContainerCollectorPlugin collects container metrics
type ContainerCollectorPlugin struct {
	containers []string
	metrics    map[string]float64
	mu         sync.RWMutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewContainerCollectorPlugin() *ContainerCollectorPlugin {
	return &ContainerCollectorPlugin{
		containers: make([]string, 0),
		metrics:    make(map[string]float64),
	}
}

func (c *ContainerCollectorPlugin) Name() string { return "container_collector" }
func (c *ContainerCollectorPlugin) Priority() int { return 35 }
func (c *ContainerCollectorPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if event.Metrics == nil {
		event.Metrics = make(map[string]float64)
	}
	
	for key, value := range c.metrics {
		event.Metrics[key] = value
	}
	
	return event, nil
}

func (c *ContainerCollectorPlugin) Initialize(config map[string]interface{}) error {
	if containers, ok := config["containers"].([]interface{}); ok {
		for _, cont := range containers {
			if name, ok := cont.(string); ok {
				c.containers = append(c.containers, name)
			}
		}
	}
	return nil
}

func (c *ContainerCollectorPlugin) Shutdown() error {
	c.running = false
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *ContainerCollectorPlugin) Health() error {
	if !c.running {
		return fmt.Errorf("container collector not running")
	}
	return nil
}

func (c *ContainerCollectorPlugin) Metrics() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range c.metrics {
		result[k] = v
	}
	return result
}

func (c *ContainerCollectorPlugin) Start(ctx context.Context) {
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running = true
	go c.collectLoop()
}

func (c *ContainerCollectorPlugin) collectLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.collectContainerMetrics()
		}
	}
}

func (c *ContainerCollectorPlugin) collectContainerMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for _, container := range c.containers {
		// In production, would read cgroup stats
		c.metrics[fmt.Sprintf("container_%s_cpu", container)] = 25.0
		c.metrics[fmt.Sprintf("container_%s_memory", container)] = 512.0
		c.metrics[fmt.Sprintf("container_%s_network_rx", container)] = 1000000.0
		c.metrics[fmt.Sprintf("container_%s_network_tx", container)] = 500000.0
	}
}

// Helper functions

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getMap(m map[string]interface{}, key string) map[string]string {
	result := make(map[string]string)
	if v, ok := m[key].(map[string]interface{}); ok {
		for k, val := range v {
			if str, ok := val.(string); ok {
				result[k] = str
			}
		}
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}