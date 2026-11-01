// Package ingestion provides batch processing
package ingestion

import (
	"sync"
	"time"
)

// BatchProcessor processes batches
type BatchProcessor struct {
	batchSize int
	flushInterval time.Duration
	buffer    []map[string]interface{}
	mu        sync.Mutex
	flushFn   func([]map[string]interface{})
	stop      chan struct{}
}

// NewBatchProcessor creates processor
func NewBatchProcessor(size int, interval time.Duration, fn func([]map[string]interface{})) *BatchProcessor {
	bp := &BatchProcessor{
		batchSize:     size,
		flushInterval: interval,
		buffer:        []map[string]interface{}{},
		flushFn:       fn,
		stop:          make(chan struct{}),
	}
	go bp.run()
	return bp
}

// Add adds record
func (bp *BatchProcessor) Add(rec map[string]interface{}) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.buffer = append(bp.buffer, rec)
	if len(bp.buffer) >= bp.batchSize {
		bp.flushLocked()
	}
}

// Flush flushes buffer
func (bp *BatchProcessor) Flush() {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.flushLocked()
}

func (bp *BatchProcessor) flushLocked() {
	if len(bp.buffer) == 0 {
		return
	}
	batch := bp.buffer
	bp.buffer = []map[string]interface{}{}
	go bp.flushFn(batch)
}

func (bp *BatchProcessor) run() {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bp.Flush()
		case <-bp.stop:
			bp.Flush()
			return
		}
	}
}

// Stop stops processor
func (bp *BatchProcessor) Stop() {
	close(bp.stop)
}

// Size returns buffer size
func (bp *BatchProcessor) Size() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.buffer)
}
