// Package ingestion provides compression for ingestion
package ingestion

import (
	"bytes"
	"compress/gzip"
	"io"
)

// IngestCompressor compresses payloads
type IngestCompressor struct{}

// NewIngestCompressor creates compressor
func NewIngestCompressor() *IngestCompressor {
	return &IngestCompressor{}
}

// Compress compresses data
func (ic *IngestCompressor) Compress(data []byte) ([]byte, error) {
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

// Decompress decompresses data
func (ic *IngestCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
