// Package api provides token utilities
package api

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Token represents a token
type Token struct {
	Value     string
	ExpiresAt time.Time
}

// GenerateToken generates token
func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// NewToken creates token
func NewToken(ttl time.Duration) Token {
	return Token{Value: GenerateToken(), ExpiresAt: time.Now().Add(ttl)}
}

// IsExpired checks expired
func (t Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
