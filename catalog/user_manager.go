// Package catalog provides user management
package catalog

import (
	"fmt"
	"sync"
	"time"
)

// User represents a user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Metadata  map[string]string `json:"metadata"`
}

// UserManager manages users
type UserManager struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewUserManager creates manager
func NewUserManager() *UserManager {
	return &UserManager{users: make(map[string]*User)}
}

// Create creates user
func (um *UserManager) Create(u User) *User {
	um.mu.Lock()
	defer um.mu.Unlock()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	if u.Metadata == nil {
		u.Metadata = map[string]string{}
	}
	um.users[u.ID] = &u
	return &u
}

// Get returns user
func (um *UserManager) Get(id string) (*User, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	u, ok := um.users[id]
	return u, ok
}

// GetByUsername returns user
func (um *UserManager) GetByUsername(username string) (*User, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	for _, u := range um.users {
		if u.Username == username {
			return u, true
		}
	}
	return nil, false
}

// Update updates user
func (um *UserManager) Update(u User) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	existing, ok := um.users[u.ID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	existing.Username = u.Username
	existing.Email = u.Email
	existing.Role = u.Role
	existing.Active = u.Active
	existing.UpdatedAt = time.Now()
	return nil
}

// Delete deletes user
func (um *UserManager) Delete(id string) bool {
	um.mu.Lock()
	defer um.mu.Unlock()
	if _, ok := um.users[id]; !ok {
		return false
	}
	delete(um.users, id)
	return true
}

// List lists users
func (um *UserManager) List() []*User {
	um.mu.RLock()
	defer um.mu.RUnlock()
	out := []*User{}
	for _, u := range um.users {
		out = append(out, u)
	}
	return out
}

// Count returns count
func (um *UserManager) Count() int {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return len(um.users)
}

// ActiveCount returns active count
func (um *UserManager) ActiveCount() int {
	um.mu.RLock()
	defer um.mu.RUnlock()
	n := 0
	for _, u := range um.users {
		if u.Active {
			n++
		}
	}
	return n
}
