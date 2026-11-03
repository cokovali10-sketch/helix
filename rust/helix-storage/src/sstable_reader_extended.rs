//! Helix Storage Engine - SSTable Reader Module (Extended)
//!
//! This module provides an extended SSTable reader with full read capabilities,
//! including block-based reading, index lookup, and range scans.

use std::collections::BTreeMap;
use std::fs::{self, File};
use std::io::{self, Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};

/// Extended SSTable reader
pub struct ExtendedSSTableReader {
    path: PathBuf,
    index: BTreeMap<Vec<u8>, BlockPointer>,
    bloom_filter: BloomFilter,
    footer: SSTableFooter,
    block_cache: BlockCache,
    metrics: Option<std::sync::Arc<crate::metrics::StorageMetrics>>,
}

/// Block pointer in SSTable
#[derive(Debug, Clone, Copy)]
pub struct BlockPointer {
    pub offset: u64,
    pub size: u32,
    pub compressed_size: u32,
}

/// SSTable footer
#[derive(Debug, Clone)]
pub struct SSTableFooter {
    pub magic: u32,
    pub version: u32,
    pub index_offset: u64,
    pub index_size: u64,
    pub bloom_offset: u64,
    pub bloom_size: u64,
    pub num_entries: u64,
    pub num_blocks: u64,
    pub compression_type: u8,
    pub checksum: u32,
}

impl SSTableFooter {
    pub fn new() -> Self {
        SSTableFooter {
            magic: 0x53535442, // "SSTB"
            version: 1,
            index_offset: 0,
            index_size: 0,
            bloom_offset: 0,
            bloom_size: 0,
            num_entries: 0,
            num_blocks: 0,
            compression_type: 0,
            checksum: 0,
        }
    }

    pub fn size() -> usize {
        4 + 4 + 8 + 8 + 8 + 8 + 8 + 8 + 1 + 4 // 57 bytes
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut buf = Vec::with_capacity(Self::size());
        buf.extend_from_slice(&self.magic.to_le_bytes());
        buf.extend_from_slice(&self.version.to_le_bytes());
        buf.extend_from_slice(&self.index_offset.to_le_bytes());
        buf.extend_from_slice(&self.index_size.to_le_bytes());
        buf.extend_from_slice(&self.bloom_offset.to_le_bytes());
        buf.extend_from_slice(&self.bloom_size.to_le_bytes());
        buf.extend_from_slice(&self.num_entries.to_le_bytes());
        buf.extend_from_slice(&self.num_blocks.to_le_bytes());
        buf.push(self.compression_type);
        buf.extend_from_slice(&self.checksum.to_le_bytes());
        buf
    }

    pub fn decode(data: &[u8]) -> Result<Self, String> {
        if data.len() < Self::size() {
            return Err("footer too small".to_string());
        }

        let mut offset = 0;
        let magic = u32::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
        ]);
        offset += 4;

        if magic != 0x53535442 {
            return Err("invalid magic number".to_string());
        }

        let version = u32::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
        ]);
        offset += 4;

        let index_offset = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let index_size = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let bloom_offset = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let bloom_size = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let num_entries = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let num_blocks = u64::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
            data[offset + 4],
            data[offset + 5],
            data[offset + 6],
            data[offset + 7],
        ]);
        offset += 8;

        let compression_type = data[offset];
        offset += 1;

        let checksum = u32::from_le_bytes([
            data[offset],
            data[offset + 1],
            data[offset + 2],
            data[offset + 3],
        ]);

        Ok(SSTableFooter {
            magic,
            version,
            index_offset,
            index_size,
            bloom_offset,
            bloom_size,
            num_entries,
            num_blocks,
            compression_type,
            checksum,
        })
    }
}

impl Default for SSTableFooter {
    fn default() -> Self {
        Self::new()
    }
}

/// Block cache
pub struct BlockCache {
    cache: std::sync::RwLock<BTreeMap<u64, Vec<u8>>>,
    capacity: usize,
    current_size: std::sync::atomic::AtomicUsize,
}

impl BlockCache {
    pub fn new(capacity: usize) -> Self {
        BlockCache {
            cache: std::sync::RwLock::new(BTreeMap::new()),
            capacity,
            current_size: std::sync::atomic::AtomicUsize::new(0),
        }
    }

    pub fn get(&self, block_id: u64) -> Option<Vec<u8>> {
        let cache = self.cache.read().unwrap();
        cache.get(&block_id).cloned()
    }

    pub fn put(&self, block_id: u64, data: Vec<u8>) {
        use std::sync::atomic::Ordering;
        let mut cache = self.cache.write().unwrap();
        if self.current_size.load(Ordering::Relaxed) + data.len() > self.capacity {
            // Evict oldest (simplified)
            if let Some(key) = cache.keys().next().cloned() {
                if let Some(val) = cache.remove(&key) {
                    self.current_size.fetch_sub(val.len(), Ordering::Relaxed);
                }
            }
        }
        self.current_size.fetch_add(data.len(), Ordering::Relaxed);
        cache.insert(block_id, data);
    }

    pub fn clear(&self) {
        let mut cache = self.cache.write().unwrap();
        cache.clear();
        self.current_size
            .store(0, std::sync::atomic::Ordering::Relaxed);
    }

    pub fn size(&self) -> usize {
        self.current_size.load(std::sync::atomic::Ordering::Relaxed)
    }
}

/// Bloom filter
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

impl ExtendedSSTableReader {
    pub fn new(path: PathBuf) -> io::Result<Self> {
        let file_size = fs::metadata(&path)?.len();

        if file_size < SSTableFooter::size() as u64 {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "file too small"));
        }

        // Read footer
        let mut file = File::open(&path)?;
        file.seek(SeekFrom::End(-(SSTableFooter::size() as i64)))?;
        let mut footer_buf = vec![0u8; SSTableFooter::size()];
        file.read_exact(&mut footer_buf)?;
        let footer = SSTableFooter::decode(&footer_buf)
            .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e))?;

        // Read index
        file.seek(SeekFrom::Start(footer.index_offset))?;
        let mut index_buf = vec![0u8; footer.index_size as usize];
        file.read_exact(&mut index_buf)?;
        let index = Self::parse_index(&index_buf)?;

        // Read bloom filter
        file.seek(SeekFrom::Start(footer.bloom_offset))?;
        let mut bloom_buf = vec![0u8; footer.bloom_size as usize];
        file.read_exact(&mut bloom_buf)?;
        let num_bits = bloom_buf.len() * 8;
        let bloom_filter = BloomFilter {
            bits: bloom_buf,
            num_hashes: 3,
            num_bits,
        };

        Ok(ExtendedSSTableReader {
            path,
            index,
            bloom_filter,
            footer,
            block_cache: BlockCache::new(1024 * 1024), // 1MB cache
            metrics: None,
        })
    }

    fn parse_index(buf: &[u8]) -> io::Result<BTreeMap<Vec<u8>, BlockPointer>> {
        let mut index = BTreeMap::new();
        let mut offset = 0;

        while offset + 16 <= buf.len() {
            let key_len = u32::from_le_bytes([
                buf[offset],
                buf[offset + 1],
                buf[offset + 2],
                buf[offset + 3],
            ]) as usize;
            offset += 4;

            if offset + key_len + 12 > buf.len() {
                break;
            }

            let key = buf[offset..offset + key_len].to_vec();
            offset += key_len;

            let block_offset = u64::from_le_bytes([
                buf[offset],
                buf[offset + 1],
                buf[offset + 2],
                buf[offset + 3],
                buf[offset + 4],
                buf[offset + 5],
                buf[offset + 6],
                buf[offset + 7],
            ]);
            offset += 8;

            let block_size = u32::from_le_bytes([
                buf[offset],
                buf[offset + 1],
                buf[offset + 2],
                buf[offset + 3],
            ]);
            offset += 4;

            index.insert(
                key,
                BlockPointer {
                    offset: block_offset,
                    size: block_size,
                    compressed_size: block_size,
                },
            );
        }

        Ok(index)
    }

    pub fn get(&self, key: &[u8]) -> io::Result<Option<Vec<u8>>> {
        // Check bloom filter
        if !self.bloom_filter.contains(key) {
            return Ok(None);
        }

        // Look up in index
        let pointer = match self.index.get(key) {
            Some(p) => p,
            None => return Ok(None),
        };

        // Check block cache
        if let Some(cached) = self.block_cache.get(pointer.offset) {
            return self.extract_value(&cached, key);
        }

        // Read from file
        let mut file = File::open(&self.path)?;
        file.seek(SeekFrom::Start(pointer.offset))?;

        let mut block_data = vec![0u8; pointer.size as usize];
        file.read_exact(&mut block_data)?;

        // Cache the block
        self.block_cache.put(pointer.offset, block_data.clone());

        self.extract_value(&block_data, key)
    }

    fn extract_value(&self, block_data: &[u8], key: &[u8]) -> io::Result<Option<Vec<u8>>> {
        let mut offset = 0;

        while offset < block_data.len() {
            if offset + 8 > block_data.len() {
                break;
            }

            let key_len = u32::from_le_bytes([
                block_data[offset],
                block_data[offset + 1],
                block_data[offset + 2],
                block_data[offset + 3],
            ]) as usize;
            offset += 4;

            if offset + key_len + 4 > block_data.len() {
                break;
            }

            let entry_key = &block_data[offset..offset + key_len];
            offset += key_len;

            let value_len = u32::from_le_bytes([
                block_data[offset],
                block_data[offset + 1],
                block_data[offset + 2],
                block_data[offset + 3],
            ]) as usize;
            offset += 4;

            if offset + value_len > block_data.len() {
                break;
            }

            let value = &block_data[offset..offset + value_len];
            offset += value_len;

            if entry_key == key {
                return Ok(Some(value.to_vec()));
            }
        }

        Ok(None)
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

    pub fn entry_count(&self) -> usize {
        self.index.len()
    }

    pub fn num_blocks(&self) -> u64 {
        self.footer.num_blocks
    }

    pub fn path(&self) -> &Path {
        &self.path
    }

    pub fn range_scan(&self, start: &[u8], end: &[u8]) -> io::Result<Vec<(Vec<u8>, Vec<u8>)>> {
        let mut results = Vec::new();

        for key in self.index.keys() {
            if key.as_slice() >= start && key.as_slice() <= end {
                if let Ok(Some(value)) = self.get(key) {
                    results.push((key.clone(), value));
                }
            }
        }

        Ok(results)
    }

    pub fn prefix_scan(&self, prefix: &[u8]) -> io::Result<Vec<(Vec<u8>, Vec<u8>)>> {
        let mut results = Vec::new();

        for key in self.index.keys() {
            if key.starts_with(prefix) {
                if let Ok(Some(value)) = self.get(key) {
                    results.push((key.clone(), value));
                }
            }
        }

        Ok(results)
    }

    pub fn with_metrics(mut self, metrics: std::sync::Arc<crate::metrics::StorageMetrics>) -> Self {
        self.metrics = Some(metrics);
        self
    }
}

/// SSTable iterator
pub struct SSTableIterator {
    reader: ExtendedSSTableReader,
    keys: Vec<Vec<u8>>,
    current: usize,
}

impl SSTableIterator {
    pub fn new(reader: ExtendedSSTableReader) -> Self {
        let keys = reader.keys();
        SSTableIterator {
            reader,
            keys,
            current: 0,
        }
    }

    #[allow(clippy::should_implement_trait)]
    pub fn next(&mut self) -> Option<(Vec<u8>, Vec<u8>)> {
        if self.current >= self.keys.len() {
            return None;
        }

        let key = self.keys[self.current].clone();
        self.current += 1;

        match self.reader.get(&key) {
            Ok(Some(value)) => Some((key, value)),
            _ => None,
        }
    }

    pub fn has_next(&self) -> bool {
        self.current < self.keys.len()
    }

    pub fn reset(&mut self) {
        self.current = 0;
    }

    pub fn position(&self) -> usize {
        self.current
    }

    pub fn remaining(&self) -> usize {
        self.keys.len() - self.current
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_sstable_footer_encode_decode() {
        let footer = SSTableFooter::new();
        let encoded = footer.encode();
        let decoded = SSTableFooter::decode(&encoded).unwrap();

        assert_eq!(decoded.magic, footer.magic);
        assert_eq!(decoded.version, footer.version);
        assert_eq!(decoded.num_entries, footer.num_entries);
    }

    #[test]
    fn test_sstable_footer_size() {
        assert_eq!(SSTableFooter::size(), 61);
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
    fn test_block_cache() {
        let cache = BlockCache::new(1024);
        cache.put(1, vec![1, 2, 3]);

        assert_eq!(cache.get(1), Some(vec![1, 2, 3]));
        assert_eq!(cache.get(2), None);
        assert_eq!(cache.size(), 3);
    }

    #[test]
    fn test_sstable_reader_invalid() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("nonexistent.sst");

        let result = ExtendedSSTableReader::new(path);
        assert!(result.is_err());
    }

    #[test]
    fn test_sstable_iterator() {
        // Would need a valid SSTable file
        // For now, just test the struct
    }
}
