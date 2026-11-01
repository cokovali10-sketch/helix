// Package api provides additional API handlers
// including authentication, authorization, and session management
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	users    map[string]*User
	sessions map[string]*Session
	mu       sync.RWMutex
}

// User represents a user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Roles     []string  `json:"roles"`
	TenantID  string    `json:"tenantId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Active    bool      `json:"active"`
}

// Session represents a session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		users:    make(map[string]*User),
		sessions: make(map[string]*Session),
	}
}

// Register registers a new user
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ah.mu.Lock()
	defer ah.mu.Unlock()

	// Check if user exists
	for _, u := range ah.users {
		if u.Username == user.Username || u.Email == user.Email {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
	}

	user.ID = fmt.Sprintf("user-%d", time.Now().UnixNano())
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Active = true
	ah.users[user.ID] = &user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login logs in a user
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ah.mu.RLock()
	var foundUser *User
	for _, u := range ah.users {
		if u.Username == creds.Username {
			foundUser = u
			break
		}
	}
	ah.mu.RUnlock()

	if foundUser == nil || foundUser.Password != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session := &Session{
		ID:        fmt.Sprintf("session-%d", time.Now().UnixNano()),
		UserID:    foundUser.ID,
		Token:     fmt.Sprintf("token-%d", time.Now().UnixNano()),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
	}

	ah.mu.Lock()
	ah.sessions[session.ID] = session
	ah.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   session.Token,
		"expires": session.ExpiresAt,
		"user": map[string]interface{}{
			"id":       foundUser.ID,
			"username": foundUser.Username,
			"email":    foundUser.Email,
			"roles":    foundUser.Roles,
			"tenantId": foundUser.TenantID,
		},
	})
}

// Logout logs out a user
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	ah.mu.Lock()
	defer ah.mu.Unlock()

	for id, session := range ah.sessions {
		if session.Token == token {
			delete(ah.sessions, id)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ValidateToken validates a token
func (ah *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	ah.mu.RLock()
	defer ah.mu.RUnlock()

	for _, session := range ah.sessions {
		if session.Token == token {
			if time.Now().After(session.ExpiresAt) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"valid":   true,
				"userId":  session.UserID,
				"expires": session.ExpiresAt,
			})
			return
		}
	}

	http.Error(w, "Invalid token", http.StatusUnauthorized)
}

// RefreshToken refreshes a token
func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	ah.mu.Lock()
	defer ah.mu.Unlock()

	for id, session := range ah.sessions {
		if session.Token == token {
			session.Token = fmt.Sprintf("token-%d", time.Now().UnixNano())
			session.ExpiresAt = time.Now().Add(24 * time.Hour)
			ah.sessions[id] = session

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"token":   session.Token,
				"expires": session.ExpiresAt,
			})
			return
		}
	}

	http.Error(w, "Invalid token", http.StatusUnauthorized)
}

// GetUser gets the current user
func (ah *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	ah.mu.RLock()
	defer ah.mu.RUnlock()

	for _, session := range ah.sessions {
		if session.Token == token {
			for _, user := range ah.users {
				if user.ID == session.UserID {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(user)
					return
				}
			}
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// UpdateUser updates a user
func (ah *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ah.mu.Lock()
	defer ah.mu.Unlock()

	for _, session := range ah.sessions {
		if session.Token == token {
			for _, user := range ah.users {
				if user.ID == session.UserID {
					if name, ok := updates["name"].(string); ok {
						user.Username = name
					}
					if email, ok := updates["email"].(string); ok {
						user.Email = email
					}
					user.UpdatedAt = time.Now()

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(user)
					return
				}
			}
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// ChangePassword changes a user's password
func (ah *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var body struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ah.mu.Lock()
	defer ah.mu.Unlock()

	for _, session := range ah.sessions {
		if session.Token == token {
			for _, user := range ah.users {
				if user.ID == session.UserID {
					if user.Password != body.OldPassword {
						http.Error(w, "Invalid old password", http.StatusBadRequest)
						return
					}
					user.Password = body.NewPassword
					user.UpdatedAt = time.Now()

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
					return
				}
			}
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// ServeHTTP serves auth endpoints
func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/v1/auth/register":
		ah.Register(w, r)
	case path == "/api/v1/auth/login":
		ah.Login(w, r)
	case path == "/api/v1/auth/logout":
		ah.Logout(w, r)
	case path == "/api/v1/auth/validate":
		ah.ValidateToken(w, r)
	case path == "/api/v1/auth/refresh":
		ah.RefreshToken(w, r)
	case path == "/api/v1/auth/user":
		if r.Method == http.MethodGet {
			ah.GetUser(w, r)
		} else {
			ah.UpdateUser(w, r)
		}
	case path == "/api/v1/auth/password":
		ah.ChangePassword(w, r)
	default:
		http.NotFound(w, r)
	}
}
