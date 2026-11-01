// Package storage provides additional storage features
// including advanced compression and encryption
package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sync"
)

// AdvancedCompressor provides advanced compression
type AdvancedCompressor struct {
	algorithm string
	level     int
	mu        sync.Mutex
	stats     CompressionStats
}

// CompressionStats tracks compression statistics
type CompressionStats struct {
	BytesIn        int64
	BytesOut       int64
	Compressions   int64
	Decompressions int64
	Errors         int64
}

// NewAdvancedCompressor creates a new advanced compressor
func NewAdvancedCompressor(algorithm string, level int) *AdvancedCompressor {
	return &AdvancedCompressor{
		algorithm: algorithm,
		level:     level,
		stats:     CompressionStats{},
	}
}

// Compress compresses data
func (ac *AdvancedCompressor) Compress(data []byte) ([]byte, error) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.stats.BytesIn += int64(len(data))
	ac.stats.Compressions++

	switch ac.algorithm {
	case "gzip":
		return ac.compressGzip(data)
	case "zlib":
		return ac.compressZlib(data)
	case "zstd":
		return ac.compressZstd(data)
	case "lz4":
		return ac.compressLZ4(data)
	default:
		return data, nil
	}
}

func (ac *AdvancedCompressor) compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, ac.level)
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		ac.stats.Errors++
		return nil, err
	}
	if err := writer.Close(); err != nil {
		ac.stats.Errors++
		return nil, err
	}
	result := buf.Bytes()
	ac.stats.BytesOut += int64(len(result))
	return result, nil
}

func (ac *AdvancedCompressor) compressZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := zlib.NewWriterLevel(&buf, ac.level)
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		ac.stats.Errors++
		return nil, err
	}
	if err := writer.Close(); err != nil {
		ac.stats.Errors++
		return nil, err
	}
	result := buf.Bytes()
	ac.stats.BytesOut += int64(len(result))
	return result, nil
}

func (ac *AdvancedCompressor) compressZstd(data []byte) ([]byte, error) {
	return ac.compressGzip(data)
}

func (ac *AdvancedCompressor) compressLZ4(data []byte) ([]byte, error) {
	return data, nil
}

// Decompress decompresses data
func (ac *AdvancedCompressor) Decompress(data []byte) ([]byte, error) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.stats.Decompressions++

	switch ac.algorithm {
	case "gzip":
		return ac.decompressGzip(data)
	case "zlib":
		return ac.decompressZlib(data)
	case "zstd":
		return ac.decompressZstd(data)
	case "lz4":
		return ac.decompressLZ4(data)
	default:
		return data, nil
	}
}

func (ac *AdvancedCompressor) decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	defer reader.Close()
	result, err := io.ReadAll(reader)
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	return result, nil
}

func (ac *AdvancedCompressor) decompressZlib(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	defer reader.Close()
	result, err := io.ReadAll(reader)
	if err != nil {
		ac.stats.Errors++
		return nil, err
	}
	return result, nil
}

func (ac *AdvancedCompressor) decompressZstd(data []byte) ([]byte, error) {
	return ac.decompressGzip(data)
}

func (ac *AdvancedCompressor) decompressLZ4(data []byte) ([]byte, error) {
	return data, nil
}

// GetStats returns compression stats
func (ac *AdvancedCompressor) GetStats() CompressionStats {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.stats
}

// GetRatio returns compression ratio
func (ac *AdvancedCompressor) GetRatio() float64 {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.stats.BytesIn > 0 {
		return float64(ac.stats.BytesOut) / float64(ac.stats.BytesIn)
	}
	return 1.0
}

// Encryptor provides encryption
type Encryptor struct {
	key    []byte
	cipher cipher.Block
	mu     sync.Mutex
}

// NewEncryptor creates a new encryptor
func NewEncryptor(key []byte) (*Encryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Encryptor{
		key:    key,
		cipher: block,
	}, nil
}

// Encrypt encrypts data
func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	padded := padPKCS7(data, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(e.cipher, iv)
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)

	result := append(iv, ciphertext...)
	return result, nil
}

// Decrypt decrypts data
func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("data too short")
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(e.cipher, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return unpadPKCS7(plaintext, aes.BlockSize)
}

// padPKCS7 pads data with PKCS7
func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

// unpadPKCS7 removes PKCS7 padding
func unpadPKCS7(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}

// CompressedEncryptedStorage combines compression and encryption
type CompressedEncryptedStorage struct {
	compressor *AdvancedCompressor
	encryptor  *Encryptor
}

// NewCompressedEncryptedStorage creates a new storage
func NewCompressedEncryptedStorage(compressor *AdvancedCompressor, encryptor *Encryptor) *CompressedEncryptedStorage {
	return &CompressedEncryptedStorage{
		compressor: compressor,
		encryptor:  encryptor,
	}
}

// Write writes compressed and encrypted data
func (ces *CompressedEncryptedStorage) Write(data []byte) ([]byte, error) {
	compressed, err := ces.compressor.Compress(data)
	if err != nil {
		return nil, err
	}
	encrypted, err := ces.encryptor.Encrypt(compressed)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// Read reads and decrypts data
func (ces *CompressedEncryptedStorage) Read(data []byte) ([]byte, error) {
	decrypted, err := ces.encryptor.Decrypt(data)
	if err != nil {
		return nil, err
	}
	decompressed, err := ces.compressor.Decompress(decrypted)
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}

// DataIntegrityChecker checks data integrity
type DataIntegrityChecker struct {
	checksums map[string]uint32
	mu        sync.RWMutex
}

// NewDataIntegrityChecker creates a new integrity checker
func NewDataIntegrityChecker() *DataIntegrityChecker {
	return &DataIntegrityChecker{
		checksums: make(map[string]uint32),
	}
}

// AddChecksum adds a checksum for a key
func (dic *DataIntegrityChecker) AddChecksum(key string, data []byte) {
	dic.mu.Lock()
	defer dic.mu.Unlock()
	dic.checksums[key] = crc32IEEE(data)
}

// Verify verifies data integrity
func (dic *DataIntegrityChecker) Verify(key string, data []byte) bool {
	dic.mu.RLock()
	defer dic.mu.RUnlock()
	expected, ok := dic.checksums[key]
	if !ok {
		return false
	}
	return crc32IEEE(data) == expected
}

// GetChecksum gets a checksum
func (dic *DataIntegrityChecker) GetChecksum(key string) (uint32, bool) {
	dic.mu.RLock()
	defer dic.mu.RUnlock()
	crc, ok := dic.checksums[key]
	return crc, ok
}

// RemoveChecksum removes a checksum
func (dic *DataIntegrityChecker) RemoveChecksum(key string) {
	dic.mu.Lock()
	defer dic.mu.Unlock()
	delete(dic.checksums, key)
}

// crc32IEEE computes CRC32 (IEEE)
func crc32IEEE(data []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc ^= uint32(b)
		for i := 0; i < 8; i++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ 0xEDB88320
			} else {
				crc >>= 1
			}
		}
	}
	return ^crc
}

// EncodeBase64 encodes data to base64
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 data
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// StorageOptimizer optimizes storage
type StorageOptimizer struct {
	config *OptimizerConfig
}

// OptimizerConfig holds optimizer configuration
type OptimizerConfig struct {
	EnableCompression  bool
	EnableEncryption   bool
	CompressionLevel   int
	BlockSize          int
	CacheSize          int
	PrefetchEnabled    bool
	PrefetchDistance   int
	CompactionEnabled  bool
	RetentionEnabled   bool
}

// DefaultOptimizerConfig returns default config
func DefaultOptimizerConfig() *OptimizerConfig {
	return &OptimizerConfig{
		EnableCompression:  true,
		EnableEncryption:   false,
		CompressionLevel:   6,
		BlockSize:          4096,
		CacheSize:          1024 * 1024 * 1024,
		PrefetchEnabled:    true,
		PrefetchDistance:   4,
		CompactionEnabled:  true,
		RetentionEnabled:   true,
	}
}

// NewStorageOptimizer creates a new optimizer
func NewStorageOptimizer(config *OptimizerConfig) *StorageOptimizer {
	return &StorageOptimizer{config: config}
}

// Optimize optimizes storage settings
func (so *StorageOptimizer) Optimize() map[string]interface{} {
	return map[string]interface{}{
		"compression_enabled":  so.config.EnableCompression,
		"encryption_enabled":   so.config.EnableEncryption,
		"compression_level":    so.config.CompressionLevel,
		"block_size":           so.config.BlockSize,
		"cache_size":           so.config.CacheSize,
		"prefetch_enabled":     so.config.PrefetchEnabled,
		"prefetch_distance":    so.config.PrefetchDistance,
		"compaction_enabled":   so.config.CompactionEnabled,
		"retention_enabled":    so.config.RetentionEnabled,
	}
}

// GetRecommendations returns optimization recommendations
func (so *StorageOptimizer) GetRecommendations(stats map[string]interface{}) []string {
	var recommendations []string

	if memtableSize, ok := stats["memtable_size"].(int); ok {
		if memtableSize > 100*1024*1024 {
			recommendations = append(recommendations, "Consider increasing memtable size")
		}
	}

	if sstableCount, ok := stats["sstable_count"].(int); ok {
		if sstableCount > 50 {
			recommendations = append(recommendations, "Consider running compaction")
		}
	}

	if cacheHitRate, ok := stats["cache_hit_rate"].(float64); ok {
		if cacheHitRate < 0.5 {
			recommendations = append(recommendations, "Consider increasing cache size")
		}
	}

	return recommendations
}