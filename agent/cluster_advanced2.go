// Package agent provides additional agent features
// including configuration, monitoring, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AgentClusterManager manages agent clusters
type AgentClusterManager struct {
	clusters map[string]*AgentCluster
	mu       sync.RWMutex
}

// NewAgentClusterManager creates a new cluster manager
func NewAgentClusterManager() *AgentClusterManager {
	return &AgentClusterManager{
		clusters: make(map[string]*AgentCluster),
	}
}

// CreateCluster creates a new cluster
func (acm *AgentClusterManager) CreateCluster(id, name string, config *ClusterConfig) error {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if _, exists := acm.clusters[id]; exists {
		return fmt.Errorf("cluster %s already exists", id)
	}

	acm.clusters[id] = &AgentCluster{
		ID:     id,
		Name:   name,
		Config: config,
		Agents: make(map[string]*Agent),
	}
	return nil
}

// GetCluster gets a cluster by ID
func (acm *AgentClusterManager) GetCluster(id string) (*AgentCluster, bool) {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	cluster, ok := acm.clusters[id]
	return cluster, ok
}

// ListClusters lists all clusters
func (acm *AgentClusterManager) ListClusters() []*AgentCluster {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	clusters := make([]*AgentCluster, 0, len(acm.clusters))
	for _, c := range acm.clusters {
		clusters = append(clusters, c)
	}
	return clusters
}

// DeleteCluster deletes a cluster
func (acm *AgentClusterManager) DeleteCluster(id string) error {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if _, exists := acm.clusters[id]; !exists {
		return fmt.Errorf("cluster %s not found", id)
	}

	delete(acm.clusters, id)
	return nil
}

// AgentCluster represents a cluster of agents
type AgentCluster struct {
	ID      string
	Name    string
	Config  *ClusterConfig
	Agents  map[string]*Agent
	mu      sync.RWMutex
	running bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// ClusterConfig holds cluster configuration
type ClusterConfig struct {
	Name            string
	LeaderID        string
	Replicas        int
	Sharding        bool
	Consensus       string
	Heartbeat       time.Duration
	ElectionTimeout time.Duration
}

// RegisterAgent registers an agent in the cluster
func (ac *AgentCluster) RegisterAgent(id string, agent *Agent) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.Agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	ac.Agents[id] = agent
	return nil
}

// UnregisterAgent unregisters an agent
func (ac *AgentCluster) UnregisterAgent(id string) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.Agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	delete(ac.Agents, id)
	return nil
}

// GetAgent gets an agent by ID
func (ac *AgentCluster) GetAgent(id string) (*Agent, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	agent, ok := ac.Agents[id]
	return agent, ok
}

// ListAgents lists all agents
func (ac *AgentCluster) ListAgents() []*Agent {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	agents := make([]*Agent, 0, len(ac.Agents))
	for _, agent := range ac.Agents {
		agents = append(agents, agent)
	}
	return agents
}

// Start starts the cluster
func (ac *AgentCluster) Start(ctx context.Context) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.ctx, ac.cancel = context.WithCancel(ctx)
	ac.running = true

	for id, agent := range ac.Agents {
		if err := agent.Start(ac.ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
	}

	return nil
}

// Stop stops the cluster
func (ac *AgentCluster) Stop() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.running = false
	if ac.cancel != nil {
		ac.cancel()
	}

	for id, agent := range ac.Agents {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
	}

	return nil
}

// GetStatus returns cluster status
func (ac *AgentCluster) GetStatus() map[string]interface{} {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	status := make(map[string]interface{})
	status["running"] = ac.running
	status["agent_count"] = len(ac.Agents)
	status["cluster_name"] = ac.Name
	status["leader"] = ac.Config.LeaderID
	status["consensus"] = ac.Config.Consensus

	agents := make(map[string]bool)
	for id, agent := range ac.Agents {
		agents[id] = agent.IsRunning()
	}
	status["agents"] = agents

	return status
}

// AgentDiscovery discovers agents
type AgentDiscovery struct {
	agents   map[string]*AgentInfo
	mu       sync.RWMutex
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// AgentInfo represents agent information
type AgentInfo struct {
	ID       string
	Endpoint string
	Version  string
	Status   string
	Metadata map[string]string
	LastSeen time.Time
}

// NewAgentDiscovery creates a new agent discovery
func NewAgentDiscovery(interval time.Duration) *AgentDiscovery {
	return &AgentDiscovery{
		agents:   make(map[string]*AgentInfo),
		interval: interval,
	}
}

// Start starts discovery
func (ad *AgentDiscovery) Start(ctx context.Context) {
	ad.ctx, ad.cancel = context.WithCancel(ctx)
	ad.running = true
	go ad.run()
}

// Stop stops discovery
func (ad *AgentDiscovery) Stop() {
	ad.running = false
	if ad.cancel != nil {
		ad.cancel()
	}
}

// run runs the discovery loop
func (ad *AgentDiscovery) run() {
	ticker := time.NewTicker(ad.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ad.ctx.Done():
			return
		case <-ticker.C:
			if !ad.running {
				return
			}
			ad.discover()
		}
	}
}

// discover discovers agents
func (ad *AgentDiscovery) discover() {
	// In production, would use mDNS, DNS-SD, or consensus
}

// Register registers an agent
func (ad *AgentDiscovery) Register(info *AgentInfo) {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	info.LastSeen = time.Now()
	ad.agents[info.ID] = info
}

// Unregister unregisters an agent
func (ad *AgentDiscovery) Unregister(id string) {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	delete(ad.agents, id)
}

// Get gets an agent info
func (ad *AgentDiscovery) Get(id string) (*AgentInfo, bool) {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	info, ok := ad.agents[id]
	return info, ok
}

// List lists all agents
func (ad *AgentDiscovery) List() []*AgentInfo {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	infos := make([]*AgentInfo, 0, len(ad.agents))
	for _, info := range ad.agents {
		infos = append(infos, info)
	}
	return infos
}

// GetActive gets active agents
func (ad *AgentDiscovery) GetActive() []*AgentInfo {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	var active []*AgentInfo
	for _, info := range ad.agents {
		if info.LastSeen.After(cutoff) {
			active = append(active, info)
		}
	}
	return active
}

// Cleanup removes stale agents
func (ad *AgentDiscovery) Cleanup() {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	for id, info := range ad.agents {
		if info.LastSeen.Before(cutoff) {
			delete(ad.agents, id)
		}
	}
}

// AgentRegistry manages agent registration
type AgentRegistry struct {
	agents map[string]*RegisteredAgent
	mu     sync.RWMutex
}

// RegisteredAgent represents a registered agent
type RegisteredAgent struct {
	ID       string
	Agent    *Agent
	Metadata map[string]string
	Status   string
	LastSeen time.Time
}

// NewAgentRegistry creates a new agent registry
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*RegisteredAgent),
	}
}

// Register registers an agent
func (ar *AgentRegistry) Register(id string, agent *Agent, metadata map[string]string) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, exists := ar.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	ar.agents[id] = &RegisteredAgent{
		ID:       id,
		Agent:    agent,
		Metadata: metadata,
		Status:   "registered",
		LastSeen: time.Now(),
	}
	return nil
}

// Unregister unregisters an agent
func (ar *AgentRegistry) Unregister(id string) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, exists := ar.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	delete(ar.agents, id)
	return nil
}

// Get gets a registered agent
func (ar *AgentRegistry) Get(id string) (*RegisteredAgent, bool) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	agent, ok := ar.agents[id]
	return agent, ok
}

// List lists all registered agents
func (ar *AgentRegistry) List() []*RegisteredAgent {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	agents := make([]*RegisteredAgent, 0, len(ar.agents))
	for _, a := range ar.agents {
		agents = append(agents, a)
	}
	return agents
}

// UpdateStatus updates agent status
func (ar *AgentRegistry) UpdateStatus(id, status string) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	agent, exists := ar.agents[id]
	if !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	agent.Status = status
	agent.LastSeen = time.Now()
	return nil
}

// GetActive gets active agents
func (ar *AgentRegistry) GetActive() []*RegisteredAgent {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	var active []*RegisteredAgent
	for _, agent := range ar.agents {
		if agent.LastSeen.After(cutoff) {
			active = append(active, agent)
		}
	}
	return active
}

// AgentConfigSync synchronizes agent configurations
type AgentConfigSync struct {
	configs map[string]*AgentConfig
	mu      sync.RWMutex
}

// NewAgentConfigSync creates a new config sync
func NewAgentConfigSync() *AgentConfigSync {
	return &AgentConfigSync{
		configs: make(map[string]*AgentConfig),
	}
}

// Sync synchronizes configs
func (acs *AgentConfigSync) Sync(configs map[string]*AgentConfig) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	for id, config := range configs {
		acs.configs[id] = config
	}
}

// Get gets a config
func (acs *AgentConfigSync) Get(id string) (*AgentConfig, bool) {
	acs.mu.RLock()
	defer acs.mu.RUnlock()
	config, ok := acs.configs[id]
	return config, ok
}

// GetAll gets all configs
func (acs *AgentConfigSync) GetAll() map[string]*AgentConfig {
	acs.mu.RLock()
	defer acs.mu.RUnlock()

	result := make(map[string]*AgentConfig)
	for k, v := range acs.configs {
		result[k] = v
	}
	return result
}

// AgentEventBus provides event bus for agents
type AgentEventBus struct {
	subscribers map[string][]chan Event
	mu          sync.RWMutex
}

// NewAgentEventBus creates a new event bus
func NewAgentEventBus() *AgentEventBus {
	return &AgentEventBus{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe subscribes to a topic
func (aeb *AgentEventBus) Subscribe(topic string) <-chan Event {
	aeb.mu.Lock()
	defer aeb.mu.Unlock()

	ch := make(chan Event, 100)
	aeb.subscribers[topic] = append(aeb.subscribers[topic], ch)
	return ch
}

// Publish publishes an event
func (aeb *AgentEventBus) Publish(topic string, event Event) {
	aeb.mu.RLock()
	defer aeb.mu.RUnlock()

	for _, ch := range aeb.subscribers[topic] {
		select {
		case ch <- event:
		default:
		}
	}
}

// Unsubscribe unsubscribes from a topic
func (aeb *AgentEventBus) Unsubscribe(topic string, ch <-chan Event) {
	aeb.mu.Lock()
	defer aeb.mu.Unlock()

	subs := aeb.subscribers[topic]
	for i, sub := range subs {
		if sub == ch {
			aeb.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// AgentMetricsCollector collects metrics from agents
type AgentMetricsCollector struct {
	metrics map[string]map[string]float64
	mu      sync.RWMutex
}

// NewAgentMetricsCollector creates a new metrics collector
func NewAgentMetricsCollector() *AgentMetricsCollector {
	return &AgentMetricsCollector{
		metrics: make(map[string]map[string]float64),
	}
}

// Collect collects metrics from an agent
func (amc *AgentMetricsCollector) Collect(agentID string, metrics map[string]float64) {
	amc.mu.Lock()
	defer amc.mu.Unlock()
	amc.metrics[agentID] = metrics
}

// Get gets metrics for an agent
func (amc *AgentMetricsCollector) Get(agentID string) (map[string]float64, bool) {
	amc.mu.RLock()
	defer amc.mu.RUnlock()
	metrics, ok := amc.metrics[agentID]
	return metrics, ok
}

// GetAll gets all metrics
func (amc *AgentMetricsCollector) GetAll() map[string]map[string]float64 {
	amc.mu.RLock()
	defer amc.mu.RUnlock()

	result := make(map[string]map[string]float64)
	for k, v := range amc.metrics {
		result[k] = make(map[string]float64)
		for mk, mv := range v {
			result[k][mk] = mv
		}
	}
	return result
}

// GetTotal gets total metrics across all agents
func (amc *AgentMetricsCollector) GetTotal() map[string]float64 {
	amc.mu.RLock()
	defer amc.mu.RUnlock()

	total := make(map[string]float64)
	for _, metrics := range amc.metrics {
		for k, v := range metrics {
			total[k] += v
		}
	}
	return total
}

// AgentAlertManager manages agent alerts
type AgentAlertManager struct {
	alerts []AgentAlert
	mu     sync.RWMutex
}

// AgentAlert represents an agent alert
type AgentAlert struct {
	ID        string
	AgentID   string
	Type      string
	Message   string
	Timestamp time.Time
	Resolved  bool
}

// NewAgentAlertManager creates a new alert manager
func NewAgentAlertManager() *AgentAlertManager {
	return &AgentAlertManager{
		alerts: make([]AgentAlert, 0),
	}
}

// AddAlert adds an alert
func (aam *AgentAlertManager) AddAlert(agentID, alertType, message string) {
	aam.mu.Lock()
	defer aam.mu.Unlock()

	alert := AgentAlert{
		ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
		AgentID:   agentID,
		Type:      alertType,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  false,
	}
	aam.alerts = append(aam.alerts, alert)
}

// GetAlerts gets all alerts
func (aam *AgentAlertManager) GetAlerts() []AgentAlert {
	aam.mu.RLock()
	defer aam.mu.RUnlock()

	alerts := make([]AgentAlert, len(aam.alerts))
	copy(alerts, aam.alerts)
	return alerts
}

// GetUnresolvedAlerts gets unresolved alerts
func (aam *AgentAlertManager) GetUnresolvedAlerts() []AgentAlert {
	aam.mu.RLock()
	defer aam.mu.RUnlock()

	var result []AgentAlert
	for _, alert := range aam.alerts {
		if !alert.Resolved {
			result = append(result, alert)
		}
	}
	return result
}

// ResolveAlert resolves an alert
func (aam *AgentAlertManager) ResolveAlert(id string) {
	aam.mu.Lock()
	defer aam.mu.Unlock()

	for i, alert := range aam.alerts {
		if alert.ID == id {
			aam.alerts[i].Resolved = true
			break
		}
	}
}

// GetAlertsByAgent gets alerts for an agent
func (aam *AgentAlertManager) GetAlertsByAgent(agentID string) []AgentAlert {
	aam.mu.RLock()
	defer aam.mu.RUnlock()

	var result []AgentAlert
	for _, alert := range aam.alerts {
		if alert.AgentID == agentID {
			result = append(result, alert)
		}
	}
	return result
}

// GetAlertsByType gets alerts by type
func (aam *AgentAlertManager) GetAlertsByType(alertType string) []AgentAlert {
	aam.mu.RLock()
	defer aam.mu.RUnlock()

	var result []AgentAlert
	for _, alert := range aam.alerts {
		if alert.Type == alertType {
			result = append(result, alert)
		}
	}
	return result
}

// ClearAlerts clears all alerts
func (aam *AgentAlertManager) ClearAlerts() {
	aam.mu.Lock()
	defer aam.mu.Unlock()
	aam.alerts = make([]AgentAlert, 0)
}

// GetAlertCount returns the number of alerts
func (aam *AgentAlertManager) GetAlertCount() int {
	aam.mu.RLock()
	defer aam.mu.RUnlock()
	return len(aam.alerts)
}

// GetUnresolvedCount returns the number of unresolved alerts
func (aam *AgentAlertManager) GetUnresolvedCount() int {
	aam.mu.RLock()
	defer aam.mu.RUnlock()

	count := 0
	for _, alert := range aam.alerts {
		if !alert.Resolved {
			count++
		}
	}
	return count
}