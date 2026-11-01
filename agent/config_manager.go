// Package agent provides config management
package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// ConfigManager manages agent config
type ConfigManager struct {
	config map[string]interface{}
	path   string
	mu     sync.RWMutex
}

// NewConfigManager creates manager
func NewConfigManager(path string) *ConfigManager {
	cm := &ConfigManager{config: make(map[string]interface{}), path: path}
	cm.Load()
	return cm
}

// Load loads config from file
func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &cm.config)
}

// Save saves config to file
func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.path, data, 0644)
}

// Get returns value
func (cm *ConfigManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	v, ok := cm.config[key]
	return v, ok
}

// Set sets value
func (cm *ConfigManager) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config[key] = value
}

// Delete deletes value
func (cm *ConfigManager) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.config, key)
}

// All returns all config
func (cm *ConfigManager) All() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	out := make(map[string]interface{}, len(cm.config))
	for k, v := range cm.config {
		out[k] = v
	}
	return out
}

// String returns summary
func (cm *ConfigManager) String() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return fmt.Sprintf("ConfigManager: %d keys", len(cm.config))
}
