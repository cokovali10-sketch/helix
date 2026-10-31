//! Helix Storage Engine - SSTable Reader Module
//! 
//! This module provides reading capabilities for SSTable files.
//! Includes index lookup, bloom filter, and block-based reading.

use std::collections::BTreeMap;
use std::fs;
use std::io::{self, Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};

/// SSTable reader for reading data from SSTable files
pub struct SSTableReader {
    path: PathBuf,
    index: BTreeMap<Vec<u8>, u64>,
    bloom_filter: BloomFilter,
    data_offset: u64,
    footer_size: usize,
    block_size: usize,
}

/// SSTable footer metadata
#[derive(Debug, Clone)]
pub struct SSTableFooter {
    pub index_offset: u64,
    pub index_size: u64,
    pub bloom_filter_offset: u64,
    pub bloom_filter_size: u64,
    pub compression_type: u8,
    pub num_entries: u64,
    pub checksum: u32,
}

impl SSTableFooter {
    pub fn new() -> Self {
        SSTableFooter {
            index_offset: 0,
            index_size: 0,
            bloom_filter_offset: 0,
            bloom_filter_size: 0,
            compression_type: 0,
            num_entries: 0,
            checksum: 0,
        }
    }
    
    pub fn size() -> usize {
        8 + 8 + 8 + 8 + 1 + 8 + 4 // 45 bytes
    }
}

/// Bloom filter for SSTable
pub struct BloomFilter {
    bits: Vec<u8>,
    num_hash_functions: usize,
    num_bits: usize,
}

impl BloomFilter {
    pub fn new(num_bits: usize, num_hash_functions: usize) -> Self {
        let num_bytes = (num_bits + 7) / 8;
        BloomFilter {
            bits: vec![0u8; num_bytes],
            num_hash_functions,
            num_bits,
        }
    }
    
    pub fn from_bytes(bits: Vec<u8>, num_hash_functions: usize) -> Self {
        let num_bits = bits.len() * 8;
        BloomFilter {
            bits,
            num_hash_functions,
            num_bits,
        }
    }
    
    pub fn add(&mut self, key: &[u8]) {
        for i in 0..self.num_hash_functions {
            let hash = self.hash(key, i);
            let bit_index = hash % self.num_bits;
            let byte_index = bit_index / 8;
            let bit_offset = bit_index % 8;
            self.bits[byte_index] |= 1 << bit_offset;
        }
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        for i in 0..self.num_hash_functions {
            let hash = self.hash(key, i);
            let bit_index = hash % self.num_bits;
            let byte_index = bit_index / 8;
            let bit_offset = bit_index % 8;
            if self.bits[byte_index] & (1 << bit_offset) == 0 {
                return false;
            }
        }
        true
    }
    
    fn hash(&self, key: &[u8], seed: usize) -> usize {
        // Simple hash function (in production, use murmur3 or xxhash)
        let mut hash: usize = seed.wrapping_mul(0x9e3779b9);
        for &byte in key {
            hash = hash.wrapping_mul(31).wrapping_add(byte as usize);
        }
        hash
    }
    
    pub fn size(&self) -> usize {
        self.bits.len()
    }
    
    pub fn false_positive_rate(&self, num_items: usize) -> f64 {
        let m = self.num_bits as f64;
        let k = self.num_hash_functions as f64;
        let n = num_items as f64;
        (1.0 - (-k * n / m).exp()).powf(k)
    }
}

impl SSTableReader {
    pub fn new(path: PathBuf) -> io::Result<Self> {
        let mut file = fs::File::open(&path)?;
        
        // Read footer
        let file_size = file.metadata()?.len();
        if file_size < SSTableFooter::size() as u64 {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "file too small"));
        }
        
        file.seek(SeekFrom::End(-(SSTableFooter::size() as i64)))?;
        let mut footer_buf = vec![0u8; SSTableFooter::size()];
        file.read_exact(&mut footer_buf)?;
        
        let footer = Self::parse_footer(&footer_buf)?;
        
        // Read index
        file.seek(SeekFrom::Start(footer.index_offset))?;
        let mut index_buf = vec![0u8; footer.index_size as usize];
        file.read_exact(&mut index_buf)?;
        let index = Self::parse_index(&index_buf)?;
        
        // Read bloom filter
        file.seek(SeekFrom::Start(footer.bloom_filter_offset))?;
        let mut bloom_buf = vec![0u8; footer.bloom_filter_size as usize];
        file.read_exact(&mut bloom_buf)?;
        let bloom_filter = BloomFilter::from_bytes(bloom_buf, 3);
        
        Ok(SSTableReader {
            path,
            index,
            bloom_filter,
            data_offset: 0,
            footer_size: SSTableFooter::size(),
            block_size: 4096,
        })
    }
    
    fn parse_footer(buf: &[u8]) -> io::Result<SSTableFooter> {
        if buf.len() < SSTableFooter::size() {
            return Err(io::Error::new(io::ErrorKind::InvalidData, "footer too small"));
        }
        
        let mut footer = SSTableFooter::new();
        let mut offset = 0;
        
        footer.index_offset = u64::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
            buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
        ]);
        offset += 8;
        
        footer.index_size = u64::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
            buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
        ]);
        offset += 8;
        
        footer.bloom_filter_offset = u64::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
            buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
        ]);
        offset += 8;
        
        footer.bloom_filter_size = u64::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
            buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
        ]);
        offset += 8;
        
        footer.compression_type = buf[offset];
        offset += 1;
        
        footer.num_entries = u64::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
            buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
        ]);
        offset += 8;
        
        footer.checksum = u32::from_le_bytes([
            buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
        ]);
        
        Ok(footer)
    }
    
    fn parse_index(buf: &[u8]) -> io::Result<BTreeMap<Vec<u8>, u64>> {
        let mut index = BTreeMap::new();
        let mut offset = 0;
        
        while offset < buf.len() {
            if offset + 8 > buf.len() {
                break;
            }
            
            let key_len = u64::from_le_bytes([
                buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
                buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
            ]) as usize;
            offset += 8;
            
            if offset + key_len + 8 > buf.len() {
                break;
            }
            
            let key = buf[offset..offset + key_len].to_vec();
            offset += key_len;
            
            let value_offset = u64::from_le_bytes([
                buf[offset], buf[offset+1], buf[offset+2], buf[offset+3],
                buf[offset+4], buf[offset+5], buf[offset+6], buf[offset+7],
            ]);
            offset += 8;
            
            index.insert(key, value_offset);
        }
        
        Ok(index)
    }
    
    /// Get a value by key
    pub fn get(&self, key: &[u8]) -> io::Result<Option<Vec<u8>>> {
        // Check bloom filter first
        if !self.bloom_filter.contains(key) {
            return Ok(None);
        }
        
        // Look up in index
        let offset = match self.index.get(key) {
            Some(&off) => off,
            None => return Ok(None),
        };
        
        // Read from file at offset
        let mut file = fs::File::open(&self.path)?;
        file.seek(SeekFrom::Start(offset))?;
        
        // Read key length
        let mut len_buf = [0u8; 8];
        file.read_exact(&mut len_buf)?;
        let key_len = u64::from_le_bytes(len_buf) as usize;
        
        // Skip key
        file.seek(SeekFrom::Current(key_len as i64))?;
        
        // Read value length
        file.read_exact(&mut len_buf)?;
        let value_len = u64::from_le_bytes(len_buf) as usize;
        
        // Read value
        let mut value = vec![0u8; value_len];
        file.read_exact(&mut value)?;
        
        Ok(Some(value))
    }
    
    /// Check if a key exists
    pub fn contains(&self, key: &[u8]) -> bool {
        if !self.bloom_filter.contains(key) {
            return false;
        }
        self.index.contains_key(key)
    }
    
    /// Get number of entries
    pub fn entry_count(&self) -> usize {
        self.index.len()
    }
    
    /// Get the path
    pub fn path(&self) -> &Path {
        &self.path
    }
    
    /// Iterate over all keys in range
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
    
    /// Get all keys
    pub fn keys(&self) -> Vec<Vec<u8>> {
        self.index.keys().cloned().collect()
    }
    
    /// Get size of bloom filter
    pub fn bloom_filter_size(&self) -> usize {
        self.bloom_filter.size()
    }
}

/// SSTable iterator for sequential reading
pub struct SSTableIterator {
    reader: SSTableReader,
    keys: Vec<Vec<u8>>,
    current: usize,
}

impl SSTableIterator {
    pub fn new(reader: SSTableReader) -> Self {
        let keys = reader.keys();
        SSTableIterator {
            reader,
            keys,
            current: 0,
        }
    }
    
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
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;
    
    #[test]
    fn test_bloom_filter_basic() {
        let mut bf = BloomFilter::new(1000, 3);
        
        bf.add(b"key1");
        bf.add(b"key2");
        bf.add(b"key3");
        
        assert!(bf.contains(b"key1"));
        assert!(bf.contains(b"key2"));
        assert!(bf.contains(b"key3"));
        
        // May or may not contain (false positive possible)
        // Just check it doesn't panic
        let _ = bf.contains(b"nonexistent");
    }
    
    #[test]
    fn test_bloom_filter_false_positive_rate() {
        let bf = BloomFilter::new(10000, 3);
        let fpr = bf.false_positive_rate(1000);
        
        assert!(fpr > 0.0 && fpr < 1.0);
    }
    
    #[test]
    fn test_sstable_footer_size() {
        assert_eq!(SSTableFooter::size(), 45);
    }
    
    #[test]
    fn test_sstable_footer_new() {
        let footer = SSTableFooter::new();
        assert_eq!(footer.num_entries, 0);
        assert_eq!(footer.compression_type, 0);
    }
    
    #[test]
    fn test_bloom_filter_size() {
        let bf = BloomFilter::new(1024, 3);
        assert_eq!(bf.size(), 128); // 1024 bits = 128 bytes
    }
    
    #[test]
    fn test_sstable_reader_invalid_file() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("nonexistent.sst");
        
        let result = SSTableReader::new(path);
        assert!(result.is_err());
    }
    
    #[test]
    fn test_sstable_iterator_empty() {
        // Create a minimal SSTable reader would require a valid file
        // For now, just test the iterator struct
        let keys = vec![b"key1".to_vec(), b"key2".to_vec()];
        // Can't easily test without a real file
    }
}