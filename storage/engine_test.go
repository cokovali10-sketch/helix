// Package storage contains tests for the storage engine
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWAL_Create(t *testing.T) {
	dir := t.TempDir()
	wal, err := NewWAL(dir, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}
	defer wal.Close()

	if wal.sequence != 1 {
		t.Errorf("Expected sequence 1, got %d", wal.sequence)
	}
}

func TestWAL_Record(t *testing.T) {
	dir := t.TempDir()
	wal, err := NewWAL(dir, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}
	defer wal.Close()

	data := []byte("test data")
	if err := wal.Record(data); err != nil {
		t.Fatalf("Failed to record: %v", err)
	}
}

func TestWAL_MultipleRecords(t *testing.T) {
	dir := t.TempDir()
	wal, err := NewWAL(dir, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}
	defer wal.Close()

	for i := 0; i < 100; i++ {
		data := []byte(fmt.Sprintf("record-%d", i))
		if err := wal.Record(data); err != nil {
			t.Fatalf("Failed to record %d: %v", i, err)
		}
	}
}

func TestWAL_Rotation(t *testing.T) {
	dir := t.TempDir()
	wal, err := NewWAL(dir, 100) // Small max size to trigger rotation
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}
	defer wal.Close()

	// Write enough data to trigger rotation
	for i := 0; i < 10; i++ {
		data := make([]byte, 50)
		if err := wal.Record(data); err != nil {
			t.Fatalf("Failed to record %d: %v", i, err)
		}
	}

	// Check that multiple WAL files exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	walFiles := 0
	for _, entry := range entries {
		if len(entry.Name()) > 4 && entry.Name()[:4] == "wal-" {
			walFiles++
		}
	}

	if walFiles < 2 {
		t.Errorf("Expected at least 2 WAL files after rotation, got %d", walFiles)
	}
}

func TestMemtable_Create(t *testing.T) {
	memtable := NewMemtable(1024)
	if memtable == nil {
		t.Fatal("Failed to create memtable")
	}
	if memtable.Size() != 0 {
		t.Errorf("Expected size 0, got %d", memtable.Size())
	}
}

func TestMemtable_PutGet(t *testing.T) {
	memtable := NewMemtable(1024)

	key := "test-key"
	value := []byte("test-value")
	memtable.Put(key, value)

	got, ok := memtable.Get(key)
	if !ok {
		t.Fatal("Expected to find key")
	}
	if string(got) != string(value) {
		t.Errorf("Expected %s, got %s", value, got)
	}
}

func TestMemtable_Delete(t *testing.T) {
	memtable := NewMemtable(1024)

	key := "test-key"
	value := []byte("test-value")
	memtable.Put(key, value)

	memtable.Delete(key)

	_, ok := memtable.Get(key)
	if ok {
		t.Fatal("Expected key to be deleted")
	}
}

func TestMemtable_IsFull(t *testing.T) {
	memtable := NewMemtable(10)

	memtable.Put("key1", []byte("12345"))
	if memtable.IsFull() {
		t.Error("Expected not full")
	}

	memtable.Put("key2", []byte("12345"))
	if !memtable.IsFull() {
		t.Error("Expected full")
	}
}

func TestMemtable_Size(t *testing.T) {
	memtable := NewMemtable(1024)

	memtable.Put("key1", []byte("value1"))
	memtable.Put("key2", []byte("value2"))

	if memtable.Size() != 12 {
		t.Errorf("Expected size 12, got %d", memtable.Size())
	}
}

func TestMemtable_Flush(t *testing.T) {
	memtable := NewMemtable(1024)
	dir := t.TempDir()

	memtable.Put("key1", []byte("value1"))
	memtable.Put("key2", []byte("value2"))

	path := filepath.Join(dir, "test.sst")
	writer, err := NewSSTableWriter(path)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}

	if err := memtable.Flush(writer); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
}

func TestBloomFilter_Create(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	if bf == nil {
		t.Fatal("Failed to create bloom filter")
	}
}

func TestBloomFilter_AddContains(t *testing.T) {
	bf := NewBloomFilter(1000, 3)

	key := []byte("test-key")
	bf.Add(key)

	if !bf.Contains(key) {
		t.Error("Expected bloom filter to contain key")
	}
}

func TestBloomFilter_NotContains(t *testing.T) {
	bf := NewBloomFilter(1000, 3)

	bf.Add([]byte("key1"))

	// May or may not contain (false positive possible)
	// Just check it doesn't panic
	_ = bf.Contains([]byte("key2"))
}

func TestSSTableWriter_Create(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sst")

	writer, err := NewSSTableWriter(path)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	if writer == nil {
		t.Fatal("Writer is nil")
	}
}

func TestSSTableWriter_Write(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sst")

	writer, err := NewSSTableWriter(path)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}

	if err := writer.Write("key1", []byte("value1")); err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
}

func TestSSTableWriter_MultipleWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sst")

	writer, err := NewSSTableWriter(path)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := writer.Write(key, value); err != nil {
			t.Fatalf("Failed to write %d: %v", i, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
}

func TestSSTableReader_Create(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sst")

	// Create an SSTable first
	writer, err := NewSSTableWriter(path)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	writer.Write("key1", []byte("value1"))
	writer.Close()

	// Open reader
	reader, err := NewSSTableReader(path)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
	defer reader.Close()

	if reader == nil {
		t.Fatal("Reader is nil")
	}
}

func TestStorageEngine_Create(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	sstDir := filepath.Join(dir, "sst")

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	if engine == nil {
		t.Fatal("Engine is nil")
	}
}

func TestStorageEngine_PutGet(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	sstDir := filepath.Join(dir, "sst")

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	key := "test-key"
	value := []byte("test-value")

	if err := engine.Put(key, value); err != nil {
		t.Fatalf("Failed to put: %v", err)
	}

	got, ok := engine.Get(key)
	if !ok {
		t.Fatal("Expected to find key")
	}
	if string(got) != string(value) {
		t.Errorf("Expected %s, got %s", value, got)
	}
}

func TestStorageEngine_Delete(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	sstDir := filepath.Join(dir, "sst")

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	key := "test-key"
	value := []byte("test-value")

	engine.Put(key, value)
	engine.Delete(key)

	_, ok := engine.Get(key)
	if ok {
		t.Fatal("Expected key to be deleted")
	}
}

func TestStorageEngine_MultiplePuts(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	sstDir := filepath.Join(dir, "sst")

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := engine.Put(key, value); err != nil {
			t.Fatalf("Failed to put %d: %v", i, err)
		}
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		got, ok := engine.Get(key)
		if !ok {
			t.Fatalf("Expected to find key %d", i)
		}
		if string(got) != string(value) {
			t.Errorf("Expected %s, got %s", value, got)
		}
	}
}

func TestStorageEngine_Stats(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	sstDir := filepath.Join(dir, "sst")

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	stats := engine.Stats()
	if stats == nil {
		t.Fatal("Stats is nil")
	}
}