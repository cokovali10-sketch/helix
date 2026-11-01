// Package storage provides logger
package storage

import (
	"fmt"
	"sync"
	"time"
)

// StorageLogger logs storage
type StorageLogger struct {
	messages []string
	mu       sync.Mutex
}

// NewStorageLogger creates logger
func NewStorageLogger() *StorageLogger {
	return &StorageLogger{messages: []string{}}
}

// Log logs message
func (l *StorageLogger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), msg))
}

// Messages returns messages
func (l *StorageLogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.messages))
	copy(out, l.messages)
	return out
}

// Clear clears
func (l *StorageLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = []string{}
}

// Count returns count
func (l *StorageLogger) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}
