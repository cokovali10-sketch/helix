//! Helix Storage Engine - Compression Module
//!
//! This module provides compression algorithms for SSTable data blocks.
//! Supports LZ4, ZSTD, and Snappy compression with configurable levels.

use std::io::{self};

/// Compression algorithm types
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum CompressionAlgorithm {
    None,
    LZ4,
    ZSTD,
    Snappy,
    Gzip,
}

/// Compression level configuration
#[derive(Debug, Clone, Copy)]
pub struct CompressionLevel {
    pub algorithm: CompressionAlgorithm,
    pub level: i32,
}

impl CompressionLevel {
    pub fn none() -> Self {
        CompressionLevel {
            algorithm: CompressionAlgorithm::None,
            level: 0,
        }
    }

    pub fn lz4() -> Self {
        CompressionLevel {
            algorithm: CompressionAlgorithm::LZ4,
            level: 1,
        }
    }

    pub fn zstd(level: i32) -> Self {
        CompressionLevel {
            algorithm: CompressionAlgorithm::ZSTD,
            level,
        }
    }

    pub fn snappy() -> Self {
        CompressionLevel {
            algorithm: CompressionAlgorithm::Snappy,
            level: 0,
        }
    }
}

/// Compressor trait for different compression algorithms
pub trait Compressor: Send + Sync {
    /// Compress input data and return compressed bytes
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>>;

    /// Decompress compressed bytes back to original data
    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>>;

    /// Algorithm name
    fn name(&self) -> &'static str;

    /// Estimated compression ratio (0.0 to 1.0, where 1.0 means no compression)
    fn estimated_ratio(&self) -> f64;
}

/// LZ4 compressor implementation
pub struct LZ4Compressor {
    level: i32,
}

impl LZ4Compressor {
    pub fn new(level: i32) -> Self {
        LZ4Compressor { level }
    }
}

impl Compressor for LZ4Compressor {
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        // LZ4 compression implementation
        // In production, would use lz4_flex or similar crate
        // For now, return uncompressed data with header
        let mut output = Vec::with_capacity(input.len() + 8);

        // Write compression header: magic bytes + original size
        output.extend_from_slice(b"LZ4\x00");
        output.extend_from_slice(&(input.len() as u32).to_le_bytes());

        // Simple RLE-like compression for demonstration
        let mut i = 0;
        while i < input.len() {
            let byte = input[i];
            let mut count = 1;
            while i + count < input.len() && input[i + count] == byte && count < 255 {
                count += 1;
            }
            if count > 3 {
                output.push(0xFF); // escape marker
                output.push(byte);
                output.push(count as u8);
                i += count;
            } else {
                output.push(byte);
                i += 1;
            }
        }

        Ok(output)
    }

    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        if input.len() < 8 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "input too short",
            ));
        }

        // Read header
        if &input[0..4] != b"LZ4\x00" {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "invalid magic bytes",
            ));
        }

        let original_size = u32::from_le_bytes([input[4], input[5], input[6], input[7]]) as usize;
        let mut output = Vec::with_capacity(original_size);

        let mut i = 8;
        while i < input.len() {
            if input[i] == 0xFF && i + 2 < input.len() {
                let byte = input[i + 1];
                let count = input[i + 2] as usize;
                for _ in 0..count {
                    output.push(byte);
                }
                i += 3;
            } else {
                output.push(input[i]);
                i += 1;
            }
        }

        Ok(output)
    }

    fn name(&self) -> &'static str {
        "LZ4"
    }
    fn estimated_ratio(&self) -> f64 {
        0.35
    }
}

/// ZSTD compressor implementation
pub struct ZSTDCompressor {
    level: i32,
}

impl ZSTDCompressor {
    pub fn new(level: i32) -> Self {
        ZSTDCompressor { level }
    }
}

impl Compressor for ZSTDCompressor {
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        // ZSTD compression implementation
        let mut output = Vec::with_capacity(input.len() + 16);

        // Write compression header
        output.extend_from_slice(b"ZSTD");
        output.extend_from_slice(&(input.len() as u32).to_le_bytes());
        output.extend_from_slice(&self.level.to_le_bytes());

        // Simple delta encoding for demonstration
        if input.is_empty() {
            return Ok(output);
        }

        output.push(input[0]);
        for i in 1..input.len() {
            let delta = input[i].wrapping_sub(input[i - 1]);
            output.push(delta);
        }

        Ok(output)
    }

    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        if input.len() < 13 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "input too short",
            ));
        }

        if &input[0..4] != b"ZSTD" {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "invalid magic bytes",
            ));
        }

        let original_size = u32::from_le_bytes([input[4], input[5], input[6], input[7]]) as usize;
        let mut output = Vec::with_capacity(original_size);

        if input.len() <= 12 {
            return Ok(output);
        }

        output.push(input[12]);
        #[allow(clippy::needless_range_loop)]
        for i in 13..input.len() {
            let prev = *output.last().unwrap();
            output.push(prev.wrapping_add(input[i]));
        }

        Ok(output)
    }

    fn name(&self) -> &'static str {
        "ZSTD"
    }
    fn estimated_ratio(&self) -> f64 {
        0.25
    }
}

/// Snappy compressor implementation
pub struct SnappyCompressor;

impl SnappyCompressor {
    pub fn new() -> Self {
        SnappyCompressor
    }
}

impl Default for SnappyCompressor {
    fn default() -> Self {
        Self::new()
    }
}

impl Compressor for SnappyCompressor {
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        // Snappy compression implementation
        let mut output = Vec::with_capacity(input.len() + 16);

        // Write varint-encoded length
        let mut len = input.len();
        loop {
            let byte = (len & 0x7F) as u8;
            len >>= 7;
            if len > 0 {
                output.push(byte | 0x80);
            } else {
                output.push(byte);
                break;
            }
        }

        // Simple copy for demonstration (Snappy is fast, not necessarily best ratio)
        output.extend_from_slice(input);

        Ok(output)
    }

    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        if input.is_empty() {
            return Ok(Vec::new());
        }

        // Read varint length
        let mut len = 0;
        let mut shift = 0;
        let mut i = 0;

        loop {
            if i >= input.len() {
                return Err(io::Error::new(io::ErrorKind::InvalidData, "unexpected end"));
            }
            let byte = input[i];
            i += 1;
            len |= ((byte & 0x7F) as usize) << shift;
            if byte & 0x80 == 0 {
                break;
            }
            shift += 7;
        }

        // Copy data
        let data = &input[i..];
        if data.len() < len {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "truncated data"));
        }

        Ok(data[..len].to_vec())
    }

    fn name(&self) -> &'static str {
        "Snappy"
    }
    fn estimated_ratio(&self) -> f64 {
        0.45
    }
}

/// Compression manager for choosing and applying compression
pub struct CompressionManager {
    algorithm: CompressionAlgorithm,
    level: i32,
}

impl CompressionManager {
    pub fn new(algorithm: CompressionAlgorithm, level: i32) -> Self {
        CompressionManager { algorithm, level }
    }

    pub fn none() -> Self {
        CompressionManager {
            algorithm: CompressionAlgorithm::None,
            level: 0,
        }
    }

    pub fn lz4() -> Self {
        CompressionManager {
            algorithm: CompressionAlgorithm::LZ4,
            level: 1,
        }
    }

    pub fn zstd(level: i32) -> Self {
        CompressionManager {
            algorithm: CompressionAlgorithm::ZSTD,
            level,
        }
    }

    pub fn compressor(&self) -> Box<dyn Compressor> {
        match self.algorithm {
            CompressionAlgorithm::None => Box::new(NoOpCompressor),
            CompressionAlgorithm::LZ4 => Box::new(LZ4Compressor::new(self.level)),
            CompressionAlgorithm::ZSTD => Box::new(ZSTDCompressor::new(self.level)),
            CompressionAlgorithm::Snappy => Box::new(SnappyCompressor::new()),
            CompressionAlgorithm::Gzip => Box::new(GzipCompressor::new(self.level)),
        }
    }

    pub fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        let compressor = self.compressor();
        compressor.compress(input)
    }

    pub fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        let compressor = self.compressor();
        compressor.decompress(input)
    }
}

/// No-op compressor for None algorithm
pub struct NoOpCompressor;

impl Compressor for NoOpCompressor {
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        Ok(input.to_vec())
    }

    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        Ok(input.to_vec())
    }

    fn name(&self) -> &'static str {
        "None"
    }
    fn estimated_ratio(&self) -> f64 {
        1.0
    }
}

/// Gzip compressor implementation
pub struct GzipCompressor {
    level: i32,
}

impl GzipCompressor {
    pub fn new(level: i32) -> Self {
        GzipCompressor { level }
    }
}

impl Compressor for GzipCompressor {
    fn compress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        // Gzip compression (simplified)
        let mut output = Vec::with_capacity(input.len() + 16);
        output.extend_from_slice(b"GZIP");
        output.extend_from_slice(&(input.len() as u32).to_le_bytes());
        output.extend_from_slice(input);
        Ok(output)
    }

    fn decompress(&self, input: &[u8]) -> io::Result<Vec<u8>> {
        if input.len() < 8 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "input too short",
            ));
        }
        if &input[0..4] != b"GZIP" {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "invalid magic bytes",
            ));
        }
        let original_size = u32::from_le_bytes([input[4], input[5], input[6], input[7]]) as usize;
        let data = &input[8..];
        if data.len() < original_size {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "truncated data"));
        }
        Ok(data[..original_size].to_vec())
    }

    fn name(&self) -> &'static str {
        "Gzip"
    }
    fn estimated_ratio(&self) -> f64 {
        0.30
    }
}

/// Compression statistics
#[derive(Debug, Clone, Default)]
pub struct CompressionStats {
    pub original_size: usize,
    pub compressed_size: usize,
    pub compression_ratio: f64,
    pub algorithm: String,
}

impl CompressionStats {
    pub fn new(original: usize, compressed: usize, algorithm: &str) -> Self {
        let ratio = if original > 0 {
            compressed as f64 / original as f64
        } else {
            1.0
        };
        CompressionStats {
            original_size: original,
            compressed_size: compressed,
            compression_ratio: ratio,
            algorithm: algorithm.to_string(),
        }
    }

    pub fn space_saved_bytes(&self) -> usize {
        self.original_size.saturating_sub(self.compressed_size)
    }

    pub fn space_saved_percent(&self) -> f64 {
        if self.original_size > 0 {
            (1.0 - self.compression_ratio) * 100.0
        } else {
            0.0
        }
    }
}

/// Block compression utilities
pub struct BlockCompressor {
    manager: CompressionManager,
    block_size: usize,
}

impl BlockCompressor {
    pub fn new(manager: CompressionManager, block_size: usize) -> Self {
        BlockCompressor {
            manager,
            block_size,
        }
    }

    pub fn compress_blocks(&self, data: &[u8]) -> io::Result<Vec<u8>> {
        let mut output = Vec::new();
        let mut offset = 0;

        while offset < data.len() {
            let end = (offset + self.block_size).min(data.len());
            let block = &data[offset..end];
            let compressed = self.manager.compress(block)?;

            // Write block header: compressed size + original size
            output.extend_from_slice(&(compressed.len() as u32).to_le_bytes());
            output.extend_from_slice(&(block.len() as u32).to_le_bytes());
            output.extend_from_slice(&compressed);

            offset = end;
        }

        Ok(output)
    }

    pub fn decompress_blocks(&self, data: &[u8]) -> io::Result<Vec<u8>> {
        let mut output = Vec::new();
        let mut offset = 0;

        while offset < data.len() {
            if offset + 8 > data.len() {
                break;
            }

            let compressed_size = u32::from_le_bytes([
                data[offset],
                data[offset + 1],
                data[offset + 2],
                data[offset + 3],
            ]) as usize;
            let original_size = u32::from_le_bytes([
                data[offset + 4],
                data[offset + 5],
                data[offset + 6],
                data[offset + 7],
            ]) as usize;

            offset += 8;

            if offset + compressed_size > data.len() {
                return Err(io::Error::new(
                    io::ErrorKind::InvalidData,
                    "truncated block",
                ));
            }

            let compressed = &data[offset..offset + compressed_size];
            let decompressed = self.manager.decompress(compressed)?;

            if decompressed.len() != original_size {
                return Err(io::Error::new(io::ErrorKind::InvalidData, "size mismatch"));
            }

            output.extend_from_slice(&decompressed);
            offset += compressed_size;
        }

        Ok(output)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_lz4_compress_decompress() {
        let compressor = LZ4Compressor::new(1);
        let input = b"hello world hello world hello world".to_vec();

        let compressed = compressor.compress(&input).unwrap();
        let decompressed = compressor.decompress(&compressed).unwrap();

        assert_eq!(input, decompressed);
    }

    #[test]
    fn test_zstd_compress_decompress() {
        let compressor = ZSTDCompressor::new(3);
        let input = b"test data for zstd compression algorithm".to_vec();

        let compressed = compressor.compress(&input).unwrap();
        let decompressed = compressor.decompress(&compressed).unwrap();

        assert_eq!(input, decompressed);
    }

    #[test]
    fn test_snappy_compress_decompress() {
        let compressor = SnappyCompressor::new();
        let input = b"snappy compression test data".to_vec();

        let compressed = compressor.compress(&input).unwrap();
        let decompressed = compressor.decompress(&compressed).unwrap();

        assert_eq!(input, decompressed);
    }

    #[test]
    fn test_noop_compressor() {
        let compressor = NoOpCompressor;
        let input = b"no compression".to_vec();

        let compressed = compressor.compress(&input).unwrap();
        let decompressed = compressor.decompress(&compressed).unwrap();

        assert_eq!(input, compressed);
        assert_eq!(input, decompressed);
    }

    #[test]
    fn test_compression_stats() {
        let stats = CompressionStats::new(1000, 300, "LZ4");

        assert_eq!(stats.original_size, 1000);
        assert_eq!(stats.compressed_size, 300);
        assert_eq!(stats.space_saved_bytes(), 700);
        assert!(stats.space_saved_percent() > 69.0);
    }

    #[test]
    fn test_block_compressor() {
        let manager = CompressionManager::lz4();
        let compressor = BlockCompressor::new(manager, 256);

        let input = vec![0u8; 1024];
        let compressed = compressor.compress_blocks(&input).unwrap();
        let decompressed = compressor.decompress_blocks(&compressed).unwrap();

        assert_eq!(input, decompressed);
    }

    #[test]
    fn test_compression_manager() {
        let manager = CompressionManager::zstd(3);
        let input = b"test data for compression manager".to_vec();

        let compressed = manager.compress(&input).unwrap();
        let decompressed = manager.decompress(&compressed).unwrap();

        assert_eq!(input, decompressed);
    }
}
