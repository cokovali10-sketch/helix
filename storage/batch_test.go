// Package storage provides storage engine tests for transaction and snapshot
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteBatch_Create(t *testing.T) {
	batch := NewWriteBatch()
	if batch == nil {
		t.Fatal("Batch is nil")
	}
	if batch.Size() != 0 {
		t.Errorf("Expected 0, got %d", batch.Size())
	}
}

func TestWriteBatch_Put(t *testing.T) {
	batch := NewWriteBatch()
	batch.Put("key1", []byte("value1"))
	batch.Put("key2", []byte("value2"))

	if batch.Size() != 2 {
		t.Errorf("Expected 2, got %d", batch.Size())
	}
}

func TestWriteBatch_Delete(t *testing.T) {
	batch := NewWriteBatch()
	batch.Delete("key1")

	if batch.Size() != 1 {
		t.Errorf("Expected 1, got %d", batch.Size())
	}
}

func TestWriteBatch_Clear(t *testing.T) {
	batch := NewWriteBatch()
	batch.Put("key1", []byte("value1"))
	batch.Clear()

	if batch.Size() != 0 {
		t.Errorf("Expected 0, got %d", batch.Size())
	}
}

func TestWriteBatch_Apply(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	batch := NewWriteBatch()
	batch.Put("key1", []byte("value1"))
	batch.Put("key2", []byte("value2"))
	batch.Delete("key3")

	if err := batch.Apply(engine); err != nil {
		t.Fatalf("Failed to apply: %v", err)
	}

	value, ok := engine.Get("key1")
	if !ok || string(value) != "value1" {
		t.Error("Expected key1 to exist")
	}

	value, ok = engine.Get("key2")
	if !ok || string(value) != "value2" {
		t.Error("Expected key2 to exist")
	}
}

func TestIterator_Create(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	it := NewIterator(engine, "", false)
	if it == nil {
		t.Fatal("Iterator is nil")
	}
}

func TestRangeScan(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	engine.Put("key1", []byte("value1"))
	engine.Put("key2", []byte("value2"))

	result, err := engine.RangeScan("key1", "key2")
	if err != nil {
		t.Fatalf("Failed to range scan: %v", err)
	}
	_ = result
}

func TestPrefixScan(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	engine.Put("prefix-key1", []byte("value1"))
	engine.Put("prefix-key2", []byte("value2"))
	engine.Put("other-key", []byte("value3"))

	result, err := engine.PrefixScan("prefix-")
	if err != nil {
		t.Fatalf("Failed to prefix scan: %v", err)
	}
	_ = result
}

func TestStorageEngine_Count(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	for i := 0; i < 100; i++ {
		engine.Put(fmt.Sprintf("key-%d", i), []byte(fmt.Sprintf("value-%d", i)))
	}

	count, err := engine.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}
	_ = count
}

func TestStorageEngine_Compact(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	if err := engine.Compact(); err != nil {
		t.Fatalf("Failed to compact: %v", err)
	}
}

func TestStorageEngine_Flush(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	engine.Put("key1", []byte("value1"))
	if err := engine.Flush(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
}

func TestStorageEngine_Health(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	health := engine.Health()
	if health == nil {
		t.Fatal("Health is nil")
	}
}

func TestStorageEngine_GetStats(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024, 1024)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	engine.Put("key1", []byte("value1"))

	stats := engine.GetStats()
	if stats == nil {
		t.Fatal("Stats is nil")
	}
}