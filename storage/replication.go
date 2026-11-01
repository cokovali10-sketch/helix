// Package storage provides replication support
package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Replica represents a storage replica
type Replica struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	Role     string    `json:"role"` // leader, follower
	Status   string    `json:"status"`
	Lag      time.Duration `json:"lag"`
	LastSync time.Time `json:"last_sync"`
}

// ReplicationManager manages replicas
type ReplicationManager struct {
	replicas map[string]*Replica
	mu       sync.RWMutex
}

// NewReplicationManager creates manager
func NewReplicationManager() *ReplicationManager {
	return &ReplicationManager{replicas: make(map[string]*Replica)}
}

// AddReplica adds a replica
func (rm *ReplicationManager) AddReplica(r Replica) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.replicas[r.ID] = &r
}

// RemoveReplica removes a replica
func (rm *ReplicationManager) RemoveReplica(id string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.replicas, id)
}

// ListReplicas returns replicas
func (rm *ReplicationManager) ListReplicas() []*Replica {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	out := []*Replica{}
	for _, r := range rm.replicas {
		out = append(out, r)
	}
	return out
}

// GetReplica returns a replica
func (rm *ReplicationManager) GetReplica(id string) (*Replica, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.replicas[id]
	return r, ok
}

// UpdateLag updates lag
func (rm *ReplicationManager) UpdateLag(id string, lag time.Duration) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if r, ok := rm.replicas[id]; ok {
		r.Lag = lag
		r.LastSync = time.Now()
	}
}

// Leader returns leader replica
func (rm *ReplicationManager) Leader() *Replica {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	for _, r := range rm.replicas {
		if r.Role == "leader" {
			return r
		}
	}
	return nil
}

// Followers returns follower replicas
func (rm *ReplicationManager) Followers() []*Replica {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	out := []*Replica{}
	for _, r := range rm.replicas {
		if r.Role == "follower" {
			out = append(out, r)
		}
	}
	return out
}

// HealthCheck checks replica health
func (rm *ReplicationManager) HealthCheck() map[string]string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	res := map[string]string{}
	for id, r := range rm.replicas {
		if r.Status == "healthy" {
			res[id] = "ok"
		} else {
			res[id] = "degraded"
		}
	}
	return res
}

// MarshalJSON for manager
func (rm *ReplicationManager) MarshalJSON() ([]byte, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return json.Marshal(rm.replicas)
}

// SyncStatus returns sync status
func (rm *ReplicationManager) SyncStatus() string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if len(rm.replicas) == 0 {
		return "no replicas"
	}
	return fmt.Sprintf("%d replicas", len(rm.replicas))
}
