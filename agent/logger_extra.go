// Package agent provides logger extra
package agent

import (
	"fmt"
	"sync"
	"time"
)

// LoggerExtra logs messages
type LoggerExtra struct {
	messages []string
	mu       sync.Mutex
}

// NewLoggerExtra creates logger
func NewLoggerExtra() *LoggerExtra {
	return &LoggerExtra{messages: []string{}}
}

// Log logs message
func (l *LoggerExtra) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), msg))
}

// Messages returns messages
func (l *LoggerExtra) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.messages))
	copy(out, l.messages)
	return out
}

// Clear clears messages
func (l *LoggerExtra) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = []string{}
}

// Count returns count
func (l *LoggerExtra) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}
