// Package agent provides plugin manager
package agent

import (
	"fmt"
	"sync"
)

// Plugin interface
type Plugin interface {
	Name() string
	Init() error
	Start() error
	Stop() error
}

// PluginManager manages plugins
type PluginManager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
	order   []string
}

// NewPluginManager creates manager
func NewPluginManager() *PluginManager {
	return &PluginManager{plugins: make(map[string]Plugin)}
}

// Register registers plugin
func (pm *PluginManager) Register(p Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.plugins[p.Name()]; exists {
		return fmt.Errorf("plugin already registered: %s", p.Name())
	}
	pm.plugins[p.Name()] = p
	pm.order = append(pm.order, p.Name())
	return nil
}

// Unregister unregisters plugin
func (pm *PluginManager) Unregister(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.plugins, name)
}

// Get returns plugin
func (pm *PluginManager) Get(name string) (Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	p, ok := pm.plugins[name]
	return p, ok
}

// List returns plugin names
func (pm *PluginManager) List() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := make([]string, len(pm.order))
	copy(out, pm.order)
	return out
}

// StartAll starts all plugins
func (pm *PluginManager) StartAll() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for _, name := range pm.order {
		if err := pm.plugins[name].Start(); err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err)
		}
	}
	return nil
}

// StopAll stops all plugins
func (pm *PluginManager) StopAll() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for i := len(pm.order) - 1; i >= 0; i-- {
		if err := pm.plugins[pm.order[i]].Stop(); err != nil {
			return fmt.Errorf("failed to stop %s: %w", pm.order[i], err)
		}
	}
	return nil
}

// InitAll inits all plugins
func (pm *PluginManager) InitAll() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for _, name := range pm.order {
		if err := pm.plugins[name].Init(); err != nil {
			return fmt.Errorf("failed to init %s: %w", name, err)
		}
	}
	return nil
}

// Count returns count
func (pm *PluginManager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.plugins)
}
