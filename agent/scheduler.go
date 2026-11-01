// Package agent provides additional agent features
// including configuration, monitoring, and diagnostics
package agent

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// AgentScheduler schedules agent tasks
type AgentScheduler struct {
	tasks   []ScheduledTask
	mu      sync.Mutex
	running bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// ScheduledTask represents a scheduled task
type ScheduledTask struct {
	ID       string
	Name     string
	AgentID  string
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
	Enabled  bool
	Run      func(ctx context.Context) error
}

// NewAgentScheduler creates a new scheduler
func NewAgentScheduler() *AgentScheduler {
	return &AgentScheduler{
		tasks: make([]ScheduledTask, 0),
	}
}

// AddTask adds a scheduled task
func (as *AgentScheduler) AddTask(task ScheduledTask) {
	as.mu.Lock()
	defer as.mu.Unlock()
	task.NextRun = time.Now().Add(task.Interval)
	as.tasks = append(as.tasks, task)
}

// RemoveTask removes a scheduled task
func (as *AgentScheduler) RemoveTask(id string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	for i, task := range as.tasks {
		if task.ID == id {
			as.tasks = append(as.tasks[:i], as.tasks[i+1:]...)
			break
		}
	}
}

// Start starts the scheduler
func (as *AgentScheduler) Start(ctx context.Context) {
	as.ctx, as.cancel = context.WithCancel(ctx)
	as.running = true
	go as.run()
}

// Stop stops the scheduler
func (as *AgentScheduler) Stop() {
	as.running = false
	if as.cancel != nil {
		as.cancel()
	}
}

// run runs the scheduler loop
func (as *AgentScheduler) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-as.ctx.Done():
			return
		case <-ticker.C:
			if !as.running {
				return
			}
			as.runDueTasks()
		}
	}
}

// runDueTasks runs tasks that are due
func (as *AgentScheduler) runDueTasks() {
	as.mu.Lock()
	defer as.mu.Unlock()

	now := time.Now()
	for i := range as.tasks {
		if as.tasks[i].Enabled && now.After(as.tasks[i].NextRun) {
			go func(task *ScheduledTask) {
				task.Run(as.ctx)
				task.LastRun = time.Now()
				task.NextRun = time.Now().Add(task.Interval)
			}(&as.tasks[i])
		}
	}
}

// ListTasks lists all tasks
func (as *AgentScheduler) ListTasks() []ScheduledTask {
	as.mu.Lock()
	defer as.mu.Unlock()
	tasks := make([]ScheduledTask, len(as.tasks))
	copy(tasks, as.tasks)
	return tasks
}

// EnableTask enables a task
func (as *AgentScheduler) EnableTask(id string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	for i := range as.tasks {
		if as.tasks[i].ID == id {
			as.tasks[i].Enabled = true
			break
		}
	}
}

// DisableTask disables a task
func (as *AgentScheduler) DisableTask(id string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	for i := range as.tasks {
		if as.tasks[i].ID == id {
			as.tasks[i].Enabled = false
			break
		}
	}
}

// AgentPool manages a pool of agents
type AgentPool struct {
	agents  []*Agent
	mu      sync.Mutex
	maxSize int
	created int
}

// NewAgentPool creates a new agent pool
func NewAgentPool(maxSize int) *AgentPool {
	return &AgentPool{
		agents:  make([]*Agent, 0),
		maxSize: maxSize,
	}
}

// Get gets an agent from the pool
func (ap *AgentPool) Get() (*Agent, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if len(ap.agents) > 0 {
		agent := ap.agents[len(ap.agents)-1]
		ap.agents = ap.agents[:len(ap.agents)-1]
		return agent, nil
	}

	if ap.created >= ap.maxSize {
		return nil, fmt.Errorf("pool exhausted")
	}

	ap.created++
	return nil, nil
}

// Put returns an agent to the pool
func (ap *AgentPool) Put(agent *Agent) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	ap.agents = append(ap.agents, agent)
}

// Size returns the pool size
func (ap *AgentPool) Size() int {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	return len(ap.agents)
}

// Created returns the number of created agents
func (ap *AgentPool) Created() int {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	return ap.created
}

// MaxSize returns the max pool size
func (ap *AgentPool) MaxSize() int {
	return ap.maxSize
}

// AgentLoadBalancer balances load across agents
type AgentLoadBalancer struct {
	agents   []*Agent
	weights  []int
	mu       sync.RWMutex
	strategy string
	counter  int
}

// NewAgentLoadBalancer creates a new load balancer
func NewAgentLoadBalancer(strategy string) *AgentLoadBalancer {
	return &AgentLoadBalancer{
		agents:   make([]*Agent, 0),
		weights:  make([]int, 0),
		strategy: strategy,
	}
}

// AddAgent adds an agent
func (alb *AgentLoadBalancer) AddAgent(agent *Agent, weight int) {
	alb.mu.Lock()
	defer alb.mu.Unlock()
	alb.agents = append(alb.agents, agent)
	alb.weights = append(alb.weights, weight)
}

// RemoveAgent removes an agent
func (alb *AgentLoadBalancer) RemoveAgent(agent *Agent) {
	alb.mu.Lock()
	defer alb.mu.Unlock()
	for i, a := range alb.agents {
		if a == agent {
			alb.agents = append(alb.agents[:i], alb.agents[i+1:]...)
			alb.weights = append(alb.weights[:i], alb.weights[i+1:]...)
			break
		}
	}
}

// Select selects an agent
func (alb *AgentLoadBalancer) Select() (*Agent, error) {
	alb.mu.RLock()
	defer alb.mu.RUnlock()

	if len(alb.agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	switch alb.strategy {
	case "round_robin":
		agent := alb.agents[alb.counter%len(alb.agents)]
		alb.counter++
		return agent, nil
	case "random":
		idx := alb.counter % len(alb.agents)
		alb.counter++
		return alb.agents[idx], nil
	case "weighted":
		return alb.weighted(), nil
	default:
		return alb.agents[0], nil
	}
}

func (alb *AgentLoadBalancer) weighted() *Agent {
	totalWeight := 0
	for _, w := range alb.weights {
		totalWeight += w
	}

	r := 0
	if totalWeight > 0 {
		r = rand.Intn(totalWeight)
	}

	current := 0
	for i, w := range alb.weights {
		current += w
		if r < current {
			return alb.agents[i]
		}
	}
	return alb.agents[0]
}

// GetAgentCount returns the number of agents
func (alb *AgentLoadBalancer) GetAgentCount() int {
	alb.mu.RLock()
	defer alb.mu.RUnlock()
	return len(alb.agents)
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

// AgentMetricsAggregator aggregates metrics from multiple agents
type AgentMetricsAggregator struct {
	agents  map[string]*Agent
	metrics map[string]map[string]float64
	mu      sync.RWMutex
}

// NewAgentMetricsAggregator creates a new metrics aggregator
func NewAgentMetricsAggregator() *AgentMetricsAggregator {
	return &AgentMetricsAggregator{
		agents:  make(map[string]*Agent),
		metrics: make(map[string]map[string]float64),
	}
}

// RegisterAgent registers an agent
func (ama *AgentMetricsAggregator) RegisterAgent(id string, agent *Agent) {
	ama.mu.Lock()
	defer ama.mu.Unlock()
	ama.agents[id] = agent
	ama.metrics[id] = make(map[string]float64)
}

// UnregisterAgent unregisters an agent
func (ama *AgentMetricsAggregator) UnregisterAgent(id string) {
	ama.mu.Lock()
	defer ama.mu.Unlock()
	delete(ama.agents, id)
	delete(ama.metrics, id)
}

// Collect collects metrics from all agents
func (ama *AgentMetricsAggregator) Collect() {
	ama.mu.RLock()
	defer ama.mu.RUnlock()

	for id, agent := range ama.agents {
		agentMetrics := agent.GetMetrics()
		ama.metrics[id]["events_collected"] = float64(agentMetrics.EventsCollected)
		ama.metrics[id]["events_processed"] = float64(agentMetrics.EventsProcessed)
		ama.metrics[id]["events_dropped"] = float64(agentMetrics.EventsDropped)
		ama.metrics[id]["events_exported"] = float64(agentMetrics.EventsExported)
	}
}

// GetMetrics returns aggregated metrics
func (ama *AgentMetricsAggregator) GetMetrics() map[string]map[string]float64 {
	ama.mu.RLock()
	defer ama.mu.RUnlock()

	result := make(map[string]map[string]float64)
	for id, metrics := range ama.metrics {
		result[id] = make(map[string]float64)
		for k, v := range metrics {
			result[id][k] = v
		}
	}
	return result
}

// GetTotalMetrics returns total metrics across all agents
func (ama *AgentMetricsAggregator) GetTotalMetrics() map[string]float64 {
	ama.mu.RLock()
	defer ama.mu.RUnlock()

	total := make(map[string]float64)
	for _, metrics := range ama.metrics {
		for k, v := range metrics {
			total[k] += v
		}
	}
	return total
}

// AgentConfigManager manages agent configurations
type AgentConfigManager struct {
	configs map[string]*AgentConfig
	mu      sync.RWMutex
}

// NewAgentConfigManager creates a new config manager
func NewAgentConfigManager() *AgentConfigManager {
	return &AgentConfigManager{
		configs: make(map[string]*AgentConfig),
	}
}

// Set sets a config for an agent
func (acm *AgentConfigManager) Set(id string, config *AgentConfig) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.configs[id] = config
}

// Get gets a config for an agent
func (acm *AgentConfigManager) Get(id string) (*AgentConfig, bool) {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	config, ok := acm.configs[id]
	return config, ok
}

// Delete deletes a config
func (acm *AgentConfigManager) Delete(id string) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	delete(acm.configs, id)
}

// List lists all configs
func (acm *AgentConfigManager) List() map[string]*AgentConfig {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	result := make(map[string]*AgentConfig)
	for k, v := range acm.configs {
		result[k] = v
	}
	return result
}