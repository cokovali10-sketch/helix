// Package agent implements the Helix event collection agent
// with backpressure, plugins, and multiple data sources
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Agent represents the complete Helix observability agent
type Agent struct {
	config       *AgentConfig
	collector    *Collector
	processors   []Processor
	sources      []DataSource
	exporters    []Exporter
	metrics      *AgentMetrics
	mu           sync.RWMutex
	running      bool
	ctx          context.Context
	cancel       context.CancelFunc
}

// AgentConfig holds agent configuration
type AgentConfig struct {
	AgentID          string
	TenantID         string
	MaxEventsPerSec  int
	BufferSize       int
	BatchSize        int
	FlushInterval    time.Duration
	EnableTLS        bool
	CertFile         string
	KeyFile          string
	CAFile           string
	IngestionEndpoint string
	Plugins          []string
	Sources          []SourceConfig
	Exporters        []ExporterConfig
}

// SourceConfig configures a data source
type SourceConfig struct {
	Type     string
	Name     string
	Interval time.Duration
	Config   map[string]interface{}
}

// ExporterConfig configures an exporter
type ExporterConfig struct {
	Type   string
	Name   string
	Config map[string]interface{}
}

// Processor processes events
type Processor interface {
	Process(ctx context.Context, event *Event) (*Event, error)
	Name() string
	Priority() int
}

// DataSource collects events from a source
type DataSource interface {
	Start(ctx context.Context) error
	Stop() error
	Events() <-chan Event
	Name() string
}

// Exporter exports events to a destination
type Exporter interface {
	Export(ctx context.Context, events []Event) error
	Name() string
}

// AgentMetrics tracks agent metrics
type AgentMetrics struct {
	EventsCollected   int64
	EventsProcessed   int64
	EventsDropped     int64
	EventsExported    int64
	ExportErrors      int64
	SourceErrors      int64
	ProcessingTimeMs  int64
	ExportTimeMs      int64
	LastFlushTime     time.Time
	mu                sync.RWMutex
}

// NewAgent creates a new agent
func NewAgent(config *AgentConfig) *Agent {
	agent := &Agent{
		config:     config,
		collector:  NewCollector(config.MaxEventsPerSec),
		processors: make([]Processor, 0),
		sources:    make([]DataSource, 0),
		exporters:  make([]Exporter, 0),
		metrics:    &AgentMetrics{},
	}
	
	// Add default processors
	agent.AddProcessor(&ValidationProcessor{})
	agent.AddProcessor(&EnrichmentProcessor{})
	agent.AddProcessor(&FilterProcessor{})
	
	return agent
}

// AddProcessor adds a processor to the pipeline
func (a *Agent) AddProcessor(p Processor) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.processors = append(a.processors, p)
}

// AddSource adds a data source
func (a *Agent) AddSource(source DataSource) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sources = append(a.sources, source)
}

// AddExporter adds an exporter
func (a *Agent) AddExporter(exporter Exporter) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.exporters = append(a.exporters, exporter)
}

// Start starts the agent
func (a *Agent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if a.running {
		return fmt.Errorf("agent already running")
	}
	
	a.ctx, a.cancel = context.WithCancel(ctx)
	a.running = true
	
	// Start collector
	a.collector.Start(a.ctx)
	
	// Start sources
	for _, source := range a.sources {
		if err := source.Start(a.ctx); err != nil {
			return fmt.Errorf("failed to start source %s: %w", source.Name(), err)
		}
	}
	
	// Start processing loop
	go a.processLoop()
	
	// Start export loop
	go a.exportLoop()
	
	return nil
}

// Stop stops the agent
func (a *Agent) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if !a.running {
		return fmt.Errorf("agent not running")
	}
	
	a.running = false
	a.cancel()
	
	// Stop sources
	for _, source := range a.sources {
		if err := source.Stop(); err != nil {
			fmt.Printf("Error stopping source %s: %v\n", source.Name(), err)
		}
	}
	
	// Stop collector
	a.collector.Stop()
	
	return nil
}

// processLoop processes events from sources
func (a *Agent) processLoop() {
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			// Process events from sources
			for _, source := range a.sources {
				select {
				case event := <-source.Events():
					a.processEvent(event)
				default:
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// processEvent processes a single event
func (a *Agent) processEvent(event Event) {
	a.metrics.mu.Lock()
	a.metrics.EventsCollected++
	a.metrics.mu.Unlock()
	
	// Process through processors
	currentEvent := event
	for _, processor := range a.processors {
		result, err := processor.Process(a.ctx, &currentEvent)
		if err != nil {
			fmt.Printf("Processor %s error: %v\n", processor.Name(), err)
			return
		}
		if result != nil {
			currentEvent = *result
		}
	}
	
	a.metrics.mu.Lock()
	a.metrics.EventsProcessed++
	a.metrics.mu.Unlock()
	
	// Send to collector
	select {
	case a.collector.events <- currentEvent:
	default:
		a.metrics.mu.Lock()
		a.metrics.EventsDropped++
		a.metrics.mu.Unlock()
	}
}

// exportLoop exports events periodically
func (a *Agent) exportLoop() {
	ticker := time.NewTicker(a.config.FlushInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.flush()
		}
	}
}

// flush flushes events to exporters
func (a *Agent) flush() {
	// In production, would collect events from collector and export
	start := time.Now()
	
	a.metrics.mu.Lock()
	a.metrics.LastFlushTime = time.Now()
	a.metrics.mu.Unlock()
	
	_ = start
}

// GetMetrics returns agent metrics
func (a *Agent) GetMetrics() *AgentMetrics {
	a.metrics.mu.RLock()
	defer a.metrics.mu.RUnlock()
	return a.metrics
}

// IsRunning returns if agent is running
func (a *Agent) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// Processor implementations

// ValidationProcessor validates events
type ValidationProcessor struct{}

func (v *ValidationProcessor) Name() string { return "validation" }
func (v *ValidationProcessor) Priority() int { return 10 }
func (v *ValidationProcessor) Process(ctx context.Context, event *Event) (*Event, error) {
	if event.Source == "" {
		event.Source = "unknown"
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Metrics == nil {
		event.Metrics = make(map[string]float64)
	}
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}
	return event, nil
}

// EnrichmentProcessor enriches events with additional metadata
type EnrichmentProcessor struct {
	hostname string
	region   string
}

func (e *EnrichmentProcessor) Name() string { return "enrichment" }
func (e *EnrichmentProcessor) Priority() int { return 20 }
func (e *EnrichmentProcessor) Process(ctx context.Context, event *Event) (*Event, error) {
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}
	if _, ok := event.Tags["host"]; !ok {
		event.Tags["host"] = e.hostname
	}
	if _, ok := event.Tags["region"]; !ok {
		event.Tags["region"] = e.region
	}
	event.Tags["agent_id"] = "helix-agent-1"
	event.Tags["processed_at"] = time.Now().Format(time.RFC3339)
	return event, nil
}

// FilterProcessor filters events based on rules
type FilterProcessor struct {
	rules []FilterRule
}

type FilterRule struct {
	Field    string
	Operator string
	Value    string
	Action   string // drop, keep, modify
}

func (f *FilterProcessor) Name() string { return "filter" }
func (f *FilterProcessor) Priority() int { return 30 }
func (f *FilterProcessor) Process(ctx context.Context, event *Event) (*Event, error) {
	for _, rule := range f.rules {
		if f.matches(event, rule) {
			if rule.Action == "drop" {
				return nil, fmt.Errorf("event dropped by filter rule")
			}
		}
	}
	return event, nil
}

func (f *FilterProcessor) matches(event *Event, rule FilterRule) bool {
	var value string
	switch rule.Field {
	case "source":
		value = event.Source
	case "tenant":
		value = event.Tenant
	default:
		if v, ok := event.Tags[rule.Field]; ok {
			value = v
		}
	}
	
	switch rule.Operator {
	case "equals":
		return value == rule.Value
	case "contains":
		return len(value) > 0 && len(rule.Value) > 0
	case "startswith":
		return len(value) >= len(rule.Value) && value[:len(rule.Value)] == rule.Value
	default:
		return false
	}
}

// Data sources

// FileSource collects events from a file
type FileSource struct {
	path     string
	interval time.Duration
	events   chan Event
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewFileSource(path string, interval time.Duration) *FileSource {
	return &FileSource{
		path:     path,
		interval: interval,
		events:   make(chan Event, 1000),
		quit:     make(chan struct{}),
	}
}

func (fs *FileSource) Start(ctx context.Context) error {
	fs.wg.Add(1)
	go fs.run(ctx)
	return nil
}

func (fs *FileSource) Stop() error {
	close(fs.quit)
	fs.wg.Wait()
	close(fs.events)
	return nil
}

func (fs *FileSource) Events() <-chan Event {
	return fs.events
}

func (fs *FileSource) Name() string {
	return "file"
}

func (fs *FileSource) run(ctx context.Context) {
	defer fs.wg.Done()
	ticker := time.NewTicker(fs.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-fs.quit:
			return
		case <-ticker.C:
			// Read file and emit events
			// Simplified implementation
		}
	}
}

// HTTPSource collects events from an HTTP endpoint
type HTTPSource struct {
	url      string
	interval time.Duration
	events   chan Event
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewHTTPSource(url string, interval time.Duration) *HTTPSource {
	return &HTTPSource{
		url:      url,
		interval: interval,
		events:   make(chan Event, 1000),
		quit:     make(chan struct{}),
	}
}

func (hs *HTTPSource) Start(ctx context.Context) error {
	hs.wg.Add(1)
	go hs.run(ctx)
	return nil
}

func (hs *HTTPSource) Stop() error {
	close(hs.quit)
	hs.wg.Wait()
	close(hs.events)
	return nil
}

func (hs *HTTPSource) Events() <-chan Event {
	return hs.events
}

func (hs *HTTPSource) Name() string {
	return "http"
}

func (hs *HTTPSource) run(ctx context.Context) {
	defer hs.wg.Done()
	ticker := time.NewTicker(hs.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-hs.quit:
			return
		case <-ticker.C:
			// Poll HTTP endpoint
		}
	}
}

// Exporters

// HTTPExporter exports events via HTTP
type HTTPExporter struct {
	endpoint string
	batchSize int
}

func NewHTTPExporter(endpoint string, batchSize int) *HTTPExporter {
	return &HTTPExporter{
		endpoint:  endpoint,
		batchSize: batchSize,
	}
}

func (he *HTTPExporter) Export(ctx context.Context, events []Event) error {
	// In production, would POST events to endpoint
	return nil
}

func (he *HTTPExporter) Name() string {
	return "http"
}

// KafkaExporter exports events to Kafka
type KafkaExporter struct {
	brokers []string
	topic   string
}

func NewKafkaExporter(brokers []string, topic string) *KafkaExporter {
	return &KafkaExporter{
		brokers: brokers,
		topic:   topic,
	}
}

func (ke *KafkaExporter) Export(ctx context.Context, events []Event) error {
	// In production, would produce to Kafka topic
	return nil
}

func (ke *KafkaExporter) Name() string {
	return "kafka"
}

// NATSExporter exports events to NATS
type NATSExporter struct {
	url   string
	topic string
}

func NewNATSExporter(url, topic string) *NATSExporter {
	return &NATSExporter{
		url:   url,
		topic: topic,
	}
}

func (ne *NATSExporter) Export(ctx context.Context, events []Event) error {
	// In production, would publish to NATS subject
	return nil
}

func (ne *NATSExporter) Name() string {
	return "nats"
}

// Plugin registry
var pluginRegistry = make(map[string]func() Processor)

func RegisterPlugin(name string, factory func() Processor) {
	pluginRegistry[name] = factory
}

func GetPlugin(name string) (Processor, bool) {
	factory, ok := pluginRegistry[name]
	if !ok {
		return nil, false
	}
	return factory(), true
}

func ListPlugins() []string {
	names := make([]string, 0, len(pluginRegistry))
	for name := range pluginRegistry {
		names = append(names, name)
	}
	return names
}

// Default agent config
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
	}
}

// CreateDefaultAgent creates an agent with default configuration
func CreateDefaultAgent() *Agent {
	config := DefaultAgentConfig()
	agent := NewAgent(config)
	
	// Add default sources
	agent.AddSource(NewFileSource("/var/log/helix/events.log", 5*time.Second))
	
	// Add default exporters
	agent.AddExporter(NewHTTPExporter(config.IngestionEndpoint+"/v1/events", config.BatchSize))
	
	return agent
}