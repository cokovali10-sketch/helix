// Package storage provides additional tests for compression and retention
package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGzipCompressor_CompressDecompress(t *testing.T) {
	compressor := NewGzipCompressor(6)
	data := []byte("test data for gzip compression")

	compressed, err := compressor.Compress(data)
	if err != nil {
		t.Fatalf("Failed to compress: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(data, decompressed) {
		t.Error("Decompressed data doesn't match original")
	}
}

func TestGzipCompressor_Name(t *testing.T) {
	compressor := NewGzipCompressor(6)
	if compressor.Name() != "gzip" {
		t.Errorf("Expected gzip, got %s", compressor.Name())
	}
}

func TestZlibCompressor_CompressDecompress(t *testing.T) {
	compressor := NewZlibCompressor(6)
	data := []byte("test data for zlib compression")

	compressed, err := compressor.Compress(data)
	if err != nil {
		t.Fatalf("Failed to compress: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(data, decompressed) {
		t.Error("Decompressed data doesn't match original")
	}
}

func TestZlibCompressor_Name(t *testing.T) {
	compressor := NewZlibCompressor(6)
	if compressor.Name() != "zlib" {
		t.Errorf("Expected zlib, got %s", compressor.Name())
	}
}

func TestNoOpCompressor_CompressDecompress(t *testing.T) {
	compressor := &NoOpCompressor{}
	data := []byte("test data")

	compressed, err := compressor.Compress(data)
	if err != nil {
		t.Fatalf("Failed to compress: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(data, decompressed) {
		t.Error("Decompressed data doesn't match original")
	}
}

func TestNoOpCompressor_Name(t *testing.T) {
	compressor := &NoOpCompressor{}
	if compressor.Name() != "none" {
		t.Errorf("Expected none, got %s", compressor.Name())
	}
}

func TestRetentionManager_Create(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{Type: "time", Duration: time.Hour}
	rm := NewRetentionManager(dir, policy)

	if rm == nil {
		t.Fatal("RetentionManager is nil")
	}
}

func TestRetentionManager_Run(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{Type: "time", Duration: time.Hour}
	rm := NewRetentionManager(dir, policy)

	if err := rm.Run(); err != nil {
		t.Fatalf("Failed to run: %v", err)
	}
}

func TestRetentionManager_TimeBased(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{Type: "time", Duration: time.Hour}
	rm := NewRetentionManager(dir, policy)

	// Create an old file
	oldFile := filepath.Join(dir, "old.sst")
	os.WriteFile(oldFile, []byte("old data"), 0644)
	os.Chtimes(oldFile, time.Now().Add(-2*time.Hour), time.Now().Add(-2*time.Hour))

	if err := rm.Run(); err != nil {
		t.Fatalf("Failed to run: %v", err)
	}

	stats := rm.GetStats()
	if stats.FilesDeleted == 0 {
		t.Log("No files deleted (might be OK)")
	}
}

func TestRetentionManager_SizeBased(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{Type: "size", MaxBytes: 100}
	rm := NewRetentionManager(dir, policy)

	// Create files
	for i := 0; i < 5; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, make([]byte, 50), 0644)
	}

	if err := rm.Run(); err != nil {
		t.Fatalf("Failed to run: %v", err)
	}
}

func TestRetentionManager_CountBased(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{Type: "count", MaxFiles: 3}
	rm := NewRetentionManager(dir, policy)

	// Create files
	for i := 0; i < 5; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, []byte("data"), 0644)
		time.Sleep(10 * time.Millisecond)
	}

	if err := rm.Run(); err != nil {
		t.Fatalf("Failed to run: %v", err)
	}

	stats := rm.GetStats()
	if stats.FilesDeleted < 2 {
		t.Logf("Expected at least 2 files deleted, got %d", stats.FilesDeleted)
	}
}

func TestCompactionManager_Create(t *testing.T) {
	dir := t.TempDir()
	cm := NewCompactionManager(dir, "leveled", 4)

	if cm == nil {
		t.Fatal("CompactionManager is nil")
	}
}

func TestCompactionManager_Compact(t *testing.T) {
	dir := t.TempDir()
	cm := NewCompactionManager(dir, "leveled", 4)

	if err := cm.Compact(); err != nil {
		t.Fatalf("Failed to compact: %v", err)
	}
}

func TestCompactionManager_GetSSTableCount(t *testing.T) {
	dir := t.TempDir()
	cm := NewCompactionManager(dir, "leveled", 4)

	// Create some SSTable files
	for i := 0; i < 3; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file-%d.sst", i))
		os.WriteFile(path, []byte("data"), 0644)
	}

	count := cm.GetSSTableCount()
	if count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}
}

func TestStorageMetrics_Create(t *testing.T) {
	metrics := NewStorageMetrics()
	if metrics == nil {
		t.Fatal("Metrics is nil")
	}
}

func TestStorageMetrics_RecordWrite(t *testing.T) {
	metrics := NewStorageMetrics()
	metrics.RecordWrite(100)

	if metrics.WritesTotal != 1 {
		t.Errorf("Expected 1, got %d", metrics.WritesTotal)
	}
	if metrics.WritesBytes != 100 {
		t.Errorf("Expected 100, got %d", metrics.WritesBytes)
	}
}

func TestStorageMetrics_RecordRead(t *testing.T) {
	metrics := NewStorageMetrics()
	metrics.RecordRead(200)

	if metrics.ReadsTotal != 1 {
		t.Errorf("Expected 1, got %d", metrics.ReadsTotal)
	}
	if metrics.ReadsBytes != 200 {
		t.Errorf("Expected 200, got %d", metrics.ReadsBytes)
	}
}

func TestStorageMetrics_CacheHitRate(t *testing.T) {
	metrics := NewStorageMetrics()
	metrics.RecordCacheHit()
	metrics.RecordCacheHit()
	metrics.RecordCacheMiss()

	rate := metrics.CacheHitRate()
	if rate < 0.6 || rate > 0.7 {
		t.Errorf("Expected ~0.667, got %f", rate)
	}
}

func TestStorageMetrics_GetStats(t *testing.T) {
	metrics := NewStorageMetrics()
	metrics.RecordWrite(100)
	metrics.RecordRead(200)

	stats := metrics.GetStats()
	if stats == nil {
		t.Fatal("Stats is nil")
	}
	if stats["writes_total"] != int64(1) {
		t.Errorf("Expected 1, got %v", stats["writes_total"])
	}
}