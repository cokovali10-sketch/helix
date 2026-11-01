// Package api provides logger
package api

import (
	"fmt"
	"sync"
	"time"
)

// APILogger logs API
type APILogger struct {
	messages []string
	mu       sync.Mutex
}

// NewAPILogger creates logger
func NewAPILogger() *APILogger {
	return &APILogger{messages: []string{}}
}

// Log logs message
func (l *APILogger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), msg))
}

// Messages returns messages
func (l *APILogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.messages))
	copy(out, l.messages)
	return out
}

// Clear clears
func (l *APILogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = []string{}
}

// Count returns count
func (l *APILogger) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}
