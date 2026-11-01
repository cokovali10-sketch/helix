//! Helix Storage Engine - Bloom Filter Module (Extended)
//! 
//! This module provides advanced bloom filter implementations
//! including counting bloom filter, scalable bloom filter, and cuckoo filter.

use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};
use std::sync::RwLock;

/// Standard bloom filter
pub struct BloomFilter {
    bits: Vec<u8>,
    num_hashes: usize,
    num_bits: usize,
}

impl BloomFilter {
    pub fn new(num_bits: usize, num_hashes: usize) -> Self {
        BloomFilter {
            bits: vec![0u8; (num_bits + 7) / 8],
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
        let mut hasher = DefaultHasher::new();
        seed.hash(&mut hasher);
        key.hash(&mut hasher);
        hasher.finish() as usize
    }
    
    pub fn size(&self) -> usize {
        self.bits.len()
    }
    
    pub fn clear(&mut self) {
        for byte in &mut self.bits {
            *byte = 0;
        }
    }
}

/// Counting bloom filter (supports deletion)
pub struct CountingBloomFilter {
    counters: Vec<u8>,
    num_hashes: usize,
    num_bits: usize,
}

impl CountingBloomFilter {
    pub fn new(num_bits: usize, num_hashes: usize) -> Self {
        CountingBloomFilter {
            counters: vec![0u8; num_bits],
            num_hashes,
            num_bits,
        }
    }
    
    pub fn add(&mut self, key: &[u8]) {
        for i in 0..self.num_hashes {
            let hash = self.hash(key, i);
            let idx = hash % self.num_bits;
            if self.counters[idx] < 255 {
                self.counters[idx] += 1;
            }
        }
    }
    
    pub fn remove(&mut self, key: &[u8]) {
        for i in 0..self.num_hashes {
            let hash = self.hash(key, i);
            let idx = hash % self.num_bits;
            if self.counters[idx] > 0 {
                self.counters[idx] -= 1;
            }
        }
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        for i in 0..self.num_hashes {
            let hash = self.hash(key, i);
            let idx = hash % self.num_bits;
            if self.counters[idx] == 0 {
                return false;
            }
        }
        true
    }
    
    fn hash(&self, key: &[u8], seed: usize) -> usize {
        let mut hasher = DefaultHasher::new();
        seed.hash(&mut hasher);
        key.hash(&mut hasher);
        hasher.finish() as usize
    }
    
    pub fn size(&self) -> usize {
        self.counters.len()
    }
}

/// Scalable bloom filter (grows as needed)
pub struct ScalableBloomFilter {
    filters: Vec<BloomFilter>,
    num_hashes: usize,
    growth_factor: f64,
    max_false_positive: f64,
    current_false_positive: f64,
}

impl ScalableBloomFilter {
    pub fn new(initial_capacity: usize, num_hashes: usize, growth_factor: f64, max_false_positive: f64) -> Self {
        let first_filter = BloomFilter::new(initial_capacity, num_hashes);
        ScalableBloomFilter {
            filters: vec![first_filter],
            num_hashes,
            growth_factor,
            max_false_positive,
            current_false_positive: max_false_positive / growth_factor,
        }
    }
    
    pub fn add(&mut self, key: &[u8]) {
        // Check if we need to add a new filter
        let last = self.filters.last().unwrap();
        // Simplified: always add to last filter
        // In production, would check fill ratio
        
        let last_filter = self.filters.last_mut().unwrap();
        last_filter.add(key);
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        for filter in &self.filters {
            if filter.contains(key) {
                return true;
            }
        }
        false
    }
    
    pub fn filter_count(&self) -> usize {
        self.filters.len()
    }
}

/// Cuckoo filter
pub struct CuckooFilter {
    buckets: Vec<CuckooBucket>,
    num_buckets: usize,
    bucket_size: usize,
    max_kicks: usize,
}

#[derive(Debug, Clone)]
pub struct CuckooBucket {
    fingerprints: Vec<u8>,
}

impl CuckooBucket {
    pub fn new(size: usize) -> Self {
        CuckooBucket {
            fingerprints: vec![0u8; size],
        }
    }
}

impl CuckooFilter {
    pub fn new(num_buckets: usize, bucket_size: usize) -> Self {
        let mut buckets = Vec::with_capacity(num_buckets);
        for _ in 0..num_buckets {
            buckets.push(CuckooBucket::new(bucket_size));
        }
        
        CuckooFilter {
            buckets,
            num_buckets,
            bucket_size,
            max_kicks: 500,
        }
    }
    
    pub fn add(&mut self, key: &[u8]) -> bool {
        let (fp, i1, i2) = self.hash(key);
        
        // Try bucket 1
        if self.insert_into_bucket(i1, fp) {
            return true;
        }
        
        // Try bucket 2
        if self.insert_into_bucket(i2, fp) {
            return true;
        }
        
        // Kick out
        let mut current_bucket = i1;
        let mut current_fp = fp;
        
        for _ in 0..self.max_kicks {
            let idx = self.random_index(current_bucket);
            let kicked_fp = self.buckets[current_bucket].fingerprints[idx];
            self.buckets[current_bucket].fingerprints[idx] = current_fp;
            
            current_fp = kicked_fp;
            current_bucket = current_bucket ^ self.hash_index(current_fp);
            
            if self.insert_into_bucket(current_bucket, current_fp) {
                return true;
            }
        }
        
        false
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        let (fp, i1, i2) = self.hash(key);
        
        self.buckets[i1].fingerprints.contains(&fp) ||
        self.buckets[i2].fingerprints.contains(&fp)
    }
    
    pub fn remove(&mut self, key: &[u8]) -> bool {
        let (fp, i1, i2) = self.hash(key);
        
        if let Some(idx) = self.buckets[i1].fingerprints.iter().position(|&f| f == fp) {
            self.buckets[i1].fingerprints[idx] = 0;
            return true;
        }
        
        if let Some(idx) = self.buckets[i2].fingerprints.iter().position(|&f| f == fp) {
            self.buckets[i2].fingerprints[idx] = 0;
            return true;
        }
        
        false
    }
    
    fn hash(&self, key: &[u8]) -> (u8, usize, usize) {
        let mut hasher = DefaultHasher::new();
        key.hash(&mut hasher);
        let hash = hasher.finish();
        
        let fp = (hash & 0xFF) as u8;
        let i1 = (hash >> 8) as usize % self.num_buckets;
        let i2 = i1 ^ self.hash_index(fp);
        
        (fp, i1, i2)
    }
    
    fn hash_index(&self, fp: u8) -> usize {
        (fp as usize).wrapping_mul(0x5bd1e995) % self.num_buckets
    }
    
    fn insert_into_bucket(&mut self, bucket: usize, fp: u8) -> bool {
        for slot in &mut self.buckets[bucket].fingerprints {
            if *slot == 0 {
                *slot = fp;
                return true;
            }
        }
        false
    }
    
    fn random_index(&self, bucket: usize) -> usize {
        // Simple hash-based index selection
        bucket % self.bucket_size
    }
}

/// Quotient filter
pub struct QuotientFilter {
    table: Vec<u64>,
    num_quotients: usize,
    remainder_bits: usize,
}

impl QuotientFilter {
    pub fn new(num_quotients: usize, remainder_bits: usize) -> Self {
        QuotientFilter {
            table: vec![0u64; num_quotients],
            num_quotients,
            remainder_bits,
        }
    }
    
    pub fn insert(&mut self, key: &[u8]) -> bool {
        let (q, r) = self.hash(key);
        // Simplified insertion
        if self.table[q] == 0 {
            self.table[q] = r;
            true
        } else {
            false
        }
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        let (q, r) = self.hash(key);
        self.table[q] == r
    }
    
    fn hash(&self, key: &[u8]) -> (usize, u64) {
        let mut hasher = DefaultHasher::new();
        key.hash(&mut hasher);
        let hash = hasher.finish();
        
        let q = (hash as usize) % self.num_quotients;
        let r = hash >> (64 - self.remainder_bits);
        
        (q, r)
    }
}

/// Thread-safe bloom filter
pub struct ThreadSafeBloomFilter {
    filter: RwLock<BloomFilter>,
}

impl ThreadSafeBloomFilter {
    pub fn new(num_bits: usize, num_hashes: usize) -> Self {
        ThreadSafeBloomFilter {
            filter: RwLock::new(BloomFilter::new(num_bits, num_hashes)),
        }
    }
    
    pub fn add(&self, key: &[u8]) {
        let mut filter = self.filter.write().unwrap();
        filter.add(key);
    }
    
    pub fn contains(&self, key: &[u8]) -> bool {
        let filter = self.filter.read().unwrap();
        filter.contains(key)
    }
    
    pub fn clear(&self) {
        let mut filter = self.filter.write().unwrap();
        filter.clear();
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_bloom_filter_basic() {
        let mut bf = BloomFilter::new(1000, 3);
        bf.add(b"key1");
        bf.add(b"key2");
        
        assert!(bf.contains(b"key1"));
        assert!(bf.contains(b"key2"));
    }
    
    #[test]
    fn test_counting_bloom_filter() {
        let mut cbf = CountingBloomFilter::new(1000, 3);
        cbf.add(b"key1");
        assert!(cbf.contains(b"key1"));
        
        cbf.remove(b"key1");
        assert!(!cbf.contains(b"key1"));
    }
    
    #[test]
    fn test_scalable_bloom_filter() {
        let mut sbf = ScalableBloomFilter::new(100, 3, 2.0, 0.01);
        sbf.add(b"key1");
        sbf.add(b"key2");
        
        assert!(sbf.contains(b"key1"));
        assert!(sbf.contains(b"key2"));
    }
    
    #[test]
    fn test_cuckoo_filter() {
        let mut cf = CuckooFilter::new(1000, 4);
        assert!(cf.add(b"key1"));
        assert!(cf.add(b"key2"));
        
        assert!(cf.contains(b"key1"));
        assert!(cf.contains(b"key2"));
        
        assert!(cf.remove(b"key1"));
        assert!(!cf.contains(b"key1"));
    }
    
    #[test]
    fn test_quotient_filter() {
        let mut qf = QuotientFilter::new(1000, 8);
        assert!(qf.insert(b"key1"));
        assert!(qf.contains(b"key1"));
    }
    
    #[test]
    fn test_thread_safe_bloom_filter() {
        let bf = ThreadSafeBloomFilter::new(1000, 3);
        bf.add(b"key1");
        assert!(bf.contains(b"key1"));
        bf.clear();
        assert!(!bf.contains(b"key1"));
    }
    
    #[test]
    fn test_bloom_filter_size() {
        let bf = BloomFilter::new(1000, 3);
        assert_eq!(bf.size(), 125); // 1000 bits = 125 bytes
    }
    
    #[test]
    fn test_counting_bloom_filter_size() {
        let cbf = CountingBloomFilter::new(1000, 3);
        assert_eq!(cbf.size(), 1000);
    }
    
    #[test]
    fn test_cuckoo_filter_remove() {
        let mut cf = CuckooFilter::new(1000, 4);
        cf.add(b"key1");
        assert!(cf.remove(b"key1"));
        assert!(!cf.remove(b"key1")); // Already removed
    }
}