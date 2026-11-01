// Package ingestion provides queue
package ingestion

import (
	"sync"
)

// Queue is a FIFO queue
type Queue struct {
	items []map[string]interface{}
	mu    sync.Mutex
}

// NewQueue creates queue
func NewQueue() *Queue {
	return &Queue{items: []map[string]interface{}{}}
}

// Enqueue enqueues item
func (q *Queue) Enqueue(item map[string]interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue dequeues item
func (q *Queue) Dequeue() (map[string]interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek peeks
func (q *Queue) Peek() (map[string]interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, false
	}
	return q.items[0], true
}

// Size returns size
func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Clear clears
func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = []map[string]interface{}{}
}

// IsEmpty returns empty
func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}
