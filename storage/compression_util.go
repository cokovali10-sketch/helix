// Package storage provides compression utilities
package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"sync"
)

// Compressor compresses data
type Compressor struct {
	mu sync.Mutex
}

// NewCompressor creates compressor
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Gzip compresses
func (c *Compressor) Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gunzip decompresses
func (c *Compressor) Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// Zlib compresses
func (c *Compressor) Zlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ZlibDecompress decompresses
func (c *Compressor) ZlibDecompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// Ratio returns compression ratio
func (c *Compressor) Ratio(original, compressed int) float64 {
	if original == 0 {
		return 0
	}
	return float64(compressed) / float64(original)
}

// EstimateSize estimates compressed size
func EstimateSize(data []byte) int {
	return len(data) / 2
}

// Validate checks gzip data
func ValidateGzip(data []byte) error {
	_, err := (&Compressor{}).Gunzip(data)
	if err != nil {
		return fmt.Errorf("invalid gzip: %w", err)
	}
	return nil
}
