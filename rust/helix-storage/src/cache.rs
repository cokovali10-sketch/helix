//! Helix Storage Engine - Cache Module
//!
//! This module provides caching capabilities for the storage engine.
//! Includes LRU cache, write-through cache, and cache statistics.

use std::collections::{BTreeMap, VecDeque};
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::{Arc, Mutex, RwLock};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Cache entry
#[derive(Debug, Clone)]
pub struct CacheEntry {
    pub key: Vec<u8>,
    pub value: Vec<u8>,
    pub created_at: i64,
    pub accessed_at: i64,
    pub access_count: u64,
    pub size: usize,
    pub ttl: Option<Duration>,
}

impl CacheEntry {
    pub fn new(key: Vec<u8>, value: Vec<u8>) -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        let size = key.len() + value.len();
        CacheEntry {
            key,
            value: value.clone(),
            created_at: now,
            accessed_at: now,
            access_count: 1,
            size,
            ttl: None,
        }
    }

    pub fn with_ttl(mut self, ttl: Duration) -> Self {
        self.ttl = Some(ttl);
        self
    }

    pub fn is_expired(&self) -> bool {
        if let Some(ttl) = self.ttl {
            let now = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_millis() as i64;
            now - self.created_at > ttl.as_millis() as i64
        } else {
            false
        }
    }

    pub fn touch(&mut self) {
        self.accessed_at = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        self.access_count += 1;
    }
}

/// LRU Cache
pub struct LRUCache {
    capacity: usize,
    current_size: AtomicUsize,
    entries: RwLock<BTreeMap<Vec<u8>, CacheEntry>>,
    order: Mutex<VecDeque<Vec<u8>>>,
    stats: Mutex<CacheStats>,
}

/// Cache statistics
#[derive(Debug, Clone, Default)]
pub struct CacheStats {
    pub hits: u64,
    pub misses: u64,
    pub evictions: u64,
    pub insertions: u64,
    pub expirations: u64,
}

impl LRUCache {
    pub fn new(capacity: usize) -> Self {
        LRUCache {
            capacity,
            current_size: AtomicUsize::new(0),
            entries: RwLock::new(BTreeMap::new()),
            order: Mutex::new(VecDeque::new()),
            stats: Mutex::new(CacheStats::default()),
        }
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        let mut entries = self.entries.write().unwrap();

        if let Some(entry) = entries.get_mut(key) {
            if entry.is_expired() {
                entries.remove(key);
                self.stats.lock().unwrap().expirations += 1;
                return None;
            }
            entry.touch();
            self.stats.lock().unwrap().hits += 1;

            // Move to front of LRU
            let mut order = self.order.lock().unwrap();
            order.retain(|k| k != key);
            order.push_front(key.to_vec());

            return Some(entry.value.clone());
        }

        self.stats.lock().unwrap().misses += 1;
        None
    }

    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) {
        let mut entries = self.entries.write().unwrap();

        // Check if key already exists
        if let Some(existing) = entries.get(&key) {
            self.current_size
                .fetch_sub(existing.size, Ordering::Relaxed);
        }

        let entry = CacheEntry::new(key.clone(), value);
        let entry_size = entry.size;

        // Evict if necessary
        while self.current_size.load(Ordering::Relaxed) + entry_size > self.capacity
            && !entries.is_empty()
        {
            self.evict_lru(&mut entries);
        }

        entries.insert(key.clone(), entry);
        self.current_size.fetch_add(entry_size, Ordering::Relaxed);
        self.stats.lock().unwrap().insertions += 1;

        let mut order = self.order.lock().unwrap();
        order.retain(|k| k != &key);
        order.push_front(key);
    }

    pub fn put_with_ttl(&self, key: Vec<u8>, value: Vec<u8>, ttl: Duration) {
        let mut entries = self.entries.write().unwrap();

        if let Some(existing) = entries.get(&key) {
            self.current_size
                .fetch_sub(existing.size, Ordering::Relaxed);
        }

        let entry = CacheEntry::new(key.clone(), value).with_ttl(ttl);
        let entry_size = entry.size;

        while self.current_size.load(Ordering::Relaxed) + entry_size > self.capacity
            && !entries.is_empty()
        {
            self.evict_lru(&mut entries);
        }

        entries.insert(key.clone(), entry);
        self.current_size.fetch_add(entry_size, Ordering::Relaxed);
        self.stats.lock().unwrap().insertions += 1;

        let mut order = self.order.lock().unwrap();
        order.retain(|k| k != &key);
        order.push_front(key);
    }

    fn evict_lru(&self, entries: &mut BTreeMap<Vec<u8>, CacheEntry>) {
        let mut order = self.order.lock().unwrap();
        if let Some(key) = order.pop_back() {
            if let Some(entry) = entries.remove(&key) {
                self.current_size.fetch_sub(entry.size, Ordering::Relaxed);
                self.stats.lock().unwrap().evictions += 1;
            }
        }
    }

    pub fn delete(&self, key: &[u8]) -> bool {
        let mut entries = self.entries.write().unwrap();
        if let Some(entry) = entries.remove(key) {
            self.current_size.fetch_sub(entry.size, Ordering::Relaxed);
            let mut order = self.order.lock().unwrap();
            order.retain(|k| k != key);
            return true;
        }
        false
    }

    pub fn clear(&self) {
        let mut entries = self.entries.write().unwrap();
        entries.clear();
        self.current_size.store(0, Ordering::Relaxed);
        let mut order = self.order.lock().unwrap();
        order.clear();
    }

    pub fn size(&self) -> usize {
        self.current_size.load(Ordering::Relaxed)
    }

    pub fn capacity(&self) -> usize {
        self.capacity
    }

    pub fn len(&self) -> usize {
        let entries = self.entries.read().unwrap();
        entries.len()
    }

    pub fn is_empty(&self) -> bool {
        self.len() == 0
    }

    pub fn utilization(&self) -> f64 {
        if self.capacity > 0 {
            self.current_size.load(Ordering::Relaxed) as f64 / self.capacity as f64
        } else {
            0.0
        }
    }

    pub fn stats(&self) -> CacheStats {
        self.stats.lock().unwrap().clone()
    }

    pub fn hit_rate(&self) -> f64 {
        let stats = self.stats.lock().unwrap();
        let total = stats.hits + stats.misses;
        if total > 0 {
            stats.hits as f64 / total as f64
        } else {
            0.0
        }
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        let entries = self.entries.read().unwrap();
        if let Some(entry) = entries.get(key) {
            !entry.is_expired()
        } else {
            false
        }
    }

    pub fn keys(&self) -> Vec<Vec<u8>> {
        let entries = self.entries.read().unwrap();
        entries.keys().cloned().collect()
    }

    pub fn cleanup_expired(&self) -> usize {
        let mut entries = self.entries.write().unwrap();
        let mut order = self.order.lock().unwrap();
        let mut removed = 0;

        let expired_keys: Vec<Vec<u8>> = entries
            .iter()
            .filter(|(_, e)| e.is_expired())
            .map(|(k, _)| k.clone())
            .collect();

        for key in expired_keys {
            if let Some(entry) = entries.remove(&key) {
                self.current_size.fetch_sub(entry.size, Ordering::Relaxed);
                order.retain(|k| k != &key);
                self.stats.lock().unwrap().expirations += 1;
                removed += 1;
            }
        }

        removed
    }
}

/// Write-through cache
pub struct WriteThroughCache {
    cache: LRUCache,
    backend: Arc<dyn CacheBackend>,
}

/// Cache backend trait
pub trait CacheBackend: Send + Sync {
    fn get(&self, key: &[u8]) -> Option<Vec<u8>>;
    fn put(&self, key: Vec<u8>, value: Vec<u8>);
    fn delete(&self, key: &[u8]) -> bool;
}

impl WriteThroughCache {
    pub fn new(capacity: usize, backend: Arc<dyn CacheBackend>) -> Self {
        WriteThroughCache {
            cache: LRUCache::new(capacity),
            backend,
        }
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        // Check cache first
        if let Some(value) = self.cache.get(key) {
            return Some(value);
        }

        // Cache miss - fetch from backend
        if let Some(value) = self.backend.get(key) {
            self.cache.put(key.to_vec(), value.clone());
            return Some(value);
        }

        None
    }

    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) {
        // Write to cache
        self.cache.put(key.clone(), value.clone());

        // Write to backend
        self.backend.put(key, value);
    }

    pub fn delete(&self, key: &[u8]) -> bool {
        self.cache.delete(key);
        self.backend.delete(key)
    }
}

/// Write-back cache
pub struct WriteBackCache {
    cache: LRUCache,
    backend: Arc<dyn CacheBackend>,
    dirty: RwLock<BTreeMap<Vec<u8>, Vec<u8>>>,
    flush_interval: Duration,
}

impl WriteBackCache {
    pub fn new(capacity: usize, backend: Arc<dyn CacheBackend>, flush_interval: Duration) -> Self {
        WriteBackCache {
            cache: LRUCache::new(capacity),
            backend,
            dirty: RwLock::new(BTreeMap::new()),
            flush_interval,
        }
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        if let Some(value) = self.cache.get(key) {
            return Some(value);
        }

        if let Some(value) = self.backend.get(key) {
            self.cache.put(key.to_vec(), value.clone());
            return Some(value);
        }

        None
    }

    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) {
        self.cache.put(key.clone(), value.clone());

        let mut dirty = self.dirty.write().unwrap();
        dirty.insert(key, value);
    }

    pub fn flush(&self) {
        let mut dirty = self.dirty.write().unwrap();
        for (key, value) in dirty.iter() {
            self.backend.put(key.clone(), value.clone());
        }
        dirty.clear();
    }

    pub fn dirty_count(&self) -> usize {
        let dirty = self.dirty.read().unwrap();
        dirty.len()
    }
}

/// Cache factory
pub struct CacheFactory;

impl CacheFactory {
    pub fn create_lru(capacity: usize) -> LRUCache {
        LRUCache::new(capacity)
    }

    pub fn create_write_through(
        capacity: usize,
        backend: Arc<dyn CacheBackend>,
    ) -> WriteThroughCache {
        WriteThroughCache::new(capacity, backend)
    }

    pub fn create_write_back(
        capacity: usize,
        backend: Arc<dyn CacheBackend>,
        flush_interval: Duration,
    ) -> WriteBackCache {
        WriteBackCache::new(capacity, backend, flush_interval)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_lru_cache_basic() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());

        assert_eq!(cache.get(b"key1"), Some(b"value1".to_vec()));
        assert_eq!(cache.get(b"key2"), None);
    }

    #[test]
    fn test_lru_cache_eviction() {
        let cache = LRUCache::new(20);
        cache.put(b"key1".to_vec(), b"value1".to_vec()); // ~10 bytes
        cache.put(b"key2".to_vec(), b"value2".to_vec()); // ~10 bytes

        // Should evict key1
        cache.put(b"key3".to_vec(), b"value3".to_vec());

        assert_eq!(cache.get(b"key1"), None);
        assert!(cache.get(b"key2").is_some() || cache.get(b"key3").is_some());
    }

    #[test]
    fn test_lru_cache_ttl() {
        let cache = LRUCache::new(100);
        cache.put_with_ttl(
            b"key1".to_vec(),
            b"value1".to_vec(),
            Duration::from_millis(1),
        );

        std::thread::sleep(Duration::from_millis(10));

        assert_eq!(cache.get(b"key1"), None);
    }

    #[test]
    fn test_lru_cache_delete() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());

        assert!(cache.delete(b"key1"));
        assert_eq!(cache.get(b"key1"), None);
        assert!(!cache.delete(b"key1"));
    }

    #[test]
    fn test_lru_cache_clear() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());
        cache.put(b"key2".to_vec(), b"value2".to_vec());

        cache.clear();
        assert!(cache.is_empty());
        assert_eq!(cache.len(), 0);
    }

    #[test]
    fn test_lru_cache_stats() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());
        cache.get(b"key1");
        cache.get(b"key2");

        let stats = cache.stats();
        assert_eq!(stats.hits, 1);
        assert_eq!(stats.misses, 1);
        assert_eq!(stats.insertions, 1);
    }

    #[test]
    fn test_lru_cache_hit_rate() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());
        cache.get(b"key1");
        cache.get(b"key1");
        cache.get(b"key2");

        let rate = cache.hit_rate();
        assert!(rate > 0.6);
    }

    #[test]
    fn test_lru_cache_utilization() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());

        assert!(cache.utilization() > 0.0);
        assert!(cache.utilization() < 1.0);
    }

    #[test]
    fn test_lru_cache_cleanup_expired() {
        let cache = LRUCache::new(100);
        cache.put_with_ttl(
            b"key1".to_vec(),
            b"value1".to_vec(),
            Duration::from_millis(1),
        );

        std::thread::sleep(Duration::from_millis(10));

        let removed = cache.cleanup_expired();
        assert_eq!(removed, 1);
    }

    #[test]
    fn test_lru_cache_keys() {
        let cache = LRUCache::new(100);
        cache.put(b"key1".to_vec(), b"value1".to_vec());
        cache.put(b"key2".to_vec(), b"value2".to_vec());

        let keys = cache.keys();
        assert_eq!(keys.len(), 2);
    }

    #[test]
    fn test_cache_entry() {
        let entry = CacheEntry::new(b"key".to_vec(), b"value".to_vec());
        assert_eq!(entry.key, b"key");
        assert_eq!(entry.value, b"value");
        assert!(!entry.is_expired());

        let entry_with_ttl = entry.with_ttl(Duration::from_secs(60));
        assert!(entry_with_ttl.ttl.is_some());
    }
}
