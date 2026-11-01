// Package storage contains tests for the LSM tree
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLSMTree_Create(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	if tree == nil {
		t.Fatal("Tree is nil")
	}
}

func TestLSMTree_PutGet(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	key := "test-key"
	value := []byte("test-value")

	if err := tree.Put(key, value); err != nil {
		t.Fatalf("Failed to put: %v", err)
	}

	got, ok := tree.Get(key)
	if !ok {
		t.Fatal("Expected to find key")
	}
	if string(got) != string(value) {
		t.Errorf("Expected %s, got %s", value, got)
	}
}

func TestLSMTree_Delete(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	key := "test-key"
	value := []byte("test-value")

	tree.Put(key, value)
	tree.Delete(key)

	_, ok := tree.Get(key)
	if ok {
		t.Fatal("Expected key to be deleted")
	}
}

func TestLSMTree_MultiplePuts(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := tree.Put(key, value); err != nil {
			t.Fatalf("Failed to put %d: %v", i, err)
		}
	}

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		got, ok := tree.Get(key)
		if !ok {
			t.Fatalf("Expected to find key %d", i)
		}
		if string(got) != string(value) {
			t.Errorf("Expected %s, got %s", value, got)
		}
	}
}

func TestLSMTree_GetStats(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	tree.Put("key1", []byte("value1"))
	tree.Put("key2", []byte("value2"))

	stats := tree.GetStats()
	if stats == nil {
		t.Fatal("Stats is nil")
	}
}

func TestLSMTree_Count(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		tree.Put(key, value)
	}

	count := tree.Count()
	if count < 100 {
		t.Errorf("Expected at least 100, got %d", count)
	}
}

func TestLSMTree_Size(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	tree.Put("key1", []byte("value1"))
	tree.Put("key2", []byte("value2"))

	size := tree.Size()
	if size <= 0 {
		t.Error("Expected size > 0")
	}
}

func TestLSMTree_Health(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	health := tree.Health()
	if health == nil {
		t.Fatal("Health is nil")
	}
}

func TestLSMTree_Flush(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	config.MemtableSize = 100 // Small memtable to trigger flush
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	// Insert enough data to trigger flush
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		tree.Put(key, value)
	}

	// Verify data is still readable
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, ok := tree.Get(key)
		if !ok {
			t.Fatalf("Expected to find key %d", i)
		}
	}
}

func TestLSMTree_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}
	defer tree.Close()

	var wg sync.WaitGroup
	goroutines := 10
	opsPerGoroutine := 100

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				key := fmt.Sprintf("key-%d-%d", g, i)
				value := []byte(fmt.Sprintf("value-%d-%d", g, i))
				tree.Put(key, value)
				tree.Get(key)
			}
		}(g)
	}

	wg.Wait()
}

func TestCompactor_Create(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	cm := NewCompactor(dir, config)
	if cm == nil {
		t.Fatal("Compactor is nil")
	}
}

func TestCompactor_Compact(t *testing.T) {
	dir := t.TempDir()
	config := DefaultLSMConfig()
	cm := NewCompactor(dir, config)

	// Create some SSTable files
	for i := 0; i < 5; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, []byte("data"), 0644)
	}

	cm.Compact(make([]*SSTableReader, 0))
}

func TestLSMConfig_Default(t *testing.T) {
	config := DefaultLSMConfig()
	if config == nil {
		t.Fatal("Config is nil")
	}
	if config.MemtableSize != 64*1024*1024 {
		t.Errorf("Expected 64MB, got %d", config.MemtableSize)
	}
}

func BenchmarkLSMTree_Put(b *testing.B) {
	dir := b.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		b.Fatal(err)
	}
	defer tree.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		tree.Put(key, value)
	}
}

func BenchmarkLSMTree_Get(b *testing.B) {
	dir := b.TempDir()
	config := DefaultLSMConfig()
	tree, err := NewLSMTree(dir+"/wal", dir+"/sst", config)
	if err != nil {
		b.Fatal(err)
	}
	defer tree.Close()

	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		tree.Put(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		tree.Get(key)
	}
}