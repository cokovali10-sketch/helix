// Package api provides additional API handlers
// including webhook, notification, and integration handlers
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// WebhookHandler handles webhook endpoints
type WebhookHandler struct {
	handlers map[string]WebhookFunc
	mu       sync.RWMutex
	logs     []WebhookLog
}

// WebhookFunc is a webhook handler function
type WebhookFunc func(payload []byte) error

// WebhookLog represents a webhook log entry
type WebhookLog struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Payload   string    `json:"payload"`
	Status    int       `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		handlers: make(map[string]WebhookFunc),
		logs:     make([]WebhookLog, 0),
	}
}

// Register registers a webhook handler
func (wh *WebhookHandler) Register(path string, handler WebhookFunc) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	wh.handlers[path] = handler
}

// Unregister unregisters a webhook handler
func (wh *WebhookHandler) Unregister(path string) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	delete(wh.handlers, path)
}

// ServeHTTP serves webhook endpoint
func (wh *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	wh.mu.RLock()
	handler, ok := wh.handlers[path]
	wh.mu.RUnlock()

	if !ok {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	start := time.Now()
	body, err := readBody(r)
	if err != nil {
		wh.logWebhook(path, string(body), http.StatusBadRequest, err.Error(), time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler(body); err != nil {
		wh.logWebhook(path, string(body), http.StatusInternalServerError, err.Error(), time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wh.logWebhook(path, string(body), http.StatusOK, "", time.Since(start))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (wh *WebhookHandler) logWebhook(path, payload string, status int, err string, duration time.Duration) {
	wh.mu.Lock()
	defer wh.mu.Unlock()

	log := WebhookLog{
		ID:        fmt.Sprintf("wh-%d", time.Now().UnixNano()),
		Path:      path,
		Payload:   payload,
		Status:    status,
		Error:     err,
		Timestamp: time.Now(),
		Duration:  duration,
	}
	wh.logs = append(wh.logs, log)

	// Keep only last 1000 logs
	if len(wh.logs) > 1000 {
		wh.logs = wh.logs[1:]
	}
}

// GetLogs returns webhook logs
func (wh *WebhookHandler) GetLogs(path string, limit int) []WebhookLog {
	wh.mu.RLock()
	defer wh.mu.RUnlock()

	var result []WebhookLog
	for _, log := range wh.logs {
		if path != "" && log.Path != path {
			continue
		}
		result = append(result, log)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

func readBody(r *http.Request) ([]byte, error) {
	// Read body
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)
	for {
		n, err := r.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
	}
	return buf, nil
}

// NotificationService handles notifications
type NotificationService struct {
	notifications []Notification
	subscribers   map[string]chan Notification
	mu            sync.RWMutex
}

// Notification represents a notification
type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
	UserID    string    `json:"userId"`
	TenantID  string    `json:"tenantId"`
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: make([]Notification, 0),
		subscribers:   make(map[string]chan Notification),
	}
}

// Subscribe subscribes to notifications
func (ns *NotificationService) Subscribe(userID string) <-chan Notification {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ch := make(chan Notification, 100)
	ns.subscribers[userID] = ch
	return ch
}

// Unsubscribe unsubscribes from notifications
func (ns *NotificationService) Unsubscribe(userID string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if ch, ok := ns.subscribers[userID]; ok {
		close(ch)
		delete(ns.subscribers, userID)
	}
}

// Send sends a notification
func (ns *NotificationService) Send(notification Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	notification.ID = fmt.Sprintf("notif-%d", time.Now().UnixNano())
	notification.Timestamp = time.Now()
	ns.notifications = append(ns.notifications, notification)

	// Send to subscriber
	if ch, ok := ns.subscribers[notification.UserID]; ok {
		select {
		case ch <- notification:
		default:
		}
	}
}

// SendToTenant sends a notification to all users in a tenant
func (ns *NotificationService) SendToTenant(tenantID string, notification Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	notification.ID = fmt.Sprintf("notif-%d", time.Now().UnixNano())
	notification.Timestamp = time.Now()
	notification.TenantID = tenantID
	ns.notifications = append(ns.notifications, notification)

	// Send to all subscribers in tenant
	for userID, ch := range ns.subscribers {
		_ = userID
		select {
		case ch <- notification:
		default:
		}
	}
}

// List lists notifications
func (ns *NotificationService) List(userID string, unreadOnly bool) []Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var result []Notification
	for _, n := range ns.notifications {
		if userID != "" && n.UserID != userID {
			continue
		}
		if unreadOnly && n.Read {
			continue
		}
		result = append(result, n)
	}
	return result
}

// MarkRead marks a notification as read
func (ns *NotificationService) MarkRead(id string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	for i, n := range ns.notifications {
		if n.ID == id {
			ns.notifications[i].Read = true
			break
		}
	}
}

// MarkAllRead marks all notifications as read
func (ns *NotificationService) MarkAllRead(userID string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	for i, n := range ns.notifications {
		if userID == "" || n.UserID == userID {
			ns.notifications[i].Read = true
		}
	}
}

// Delete deletes a notification
func (ns *NotificationService) Delete(id string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	for i, n := range ns.notifications {
		if n.ID == id {
			ns.notifications = append(ns.notifications[:i], ns.notifications[i+1:]...)
			break
		}
	}
}

// Clear clears all notifications
func (ns *NotificationService) Clear(userID string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if userID == "" {
		ns.notifications = make([]Notification, 0)
	} else {
		var result []Notification
		for _, n := range ns.notifications {
			if n.UserID != userID {
				result = append(result, n)
			}
		}
		ns.notifications = result
	}
}

// GetUnreadCount returns unread count
func (ns *NotificationService) GetUnreadCount(userID string) int {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	count := 0
	for _, n := range ns.notifications {
		if (userID == "" || n.UserID == userID) && !n.Read {
			count++
		}
	}
	return count
}

// IntegrationManager manages integrations
type IntegrationManager struct {
	integrations map[string]*Integration
	mu           sync.RWMutex
}

// Integration represents an integration
type Integration struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Config      map[string]string `json:"config"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	LastSync    time.Time         `json:"lastSync"`
	Status      string            `json:"status"`
}

// NewIntegrationManager creates a new integration manager
func NewIntegrationManager() *IntegrationManager {
	return &IntegrationManager{
		integrations: make(map[string]*Integration),
	}
}

// Register registers an integration
func (im *IntegrationManager) Register(integration *Integration) {
	im.mu.Lock()
	defer im.mu.Unlock()

	integration.ID = fmt.Sprintf("int-%d", time.Now().UnixNano())
	integration.CreatedAt = time.Now()
	integration.UpdatedAt = time.Now()
	integration.Status = "connected"
	im.integrations[integration.ID] = integration
}

// Get gets an integration
func (im *IntegrationManager) Get(id string) (*Integration, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	integration, ok := im.integrations[id]
	return integration, ok
}

// List lists all integrations
func (im *IntegrationManager) List() []*Integration {
	im.mu.RLock()
	defer im.mu.RUnlock()

	integrations := make([]*Integration, 0, len(im.integrations))
	for _, i := range im.integrations {
		integrations = append(integrations, i)
	}
	return integrations
}

// Update updates an integration
func (im *IntegrationManager) Update(id string, updates map[string]interface{}) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	integration, ok := im.integrations[id]
	if !ok {
		return fmt.Errorf("integration %s not found", id)
	}

	if name, ok := updates["name"].(string); ok {
		integration.Name = name
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		integration.Enabled = enabled
	}
	if config, ok := updates["config"].(map[string]string); ok {
		integration.Config = config
	}

	integration.UpdatedAt = time.Now()
	return nil
}

// Delete deletes an integration
func (im *IntegrationManager) Delete(id string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, ok := im.integrations[id]; !ok {
		return fmt.Errorf("integration %s not found", id)
	}

	delete(im.integrations, id)
	return nil
}

// Test tests an integration
func (im *IntegrationManager) Test(id string) error {
	im.mu.RLock()
	defer im.mu.RUnlock()

	integration, ok := im.integrations[id]
	if !ok {
		return fmt.Errorf("integration %s not found", id)
	}

	// In production, would test the integration
	_ = integration
	return nil
}

// Sync syncs an integration
func (im *IntegrationManager) Sync(id string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	integration, ok := im.integrations[id]
	if !ok {
		return fmt.Errorf("integration %s not found", id)
	}

	integration.LastSync = time.Now()
	integration.Status = "connected"
	return nil
}