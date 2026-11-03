//! Helix Storage Engine - SSTable Writer Module (Extended)
//!
//! This module provides an extended SSTable writer with bloom filter,
//! compression, and block-based writing.

use std::collections::BTreeMap;
use std::fs::File;
use std::io::{self, BufWriter, Write};
use std::path::PathBuf;
use std::sync::Arc;

use super::compression::{CompressionAlgorithm, CompressionManager};
use super::metrics::StorageMetrics;

/// SSTable block
#[derive(Debug, Clone)]
pub struct Block {
    pub data: Vec<u8>,
    pub offset: u64,
    pub size: usize,
    pub compressed_size: usize,
    pub num_entries: usize,
    pub min_key: Vec<u8>,
    pub max_key: Vec<u8>,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
}

impl Block {
    pub fn new() -> Self {
        Block {
            data: Vec::new(),
            offset: 0,
            size: 0,
            compressed_size: 0,
            num_entries: 0,
            min_key: Vec::new(),
            max_key: Vec::new(),
            min_timestamp: i64::MAX,
            max_timestamp: i64::MIN,
        }
    }

    pub fn add_entry(&mut self, key: &[u8], value: &[u8], timestamp: i64) {
        // Write key length
        self.data
            .extend_from_slice(&(key.len() as u32).to_le_bytes());
        self.data.extend_from_slice(key);

        // Write value length
        self.data
            .extend_from_slice(&(value.len() as u32).to_le_bytes());
        self.data.extend_from_slice(value);

        // Write timestamp
        self.data.extend_from_slice(&timestamp.to_le_bytes());

        if self.min_key.is_empty() || key < self.min_key.as_slice() {
            self.min_key = key.to_vec();
        }
        if self.max_key.is_empty() || key > self.max_key.as_slice() {
            self.max_key = key.to_vec();
        }
        self.min_timestamp = self.min_timestamp.min(timestamp);
        self.max_timestamp = self.max_timestamp.max(timestamp);

        self.size += key.len() + value.len() + 12;
        self.num_entries += 1;
    }

    pub fn is_full(&self, block_size: usize) -> bool {
        self.size >= block_size
    }
}

impl Default for Block {
    fn default() -> Self {
        Self::new()
    }
}

/// Extended SSTable writer
pub struct ExtendedSSTableWriter {
    path: PathBuf,
    config: SSTableWriterConfig,
    blocks: Vec<Block>,
    current_block: Block,
    index: BTreeMap<Vec<u8>, (u64, u32)>, // key -> (offset, block_id)
    bloom_filter: BloomFilter,
    compression: CompressionManager,
    metrics: Option<Arc<StorageMetrics>>,
    entry_count: usize,
}

/// SSTable writer configuration
#[derive(Debug, Clone)]
pub struct SSTableWriterConfig {
    pub block_size: usize,
    pub compression: CompressionAlgorithm,
    pub compression_level: i32,
    pub bloom_filter_bits: usize,
    pub bloom_filter_hashes: usize,
    pub enable_checksums: bool,
}

impl Default for SSTableWriterConfig {
    fn default() -> Self {
        SSTableWriterConfig {
            block_size: 4096,
            compression: CompressionAlgorithm::ZSTD,
            compression_level: 3,
            bloom_filter_bits: 10000,
            bloom_filter_hashes: 3,
            enable_checksums: true,
        }
    }
}

/// Bloom filter (simplified)
#[derive(Debug)]
pub struct BloomFilter {
    bits: Vec<u8>,
    num_hashes: usize,
    num_bits: usize,
}

impl BloomFilter {
    pub fn new(num_bits: usize, num_hashes: usize) -> Self {
        BloomFilter {
            bits: vec![0u8; num_bits.div_ceil(8)],
            num_hashes,
            num_bits,
        }
    }

    pub fn add(&mut self, key: &[u8]) {
        for i in 0..self.num_hashes {
            let hash = self.hash(key, i);
            let bit = hash % self.num_bits;
            self.bits[bit / 8] |= 1 << (bit % 8);
        }
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        for i in 0..self.num_hashes {
            let hash = self.hash(key, i);
            let bit = hash % self.num_bits;
            if self.bits[bit / 8] & (1 << (bit % 8)) == 0 {
                return false;
            }
        }
        true
    }

    fn hash(&self, key: &[u8], seed: usize) -> usize {
        let mut h: usize = seed.wrapping_mul(0x9e3779b9);
        for &b in key {
            h = h.wrapping_mul(31).wrapping_add(b as usize);
        }
        h
    }
}

impl ExtendedSSTableWriter {
    pub fn new(path: PathBuf, config: SSTableWriterConfig) -> io::Result<Self> {
        let compression = CompressionManager::new(config.compression, config.compression_level);
        let bloom_filter = BloomFilter::new(config.bloom_filter_bits, config.bloom_filter_hashes);

        Ok(ExtendedSSTableWriter {
            path,
            config,
            blocks: Vec::new(),
            current_block: Block::new(),
            index: BTreeMap::new(),
            bloom_filter,
            compression,
            metrics: None,
            entry_count: 0,
        })
    }

    pub fn with_metrics(mut self, metrics: Arc<StorageMetrics>) -> Self {
        self.metrics = Some(metrics);
        self
    }

    pub fn add(&mut self, key: &[u8], value: &[u8], timestamp: i64) -> io::Result<()> {
        // Check if current block is full
        if self.current_block.is_full(self.config.block_size) && !self.current_block.data.is_empty()
        {
            self.flush_block()?;
        }

        self.current_block.add_entry(key, value, timestamp);
        self.index
            .insert(key.to_vec(), (0, self.blocks.len() as u32)); // placeholder offset
        self.bloom_filter.add(key);
        self.entry_count += 1;

        if let Some(ref metrics) = self.metrics {
            metrics
                .writes_total
                .fetch_add(1, std::sync::atomic::Ordering::Relaxed);
        }

        Ok(())
    }

    fn flush_block(&mut self) -> io::Result<()> {
        if self.current_block.data.is_empty() {
            return Ok(());
        }

        // Compress block
        let compressed = self.compression.compress(&self.current_block.data)?;
        self.current_block.compressed_size = compressed.len();

        // Store block
        self.blocks.push(self.current_block.clone());
        self.current_block = Block::new();

        Ok(())
    }

    pub fn finish(mut self) -> io::Result<ExtendedSSTable> {
        // Flush remaining block
        if !self.current_block.data.is_empty() {
            self.flush_block()?;
        }

        // Write to file
        let mut file = File::create(&self.path)?;
        let mut writer = BufWriter::new(&mut file);

        // Write blocks
        let mut offset = 0u64;
        for block in self.blocks.iter() {
            let _block_offset = offset;
            writer.write_all(&block.data)?;
            offset += block.data.len() as u64;

            // Update index with actual offset
            // (simplified - in production would track per-key offsets)
        }

        // Write index
        let index_offset = offset;
        for (key, (key_offset, block_id)) in &self.index {
            writer.write_all(&(key.len() as u32).to_le_bytes())?;
            writer.write_all(key)?;
            writer.write_all(&key_offset.to_le_bytes())?;
            writer.write_all(&block_id.to_le_bytes())?;
            offset += 4 + key.len() as u64 + 8 + 4;
        }

        // Write bloom filter
        let bloom_offset = offset;
        writer.write_all(&self.bloom_filter.bits)?;
        offset += self.bloom_filter.bits.len() as u64;

        // Write footer
        self.write_footer(&mut writer, index_offset, bloom_offset, offset)?;

        writer.flush()?;

        Ok(ExtendedSSTable {
            path: self.path,
            entry_count: self.entry_count,
            block_count: self.blocks.len(),
            index: self.index,
            bloom_filter: self.bloom_filter,
        })
    }

    fn write_footer(
        &self,
        writer: &mut BufWriter<&mut File>,
        index_offset: u64,
        bloom_offset: u64,
        _footer_offset: u64,
    ) -> io::Result<()> {
        // Write magic number
        writer.write_all(b"SSTB")?;

        // Write index offset
        writer.write_all(&index_offset.to_le_bytes())?;

        // Write bloom offset
        writer.write_all(&bloom_offset.to_le_bytes())?;

        // Write entry count
        writer.write_all(&(self.entry_count as u64).to_le_bytes())?;

        // Write block count
        writer.write_all(&(self.blocks.len() as u64).to_le_bytes())?;

        // Write checksum (simplified)
        writer.write_all(&0u32.to_le_bytes())?;

        Ok(())
    }

    pub fn entry_count(&self) -> usize {
        self.entry_count
    }

    pub fn block_count(&self) -> usize {
        self.blocks.len()
            + if self.current_block.data.is_empty() {
                0
            } else {
                1
            }
    }
}

/// Extended SSTable
#[derive(Debug)]
pub struct ExtendedSSTable {
    pub path: PathBuf,
    pub entry_count: usize,
    pub block_count: usize,
    pub index: BTreeMap<Vec<u8>, (u64, u32)>,
    pub bloom_filter: BloomFilter,
}

impl ExtendedSSTable {
    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        // Check bloom filter
        if !self.bloom_filter.contains(key) {
            return None;
        }

        // Look up in index
        let (_offset, _block_id) = self.index.get(key)?;

        // Read from file at offset
        // Simplified - in production would read and decompress
        None
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        if !self.bloom_filter.contains(key) {
            return false;
        }
        self.index.contains_key(key)
    }

    pub fn keys(&self) -> Vec<Vec<u8>> {
        self.index.keys().cloned().collect()
    }

    pub fn size(&self) -> usize {
        self.index.len()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_block_add_entry() {
        let mut block = Block::new();
        block.add_entry(b"key1", b"value1", 1000);

        assert_eq!(block.num_entries, 1);
        assert_eq!(block.min_key, b"key1");
        assert_eq!(block.max_key, b"key1");
    }

    #[test]
    fn test_block_is_full() {
        let mut block = Block::new();
        block.size = 4096;
        assert!(block.is_full(4096));
    }

    #[test]
    fn test_bloom_filter() {
        let mut bf = BloomFilter::new(1000, 3);
        bf.add(b"key1");
        bf.add(b"key2");

        assert!(bf.contains(b"key1"));
        assert!(bf.contains(b"key2"));
    }

    #[test]
    fn test_sstable_writer_config_default() {
        let config = SSTableWriterConfig::default();
        assert_eq!(config.block_size, 4096);
        assert_eq!(config.compression, CompressionAlgorithm::ZSTD);
    }

    #[test]
    fn test_extended_sstable_writer() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");
        let config = SSTableWriterConfig::default();

        let mut writer = ExtendedSSTableWriter::new(path, config).unwrap();
        writer.add(b"key1", b"value1", 1000).unwrap();
        writer.add(b"key2", b"value2", 1001).unwrap();

        let sstable = writer.finish().unwrap();
        assert_eq!(sstable.entry_count, 2);
        assert!(sstable.contains(b"key1"));
        assert!(sstable.contains(b"key2"));
    }

    #[test]
    fn test_extended_sstable_keys() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");
        let config = SSTableWriterConfig::default();

        let mut writer = ExtendedSSTableWriter::new(path, config).unwrap();
        writer.add(b"key1", b"value1", 1000).unwrap();
        writer.add(b"key2", b"value2", 1001).unwrap();
        writer.add(b"key3", b"value3", 1002).unwrap();

        let sstable = writer.finish().unwrap();
        let keys = sstable.keys();
        assert_eq!(keys.len(), 3);
    }

    #[test]
    fn test_extended_sstable_size() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");
        let config = SSTableWriterConfig::default();

        let mut writer = ExtendedSSTableWriter::new(path, config).unwrap();
        for i in 0..100 {
            let key = format!("key-{:03}", i);
            let value = format!("value-{}", i);
            writer
                .add(key.as_bytes(), value.as_bytes(), 1000 + i as i64)
                .unwrap();
        }

        let sstable = writer.finish().unwrap();
        assert_eq!(sstable.size(), 100);
    }
}
