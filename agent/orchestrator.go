// Package agent provides additional agent features
// including configuration, monitoring, and diagnostics
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// AgentOrchestrator orchestrates multiple agents
type AgentOrchestrator struct {
	agents   map[string]*Agent
	groups   map[string]*AgentGroup
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// AgentGroup represents a group of agents
type AgentGroup struct {
	ID       string
	Name     string
	Agents   []string
	Strategy string // round_robin, random, least_loaded
	mu       sync.RWMutex
}

// NewAgentOrchestrator creates a new orchestrator
func NewAgentOrchestrator() *AgentOrchestrator {
	return &AgentOrchestrator{
		agents: make(map[string]*Agent),
		groups: make(map[string]*AgentGroup),
	}
}

// RegisterAgent registers an agent
func (ao *AgentOrchestrator) RegisterAgent(id string, agent *Agent) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	if _, exists := ao.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	ao.agents[id] = agent
	return nil
}

// UnregisterAgent unregisters an agent
func (ao *AgentOrchestrator) UnregisterAgent(id string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	if _, exists := ao.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	delete(ao.agents, id)
	return nil
}

// CreateGroup creates an agent group
func (ao *AgentOrchestrator) CreateGroup(id, name, strategy string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	if _, exists := ao.groups[id]; exists {
		return fmt.Errorf("group %s already exists", id)
	}

	ao.groups[id] = &AgentGroup{
		ID:       id,
		Name:     name,
		Agents:   make([]string, 0),
		Strategy: strategy,
	}
	return nil
}

// AddToGroup adds an agent to a group
func (ao *AgentOrchestrator) AddToGroup(groupID, agentID string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	group, ok := ao.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}

	if _, ok := ao.agents[agentID]; !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	group.mu.Lock()
	defer group.mu.Unlock()
	group.Agents = append(group.Agents, agentID)
	return nil
}

// RemoveFromGroup removes an agent from a group
func (ao *AgentOrchestrator) RemoveFromGroup(groupID, agentID string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	group, ok := ao.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	for i, id := range group.Agents {
		if id == agentID {
			group.Agents = append(group.Agents[:i], group.Agents[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("agent %s not in group %s", agentID, groupID)
}

// SelectAgent selects an agent from a group
func (ao *AgentOrchestrator) SelectAgent(groupID string) (*Agent, error) {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	group, ok := ao.groups[groupID]
	if !ok {
		return nil, fmt.Errorf("group %s not found", groupID)
	}

	group.mu.RLock()
	defer group.mu.RUnlock()

	if len(group.Agents) == 0 {
		return nil, fmt.Errorf("no agents in group %s", groupID)
	}

	switch group.Strategy {
	case "round_robin":
		return ao.roundRobin(group)
	case "random":
		return ao.random(group)
	case "least_loaded":
		return ao.leastLoaded(group)
	default:
		return ao.roundRobin(group)
	}
}

func (ao *AgentOrchestrator) roundRobin(group *AgentGroup) (*Agent, error) {
	// Simplified round robin
	agentID := group.Agents[0]
	return ao.agents[agentID], nil
}

func (ao *AgentOrchestrator) random(group *AgentGroup) (*Agent, error) {
	idx := len(group.Agents) % len(group.Agents)
	agentID := group.Agents[idx]
	return ao.agents[agentID], nil
}

func (ao *AgentOrchestrator) leastLoaded(group *AgentGroup) (*Agent, error) {
	// Simplified - would check load metrics
	agentID := group.Agents[0]
	return ao.agents[agentID], nil
}

// Start starts all agents
func (ao *AgentOrchestrator) Start(ctx context.Context) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.ctx, ao.cancel = context.WithCancel(ctx)
	ao.running = true

	for id, agent := range ao.agents {
		if err := agent.Start(ao.ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
	}

	return nil
}

// Stop stops all agents
func (ao *AgentOrchestrator) Stop() error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.running = false
	if ao.cancel != nil {
		ao.cancel()
	}

	for id, agent := range ao.agents {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
	}

	return nil
}

// GetStatus returns orchestrator status
func (ao *AgentOrchestrator) GetStatus() map[string]interface{} {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	status := make(map[string]interface{})
	status["running"] = ao.running
	status["agent_count"] = len(ao.agents)
	status["group_count"] = len(ao.groups)

	agents := make(map[string]bool)
	for id, agent := range ao.agents {
		agents[id] = agent.IsRunning()
	}
	status["agents"] = agents

	return status
}

// AgentScheduler schedules agent tasks
type AgentScheduler struct {
	tasks    []ScheduledTask
	mu       sync.Mutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
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

// AgentPool manages a pool of agents
type AgentPool struct {
	agents   []*Agent
	mu       sync.Mutex
	maxSize  int
	created  int
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
	return nil, nil // Would create new agent
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
	agents  []*Agent
	mu      sync.RWMutex
	strategy string
	counter int
}

// NewAgentLoadBalancer creates a new load balancer
func NewAgentLoadBalancer(strategy string) *AgentLoadBalancer {
	return &AgentLoadBalancer{
		agents:   make([]*Agent, 0),
		strategy: strategy,
	}
}

// AddAgent adds an agent
func (alb *AgentLoadBalancer) AddAgent(agent *Agent) {
	alb.mu.Lock()
	defer alb.mu.Unlock()
	alb.agents = append(alb.agents, agent)
}

// RemoveAgent removes an agent
func (alb *AgentLoadBalancer) RemoveAgent(agent *Agent) {
	alb.mu.Lock()
	defer alb.mu.Unlock()
	for i, a := range alb.agents {
		if a == agent {
			alb.agents = append(alb.agents[:i], alb.agents[i+1:]...)
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
		idx := len(alb.agents) % len(alb.agents)
		return alb.agents[idx], nil
	default:
		return alb.agents[0], nil
	}
}

// GetAgentCount returns the number of agents
func (alb *AgentLoadBalancer) GetAgentCount() int {
	alb.mu.RLock()
	defer alb.mu.RUnlock()
	return len(alb.agents)
}