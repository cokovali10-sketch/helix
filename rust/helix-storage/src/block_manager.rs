//! Helix Storage Engine - Block Manager Module
//!
//! This module provides block management for SSTable files.
//! Includes block allocation, caching, and prefetching.

use std::collections::{BTreeMap, VecDeque};
use std::sync::atomic::{AtomicU64, AtomicUsize, Ordering};
use std::sync::{Mutex, RwLock};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Block ID
pub type BlockId = u64;

/// Block type
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum BlockType {
    Data,
    Index,
    BloomFilter,
    Metadata,
    Footer,
}

/// Block metadata
#[derive(Debug, Clone)]
pub struct BlockMetadata {
    pub id: BlockId,
    pub block_type: BlockType,
    pub offset: u64,
    pub size: usize,
    pub compressed_size: usize,
    pub checksum: u32,
    pub created_at: i64,
    pub last_accessed: i64,
    pub access_count: u64,
}

impl BlockMetadata {
    pub fn new(id: BlockId, block_type: BlockType, offset: u64, size: usize) -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        BlockMetadata {
            id,
            block_type,
            offset,
            size,
            compressed_size: size,
            checksum: 0,
            created_at: now,
            last_accessed: now,
            access_count: 0,
        }
    }

    pub fn touch(&mut self) {
        self.last_accessed = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        self.access_count += 1;
    }

    pub fn age(&self) -> Duration {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        Duration::from_secs((now - self.created_at).max(0) as u64)
    }
}

/// Block manager
pub struct BlockManager {
    blocks: RwLock<BTreeMap<BlockId, BlockMetadata>>,
    free_blocks: Mutex<VecDeque<BlockId>>,
    next_id: Mutex<BlockId>,
    cache: BlockCache,
    config: BlockManagerConfig,
    stats: Mutex<BlockManagerStats>,
}

/// Block manager configuration
#[derive(Debug, Clone)]
pub struct BlockManagerConfig {
    pub block_size: usize,
    pub cache_size: usize,
    pub max_blocks: usize,
    pub prefetch_enabled: bool,
    pub prefetch_distance: usize,
    pub write_buffer_size: usize,
}

impl Default for BlockManagerConfig {
    fn default() -> Self {
        BlockManagerConfig {
            block_size: 4096,
            cache_size: 1024 * 1024 * 1024, // 1GB
            max_blocks: 1000000,
            prefetch_enabled: true,
            prefetch_distance: 4,
            write_buffer_size: 64 * 1024, // 64KB
        }
    }
}

/// Block manager statistics
#[derive(Debug, Clone, Default)]
pub struct BlockManagerStats {
    pub blocks_allocated: u64,
    pub blocks_freed: u64,
    pub cache_hits: u64,
    pub cache_misses: u64,
    pub prefetch_hits: u64,
    pub prefetch_misses: u64,
    pub bytes_read: u64,
    pub bytes_written: u64,
}

/// Block cache
pub struct BlockCache {
    cache: RwLock<BTreeMap<BlockId, Vec<u8>>>,
    capacity: usize,
    current_size: AtomicUsize,
    hits: AtomicU64,
    misses: AtomicU64,
    evictions: AtomicU64,
}

#[derive(Debug, Clone, Default)]
pub struct CacheStats {
    pub hits: u64,
    pub misses: u64,
    pub evictions: u64,
}

impl BlockCache {
    pub fn new(capacity: usize) -> Self {
        BlockCache {
            cache: RwLock::new(BTreeMap::new()),
            capacity,
            current_size: AtomicUsize::new(0),
            hits: AtomicU64::new(0),
            misses: AtomicU64::new(0),
            evictions: AtomicU64::new(0),
        }
    }

    pub fn get(&self, block_id: BlockId) -> Option<Vec<u8>> {
        let cache = self.cache.read().unwrap();
        if let Some(data) = cache.get(&block_id) {
            self.hits.fetch_add(1, Ordering::Relaxed);
            Some(data.clone())
        } else {
            self.misses.fetch_add(1, Ordering::Relaxed);
            None
        }
    }

    pub fn put(&self, block_id: BlockId, data: Vec<u8>) {
        let mut cache = self.cache.write().unwrap();

        // Evict if necessary
        while self.current_size.load(Ordering::Relaxed) + data.len() > self.capacity
            && !cache.is_empty()
        {
            if let Some(key) = cache.keys().next().cloned() {
                if let Some(val) = cache.remove(&key) {
                    self.current_size.fetch_sub(val.len(), Ordering::Relaxed);
                    self.evictions.fetch_add(1, Ordering::Relaxed);
                }
            }
        }

        self.current_size.fetch_add(data.len(), Ordering::Relaxed);
        cache.insert(block_id, data);
    }

    pub fn remove(&self, block_id: BlockId) {
        let mut cache = self.cache.write().unwrap();
        if let Some(data) = cache.remove(&block_id) {
            self.current_size.fetch_sub(data.len(), Ordering::Relaxed);
        }
    }

    pub fn clear(&self) {
        let mut cache = self.cache.write().unwrap();
        cache.clear();
        self.current_size.store(0, Ordering::Relaxed);
    }

    pub fn size(&self) -> usize {
        self.current_size.load(Ordering::Relaxed)
    }

    pub fn capacity(&self) -> usize {
        self.capacity
    }

    pub fn stats(&self) -> CacheStats {
        CacheStats {
            hits: self.hits.load(Ordering::Relaxed),
            misses: self.misses.load(Ordering::Relaxed),
            evictions: self.evictions.load(Ordering::Relaxed),
        }
    }
}

impl BlockManager {
    pub fn new(config: BlockManagerConfig) -> Self {
        BlockManager {
            blocks: RwLock::new(BTreeMap::new()),
            free_blocks: Mutex::new(VecDeque::new()),
            next_id: Mutex::new(1),
            cache: BlockCache::new(config.cache_size),
            config,
            stats: Mutex::new(BlockManagerStats::default()),
        }
    }

    pub fn allocate_block(&self, block_type: BlockType, size: usize) -> BlockId {
        let id = {
            let mut next_id = self.next_id.lock().unwrap();
            let id = *next_id;
            *next_id += 1;
            id
        };

        let metadata = BlockMetadata::new(id, block_type, 0, size);

        let mut blocks = self.blocks.write().unwrap();
        blocks.insert(id, metadata);

        self.stats.lock().unwrap().blocks_allocated += 1;

        id
    }

    pub fn free_block(&self, block_id: BlockId) {
        let mut blocks = self.blocks.write().unwrap();
        if blocks.remove(&block_id).is_some() {
            self.stats.lock().unwrap().blocks_freed += 1;
            let mut free = self.free_blocks.lock().unwrap();
            free.push_back(block_id);
        }
        self.cache.remove(block_id);
    }

    pub fn get_block(&self, block_id: BlockId) -> Option<BlockMetadata> {
        let blocks = self.blocks.read().unwrap();
        blocks.get(&block_id).cloned()
    }

    pub fn get_block_data(&self, block_id: BlockId) -> Option<Vec<u8>> {
        // Check cache first
        if let Some(data) = self.cache.get(block_id) {
            self.stats.lock().unwrap().cache_hits += 1;
            return Some(data);
        }

        self.stats.lock().unwrap().cache_misses += 1;
        None
    }

    pub fn put_block_data(&self, block_id: BlockId, data: Vec<u8>) {
        self.cache.put(block_id, data);
    }

    pub fn prefetch(&self, block_id: BlockId) {
        if !self.config.prefetch_enabled {
            return;
        }

        // Prefetch adjacent blocks
        for i in 1..=self.config.prefetch_distance {
            let prefetch_id = block_id + i as u64;
            if self.get_block_data(prefetch_id).is_none() {
                // Would read from disk in production
                self.stats.lock().unwrap().prefetch_misses += 1;
            } else {
                self.stats.lock().unwrap().prefetch_hits += 1;
            }
        }
    }

    pub fn block_count(&self) -> usize {
        let blocks = self.blocks.read().unwrap();
        blocks.len()
    }

    pub fn stats(&self) -> BlockManagerStats {
        self.stats.lock().unwrap().clone()
    }

    pub fn cache_stats(&self) -> CacheStats {
        self.cache.stats()
    }

    pub fn config(&self) -> &BlockManagerConfig {
        &self.config
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_block_manager_create() {
        let config = BlockManagerConfig::default();
        let manager = BlockManager::new(config);
        assert_eq!(manager.block_count(), 0);
    }

    #[test]
    fn test_block_allocate() {
        let config = BlockManagerConfig::default();
        let manager = BlockManager::new(config);

        let id = manager.allocate_block(BlockType::Data, 4096);
        assert!(id > 0);
        assert_eq!(manager.block_count(), 1);
    }

    #[test]
    fn test_block_free() {
        let config = BlockManagerConfig::default();
        let manager = BlockManager::new(config);

        let id = manager.allocate_block(BlockType::Data, 4096);
        manager.free_block(id);
        assert_eq!(manager.block_count(), 0);
    }

    #[test]
    fn test_block_cache() {
        let cache = BlockCache::new(1024);
        cache.put(1, vec![1, 2, 3]);

        assert_eq!(cache.get(1), Some(vec![1, 2, 3]));
        assert_eq!(cache.get(2), None);
    }

    #[test]
    fn test_block_cache_eviction() {
        let cache = BlockCache::new(10);
        cache.put(1, vec![1, 2, 3, 4, 5]); // 5 bytes
        cache.put(2, vec![6, 7, 8, 9, 10]); // 5 bytes
        cache.put(3, vec![11, 12, 13, 14, 15]); // 5 bytes - should evict

        // At least one should be evicted
        assert!(cache.get(1).is_none() || cache.get(2).is_none() || cache.get(3).is_none());
    }

    #[test]
    fn test_block_metadata() {
        let meta = BlockMetadata::new(1, BlockType::Data, 0, 4096);
        assert_eq!(meta.id, 1);
        assert_eq!(meta.block_type, BlockType::Data);
        assert_eq!(meta.size, 4096);
    }

    #[test]
    fn test_block_metadata_touch() {
        let mut meta = BlockMetadata::new(1, BlockType::Data, 0, 4096);
        let initial_access = meta.access_count;
        meta.touch();
        assert_eq!(meta.access_count, initial_access + 1);
    }

    #[test]
    fn test_block_manager_stats() {
        let config = BlockManagerConfig::default();
        let manager = BlockManager::new(config);

        manager.allocate_block(BlockType::Data, 4096);
        manager.allocate_block(BlockType::Index, 4096);

        let stats = manager.stats();
        assert_eq!(stats.blocks_allocated, 2);
    }

    #[test]
    fn test_block_manager_prefetch() {
        let config = BlockManagerConfig::default();
        let manager = BlockManager::new(config);

        manager.allocate_block(BlockType::Data, 4096);
        manager.prefetch(1); // Should not panic
    }

    #[test]
    fn test_block_cache_clear() {
        let cache = BlockCache::new(1024);
        cache.put(1, vec![1, 2, 3]);
        cache.clear();
        assert_eq!(cache.get(1), None);
        assert_eq!(cache.size(), 0);
    }

    #[test]
    fn test_block_cache_stats() {
        let cache = BlockCache::new(1024);
        cache.put(1, vec![1, 2, 3]);
        cache.get(1); // hit
        cache.get(2); // miss

        let stats = cache.stats();
        assert_eq!(stats.hits, 1);
        assert_eq!(stats.misses, 1);
    }
}
