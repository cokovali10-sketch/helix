// Package agent provides additional agent features
// including log collection, metric scraping, and event processing
package agent

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// LogCollector collects logs from files
type LogCollector struct {
	paths     []string
	patterns  []string
	events    chan Event
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
	stats     LogCollectorStats
}

// LogCollectorStats tracks log collector stats
type LogCollectorStats struct {
	FilesWatched    int
	LinesProcessed  int64
	EventsGenerated int64
	Errors          int64
	mu              sync.RWMutex
}

// NewLogCollector creates a new log collector
func NewLogCollector(paths []string) *LogCollector {
	return &LogCollector{
		paths:    paths,
		patterns: make([]string, 0),
		events:   make(chan Event, 1000),
		stats:    LogCollectorStats{},
	}
}

// Start starts the log collector
func (lc *LogCollector) Start(ctx context.Context) {
	lc.ctx, lc.cancel = context.WithCancel(ctx)
	lc.running = true

	for _, path := range lc.paths {
		go lc.watchFile(path)
	}
}

// Stop stops the log collector
func (lc *LogCollector) Stop() {
	lc.running = false
	if lc.cancel != nil {
		lc.cancel()
	}
	close(lc.events)
}

// Events returns the events channel
func (lc *LogCollector) Events() <-chan Event {
	return lc.events
}

// watchFile watches a file for changes
func (lc *LogCollector) watchFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		lc.stats.mu.Lock()
		lc.stats.Errors++
		lc.stats.mu.Unlock()
		return
	}
	defer file.Close()

	// Seek to end
	file.Seek(0, 2)

	reader := bufio.NewReader(file)
	buf := make([]byte, 4096)

	for {
		select {
		case <-lc.ctx.Done():
			return
		default:
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				lc.stats.mu.Lock()
				lc.stats.Errors++
				lc.stats.mu.Unlock()
				return
			}
			if n > 0 {
				lines := string(buf[:n])
				for _, line := range splitLines(lines) {
					lc.processLine(path, line)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// processLine processes a log line
func (lc *LogCollector) processLine(path, line string) {
	lc.stats.mu.Lock()
	lc.stats.LinesProcessed++
	lc.stats.mu.Unlock()

	// Check patterns
	for _, pattern := range lc.patterns {
		matched, _ := regexp.MatchString(pattern, line)
		if matched {
			event := Event{
				Timestamp: time.Now(),
				Source:    "log",
				Metrics:   make(map[string]float64),
				Tags: map[string]string{
					"file":    path,
					"pattern": pattern,
					"message": line,
				},
			}
			select {
			case lc.events <- event:
				lc.stats.mu.Lock()
				lc.stats.EventsGenerated++
				lc.stats.mu.Unlock()
			default:
			}
			return
		}
	}

	// Default: emit event for all lines
	event := Event{
		Timestamp: time.Now(),
		Source:    "log",
		Metrics:   make(map[string]float64),
		Tags: map[string]string{
			"file":    path,
			"message": line,
		},
	}
	select {
	case lc.events <- event:
		lc.stats.mu.Lock()
		lc.stats.EventsGenerated++
		lc.stats.mu.Unlock()
	default:
	}
}

// AddPattern adds a pattern to match
func (lc *LogCollector) AddPattern(pattern string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.patterns = append(lc.patterns, pattern)
}

// GetStats returns log collector stats
func (lc *LogCollector) GetStats() LogCollectorStats {
	lc.stats.mu.RLock()
	defer lc.stats.mu.RUnlock()
	return lc.stats
}

// MetricScraper scrapes metrics from endpoints
type MetricScraper struct {
	endpoints []string
	interval  time.Duration
	metrics   map[string]float64
	mu        sync.RWMutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewMetricScraper creates a new metric scraper
func NewMetricScraper(endpoints []string, interval time.Duration) *MetricScraper {
	return &MetricScraper{
		endpoints: endpoints,
		interval:  interval,
		metrics:   make(map[string]float64),
	}
}

// Start starts the metric scraper
func (ms *MetricScraper) Start(ctx context.Context) {
	ms.ctx, ms.cancel = context.WithCancel(ctx)
	ms.running = true
	go ms.scrapeLoop()
}

// Stop stops the metric scraper
func (ms *MetricScraper) Stop() {
	ms.running = false
	if ms.cancel != nil {
		ms.cancel()
	}
}

// scrapeLoop runs the scraping loop
func (ms *MetricScraper) scrapeLoop() {
	ticker := time.NewTicker(ms.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			ms.scrape()
		}
	}
}

// scrape scrapes all endpoints
func (ms *MetricScraper) scrape() {
	for _, endpoint := range ms.endpoints {
		ms.scrapeEndpoint(endpoint)
	}
}

// scrapeEndpoint scrapes a single endpoint
func (ms *MetricScraper) scrapeEndpoint(endpoint string) {
	// In production, would HTTP GET the endpoint and parse Prometheus format
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.metrics[endpoint+"_up"] = 1.0
}

// GetMetrics returns scraped metrics
func (ms *MetricScraper) GetMetrics() map[string]float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]float64)
	for k, v := range ms.metrics {
		result[k] = v
	}
	return result
}

// EventProcessor processes events
type EventProcessor struct {
	processors []Processor
	mu         sync.RWMutex
	stats      EventProcessorStats
}

// EventProcessorStats tracks processor stats
type EventProcessorStats struct {
	EventsProcessed int64
	EventsDropped   int64
	ProcessingTime  time.Duration
	Errors          int64
	mu              sync.RWMutex
}

// NewEventProcessor creates a new event processor
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		processors: make([]Processor, 0),
	}
}

// AddProcessor adds a processor
func (ep *EventProcessor) AddProcessor(p Processor) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.processors = append(ep.processors, p)
	// Sort by priority
	for i := len(ep.processors) - 1; i > 0; i-- {
		if ep.processors[i].Priority() < ep.processors[i-1].Priority() {
			ep.processors[i], ep.processors[i-1] = ep.processors[i-1], ep.processors[i]
		} else {
			break
		}
	}
}

// RemoveProcessor removes a processor
func (ep *EventProcessor) RemoveProcessor(name string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	for i, p := range ep.processors {
		if p.Name() == name {
			ep.processors = append(ep.processors[:i], ep.processors[i+1:]...)
			break
		}
	}
}

// Process processes an event through all processors
func (ep *EventProcessor) Process(ctx context.Context, event *Event) (*Event, error) {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	start := time.Now()
	current := event

	for _, processor := range ep.processors {
		result, err := processor.Process(ctx, current)
		if err != nil {
			ep.stats.mu.Lock()
			ep.stats.Errors++
			ep.stats.mu.Unlock()
			return nil, err
		}
		if result != nil {
			current = result
		} else {
			ep.stats.mu.Lock()
			ep.stats.EventsDropped++
			ep.stats.mu.Unlock()
			return nil, fmt.Errorf("event dropped by %s", processor.Name())
		}
	}

	ep.stats.mu.Lock()
	ep.stats.EventsProcessed++
	ep.stats.ProcessingTime += time.Since(start)
	ep.stats.mu.Unlock()

	return current, nil
}

// GetStats returns processor stats
func (ep *EventProcessor) GetStats() EventProcessorStats {
	ep.stats.mu.RLock()
	defer ep.stats.mu.RUnlock()
	return ep.stats
}

// ListProcessors lists all processors
func (ep *EventProcessor) ListProcessors() []string {
	ep.mu.RLock()
	defer ep.mu.RUnlock()
	names := make([]string, 0, len(ep.processors))
	for _, p := range ep.processors {
		names = append(names, p.Name())
	}
	return names
}

// EventBatcher batches events for efficient processing
type EventBatcher struct {
	batchSize    int
	flushInterval time.Duration
	batch        []Event
	mu           sync.Mutex
	flushChan    chan []Event
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
}

// NewEventBatcher creates a new event batcher
func NewEventBatcher(batchSize int, flushInterval time.Duration) *EventBatcher {
	return &EventBatcher{
		batchSize:     batchSize,
		flushInterval: flushInterval,
		batch:         make([]Event, 0, batchSize),
		flushChan:     make(chan []Event, 100),
	}
}

// Start starts the batcher
func (eb *EventBatcher) Start(ctx context.Context) {
	eb.ctx, eb.cancel = context.WithCancel(ctx)
	eb.running = true
	go eb.run()
}

// Stop stops the batcher
func (eb *EventBatcher) Stop() {
	eb.running = false
	if eb.cancel != nil {
		eb.cancel()
	}
	eb.flush()
	close(eb.flushChan)
}

// Add adds an event to the batch
func (eb *EventBatcher) Add(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.batch = append(eb.batch, event)
	if len(eb.batch) >= eb.batchSize {
		eb.flushLocked()
	}
}

// run runs the batcher loop
func (eb *EventBatcher) run() {
	ticker := time.NewTicker(eb.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-eb.ctx.Done():
			eb.flush()
			return
		case <-ticker.C:
			eb.mu.Lock()
			eb.flushLocked()
			eb.mu.Unlock()
		}
	}
}

// flush flushes the batch
func (eb *EventBatcher) flush() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.flushLocked()
}

// flushLocked flushes the batch (must hold lock)
func (eb *EventBatcher) flushLocked() {
	if len(eb.batch) == 0 {
		return
	}

	batch := make([]Event, len(eb.batch))
	copy(batch, eb.batch)
	eb.batch = eb.batch[:0]

	select {
	case eb.flushChan <- batch:
	default:
		// Drop batch if channel full
	}
}

// FlushChan returns the flush channel
func (eb *EventBatcher) FlushChan() <-chan []Event {
	return eb.flushChan
}

// Size returns the current batch size
func (eb *EventBatcher) Size() int {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	return len(eb.batch)
}

// IsRunning returns if the batcher is running
func (eb *EventBatcher) IsRunning() bool {
	return eb.running
}

// splitLines splits a string into lines
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