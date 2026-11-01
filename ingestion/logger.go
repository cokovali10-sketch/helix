// Package ingestion provides logger
package ingestion

import (
	"fmt"
	"sync"
	"time"
)

// IngestLogger logs ingestion
type IngestLogger struct {
	messages []string
	mu       sync.Mutex
}

// NewIngestLogger creates logger
func NewIngestLogger() *IngestLogger {
	return &IngestLogger{messages: []string{}}
}

// Log logs message
func (l *IngestLogger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), msg))
}

// Messages returns messages
func (l *IngestLogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.messages))
	copy(out, l.messages)
	return out
}

// Clear clears
func (l *IngestLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = []string{}
}

// Count returns count
func (l *IngestLogger) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}
