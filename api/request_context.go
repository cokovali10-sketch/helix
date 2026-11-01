// Package api provides request context
package api

import (
	"context"
	"net/http"
)

// ContextKey type
type ContextKey string

const (
	TenantKey ContextKey = "tenant"
	UserKey   ContextKey = "user"
)

// WithTenant adds tenant to context
func WithTenant(r *http.Request, tenant string) *http.Request {
	ctx := context.WithValue(r.Context(), TenantKey, tenant)
	return r.WithContext(ctx)
}

// TenantFromRequest gets tenant
func TenantFromRequest(r *http.Request) string {
	if v := r.Context().Value(TenantKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return "default"
}

// WithUser adds user to context
func WithUser(r *http.Request, user string) *http.Request {
	ctx := context.WithValue(r.Context(), UserKey, user)
	return r.WithContext(ctx)
}

// UserFromRequest gets user
func UserFromRequest(r *http.Request) string {
	if v := r.Context().Value(UserKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
