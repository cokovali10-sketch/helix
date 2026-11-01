// Package storage provides storage engine tests and benchmarks
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"sync"
	"time"
)

func BenchmarkWAL_Record(b *testing.B) {
	dir := b.TempDir()
	wal, err := NewWAL(dir, 1024*1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer wal.Close()

	data := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wal.Record(data)
	}
}

func BenchmarkWAL_RecordLarge(b *testing.B) {
	dir := b.TempDir()
	wal, err := NewWAL(dir, 1024*1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer wal.Close()

	data := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wal.Record(data)
	}
}

func BenchmarkMemtable_Put(b *testing.B) {
	memtable := NewMemtable(1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		memtable.Put(key, value)
	}
}

func BenchmarkMemtable_Get(b *testing.B) {
	memtable := NewMemtable(1000000)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		memtable.Put(key, value)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		memtable.Get(key)
	}
}

func BenchmarkMemtable_GetMiss(b *testing.B) {
	memtable := NewMemtable(1000000)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		memtable.Put(key, value)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("missing-%d", i)
		memtable.Get(key)
	}
}

func BenchmarkBloomFilter_Add(b *testing.B) {
	bf := NewBloomFilter(100000, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		bf.Add(key)
	}
}

func BenchmarkBloomFilter_Contains(b *testing.B) {
	bf := NewBloomFilter(100000, 3)
	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		bf.Add(key)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key-%d", i%10000))
		bf.Contains(key)
	}
}

func BenchmarkStorageEngine_Put(b *testing.B) {
	dir := b.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"
	
	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024*1024, 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		engine.Put(key, value)
	}
}

func BenchmarkStorageEngine_Get(b *testing.B) {
	dir := b.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"
	
	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024*1024, 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		engine.Put(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		engine.Get(key)
	}
}

func BenchmarkStorageEngine_GetMiss(b *testing.B) {
	dir := b.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"
	
	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024*1024, 1024*1024)
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		engine.Put(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("missing-%d", i)
		engine.Get(key)
	}
}

func BenchmarkSSTableWriter_Write(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "test.sst")

	writer, err := NewSSTableWriter(path)
	if err != nil {
		b.Fatal(err)
	}
	defer writer.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		writer.Write(key, value)
	}
}

func BenchmarkGzipCompressor_Compress(b *testing.B) {
	compressor := NewGzipCompressor(6)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compressor.Compress(data)
	}
}

func BenchmarkGzipCompressor_Decompress(b *testing.B) {
	compressor := NewGzipCompressor(6)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	compressed, _ := compressor.Compress(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compressor.Decompress(compressed)
	}
}

func BenchmarkRetentionManager_Run(b *testing.B) {
	dir := b.TempDir()
	policy := RetentionPolicy{Type: "time", Duration: time.Hour}
	rm := NewRetentionManager(dir, policy)

	// Create files
	for i := 0; i < 100; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, make([]byte, 1024), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.Run()
	}
}

func BenchmarkCompactionManager_Compact(b *testing.B) {
	dir := b.TempDir()
	cm := NewCompactionManager(dir, "leveled", 4)

	// Create files
	for i := 0; i < 10; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, make([]byte, 1024), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Compact()
	}
}

func TestStorageEngine_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

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
				engine.Put(key, value)
				engine.Get(key)
			}
		}(g)
	}

	wg.Wait()
}

func TestStorageEngine_LargeDataset(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024*1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	// Insert 10000 entries
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%05d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := engine.Put(key, value); err != nil {
			t.Fatalf("Failed to put %d: %v", i, err)
		}
	}

	// Verify all entries
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%05d", i)
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

func TestStorageEngine_FlushAndRead(t *testing.T) {
	dir := t.TempDir()
	walDir := dir + "/wal"
	sstDir := dir + "/sst"

	engine, err := NewStorageEngine(walDir, sstDir, 1024, 100) // Small memtable to trigger flush
	if err != nil {
		t.Fatalf("Failed to create storage engine: %v", err)
	}
	defer engine.Close()

	// Insert enough entries to trigger flush
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		engine.Put(key, value)
	}

	// Verify entries are still readable
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, ok := engine.Get(key)
		if !ok {
			t.Fatalf("Expected to find key %d", i)
		}
	}
}

func TestBloomFilter_FalsePositiveRate(t *testing.T) {
	bf := NewBloomFilter(10000, 3)

	// Add 1000 items
	for i := 0; i < 1000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		bf.Add(key)
	}

	// Check false positive rate
	falsePositives := 0
	for i := 1000; i < 2000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		if bf.Contains(key) {
			falsePositives++
		}
	}

	rate := float64(falsePositives) / 1000.0
	if rate > 0.1 {
		t.Errorf("False positive rate too high: %f", rate)
	}
}

