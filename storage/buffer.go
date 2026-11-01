// Package storage provides additional storage features
package storage

import (
	"fmt"
	"sync"
	"time"
)

// WriteBuffer buffers writes for efficient storage
type WriteBuffer struct {
	buffer   []byte
	capacity int
	mu       sync.Mutex
	flushFn  func([]byte) error
}

// NewWriteBuffer creates a new write buffer
func NewWriteBuffer(capacity int, flushFn func([]byte) error) *WriteBuffer {
	return &WriteBuffer{
		buffer:   make([]byte, 0, capacity),
		capacity: capacity,
		flushFn:  flushFn,
	}
}

// Write writes data to the buffer
func (wb *WriteBuffer) Write(data []byte) error {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	if len(wb.buffer)+len(data) > wb.capacity {
		if err := wb.flush(); err != nil {
			return err
		}
	}

	wb.buffer = append(wb.buffer, data...)
	return nil
}

// Flush flushes the buffer
func (wb *WriteBuffer) Flush() error {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	return wb.flush()
}

func (wb *WriteBuffer) flush() error {
	if len(wb.buffer) == 0 {
		return nil
	}

	if err := wb.flushFn(wb.buffer); err != nil {
		return err
	}

	wb.buffer = wb.buffer[:0]
	return nil
}

// Size returns the buffer size
func (wb *WriteBuffer) Size() int {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	return len(wb.buffer)
}

// Capacity returns the buffer capacity
func (wb *WriteBuffer) Capacity() int {
	return wb.capacity
}

// Utilization returns the buffer utilization
func (wb *WriteBuffer) Utilization() float64 {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	if wb.capacity == 0 {
		return 0
	}
	return float64(len(wb.buffer)) / float64(wb.capacity)
}

// ReadBuffer buffers reads for efficient storage
type ReadBuffer struct {
	buffer   []byte
	position int
	mu       sync.Mutex
}

// NewReadBuffer creates a new read buffer
func NewReadBuffer(data []byte) *ReadBuffer {
	return &ReadBuffer{
		buffer:   data,
		position: 0,
	}
}

// Read reads data from the buffer
func (rb *ReadBuffer) Read(p []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.position >= len(rb.buffer) {
		return 0, fmt.Errorf("EOF")
	}

	n := copy(p, rb.buffer[rb.position:])
	rb.position += n
	return n, nil
}

// Seek seeks to a position
func (rb *ReadBuffer) Seek(pos int) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if pos < 0 || pos > len(rb.buffer) {
		return fmt.Errorf("invalid position")
	}

	rb.position = pos
	return nil
}

// Position returns the current position
func (rb *ReadBuffer) Position() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.position
}

// Size returns the buffer size
func (rb *ReadBuffer) Size() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer)
}

// Remaining returns the remaining bytes
func (rb *ReadBuffer) Remaining() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return len(rb.buffer) - rb.position
}

// Reset resets the buffer
func (rb *ReadBuffer) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.position = 0
}

// BufferPool manages a pool of buffers
type BufferPool struct {
	buffers  [][]byte
	mu       sync.Mutex
	capacity int
	bufferSize int
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(capacity, bufferSize int) *BufferPool {
	return &BufferPool{
		buffers:    make([][]byte, 0, capacity),
		capacity:   capacity,
		bufferSize: bufferSize,
	}
}

// Get gets a buffer from the pool
func (bp *BufferPool) Get() []byte {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if len(bp.buffers) > 0 {
		buf := bp.buffers[len(bp.buffers)-1]
		bp.buffers = bp.buffers[:len(bp.buffers)-1]
		return buf
	}

	return make([]byte, bp.bufferSize)
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf []byte) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if len(bp.buffers) < bp.capacity {
		bp.buffers = append(bp.buffers, buf[:0])
	}
}

// Size returns the pool size
func (bp *BufferPool) Size() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.buffers)
}

// Capacity returns the pool capacity
func (bp *BufferPool) Capacity() int {
	return bp.capacity
}

// CacheManager manages multiple caches
type CacheManager struct {
	caches map[string]*LRUCache
	mu     sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*LRUCache),
	}
}

// CreateCache creates a new cache
func (cm *CacheManager) CreateCache(name string, capacity int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.caches[name] = NewLRUCache(capacity)
}

// GetCache gets a cache by name
func (cm *CacheManager) GetCache(name string) (*LRUCache, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cache, ok := cm.caches[name]
	return cache, ok
}

// DeleteCache deletes a cache
func (cm *CacheManager) DeleteCache(name string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.caches, name)
}

// ListCaches lists all caches
func (cm *CacheManager) ListCaches() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	names := make([]string, 0, len(cm.caches))
	for name := range cm.caches {
		names = append(names, name)
	}
	return names
}

// ClearAll clears all caches
func (cm *CacheManager) ClearAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, cache := range cm.caches {
		cache.Clear()
	}
}

// GetTotalSize returns total cache size
func (cm *CacheManager) GetTotalSize() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	total := 0
	for _, cache := range cm.caches {
		total += cache.Size()
	}
	return total
}

// GetTotalCapacity returns total cache capacity
func (cm *CacheManager) GetTotalCapacity() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	total := 0
	for _, cache := range cm.caches {
		total += cache.Capacity()
	}
	return total
}