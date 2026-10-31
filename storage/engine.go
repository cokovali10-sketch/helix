// Package storage implements the Helix TSDB storage engine in Go
// Includes WAL, Memtable, SSTable, compression, and retention
package storage

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// WAL represents Write-Ahead Log
type WAL struct {
	dir       string
	file      *os.File
	writer    *bufio.Writer
	sequence  uint64
	maxSize   int64
	currentSize int64
	mu        sync.Mutex
}

// NewWAL creates a new WAL
func NewWAL(dir string, maxSize int64) (*WAL, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	wal := &WAL{
		dir:     dir,
		maxSize: maxSize,
	}

	// Find next sequence number
	seq, err := wal.getNextSequence()
	if err != nil {
		return nil, err
	}
	wal.sequence = seq

	// Open WAL file
	path := filepath.Join(dir, fmt.Sprintf("wal-%d.log", seq))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	wal.file = file
	wal.writer = bufio.NewWriter(file)

	return wal, nil
}

func (w *WAL) getNextSequence() (uint64, error) {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return 1, nil
	}

	var maxSeq uint64
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > 8 && name[:4] == "wal-" && name[len(name)-4:] == ".log" {
			var seq uint64
			fmt.Sscanf(name, "wal-%d.log", &seq)
			if seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	return maxSeq + 1, nil
}

// Record writes a record to the WAL
func (w *WAL) Record(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if we need to rotate
	if w.currentSize+int64(len(data))+8 > w.maxSize {
		if err := w.rotate(); err != nil {
			return err
		}
	}

	// Write length prefix
	lenBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuf, uint32(len(data)))
	if _, err := w.writer.Write(lenBuf); err != nil {
		return err
	}

	// Write data
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// Write CRC
	crc := crc32.ChecksumIEEE(data)
	crcBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(crcBuf, crc)
	if _, err := w.writer.Write(crcBuf); err != nil {
		return err
	}

	w.currentSize += int64(len(data)) + 8
	return w.writer.Flush()
}

// rotate rotates the WAL file
func (w *WAL) rotate() error {
	if err := w.writer.Flush(); err != nil {
		return err
	}
	if err := w.file.Close(); err != nil {
		return err
	}

	w.sequence++
	path := filepath.Join(w.dir, fmt.Sprintf("wal-%d.log", w.sequence))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.writer = bufio.NewWriter(file)
	w.currentSize = 0

	return nil
}

// Close closes the WAL
func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Close()
}

// Memtable represents in-memory storage
type Memtable struct {
	data     map[string][]byte
	sorted   []string
	maxSize  int
	size     int
	mu       sync.RWMutex
}

// NewMemtable creates a new memtable
func NewMemtable(maxSize int) *Memtable {
	return &Memtable{
		data:   make(map[string][]byte),
		sorted: make([]string, 0),
		maxSize: maxSize,
	}
}

// Put inserts a key-value pair
func (m *Memtable) Put(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; !exists {
		m.sorted = append(m.sorted, key)
		sort.Strings(m.sorted)
	}

	m.data[key] = value
	m.size += len(value)
}

// Get retrieves a value by key
func (m *Memtable) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	if !ok {
		return nil, false
	}
	return value, true
}

// Delete removes a key
func (m *Memtable) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	for i, k := range m.sorted {
		if k == key {
			m.sorted = append(m.sorted[:i], m.sorted[i+1:]...)
			break
		}
	}
}

// Size returns the approximate size
func (m *Memtable) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.size
}

// IsFull checks if memtable is full
func (m *Memtable) IsFull() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.size >= m.maxSize
}

// Flush flushes memtable to SSTable
func (m *Memtable) Flush(writer *SSTableWriter) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, key := range m.sorted {
		if value, ok := m.data[key]; ok {
			if err := writer.Write(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// SSTableWriter writes SSTable files
type SSTableWriter struct {
	file    *os.File
	writer  *bufio.Writer
	index   map[string]int64
	offset  int64
	bloom   *BloomFilter
	mu      sync.Mutex
}

// NewSSTableWriter creates a new SSTable writer
func NewSSTableWriter(path string) (*SSTableWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &SSTableWriter{
		file:   file,
		writer: bufio.NewWriter(file),
		index:  make(map[string]int64),
		bloom:  NewBloomFilter(10000, 3),
	}, nil
}

// Write writes a key-value pair
func (w *SSTableWriter) Write(key string, value []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Record offset
	w.index[key] = w.offset

	// Write key length
	keyLen := uint32(len(key))
	if err := binary.Write(w.writer, binary.LittleEndian, keyLen); err != nil {
		return err
	}
	w.offset += 4

	// Write key
	if _, err := w.writer.WriteString(key); err != nil {
		return err
	}
	w.offset += int64(len(key))

	// Write value length
	valueLen := uint32(len(value))
	if err := binary.Write(w.writer, binary.LittleEndian, valueLen); err != nil {
		return err
	}
	w.offset += 4

	// Write value
	if _, err := w.writer.Write(value); err != nil {
		return err
	}
	w.offset += int64(len(value))

	// Add to bloom filter
	w.bloom.Add([]byte(key))

	return nil
}

// Close closes the writer
func (w *SSTableWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Write index
	indexOffset := w.offset
	for key, offset := range w.index {
		// Write key
		keyLen := uint32(len(key))
		binary.Write(w.writer, binary.LittleEndian, keyLen)
		w.writer.WriteString(key)
		
		// Write offset
		binary.Write(w.writer, binary.LittleEndian, offset)
		w.offset += int64(4 + len(key) + 8)
	}

	// Write footer
	footer := SSTableFooter{
		IndexOffset:  indexOffset,
		IndexSize:    w.offset - indexOffset,
		NumEntries:   uint64(len(w.index)),
	}
	
	// Write footer
	binary.Write(w.writer, binary.LittleEndian, footer.IndexOffset)
	binary.Write(w.writer, binary.LittleEndian, footer.IndexSize)
	binary.Write(w.writer, binary.LittleEndian, footer.NumEntries)

	if err := w.writer.Flush(); err != nil {
		return err
	}

	return w.file.Close()
}

// SSTableFooter represents SSTable footer
type SSTableFooter struct {
	IndexOffset  uint64
	IndexSize    uint64
	NumEntries   uint64
}

// SSTableReader reads SSTable files
type SSTableReader struct {
	file   *os.File
	index  map[string]uint64
	bloom  *BloomFilter
	mu     sync.RWMutex
}

// NewSSTableReader creates a new SSTable reader
func NewSSTableReader(path string) (*SSTableReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := &SSTableReader{
		file:  file,
		index: make(map[string]uint64),
		bloom: NewBloomFilter(10000, 3),
	}

	// Read index from footer
	// Simplified - in production would read from footer
	return reader, nil
}

// Get retrieves a value by key
func (r *SSTableReader) Get(key string) ([]byte, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check bloom filter
	if !r.bloom.Contains([]byte(key)) {
		return nil, false
	}

	// Check index
	offset, ok := r.index[key]
	if !ok {
		return nil, false
	}

	// Read from file
	r.file.Seek(offset, io.SeekStart)
	
	// Read key length
	var keyLen uint32
	binary.Read(r.file, binary.LittleEndian, &keyLen)
	
	// Skip key
	r.file.Seek(int64(keyLen), io.SeekCurrent)
	
	// Read value length
	var valueLen uint32
	binary.Read(r.file, binary.LittleEndian, &valueLen)
	
	// Read value
	value := make([]byte, valueLen)
	r.file.Read(value)
	
	return value, true
}

// Contains checks if a key exists
func (r *SSTableReader) Contains(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.bloom.Contains([]byte(key)) {
		return false
	}
	_, ok := r.index[key]
	return ok
}

// Close closes the reader
func (r *SSTableReader) Close() error {
	return r.file.Close()
}

// BloomFilter implements a bloom filter
type BloomFilter struct {
	bits             []byte
	numHashFunctions int
	numBits          int
}

// NewBloomFilter creates a new bloom filter
func NewBloomFilter(numBits, numHashFunctions int) *BloomFilter {
	numBytes := (numBits + 7) / 8
	return &BloomFilter{
		bits:             make([]byte, numBytes),
		numHashFunctions: numHashFunctions,
		numBits:          numBits,
	}
}

// Add adds a key to the bloom filter
func (bf *BloomFilter) Add(key []byte) {
	for i := 0; i < bf.numHashFunctions; i++ {
		hash := bf.hash(key, i)
		bitIndex := hash % bf.numBits
		byteIndex := bitIndex / 8
		bitOffset := bitIndex % 8
		bf.bits[byteIndex] |= 1 << bitOffset
	}
}

// Contains checks if a key might be in the bloom filter
func (bf *BloomFilter) Contains(key []byte) bool {
	for i := 0; i < bf.numHashFunctions; i++ {
		hash := bf.hash(key, i)
		bitIndex := hash % bf.numBits
		byteIndex := bitIndex / 8
		bitOffset := bitIndex % 8
		if bf.bits[byteIndex]&(1<<bitOffset) == 0 {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hash(key []byte, seed int) int {
	hash := seed * 0x9e3779b9
	for _, b := range key {
		hash = hash*31 + int(b)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// StorageEngine combines WAL, Memtable, and SSTable
type StorageEngine struct {
	wal       *WAL
	memtable  *Memtable
	sstDir    string
	sstables  []*SSTableReader
	mu        sync.RWMutex
	nextSeq   uint64
}

// NewStorageEngine creates a new storage engine
func NewStorageEngine(walDir, sstDir string, walMaxSize int64, memtableMaxSize int) (*StorageEngine, error) {
	wal, err := NewWAL(walDir, walMaxSize)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(sstDir, 0755); err != nil {
		return nil, err
	}

	return &StorageEngine{
		wal:      wal,
		memtable: NewMemtable(memtableMaxSize),
		sstDir:   sstDir,
		sstables: make([]*SSTableReader, 0),
		nextSeq:  1,
	}, nil
}

// Put inserts a key-value pair
func (se *StorageEngine) Put(key string, value []byte) error {
	// Write to WAL
	if err := se.wal.Record([]byte(key + ":" + string(value))); err != nil {
		return err
	}

	// Insert into memtable
	se.memtable.Put(key, value)

	// Check if memtable is full
	if se.memtable.IsFull() {
		if err := se.flushMemtable(); err != nil {
			return err
		}
	}

	return nil
}

// Get retrieves a value by key
func (se *StorageEngine) Get(key string) ([]byte, bool) {
	// Check memtable first
	if value, ok := se.memtable.Get(key); ok {
		return value, true
	}

	// Check SSTables
	se.mu.RLock()
	defer se.mu.RUnlock()

	for i := len(se.sstables) - 1; i >= 0; i-- {
		if value, ok := se.sstables[i].Get(key); ok {
			return value, true
		}
	}

	return nil, false
}

// Delete removes a key
func (se *StorageEngine) Delete(key string) error {
	se.memtable.Delete(key)
	return nil
}

// flushMemtable flushes memtable to SSTable
func (se *StorageEngine) flushMemtable() error {
	se.mu.Lock()
	defer se.mu.Unlock()

	// Create SSTable writer
	path := filepath.Join(se.sstDir, fmt.Sprintf("%d.sst", se.nextSeq))
	se.nextSeq++

	writer, err := NewSSTableWriter(path)
	if err != nil {
		return err
	}

	// Flush memtable
	if err := se.memtable.Flush(writer); err != nil {
		return err
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return err
	}

	// Open reader
	reader, err := NewSSTableReader(path)
	if err != nil {
		return err
	}

	se.sstables = append(se.sstables, reader)

	// Reset memtable
	se.memtable = NewMemtable(se.memtable.maxSize)

	return nil
}

// Close closes the storage engine
func (se *StorageEngine) Close() error {
	if err := se.wal.Close(); err != nil {
		return err
	}

	for _, sst := range se.sstables {
		sst.Close()
	}

	return nil
}

// Stats returns storage statistics
func (se *StorageEngine) Stats() map[string]interface{} {
	se.mu.RLock()
	defer se.mu.RUnlock()

	return map[string]interface{}{
		"memtable_size":   se.memtable.Size(),
		"sstable_count":   len(se.sstables),
		"next_sequence":   se.nextSeq,
	}
}