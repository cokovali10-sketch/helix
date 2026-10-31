// Package catalog implements the catalog and metadata system for the Helix platform
// Includes metadata storage, RBAC, and API definitions
package catalog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Metadata represents metadata for an entity
type Metadata struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"` // metric, dashboard, alert, tenant
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	CreatedBy   string            `json:"createdBy"`
	UpdatedBy   string            `json:"updatedBy"`
	Version     int               `json:"version"`
}

// MetricMetadata represents metadata for a metric
type MetricMetadata struct {
	Metadata
	MetricName   string            `json:"metricName"`
	Type         string            `json:"type"` // counter, gauge, histogram, summary
	Unit         string            `json:"unit"`
	Help         string            `json:"help"`
	Labels       []string          `json:"labels"`
	Cardinality  int               `json:"cardinality"`
	LastSeen     time.Time         `json:"lastSeen"`
	SampleCount  int64             `json:"sampleCount"`
}

// DashboardMetadata represents metadata for a dashboard
type DashboardMetadata struct {
	Metadata
	Title       string   `json:"title"`
	Panels      int      `json:"panels"`
	IsPublic    bool     `json:"isPublic"`
	Owner       string   `json:"owner"`
	Tags        []string `json:"tags"`
	LastViewed  time.Time `json:"lastViewed"`
	ViewCount   int      `json:"viewCount"`
}

// AlertMetadata represents metadata for an alert
type AlertMetadata struct {
	Metadata
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	RuleID      string `json:"ruleId"`
	FiringCount int    `json:"firingCount"`
	LastFired   time.Time `json:"lastFired"`
}

// TenantMetadata represents metadata for a tenant
type TenantMetadata struct {
	Metadata
	Quota       int64  `json:"quota"`
	RateLimit   int    `json:"rateLimit"`
	UserCount   int    `json:"userCount"`
	MetricCount int    `json:"metricCount"`
	StorageUsed int64  `json:"storageUsed"`
	Active      bool   `json:"active"`
}

// Catalog manages all metadata
type Catalog struct {
	metrics    map[string]*MetricMetadata
	dashboards map[string]*DashboardMetadata
	alerts     map[string]*AlertMetadata
	tenants    map[string]*TenantMetadata
	generic    map[string]*Metadata
	mu         sync.RWMutex
	store      MetadataStore
}

// MetadataStore interface for persistence
type MetadataStore interface {
	Get(ctx context.Context, id string) (*Metadata, error)
	Put(ctx context.Context, metadata *Metadata) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *MetadataFilter) ([]*Metadata, error)
	Search(ctx context.Context, query string) ([]*Metadata, error)
}

// MetadataFilter for querying metadata
type MetadataFilter struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
	Owner  string            `json:"owner"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

// NewCatalog creates a new catalog
func NewCatalog(store MetadataStore) *Catalog {
	return &Catalog{
		metrics:    make(map[string]*MetricMetadata),
		dashboards: make(map[string]*DashboardMetadata),
		alerts:     make(map[string]*AlertMetadata),
		tenants:    make(map[string]*TenantMetadata),
		generic:    make(map[string]*Metadata),
		store:      store,
	}
}

// RegisterMetric registers a metric
func (c *Catalog) RegisterMetric(metric *MetricMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics[metric.ID] = metric
}

// GetMetric gets a metric by ID
func (c *Catalog) GetMetric(id string) (*MetricMetadata, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	metric, ok := c.metrics[id]
	return metric, ok
}

// ListMetrics lists all metrics
func (c *Catalog) ListMetrics(filter *MetadataFilter) []*MetricMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var result []*MetricMetadata
	for _, metric := range c.metrics {
		if filter != nil && filter.Type != "" && metric.Type != filter.Type {
			continue
		}
		result = append(result, metric)
		if filter != nil && filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}
	return result
}

// RegisterDashboard registers a dashboard
func (c *Catalog) RegisterDashboard(dashboard *DashboardMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dashboards[dashboard.ID] = dashboard
}

// GetDashboard gets a dashboard by ID
func (c *Catalog) GetDashboard(id string) (*DashboardMetadata, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dashboard, ok := c.dashboards[id]
	return dashboard, ok
}

// ListDashboards lists all dashboards
func (c *Catalog) ListDashboards(filter *MetadataFilter) []*DashboardMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var result []*DashboardMetadata
	for _, dashboard := range c.dashboards {
		if filter != nil && filter.Owner != "" && dashboard.Owner != filter.Owner {
			continue
		}
		result = append(result, dashboard)
		if filter != nil && filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}
	return result
}

// RegisterAlert registers an alert
func (c *Catalog) RegisterAlert(alert *AlertMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alerts[alert.ID] = alert
}

// GetAlert gets an alert by ID
func (c *Catalog) GetAlert(id string) (*AlertMetadata, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	alert, ok := c.alerts[id]
	return alert, ok
}

// ListAlerts lists all alerts
func (c *Catalog) ListAlerts(filter *MetadataFilter) []*AlertMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var result []*AlertMetadata
	for _, alert := range c.alerts {
		if filter != nil && filter.Type != "" && alert.Status != filter.Type {
			continue
		}
		result = append(result, alert)
		if filter != nil && filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}
	return result
}

// RegisterTenant registers a tenant
func (c *Catalog) RegisterTenant(tenant *TenantMetadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tenants[tenant.ID] = tenant
}

// GetTenant gets a tenant by ID
func (c *Catalog) GetTenant(id string) (*TenantMetadata, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tenant, ok := c.tenants[id]
	return tenant, ok
}

// ListTenants lists all tenants
func (c *Catalog) ListTenants() []*TenantMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	tenants := make([]*TenantMetadata, 0, len(c.tenants))
	for _, tenant := range c.tenants {
		tenants = append(tenants, tenant)
	}
	return tenants
}

// Search searches metadata
func (c *Catalog) Search(query string) []*Metadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var result []*Metadata
	for _, metric := range c.metrics {
		if contains(metric.Name, query) || contains(metric.Description, query) {
			result = append(result, &metric.Metadata)
		}
	}
	for _, dashboard := range c.dashboards {
		if contains(dashboard.Name, query) || contains(dashboard.Description, query) {
			result = append(result, &dashboard.Metadata)
		}
	}
	return result
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RBAC System
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Permission struct {
	ID          string `json:"id"`
	Resource    string `json:"resource"` // metrics, dashboards, alerts, tenants, users
	Action      string `json:"action"`   // read, write, delete, admin
	Description string `json:"description"`
}

type User struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TenantID  string   `json:"tenantId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Active    bool     `json:"active"`
}

type RBAC struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	users       map[string]*User
	rolePerms   map[string][]string // role ID -> permission IDs
	userRoles   map[string][]string // user ID -> role IDs
	mu          sync.RWMutex
}

func NewRBAC() *RBAC {
	rbac := &RBAC{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
		users:       make(map[string]*User),
		rolePerms:   make(map[string][]string),
		userRoles:   make(map[string][]string),
	}
	rbac.initDefaultRoles()
	rbac.initDefaultPermissions()
	return rbac
}

func (r *RBAC) initDefaultRoles() {
	roles := []*Role{
		{ID: "admin", Name: "Administrator", Description: "Full access", Permissions: []string{"*"}, CreatedAt: time.Now()},
		{ID: "editor", Name: "Editor", Description: "Read and write access", Permissions: []string{"metrics:read", "metrics:write", "dashboards:read", "dashboards:write", "alerts:read", "alerts:write"}, CreatedAt: time.Now()},
		{ID: "viewer", Name: "Viewer", Description: "Read-only access", Permissions: []string{"metrics:read", "dashboards:read", "alerts:read"}, CreatedAt: time.Now()},
		{ID: "operator", Name: "Operator", Description: "Operational access", Permissions: []string{"metrics:read", "alerts:read", "alerts:write"}, CreatedAt: time.Now()},
	}
	for _, role := range roles {
		r.roles[role.ID] = role
	}
}

func (r *RBAC) initDefaultPermissions() {
	permissions := []*Permission{
		{ID: "metrics:read", Resource: "metrics", Action: "read", Description: "Read metrics"},
		{ID: "metrics:write", Resource: "metrics", Action: "write", Description: "Write metrics"},
		{ID: "metrics:delete", Resource: "metrics", Action: "delete", Description: "Delete metrics"},
		{ID: "dashboards:read", Resource: "dashboards", Action: "read", Description: "Read dashboards"},
		{ID: "dashboards:write", Resource: "dashboards", Action: "write", Description: "Write dashboards"},
		{ID: "dashboards:delete", Resource: "dashboards", Action: "delete", Description: "Delete dashboards"},
		{ID: "alerts:read", Resource: "alerts", Action: "read", Description: "Read alerts"},
		{ID: "alerts:write", Resource: "alerts", Action: "write", Description: "Write alerts"},
		{ID: "alerts:delete", Resource: "alerts", Action: "delete", Description: "Delete alerts"},
		{ID: "tenants:read", Resource: "tenants", Action: "read", Description: "Read tenants"},
		{ID: "tenants:write", Resource: "tenants", Action: "write", Description: "Write tenants"},
		{ID: "tenants:delete", Resource: "tenants", Action: "delete", Description: "Delete tenants"},
		{ID: "users:read", Resource: "users", Action: "read", Description: "Read users"},
		{ID: "users:write", Resource: "users", Action: "write", Description: "Write users"},
		{ID: "users:delete", Resource: "users", Action: "delete", Description: "Delete users"},
		{ID: "admin", Resource: "*", Action: "admin", Description: "Full admin access"},
	}
	for _, perm := range permissions {
		r.permissions[perm.ID] = perm
	}
}

// CreateRole creates a new role
func (r *RBAC) CreateRole(role *Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.roles[role.ID]; exists {
		return fmt.Errorf("role %s already exists", role.ID)
	}
	r.roles[role.ID] = role
	return nil
}

// GetRole gets a role by ID
func (r *RBAC) GetRole(id string) (*Role, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	role, ok := r.roles[id]
	return role, ok
}

// ListRoles lists all roles
func (r *RBAC) ListRoles() []*Role {
	r.mu.RLock()
	defer r.mu.RUnlock()
	roles := make([]*Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles
}

// DeleteRole deletes a role
func (r *RBAC) DeleteRole(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.roles[id]; !exists {
		return fmt.Errorf("role %s not found", id)
	}
	delete(r.roles, id)
	return nil
}

// AssignRole assigns a role to a user
func (r *RBAC) AssignRole(userID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.users[userID]; !exists {
		return fmt.Errorf("user %s not found", userID)
	}
	if _, exists := r.roles[roleID]; !exists {
		return fmt.Errorf("role %s not found", roleID)
	}
	
	r.userRoles[userID] = append(r.userRoles[userID], roleID)
	return nil
}

// RevokeRole revokes a role from a user
func (r *RBAC) RevokeRole(userID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	roles := r.userRoles[userID]
	for i, rid := range roles {
		if rid == roleID {
			r.userRoles[userID] = append(roles[:i], roles[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("role %s not assigned to user %s", roleID, userID)
}

// HasPermission checks if a user has a permission
func (r *RBAC) HasPermission(userID, permission string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	roles := r.userRoles[userID]
	for _, roleID := range roles {
		role, ok := r.roles[roleID]
		if !ok {
			continue
		}
		for _, perm := range role.Permissions {
			if perm == "*" || perm == permission {
				return true
			}
		}
	}
	return false
}

// CreateUser creates a new user
func (r *RBAC) CreateUser(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.users[user.ID] = user
	r.userRoles[user.ID] = make([]string, 0)
	return nil
}

// GetUser gets a user by ID
func (r *RBAC) GetUser(id string) (*User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	return user, ok
}

// ListUsers lists all users
func (r *RBAC) ListUsers() []*User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// DeleteUser deletes a user
func (r *RBAC) DeleteUser(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.users, id)
	delete(r.userRoles, id)
	return nil
}

// GetRolePermissions gets permissions for a role
func (r *RBAC) GetRolePermissions(roleID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rolePerms[roleID]
}

// CheckAccess checks if a user can access a resource
func (r *RBAC) CheckAccess(userID, resource, action string) bool {
	permission := fmt.Sprintf("%s:%s", resource, action)
	return r.HasPermission(userID, permission)
}

// InMemoryMetadataStore implements MetadataStore
type InMemoryMetadataStore struct {
	data map[string]*Metadata
	mu   sync.RWMutex
}

func NewInMemoryMetadataStore() *InMemoryMetadataStore {
	return &InMemoryMetadataStore{
		data: make(map[string]*Metadata),
	}
}

func (s *InMemoryMetadataStore) Get(ctx context.Context, id string) (*Metadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if metadata, ok := s.data[id]; ok {
		return metadata, nil
	}
	return nil, fmt.Errorf("metadata %s not found", id)
}

func (s *InMemoryMetadataStore) Put(ctx context.Context, metadata *Metadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[metadata.ID] = metadata
	return nil
}

func (s *InMemoryMetadataStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
	return nil
}

func (s *InMemoryMetadataStore) List(ctx context.Context, filter *MetadataFilter) ([]*Metadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []*Metadata
	for _, metadata := range s.data {
		if filter != nil && filter.Type != "" && metadata.Type != filter.Type {
			continue
		}
		result = append(result, metadata)
		if filter != nil && filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}
	return result, nil
}

func (s *InMemoryMetadataStore) Search(ctx context.Context, query string) ([]*Metadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []*Metadata
	for _, metadata := range s.data {
		if contains(metadata.Name, query) || contains(metadata.Description, query) {
			result = append(result, metadata)
		}
	}
	return result, nil
}