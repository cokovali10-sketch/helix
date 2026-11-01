// Package storage provides buffer pool
package storage

import (
	"sync"
)

// BufferPool pools byte slices
type BufferPool struct {
	pool sync.Pool
	size int
}

// NewBufferPool creates pool
func NewBufferPool(size int) *BufferPool {
	return &BufferPool{
		size: size,
		pool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, size)
				return &b
			},
		},
	}
}

// Get gets buffer
func (bp *BufferPool) Get() []byte {
	b := bp.pool.Get().(*[]byte)
	return *b
}

// Put returns buffer
func (bp *BufferPool) Put(b []byte) {
	if cap(b) < bp.size {
		return
	}
	b = b[:bp.size]
	bp.pool.Put(&b)
}

// Size returns size
func (bp *BufferPool) Size() int {
	return bp.size
}
