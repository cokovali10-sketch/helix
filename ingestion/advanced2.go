// Package ingestion provides additional ingestion features
// including advanced filtering, transformation, and routing
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// AdvancedFilter provides advanced event filtering
type AdvancedFilter struct {
	rules    []AdvancedFilterRule
	mu       sync.RWMutex
	stats    FilterStats
}

// AdvancedFilterRule represents a filter rule
type AdvancedFilterRule struct {
	ID        string
	Name      string
	Condition string
	Action    string // drop, keep, modify, route
	Priority  int
	Enabled   bool
	Modified  map[string]string
}

// FilterStats tracks filter statistics
type FilterStats struct {
	TotalProcessed int64
	TotalDropped   int64
	TotalKept      int64
	TotalModified  int64
	RuleHits       map[string]int64
	mu             sync.RWMutex
}

// NewAdvancedFilter creates a new advanced filter
func NewAdvancedFilter() *AdvancedFilter {
	return &AdvancedFilter{
		rules: make([]AdvancedFilterRule, 0),
		stats: FilterStats{
			RuleHits: make(map[string]int64),
		},
	}
}

// AddRule adds a filter rule
func (af *AdvancedFilter) AddRule(rule AdvancedFilterRule) {
	af.mu.Lock()
	defer af.mu.Unlock()
	af.rules = append(af.rules, rule)
	// Sort by priority
	for i := len(af.rules) - 1; i > 0; i-- {
		if af.rules[i].Priority > af.rules[i-1].Priority {
			af.rules[i], af.rules[i-1] = af.rules[i-1], af.rules[i]
		} else {
			break
		}
	}
}

// RemoveRule removes a filter rule
func (af *AdvancedFilter) RemoveRule(id string) {
	af.mu.Lock()
	defer af.mu.Unlock()
	for i, rule := range af.rules {
		if rule.ID == id {
			af.rules = append(af.rules[:i], af.rules[i+1:]...)
			break
		}
	}
}

// Filter applies filter rules to an event
func (af *AdvancedFilter) Filter(event Event) (Event, bool) {
	af.mu.RLock()
	defer af.mu.RUnlock()

	af.stats.mu.Lock()
	af.stats.TotalProcessed++
	af.stats.mu.Unlock()

	for _, rule := range af.rules {
		if !rule.Enabled {
			continue
		}

		if af.matches(event, rule) {
			af.stats.mu.Lock()
			af.stats.RuleHits[rule.ID]++
			af.stats.mu.Unlock()

			switch rule.Action {
			case "drop":
				af.stats.mu.Lock()
				af.stats.TotalDropped++
				af.stats.mu.Unlock()
				return event, false
			case "keep":
				af.stats.mu.Lock()
				af.stats.TotalKept++
				af.stats.mu.Unlock()
				return event, true
			case "modify":
				event = af.modify(event, rule)
				af.stats.mu.Lock()
				af.stats.TotalModified++
				af.stats.mu.Unlock()
			case "route":
				// Route to different target
			}
		}
	}

	af.stats.mu.Lock()
	af.stats.TotalKept++
	af.stats.mu.Unlock()
	return event, true
}

func (af *AdvancedFilter) matches(event Event, rule AdvancedFilterRule) bool {
	// Simple condition evaluation
	// In production, would use a proper expression evaluator
	condition := rule.Condition

	// Check source
	if strings.Contains(condition, "source=") {
		parts := strings.Split(condition, "source=")
		if len(parts) > 1 {
			expected := strings.Trim(parts[1], " ")
			if event.Source != expected {
				return false
			}
		}
	}

	// Check tenant
	if strings.Contains(condition, "tenant=") {
		parts := strings.Split(condition, "tenant=")
		if len(parts) > 1 {
			expected := strings.Trim(parts[1], " ")
			if event.Tenant != expected {
				return false
			}
		}
	}

	// Check metric threshold
	if strings.Contains(condition, ">") {
		// Parse "metric > value"
		re := regexp.MustCompile(`(\w+)\s*>\s*(\d+\.?\d*)`)
		matches := re.FindStringSubmatch(condition)
		if len(matches) == 3 {
			metric := matches[1]
			threshold := parseFloat(matches[2])
			if val, ok := event.Metrics[metric]; ok {
				if !(val > threshold) {
					return false
				}
			}
		}
	}

	if strings.Contains(condition, "<") {
		re := regexp.MustCompile(`(\w+)\s*<\s*(\d+\.?\d*)`)
		matches := re.FindStringSubmatch(condition)
		if len(matches) == 3 {
			metric := matches[1]
			threshold := parseFloat(matches[2])
			if val, ok := event.Metrics[metric]; ok {
				if !(val < threshold) {
					return false
				}
			}
		}
	}

	// Check tag
	if strings.Contains(condition, "tag:") {
		parts := strings.Split(condition, "tag:")
		if len(parts) > 1 {
			tagParts := strings.Split(strings.Trim(parts[1], " "), "=")
			if len(tagParts) == 2 {
				key := strings.Trim(tagParts[0], " ")
				value := strings.Trim(tagParts[1], " ")
				if event.Tags[key] != value {
					return false
				}
			}
		}
	}

	return true
}

func (af *AdvancedFilter) modify(event Event, rule AdvancedFilterRule) Event {
	result := event
	if result.Tags == nil {
		result.Tags = make(map[string]string)
	}
	for key, value := range rule.Modified {
		result.Tags[key] = value
	}
	return result
}

func (af *AdvancedFilter) GetStats() FilterStats {
	af.stats.mu.RLock()
	defer af.stats.mu.RUnlock()
	return FilterStats{
		TotalProcessed: af.stats.TotalProcessed,
		TotalDropped:   af.stats.TotalDropped,
		TotalKept:      af.stats.TotalKept,
		TotalModified:  af.stats.TotalModified,
		RuleHits:       af.stats.RuleHits,
	}
}

func (af *AdvancedFilter) ListRules() []AdvancedFilterRule {
	af.mu.RLock()
	defer af.mu.RUnlock()
	rules := make([]AdvancedFilterRule, len(af.rules))
	copy(rules, af.rules)
	return rules
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// EventTransformer transforms events
type EventTransformer struct {
	rules []TransformRule
	mu    sync.RWMutex
}

type TransformRule struct {
	ID        string
	Name      string
	Operation string // rename, set, delete, convert, compute
	Source    string
	Target    string
	Value     interface{}
	Enabled   bool
}

func NewEventTransformer() *EventTransformer {
	return &EventTransformer{
		rules: make([]TransformRule, 0),
	}
}

func (et *EventTransformer) AddRule(rule TransformRule) {
	et.mu.Lock()
	defer et.mu.Unlock()
	et.rules = append(et.rules, rule)
}

func (et *EventTransformer) RemoveRule(id string) {
	et.mu.Lock()
	defer et.mu.Unlock()
	for i, rule := range et.rules {
		if rule.ID == id {
			et.rules = append(et.rules[:i], et.rules[i+1:]...)
			break
		}
	}
}

func (et *EventTransformer) Transform(event Event) Event {
	et.mu.RLock()
	defer et.mu.RUnlock()

	result := event
	if result.Tags == nil {
		result.Tags = make(map[string]string)
	}
	if result.Metrics == nil {
		result.Metrics = make(map[string]float64)
	}

	for _, rule := range et.rules {
		if !rule.Enabled {
			continue
		}

		switch rule.Operation {
		case "rename":
			if val, ok := result.Tags[rule.Source]; ok {
				delete(result.Tags, rule.Source)
				result.Tags[rule.Target] = val
			}
			if val, ok := result.Metrics[rule.Source]; ok {
				delete(result.Metrics, rule.Source)
				result.Metrics[rule.Target] = val
			}
		case "set":
			if strVal, ok := rule.Value.(string); ok {
				result.Tags[rule.Target] = strVal
			}
			if floatVal, ok := rule.Value.(float64); ok {
				result.Metrics[rule.Target] = floatVal
			}
		case "delete":
			delete(result.Tags, rule.Target)
			delete(result.Metrics, rule.Target)
		case "convert":
			// Type conversion
		case "compute":
			// Compute new value
		}
	}

	return result
}

func (et *EventTransformer) ListRules() []TransformRule {
	et.mu.RLock()
	defer et.mu.RUnlock()
	rules := make([]TransformRule, len(et.rules))
	copy(rules, et.rules)
	return rules
}

// EventRouter routes events to different destinations
type EventRouter struct {
	routes []EventRoute
	mu     sync.RWMutex
}

type EventRoute struct {
	ID        string
	Name      string
	Condition func(Event) bool
	Target    string
	Priority  int
	Enabled   bool
}

func NewEventRouter() *EventRouter {
	return &EventRouter{
		routes: make([]EventRoute, 0),
	}
}

func (er *EventRouter) AddRoute(route EventRoute) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.routes = append(er.routes, route)
	// Sort by priority
	for i := len(er.routes) - 1; i > 0; i-- {
		if er.routes[i].Priority > er.routes[i-1].Priority {
			er.routes[i], er.routes[i-1] = er.routes[i-1], er.routes[i]
		} else {
			break
		}
	}
}

func (er *EventRouter) RemoveRoute(id string) {
	er.mu.Lock()
	defer er.mu.Unlock()
	for i, route := range er.routes {
		if route.ID == id {
			er.routes = append(er.routes[:i], er.routes[i+1:]...)
			break
		}
	}
}

func (er *EventRouter) Route(event Event) []string {
	er.mu.RLock()
	defer er.mu.RUnlock()

	var targets []string
	for _, route := range er.routes {
		if !route.Enabled {
			continue
		}
		if route.Condition(event) {
			targets = append(targets, route.Target)
		}
	}
	return targets
}

func (er *EventRouter) ListRoutes() []EventRoute {
	er.mu.RLock()
	defer er.mu.RUnlock()
	routes := make([]EventRoute, len(er.routes))
	copy(routes, er.routes)
	return routes
}

// EventEnricher enriches events with additional data
type EventEnricher struct {
	enrichers []Enricher
	mu        sync.RWMutex
}

type Enricher interface {
	Enrich(event Event) Event
	Name() string
}

func NewEventEnricher() *EventEnricher {
	return &EventEnricher{
		enrichers: make([]Enricher, 0),
	}
}

func (ee *EventEnricher) AddEnricher(enricher Enricher) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	ee.enrichers = append(ee.enrichers, enricher)
}

func (ee *EventEnricher) Enrich(event Event) Event {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	result := event
	for _, enricher := range ee.enrichers {
		result = enricher.Enrich(result)
	}
	return result
}

// HostEnricher enriches events with host information
type HostEnricher struct {
	hostname string
	region   string
	zone     string
}

func NewHostEnricher(hostname, region, zone string) *HostEnricher {
	return &HostEnricher{
		hostname: hostname,
		region:   region,
		zone:     zone,
	}
}

func (he *HostEnricher) Name() string { return "host" }
func (he *HostEnricher) Enrich(event Event) Event {
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}
	if _, ok := event.Tags["host"]; !ok {
		event.Tags["host"] = he.hostname
	}
	if _, ok := event.Tags["region"]; !ok {
		event.Tags["region"] = he.region
	}
	if _, ok := event.Tags["zone"]; !ok {
		event.Tags["zone"] = he.zone
	}
	return event
}

// TimeEnricher enriches events with time information
type TimeEnricher struct{}

func (te *TimeEnricher) Name() string { return "time" }
func (te *TimeEnricher) Enrich(event Event) Event {
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}
	event.Tags["processed_at"] = time.Now().Format(time.RFC3339)
	event.Tags["day_of_week"] = time.Now().Weekday().String()
	event.Tags["hour"] = fmt.Sprintf("%d", time.Now().Hour())
	return event
}

// GeoEnricher enriches events with geo information
type GeoEnricher struct {
	country string
	city    string
	lat     float64
	lon     float64
}

func NewGeoEnricher(country, city string, lat, lon float64) *GeoEnricher {
	return &GeoEnricher{
		country: country,
		city:    city,
		lat:     lat,
		lon:     lon,
	}
}

func (ge *GeoEnricher) Name() string { return "geo" }
func (ge *GeoEnricher) Enrich(event Event) Event {
	if event.Tags == nil {
		event.Tags = make(map[string]string)
	}
	event.Tags["country"] = ge.country
	event.Tags["city"] = ge.city
	event.Tags["lat"] = fmt.Sprintf("%f", ge.lat)
	event.Tags["lon"] = fmt.Sprintf("%f", ge.lon)
	return event
}

// EventValidator validates events
type EventValidator struct {
	rules []ValidationRule
	mu    sync.RWMutex
}

type ValidationRule struct {
	Field    string
	Required bool
	Type     string
	Min      float64
	Max      float64
	Pattern  string
}

func NewEventValidator() *EventValidator {
	return &EventValidator{
		rules: []ValidationRule{
			{Field: "source", Required: true, Type: "string"},
			{Field: "tenant", Required: true, Type: "string"},
			{Field: "timestamp", Required: true, Type: "int64"},
		},
	}
}

func (ev *EventValidator) AddRule(rule ValidationRule) {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	ev.rules = append(ev.rules, rule)
}

func (ev *EventValidator) Validate(event Event) []string {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	var errors []string

	for _, rule := range ev.rules {
		switch rule.Field {
		case "source":
			if rule.Required && event.Source == "" {
				errors = append(errors, "source is required")
			}
		case "tenant":
			if rule.Required && event.Tenant == "" {
				errors = append(errors, "tenant is required")
			}
		case "timestamp":
			if rule.Required && event.Timestamp.IsZero() {
				errors = append(errors, "timestamp is required")
			}
		default:
			// Check metrics
			if val, ok := event.Metrics[rule.Field]; ok {
				if rule.Min > 0 && val < rule.Min {
					errors = append(errors, fmt.Sprintf("%s must be >= %f", rule.Field, rule.Min))
				}
				if rule.Max > 0 && val > rule.Max {
					errors = append(errors, fmt.Sprintf("%s must be <= %f", rule.Field, rule.Max))
				}
			} else if rule.Required {
				errors = append(errors, fmt.Sprintf("%s is required", rule.Field))
			}
		}
	}

	return errors
}

func (ev *EventValidator) IsValid(event Event) bool {
	return len(ev.Validate(event)) == 0
}

// EventSerializer serializes events
type EventSerializer struct {
	format string
}

func NewEventSerializer(format string) *EventSerializer {
	return &EventSerializer{format: format}
}

func (es *EventSerializer) Serialize(event Event) ([]byte, error) {
	switch es.format {
	case "json":
		return json.Marshal(event)
	case "protobuf":
		// In production, would use protobuf
		return json.Marshal(event)
	case "msgpack":
		// In production, would use msgpack
		return json.Marshal(event)
	default:
		return json.Marshal(event)
	}
}

func (es *EventSerializer) Deserialize(data []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	return event, err
}

// EventBatcher batches events
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

func NewEventBatcher(batchSize int, flushInterval time.Duration) *EventBatcher {
	return &EventBatcher{
		batchSize:     batchSize,
		flushInterval: flushInterval,
		batch:         make([]Event, 0, batchSize),
		flushChan:     make(chan []Event, 100),
	}
}

func (eb *EventBatcher) Start(ctx context.Context) {
	eb.ctx, eb.cancel = context.WithCancel(ctx)
	eb.running = true
	go eb.run()
}

func (eb *EventBatcher) Stop() {
	eb.running = false
	if eb.cancel != nil {
		eb.cancel()
	}
	eb.flush()
	close(eb.flushChan)
}

func (eb *EventBatcher) Add(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.batch = append(eb.batch, event)
	if len(eb.batch) >= eb.batchSize {
		eb.flushLocked()
	}
}

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

func (eb *EventBatcher) flush() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.flushLocked()
}

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

func (eb *EventBatcher) FlushChan() <-chan []Event {
	return eb.flushChan
}

func (eb *EventBatcher) Size() int {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	return len(eb.batch)
}

func (eb *EventBatcher) IsRunning() bool {
	return eb.running
}