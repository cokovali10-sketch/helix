// Package storage provides bloom filter
package storage

import (
	"hash/fnv"
	"sync"
)

// BloomFilter is a bloom filter
type BloomFilter struct {
	bits []bool
	size int
	mu   sync.RWMutex
}

// NewBloomFilter creates filter
func NewBloomFilter(size int) *BloomFilter {
	return &BloomFilter{bits: make([]bool, size), size: size}
}

// Add adds item
func (bf *BloomFilter) Add(item string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	h := fnv.New32a()
	h.Write([]byte(item))
	idx := int(h.Sum32()) % bf.size
	bf.bits[idx] = true
}

// Contains checks membership
func (bf *BloomFilter) Contains(item string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	h := fnv.New32a()
	h.Write([]byte(item))
	idx := int(h.Sum32()) % bf.size
	return bf.bits[idx]
}

// Reset resets filter
func (bf *BloomFilter) Reset() {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	bf.bits = make([]bool, bf.size)
}

// Size returns size
func (bf *BloomFilter) Size() int {
	return bf.size
}
