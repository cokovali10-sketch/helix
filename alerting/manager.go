// Package alerting implements the alerting system for the Helix observability platform
// Includes DSL, evaluation, de-duplication, and routing
package alerting

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Alert represents a single alert
type Alert struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Expression   string            `json:"expression"`
	Severity     string            `json:"severity"`
	Status       string            `json:"status"` // firing, resolved, acknowledged, suppressed
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       *time.Time        `json:"endsAt,omitempty"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	Fingerprint  string            `json:"fingerprint"`
	GeneratorURL string            `json:"generatorURL"`
}

// AlertRule represents an alerting rule
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Expression  string            `json:"expression"`
	Duration    time.Duration     `json:"duration"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Severity    string            `json:"severity"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// AlertManager manages alerts and rules
type AlertManager struct {
	rules       map[string]*AlertRule
	alerts      map[string]*Alert
	alertGroups map[string]*AlertGroup
	routes      []*Route
	receivers   map[string]*Receiver
	mu          sync.RWMutex
	evaluator   *Evaluator
	deduper     *Deduplicator
	router      *Router
	quit        chan struct{}
	wg          sync.WaitGroup
}

// AlertGroup groups related alerts
type AlertGroup struct {
	Name    string   `json:"name"`
	Alerts  []*Alert `json:"alerts"`
	Labels  map[string]string `json:"labels"`
}

// Route defines alert routing
type Route struct {
	ID         string            `json:"id"`
	Match      map[string]string `json:"match"`
	Receiver   string            `json:"receiver"`
	Continue   bool              `json:"continue"`
	GroupBy    []string          `json:"groupBy"`
	GroupWait  time.Duration     `json:"groupWait"`
	GroupInterval time.Duration  `json:"groupInterval"`
	RepeatInterval time.Duration  `json:"repeatInterval"`
	Routes     []*Route          `json:"routes,omitempty"`
}

// Receiver defines an alert receiver
type Receiver struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"` // email, slack, webhook, pagerduty
	Config map[string]string `json:"config"`
}

// Evaluator evaluates alert rules
type Evaluator struct {
	rules   map[string]*AlertRule
	mu      sync.RWMutex
}

// Deduplicator deduplicates alerts
type Deduplicator struct {
	alerts map[string]*Alert
	mu     sync.RWMutex
}

// Router routes alerts to receivers
type Router struct {
	routes    []*Route
	receivers map[string]*Receiver
	mu        sync.RWMutex
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	am := &AlertManager{
		rules:       make(map[string]*AlertRule),
		alerts:      make(map[string]*Alert),
		alertGroups: make(map[string]*AlertGroup),
		routes:      make([]*Route, 0),
		receivers:   make(map[string]*Receiver),
		quit:        make(chan struct{}),
	}
	am.evaluator = NewEvaluator()
	am.deduper = NewDeduplicator()
	am.router = NewRouter()
	return am
}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{
		rules: make(map[string]*AlertRule),
	}
}

// NewDeduplicator creates a new deduplicator
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		alerts: make(map[string]*Alert),
	}
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes:    make([]*Route, 0),
		receivers: make(map[string]*Receiver),
	}
}

// AddRule adds an alert rule
func (am *AlertManager) AddRule(rule *AlertRule) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.rules[rule.ID] = rule
	am.evaluator.AddRule(rule)
}

// RemoveRule removes an alert rule
func (am *AlertManager) RemoveRule(id string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.rules, id)
	am.evaluator.RemoveRule(id)
}

// GetRule gets an alert rule by ID
func (am *AlertManager) GetRule(id string) (*AlertRule, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	rule, ok := am.rules[id]
	return rule, ok
}

// ListRules lists all alert rules
func (am *AlertManager) ListRules() []*AlertRule {
	am.mu.RLock()
	defer am.mu.RUnlock()
	rules := make([]*AlertRule, 0, len(am.rules))
	for _, rule := range am.rules {
		rules = append(rules, rule)
	}
	return rules
}

// AddReceiver adds a receiver
func (am *AlertManager) AddReceiver(receiver *Receiver) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.receivers[receiver.ID] = receiver
	am.router.AddReceiver(receiver)
}

// AddRoute adds a route
func (am *AlertManager) AddRoute(route *Route) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.routes = append(am.routes, route)
	am.router.AddRoute(route)
}

// GetAlerts returns all alerts
func (am *AlertManager) GetAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()
	alerts := make([]*Alert, 0, len(am.alerts))
	for _, alert := range am.alerts {
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetAlert returns an alert by ID
func (am *AlertManager) GetAlert(id string) (*Alert, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	alert, ok := am.alerts[id]
	return alert, ok
}

// Start starts the alert manager
func (am *AlertManager) Start(ctx context.Context) {
	am.wg.Add(1)
	go am.run(ctx)
}

// Stop stops the alert manager
func (am *AlertManager) Stop() {
	close(am.quit)
	am.wg.Wait()
}

// run runs the alert manager evaluation loop
func (am *AlertManager) run(ctx context.Context) {
	defer am.wg.Done()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-am.quit:
			return
		case <-ticker.C:
			am.evaluateRules(ctx)
		}
	}
}

// evaluateRules evaluates all alert rules
func (am *AlertManager) evaluateRules(ctx context.Context) {
	rules := am.ListRules()
	
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		
		// Evaluate the rule expression
		firing, err := am.evaluator.Evaluate(ctx, rule)
		if err != nil {
			fmt.Printf("Error evaluating rule %s: %v\n", rule.Name, err)
			continue
		}
		
		if firing {
			// Create or update alert
			alert := am.createAlert(rule)
			am.deduper.Deduplicate(alert)
			am.router.Route(alert)
		} else {
			// Resolve alert if it was firing
			am.resolveAlert(rule.ID)
		}
	}
}

// createAlert creates an alert from a rule
func (am *AlertManager) createAlert(rule *AlertRule) *Alert {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	alert := &Alert{
		ID:          fmt.Sprintf("alert-%s-%d", rule.ID, time.Now().Unix()),
		Name:        rule.Name,
		Expression:  rule.Expression,
		Severity:    rule.Severity,
		Status:      "firing",
		Labels:      rule.Labels,
		Annotations: rule.Annotations,
		StartsAt:    time.Now(),
		UpdatedAt:   time.Now(),
		Fingerprint:  fmt.Sprintf("%s-%v", rule.ID, rule.Labels),
		GeneratorURL: fmt.Sprintf("/rules/%s", rule.ID),
	}
	
	am.alerts[alert.ID] = alert
	return alert
}

// resolveAlert resolves an alert
func (am *AlertManager) resolveAlert(ruleID string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	for _, alert := range am.alerts {
		if alert.Name == ruleID && alert.Status == "firing" {
			now := time.Now()
			alert.EndsAt = &now
			alert.Status = "resolved"
			alert.UpdatedAt = now
		}
	}
}

// Evaluator methods
func (e *Evaluator) AddRule(rule *AlertRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[rule.ID] = rule
}

func (e *Evaluator) RemoveRule(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, id)
}

func (e *Evaluator) Evaluate(ctx context.Context, rule *AlertRule) (bool, error) {
	// Simplified evaluation - in production, this would evaluate the PromQL-like expression
	// against the TSDB
	return true, nil
}

// Deduplicator methods
func (d *Deduplicator) Deduplicate(alert *Alert) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check if alert already exists
	for _, existing := range d.alerts {
		if existing.Fingerprint == alert.Fingerprint && existing.Status == "firing" {
			// Update existing alert
			existing.UpdatedAt = time.Now()
			return false
		}
	}
	
	d.alerts[alert.ID] = alert
	return true
}

// Router methods
func (r *Router) AddReceiver(receiver *Receiver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.receivers[receiver.ID] = receiver
}

func (r *Router) AddRoute(route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
}

func (r *Router) Route(alert *Alert) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, route := range r.routes {
		if r.matches(alert, route) {
			// Route to receiver
			if receiver, ok := r.receivers[route.Receiver]; ok {
				r.send(receiver, alert)
			}
			
			if !route.Continue {
				break
			}
		}
	}
}

func (r *Router) matches(alert *Alert, route *Route) bool {
	for key, value := range route.Match {
		if alert.Labels[key] != value {
			return false
		}
	}
	return true
}

func (r *Router) send(receiver *Receiver, alert *Alert) {
	// In production, this would send to the appropriate channel
	fmt.Printf("Sending alert %s to receiver %s (%s)\n", alert.Name, receiver.Name, receiver.Type)
}

// Alert DSL
type AlertDSL struct {
	evaluator *Evaluator
}

func NewAlertDSL(evaluator *Evaluator) *AlertDSL {
	return &AlertDSL{evaluator: evaluator}
}

// ParseExpression parses an alert expression
func (dsl *AlertDSL) ParseExpression(expr string) (interface{}, error) {
	// Simplified parser - in production, would use full PromQL-like parser
	return nil, nil
}

// EvaluateExpression evaluates an alert expression
func (dsl *AlertDSL) EvaluateExpression(expr string, data interface{}) (bool, error) {
	// Simplified evaluation
	return true, nil
}

// Alert template for common alert types
var DefaultAlertTemplates = map[string]*AlertRule{
	"high_cpu": {
		ID:          "high_cpu",
		Name:        "HighCPUUsage",
		Expression:  "cpu_usage > 90",
		Duration:    5 * time.Minute,
		Labels:      map[string]string{"severity": "warning"},
		Annotations: map[string]string{"summary": "High CPU usage detected"},
		Severity:    "warning",
		Enabled:     true,
	},
	"high_memory": {
		ID:          "high_memory",
		Name:        "HighMemoryUsage",
		Expression:  "memory_usage > 90",
		Duration:    5 * time.Minute,
		Labels:      map[string]string{"severity": "warning"},
		Annotations: map[string]string{"summary": "High memory usage detected"},
		Severity:    "warning",
		Enabled:     true,
	},
	"disk_full": {
		ID:          "disk_full",
		Name:        "DiskSpaceLow",
		Expression:  "disk_usage > 95",
		Duration:    10 * time.Minute,
		Labels:      map[string]string{"severity": "critical"},
		Annotations: map[string]string{"summary": "Disk space critically low"},
		Severity:    "critical",
		Enabled:     true,
	},
	"service_down": {
		ID:          "service_down",
		Name:        "ServiceDown",
		Expression:  "up == 0",
		Duration:    1 * time.Minute,
		Labels:      map[string]string{"severity": "critical"},
		Annotations: map[string]string{"summary": "Service is down"},
		Severity:    "critical",
		Enabled:     true,
	},
}

// Alert history for tracking
type AlertHistory struct {
	Alerts []*Alert `json:"alerts"`
	mu     sync.RWMutex
}

func NewAlertHistory() *AlertHistory {
	return &AlertHistory{
		Alerts: make([]*Alert, 0),
	}
}

func (h *AlertHistory) Add(alert *Alert) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Alerts = append(h.Alerts, alert)
}

func (h *AlertHistory) GetRecent(limit int) []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if len(h.Alerts) <= limit {
		return h.Alerts
	}
	return h.Alerts[len(h.Alerts)-limit:]
}

func (h *AlertHistory) GetBySeverity(severity string) []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var result []*Alert
	for _, alert := range h.Alerts {
		if alert.Severity == severity {
			result = append(result, alert)
		}
	}
	return result
}

func (h *AlertHistory) GetFiring() []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var result []*Alert
	for _, alert := range h.Alerts {
		if alert.Status == "firing" {
			result = append(result, alert)
		}
	}
	return result
}

func (h *AlertHistory) GetResolved() []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var result []*Alert
	for _, alert := range h.Alerts {
		if alert.Status == "resolved" {
			result = append(result, alert)
		}
	}
	return result
}