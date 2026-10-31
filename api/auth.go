// Package api provides authentication and authorization middleware
package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	secretKey  string
	tokenExpiry time.Duration
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(secretKey string, tokenExpiry time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

// Authenticate authenticates a request
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")

		// Validate token
		claims, err := am.validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateToken validates a JWT-like token
func (am *AuthMiddleware) validateToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Verify signature
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha256.New, []byte(am.secretKey))
	h.Write([]byte(parts[0] + "." + parts[1]))
	expectedSignature := h.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}

// GenerateToken generates a new token
func (am *AuthMiddleware) GenerateToken(userID, tenantID string, roles []string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	payload := map[string]interface{}{
		"sub":    userID,
		"tenant": tenantID,
		"roles":  roles,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(am.tokenExpiry).Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	h := hmac.New(sha256.New, []byte(am.secretKey))
	h.Write([]byte(headerEncoded + "." + payloadEncoded))
	signature := h.Sum(nil)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	return headerEncoded + "." + payloadEncoded + "." + signatureEncoded, nil
}

// RBACMiddleware handles role-based access control
type RBACMiddleware struct {
	rbac *RBAC
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(rbac *RBAC) *RBACMiddleware {
	return &RBACMiddleware{rbac: rbac}
}

// RequirePermission checks if user has required permission
func (rm *RBACMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(map[string]interface{})
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid user", http.StatusUnauthorized)
				return
			}

			permission := fmt.Sprintf("%s:%s", resource, action)
			if !rm.rbac.HasPermission(userID, permission) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole checks if user has required role
func (rm *RBACMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(map[string]interface{})
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			roles, ok := claims["roles"].([]interface{})
			if !ok {
				http.Error(w, "Invalid roles", http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, r := range roles {
				if r.(string) == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireTenant checks if user belongs to tenant
func (rm *RBACMiddleware) RequireTenant(tenantID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(map[string]interface{})
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userTenant, ok := claims["tenant"].(string)
			if !ok || userTenant != tenantID {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements rate limiting
type RateLimitMiddleware struct {
	limiter *RateLimiter
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

// TokenBucket implements token bucket
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: &RateLimiter{
			buckets: make(map[string]*TokenBucket),
		},
	}
}

// Limit applies rate limiting
func (rlm *RateLimitMiddleware) Limit(requestsPerSecond int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client identifier (IP or user ID)
			clientID := r.RemoteAddr
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				clientID = userID
			}

			bucket := rlm.limiter.getBucket(clientID, requestsPerSecond)
			if !bucket.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) getBucket(clientID string, capacity int) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, ok := rl.buckets[clientID]; ok {
		return bucket
	}

	bucket := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: capacity,
		lastRefill: time.Now(),
	}
	rl.buckets[clientID] = bucket
	return bucket
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// TenantMiddleware handles tenant isolation
type TenantMiddleware struct{}

// NewTenantMiddleware creates a new tenant middleware
func NewTenantMiddleware() *TenantMiddleware {
	return &TenantMiddleware{}
}

// ExtractTenant extracts tenant from request
func (tm *TenantMiddleware) ExtractTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = r.URL.Query().Get("tenant")
		}
		if tenantID == "" {
			tenantID = "default"
		}

		ctx := context.WithValue(r.Context(), "tenant", tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTenant gets tenant from context
func GetTenant(ctx context.Context) string {
	if tenant, ok := ctx.Value("tenant").(string); ok {
		return tenant
	}
	return "default"
}

// GetClaims gets claims from context
func GetClaims(ctx context.Context) map[string]interface{} {
	if claims, ok := ctx.Value("claims").(map[string]interface{}); ok {
		return claims
	}
	return nil
}

// GetUser gets user ID from context
func GetUser(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims != nil {
		if userID, ok := claims["sub"].(string); ok {
			return userID
		}
	}
	return ""
}

// LoggingMiddleware logs requests
type LoggingMiddleware struct {
	logger func(string)
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger func(string)) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// Log logs requests
func (lm *LoggingMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start)
		lm.logger(fmt.Sprintf("%s %s %d %v %s", r.Method, r.URL.Path, wrapped.statusCode, duration, r.RemoteAddr))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RecoveryMiddleware recovers from panics
type RecoveryMiddleware struct{}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

// Recover recovers from panics
func (rm *RecoveryMiddleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				fmt.Printf("Panic recovered: %v\n", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CompressionMiddleware compresses responses
type CompressionMiddleware struct{}

// NewCompressionMiddleware creates a new compression middleware
func NewCompressionMiddleware() *CompressionMiddleware {
	return &CompressionMiddleware{}
}

// Compress compresses responses
func (cm *CompressionMiddleware) Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		
		// In production, would wrap with gzip writer
		next.ServeHTTP(w, r)
	})
}

// CacheMiddleware adds caching headers
type CacheMiddleware struct {
	maxAge int
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(maxAge int) *CacheMiddleware {
	return &CacheMiddleware{maxAge: maxAge}
}

// Cache adds cache headers
func (cm *CacheMiddleware) Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cm.maxAge))
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
type SecurityHeadersMiddleware struct{}

// NewSecurityHeadersMiddleware creates a new security headers middleware
func NewSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

// Secure adds security headers
func (shm *SecurityHeadersMiddleware) Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// RBAC is a placeholder for the RBAC system
type RBAC struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	users       map[string]*User
	userRoles   map[string][]string
	mu          sync.RWMutex
}

type Role struct {
	ID          string
	Name        string
	Permissions []string
}

type Permission struct {
	ID       string
	Resource string
	Action   string
}

type User struct {
	ID       string
	Username string
	Roles    []string
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}