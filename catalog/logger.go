// Package catalog provides logger
package catalog

import (
	"fmt"
	"sync"
	"time"
)

// CatalogLogger logs catalog
type CatalogLogger struct {
	messages []string
	mu       sync.Mutex
}

// NewCatalogLogger creates logger
func NewCatalogLogger() *CatalogLogger {
	return &CatalogLogger{messages: []string{}}
}

// Log logs message
func (l *CatalogLogger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), msg))
}

// Messages returns messages
func (l *CatalogLogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.messages))
	copy(out, l.messages)
	return out
}

// Clear clears
func (l *CatalogLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = []string{}
}

// Count returns count
func (l *CatalogLogger) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}
