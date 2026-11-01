// Package agent provides additional agent features
// including configuration, lifecycle, and service discovery
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AgentService provides agent service discovery
type AgentService struct {
	agentID   string
	endpoint  string
	metadata  map[string]string
	heartbeat time.Duration
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
}

// NewAgentService creates a new agent service
func NewAgentService(agentID, endpoint string) *AgentService {
	return &AgentService{
		agentID:   agentID,
		endpoint:  endpoint,
		metadata:  make(map[string]string),
		heartbeat: 30 * time.Second,
	}
}

// Start starts the agent service
func (as *AgentService) Start(ctx context.Context) {
	as.ctx, as.cancel = context.WithCancel(ctx)
	as.running = true
	go as.heartbeatLoop()
}

// Stop stops the agent service
func (as *AgentService) Stop() {
	as.running = false
	if as.cancel != nil {
		as.cancel()
	}
}

// heartbeatLoop sends heartbeats
func (as *AgentService) heartbeatLoop() {
	ticker := time.NewTicker(as.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-as.ctx.Done():
			return
		case <-ticker.C:
			as.sendHeartbeat()
		}
	}
}

// sendHeartbeat sends a heartbeat
func (as *AgentService) sendHeartbeat() {
	// In production, would send to service discovery
}

// SetMetadata sets metadata
func (as *AgentService) SetMetadata(key, value string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.metadata[key] = value
}

// GetMetadata gets metadata
func (as *AgentService) GetMetadata(key string) (string, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	v, ok := as.metadata[key]
	return v, ok
}

// GetAllMetadata returns all metadata
func (as *AgentService) GetAllMetadata() map[string]string {
	as.mu.RLock()
	defer as.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range as.metadata {
		result[k] = v
	}
	return result
}

// AgentConfigLoader loads agent configuration
type AgentConfigLoader struct {
	configPath string
	watch      bool
}

// NewAgentConfigLoader creates a new config loader
func NewAgentConfigLoader(configPath string) *AgentConfigLoader {
	return &AgentConfigLoader{
		configPath: configPath,
	}
}

// Load loads configuration from file
func (acl *AgentConfigLoader) Load() (*AgentConfig, error) {
	data, err := os.ReadFile(acl.configPath)
	if err != nil {
		return nil, err
	}

	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves configuration to file
func (acl *AgentConfigLoader) Save(config *AgentConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(acl.configPath)
	os.MkdirAll(dir, 0755)

	return os.WriteFile(acl.configPath, data, 0644)
}

// Watch watches for configuration changes
func (acl *AgentConfigLoader) Watch(ctx context.Context, onChange func(*AgentConfig)) {
	if !acl.watch {
		return
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		var lastMod time.Time
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				info, err := os.Stat(acl.configPath)
				if err != nil {
					continue
				}
				if info.ModTime().After(lastMod) {
					lastMod = info.ModTime()
					config, err := acl.Load()
					if err == nil {
						onChange(config)
					}
				}
			}
		}
	}()
}

// AgentLifecycle manages agent lifecycle
type AgentLifecycle struct {
	agent    *Agent
	state    LifecycleState
	hooks    map[LifecycleEvent][]LifecycleHook
	mu       sync.RWMutex
}

// LifecycleState represents agent state
type LifecycleState string

const (
	StateInitializing LifecycleState = "initializing"
	StateStarting     LifecycleState = "starting"
	StateRunning      LifecycleState = "running"
	StateStopping     LifecycleState = "stopping"
	StateStopped      LifecycleState = "stopped"
	StateError        LifecycleState = "error"
)

// LifecycleEvent represents a lifecycle event
type LifecycleEvent string

const (
	EventBeforeStart  LifecycleEvent = "before_start"
	EventAfterStart   LifecycleEvent = "after_start"
	EventBeforeStop   LifecycleEvent = "before_stop"
	EventAfterStop    LifecycleEvent = "after_stop"
	EventError        LifecycleEvent = "error"
)

// LifecycleHook is a lifecycle hook function
type LifecycleHook func(ctx context.Context) error

// NewAgentLifecycle creates a new lifecycle manager
func NewAgentLifecycle(agent *Agent) *AgentLifecycle {
	return &AgentLifecycle{
		agent: agent,
		state: StateInitializing,
		hooks: make(map[LifecycleEvent][]LifecycleHook),
	}
}

// RegisterHook registers a lifecycle hook
func (al *AgentLifecycle) RegisterHook(event LifecycleEvent, hook LifecycleHook) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.hooks[event] = append(al.hooks[event], hook)
}

// Start starts the agent
func (al *AgentLifecycle) Start(ctx context.Context) error {
	al.setState(StateStarting)

	// Run before start hooks
	if err := al.runHooks(ctx, EventBeforeStart); err != nil {
		al.setState(StateError)
		return err
	}

	// Start agent
	if err := al.agent.Start(ctx); err != nil {
		al.setState(StateError)
		al.runHooks(ctx, EventError)
		return err
	}

	// Run after start hooks
	if err := al.runHooks(ctx, EventAfterStart); err != nil {
		al.setState(StateError)
		return err
	}

	al.setState(StateRunning)
	return nil
}

// Stop stops the agent
func (al *AgentLifecycle) Stop() error {
	al.setState(StateStopping)

	// Run before stop hooks
	ctx := context.Background()
	if err := al.runHooks(ctx, EventBeforeStop); err != nil {
		// Log error but continue
	}

	// Stop agent
	if err := al.agent.Stop(); err != nil {
		al.setState(StateError)
		al.runHooks(ctx, EventError)
		return err
	}

	// Run after stop hooks
	if err := al.runHooks(ctx, EventAfterStop); err != nil {
		// Log error
	}

	al.setState(StateStopped)
	return nil
}

// GetState returns the current state
func (al *AgentLifecycle) GetState() LifecycleState {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.state
}

func (al *AgentLifecycle) setState(state LifecycleState) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.state = state
}

func (al *AgentLifecycle) runHooks(ctx context.Context, event LifecycleEvent) error {
	al.mu.RLock()
	hooks := al.hooks[event]
	al.mu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

// AgentRegistry manages multiple agents
type AgentRegistry struct {
	agents map[string]*Agent
	mu     sync.RWMutex
}

// NewAgentRegistry creates a new agent registry
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*Agent),
	}
}

// Register registers an agent
func (ar *AgentRegistry) Register(id string, agent *Agent) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, exists := ar.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	ar.agents[id] = agent
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

// Get gets an agent by ID
func (ar *AgentRegistry) Get(id string) (*Agent, bool) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	agent, ok := ar.agents[id]
	return agent, ok
}

// List lists all agents
func (ar *AgentRegistry) List() []*Agent {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	agents := make([]*Agent, 0, len(ar.agents))
	for _, agent := range ar.agents {
		agents = append(agents, agent)
	}
	return agents
}

// Count returns the number of agents
func (ar *AgentRegistry) Count() int {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return len(ar.agents)
}

// StartAll starts all agents
func (ar *AgentRegistry) StartAll(ctx context.Context) error {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	for id, agent := range ar.agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
	}
	return nil
}

// StopAll stops all agents
func (ar *AgentRegistry) StopAll() error {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	for id, agent := range ar.agents {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
	}
	return nil
}

// AgentDiscovery discovers agents
type AgentDiscovery struct {
	agents   map[string]*AgentInfo
	mu       sync.RWMutex
	registry *AgentRegistry
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
func NewAgentDiscovery(registry *AgentRegistry) *AgentDiscovery {
	return &AgentDiscovery{
		agents:   make(map[string]*AgentInfo),
		registry: registry,
	}
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

// GetActive gets active agents (seen in last 5 minutes)
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