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

// AgentCluster manages a cluster of agents
type AgentCluster struct {
	agents   map[string]*Agent
	cluster  *ClusterConfig
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// ClusterConfig holds cluster configuration
type ClusterConfig struct {
	Name          string
	LeaderID      string
	Replicas      int
	Sharding      bool
	Consensus     string // raft, paxos, gossip
	Heartbeat     time.Duration
	ElectionTimeout time.Duration
}

// NewAgentCluster creates a new agent cluster
func NewAgentCluster(config *ClusterConfig) *AgentCluster {
	return &AgentCluster{
		agents:  make(map[string]*Agent),
		cluster: config,
	}
}

// RegisterAgent registers an agent in the cluster
func (ac *AgentCluster) RegisterAgent(id string, agent *Agent) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	ac.agents[id] = agent
	return nil
}

// UnregisterAgent unregisters an agent
func (ac *AgentCluster) UnregisterAgent(id string) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	delete(ac.agents, id)
	return nil
}

// GetAgent gets an agent by ID
func (ac *AgentCluster) GetAgent(id string) (*Agent, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	agent, ok := ac.agents[id]
	return agent, ok
}

// ListAgents lists all agents
func (ac *AgentCluster) ListAgents() []*Agent {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	agents := make([]*Agent, 0, len(ac.agents))
	for _, agent := range ac.agents {
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

	for id, agent := range ac.agents {
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

	for id, agent := range ac.agents {
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
	status["agent_count"] = len(ac.agents)
	status["cluster_name"] = ac.cluster.Name
	status["leader"] = ac.cluster.LeaderID
	status["consensus"] = ac.cluster.Consensus

	agents := make(map[string]bool)
	for id, agent := range ac.agents {
		agents[id] = agent.IsRunning()
	}
	status["agents"] = agents

	return status
}

// AgentReplicator replicates data across agents
type AgentReplicator struct {
	agents   map[string]*Agent
	mu       sync.RWMutex
	strategy string // sync, async, semi_sync
}

// NewAgentReplicator creates a new replicator
func NewAgentReplicator(strategy string) *AgentReplicator {
	return &AgentReplicator{
		agents:   make(map[string]*Agent),
		strategy: strategy,
	}
}

// AddAgent adds an agent for replication
func (ar *AgentReplicator) AddAgent(id string, agent *Agent) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.agents[id] = agent
}

// RemoveAgent removes an agent
func (ar *AgentReplicator) RemoveAgent(id string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	delete(ar.agents, id)
}

// Replicate replicates an event to all agents
func (ar *AgentReplicator) Replicate(event Event) error {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	switch ar.strategy {
	case "sync":
		return ar.replicateSync(event)
	case "async":
		return ar.replicateAsync(event)
	case "semi_sync":
		return ar.replicateSemiSync(event)
	default:
		return ar.replicateSync(event)
	}
}

func (ar *AgentReplicator) replicateSync(event Event) error {
	for _, agent := range ar.agents {
		if err := agent.ProcessEvent(context.Background(), &event); err != nil {
			return err
		}
	}
	return nil
}

func (ar *AgentReplicator) replicateAsync(event Event) error {
	for _, agent := range ar.agents {
		go agent.ProcessEvent(context.Background(), &event)
	}
	return nil
}

func (ar *AgentReplicator) replicateSemiSync(event Event) error {
	// Replicate to at least one agent synchronously
	count := 0
	for _, agent := range ar.agents {
		if err := agent.ProcessEvent(context.Background(), &event); err != nil {
			return err
		}
		count++
		if count >= 1 {
			break
		}
	}
	// Replicate to rest asynchronously
	for _, agent := range ar.agents {
		go agent.ProcessEvent(context.Background(), &event)
	}
	return nil
}

// AgentSharder shards data across agents
type AgentSharder struct {
	shards   map[int]*Agent
	numShards int
	mu       sync.RWMutex
}

// NewAgentSharder creates a new sharder
func NewAgentSharder(numShards int) *AgentSharder {
	return &AgentSharder{
		shards:    make(map[int]*Agent),
		numShards: numShards,
	}
}

// AddShard adds a shard
func (as *AgentSharder) AddShard(shardID int, agent *Agent) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.shards[shardID] = agent
}

// GetShard gets the shard for a key
func (as *AgentSharder) GetShard(key string) (*Agent, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	shardID := as.hashKey(key) % as.numShards
	agent, ok := as.shards[shardID]
	if !ok {
		return nil, fmt.Errorf("shard %d not found", shardID)
	}
	return agent, nil
}

func (as *AgentSharder) hashKey(key string) int {
	hash := 0
	for _, c := range key {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// GetShardCount returns the number of shards
func (as *AgentSharder) GetShardCount() int {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return len(as.shards)
}

// AgentBalancer balances load across agents
type AgentBalancer struct {
	agents   []*Agent
	weights  []int
	mu       sync.RWMutex
	strategy string
}

// NewAgentBalancer creates a new load balancer
func NewAgentBalancer(strategy string) *AgentBalancer {
	return &AgentBalancer{
		agents:   make([]*Agent, 0),
		weights:  make([]int, 0),
		strategy: strategy,
	}
}

// AddAgent adds an agent with weight
func (ab *AgentBalancer) AddAgent(agent *Agent, weight int) {
	ab.mu.Lock()
	defer ab.mu.Unlock()
	ab.agents = append(ab.agents, agent)
	ab.weights = append(ab.weights, weight)
}

// RemoveAgent removes an agent
func (ab *AgentBalancer) RemoveAgent(agent *Agent) {
	ab.mu.Lock()
	defer ab.mu.Unlock()
	for i, a := range ab.agents {
		if a == agent {
			ab.agents = append(ab.agents[:i], ab.agents[i+1:]...)
			ab.weights = append(ab.weights[:i], ab.weights[i+1:]...)
			break
		}
	}
}

// Select selects an agent
func (ab *AgentBalancer) Select() (*Agent, error) {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	if len(ab.agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	switch ab.strategy {
	case "round_robin":
		return ab.roundRobin(), nil
	case "weighted":
		return ab.weighted(), nil
	case "least_connections":
		return ab.leastConnections(), nil
	default:
		return ab.roundRobin(), nil
	}
}

func (ab *AgentBalancer) roundRobin() *Agent {
	// Simplified round robin
	return ab.agents[0]
}

func (ab *AgentBalancer) weighted() *Agent {
	totalWeight := 0
	for _, w := range ab.weights {
		totalWeight += w
	}

	r := 0
	if totalWeight > 0 {
		r = rand.Intn(totalWeight)
	}

	current := 0
	for i, w := range ab.weights {
		current += w
		if r < current {
			return ab.agents[i]
		}
	}
	return ab.agents[0]
}

func (ab *AgentBalancer) leastConnections() *Agent {
	// Simplified - would track connections
	return ab.agents[0]
}

// GetAgentCount returns the number of agents
func (ab *AgentBalancer) GetAgentCount() int {
	ab.mu.RLock()
	defer ab.mu.RUnlock()
	return len(ab.agents)
}

// AgentHealthMonitor monitors agent health
type AgentHealthMonitor struct {
	agents   map[string]*Agent
	health   map[string]HealthStatus
	mu       sync.RWMutex
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// HealthStatus represents health status
type HealthStatus struct {
	Healthy   bool
	LastCheck time.Time
	Error     error
	Metrics   map[string]interface{}
}

// NewAgentHealthMonitor creates a new health monitor
func NewAgentHealthMonitor(interval time.Duration) *AgentHealthMonitor {
	return &AgentHealthMonitor{
		agents:   make(map[string]*Agent),
		health:   make(map[string]HealthStatus),
		interval: interval,
	}
}

// RegisterAgent registers an agent for monitoring
func (ahm *AgentHealthMonitor) RegisterAgent(id string, agent *Agent) {
	ahm.mu.Lock()
	defer ahm.mu.Unlock()
	ahm.agents[id] = agent
	ahm.health[id] = HealthStatus{Healthy: true, LastCheck: time.Now()}
}

// UnregisterAgent unregisters an agent
func (ahm *AgentHealthMonitor) UnregisterAgent(id string) {
	ahm.mu.Lock()
	defer ahm.mu.Unlock()
	delete(ahm.agents, id)
	delete(ahm.health, id)
}

// Start starts monitoring
func (ahm *AgentHealthMonitor) Start(ctx context.Context) {
	ahm.ctx, ahm.cancel = context.WithCancel(ctx)
	ahm.running = true
	go ahm.run()
}

// Stop stops monitoring
func (ahm *AgentHealthMonitor) Stop() {
	ahm.running = false
	if ahm.cancel != nil {
		ahm.cancel()
	}
}

// run runs the monitoring loop
func (ahm *AgentHealthMonitor) run() {
	ticker := time.NewTicker(ahm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ahm.ctx.Done():
			return
		case <-ticker.C:
			if !ahm.running {
				return
			}
			ahm.checkAll()
		}
	}
}

// checkAll checks all agents
func (ahm *AgentHealthMonitor) checkAll() {
	ahm.mu.RLock()
	agents := make(map[string]*Agent)
	for id, agent := range ahm.agents {
		agents[id] = agent
	}
	ahm.mu.RUnlock()

	for id, agent := range agents {
		healthy := agent.IsRunning()
		ahm.mu.Lock()
		ahm.health[id] = HealthStatus{
			Healthy:   healthy,
			LastCheck: time.Now(),
			Metrics:   map[string]interface{}{"running": healthy},
		}
		ahm.mu.Unlock()
	}
}

// GetHealth gets health status for an agent
func (ahm *AgentHealthMonitor) GetHealth(id string) (HealthStatus, bool) {
	ahm.mu.RLock()
	defer ahm.mu.RUnlock()
	health, ok := ahm.health[id]
	return health, ok
}

// GetAllHealth gets health status for all agents
func (ahm *AgentHealthMonitor) GetAllHealth() map[string]HealthStatus {
	ahm.mu.RLock()
	defer ahm.mu.RUnlock()
	result := make(map[string]HealthStatus)
	for id, health := range ahm.health {
		result[id] = health
	}
	return result
}

// GetHealthyAgents returns healthy agents
func (ahm *AgentHealthMonitor) GetHealthyAgents() []string {
	ahm.mu.RLock()
	defer ahm.mu.RUnlock()

	var healthy []string
	for id, health := range ahm.health {
		if health.Healthy {
			healthy = append(healthy, id)
		}
	}
	return healthy
}

// GetUnhealthyAgents returns unhealthy agents
func (ahm *AgentHealthMonitor) GetUnhealthyAgents() []string {
	ahm.mu.RLock()
	defer ahm.mu.RUnlock()

	var unhealthy []string
	for id, health := range ahm.health {
		if !health.Healthy {
			unhealthy = append(unhealthy, id)
		}
	}
	return unhealthy
}