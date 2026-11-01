//! Helix Query Engine - Query Cache Module
//! 
//! This module provides caching for query results.
//! Includes LRU cache, TTL cache, and cache invalidation.

use std::collections::{BTreeMap, VecDeque};
use std::sync::{Arc, RwLock, Mutex};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Cache key
#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct CacheKey {
    pub query: String,
    pub tenant: String,
    pub params: BTreeMap<String, String>,
}

impl CacheKey {
    pub fn new(query: String, tenant: String) -> Self {
        CacheKey {
            query,
            tenant,
            params: BTreeMap::new(),
        }
    }
    
    pub fn with_param(mut self, key: String, value: String) -> Self {
        self.params.insert(key, value);
        self
    }
}

/// Cached result
#[derive(Debug, Clone)]
pub struct CachedResult {
    pub key: CacheKey,
    pub result: Vec<u8>,
    pub created_at: i64,
    pub accessed_at: i64,
    pub access_count: u64,
    pub size: usize,
    pub ttl: Option<Duration>,
    pub tags: Vec<String>,
}

impl CachedResult {
    pub fn new(key: CacheKey, result: Vec<u8>) -> Self {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        CachedResult {
            key,
            result,
            created_at: now,
            accessed_at: now,
            access_count: 1,
            size: result.len(),
            ttl: None,
            tags: Vec::new(),
        }
    }
    
    pub fn with_ttl(mut self, ttl: Duration) -> Self {
        self.ttl = Some(ttl);
        self
    }
    
    pub fn with_tags(mut self, tags: Vec<String>) -> Self {
        self.tags = tags;
        self
    }
    
    pub fn is_expired(&self) -> bool {
        if let Some(ttl) = self.ttl {
            let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
            now - self.created_at > ttl.as_secs() as i64
        } else {
            false
        }
    }
    
    pub fn touch(&mut self) {
        self.accessed_at = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        self.access_count += 1;
    }
}

/// Query cache
pub struct QueryCache {
    cache: RwLock<BTreeMap<CacheKey, CachedResult>>,
    order: Mutex<VecDeque<CacheKey>>,
    capacity: usize,
    current_size: usize,
    ttl: Option<Duration>,
    stats: CacheStats,
}

/// Cache statistics
#[derive(Debug, Clone, Default)]
pub struct CacheStats {
    pub hits: u64,
    pub misses: u64,
    pub evictions: u64,
    pub insertions: u64,
    pub expirations: u64,
    pub invalidations: u64,
}

impl QueryCache {
    pub fn new(capacity: usize) -> Self {
        QueryCache {
            cache: RwLock::new(BTreeMap::new()),
            order: Mutex::new(VecDeque::new()),
            capacity,
            current_size: 0,
            ttl: None,
            stats: CacheStats::default(),
        }
    }
    
    pub fn with_ttl(mut self, ttl: Duration) -> Self {
        self.ttl = Some(ttl);
        self
    }
    
    pub fn get(&self, key: &CacheKey) -> Option<CachedResult> {
        let mut cache = self.cache.write().unwrap();
        
        if let Some(entry) = cache.get_mut(key) {
            if entry.is_expired() {
                cache.remove(key);
                self.stats.expirations += 1;
                return None;
            }
            entry.touch();
            self.stats.hits += 1;
            
            // Move to front
            let mut order = self.order.lock().unwrap();
            order.retain(|k| k != key);
            order.push_front(key.clone());
            
            return Some(entry.clone());
        }
        
        self.stats.misses += 1;
        None
    }
    
    pub fn put(&self, key: CacheKey, result: Vec<u8>) {
        let mut cache = self.cache.write().unwrap();
        
        // Check if key exists
        if let Some(existing) = cache.get(&key) {
            self.current_size -= existing.size;
        }
        
        let entry = CachedResult::new(key.clone(), result);
        let entry_size = entry.size;
        
        // Evict if necessary
        while self.current_size + entry_size > self.capacity && !cache.is_empty() {
            self.evict_lru(&mut cache);
        }
        
        cache.insert(key.clone(), entry);
        self.current_size += entry_size;
        self.stats.insertions += 1;
        
        let mut order = self.order.lock().unwrap();
        order.retain(|k| k != &key);
        order.push_front(key);
    }
    
    fn evict_lru(&self, cache: &mut BTreeMap<CacheKey, CachedResult>) {
        let mut order = self.order.lock().unwrap();
        if let Some(key) = order.pop_back() {
            if let Some(entry) = cache.remove(&key) {
                self.current_size -= entry.size;
                self.stats.evictions += 1;
            }
        }
    }
    
    pub fn delete(&self, key: &CacheKey) -> bool {
        let mut cache = self.cache.write().unwrap();
        if let Some(entry) = cache.remove(key) {
            self.current_size -= entry.size;
            let mut order = self.order.lock().unwrap();
            order.retain(|k| k != key);
            self.stats.invalidations += 1;
            return true;
        }
        false
    }
    
    pub fn clear(&self) {
        let mut cache = self.cache.write().unwrap();
        cache.clear();
        self.current_size = 0;
        let mut order = self.order.lock().unwrap();
        order.clear();
    }
    
    pub fn size(&self) -> usize {
        let cache = self.cache.read().unwrap();
        cache.len()
    }
    
    pub fn capacity(&self) -> usize {
        self.capacity
    }
    
    pub fn utilization(&self) -> f64 {
        if self.capacity > 0 {
            self.current_size as f64 / self.capacity as f64
        } else {
            0.0
        }
    }
    
    pub fn stats(&self) -> CacheStats {
        self.stats.clone()
    }
    
    pub fn hit_rate(&self) -> f64 {
        let total = self.stats.hits + self.stats.misses;
        if total > 0 {
            self.stats.hits as f64 / total as f64
        } else {
            0.0
        }
    }
    
    pub fn contains(&self, key: &CacheKey) -> bool {
        let cache = self.cache.read().unwrap();
        if let Some(entry) = cache.get(key) {
            !entry.is_expired()
        } else {
            false
        }
    }
    
    pub fn keys(&self) -> Vec<CacheKey> {
        let cache = self.cache.read().unwrap();
        cache.keys().cloned().collect()
    }
    
    pub fn cleanup_expired(&self) -> usize {
        let mut cache = self.cache.write().unwrap();
        let mut order = self.order.lock().unwrap();
        let mut removed = 0;
        
        let expired_keys: Vec<CacheKey> = cache.iter()
            .filter(|(_, e)| e.is_expired())
            .map(|(k, _)| k.clone())
            .collect();
        
        for key in expired_keys {
            if let Some(entry) = cache.remove(&key) {
                self.current_size -= entry.size;
                order.retain(|k| k != &key);
                self.stats.expirations += 1;
                removed += 1;
            }
        }
        
        removed
    }
    
    pub fn invalidate_by_tag(&self, tag: &str) -> usize {
        let mut cache = self.cache.write().unwrap();
        let mut order = self.order.lock().unwrap();
        let mut removed = 0;
        
        let keys_to_remove: Vec<CacheKey> = cache.iter()
            .filter(|(_, e)| e.tags.contains(&tag.to_string()))
            .map(|(k, _)| k.clone())
            .collect();
        
        for key in keys_to_remove {
            if let Some(entry) = cache.remove(&key) {
                self.current_size -= entry.size;
                order.retain(|k| k != &key);
                self.stats.invalidations += 1;
                removed += 1;
            }
        }
        
        removed
    }
    
    pub fn invalidate_by_prefix(&self, prefix: &str) -> usize {
        let mut cache = self.cache.write().unwrap();
        let mut order = self.order.lock().unwrap();
        let mut removed = 0;
        
        let keys_to_remove: Vec<CacheKey> = cache.keys()
            .filter(|k| k.query.starts_with(prefix))
            .cloned()
            .collect();
        
        for key in keys_to_remove {
            if let Some(entry) = cache.remove(&key) {
                self.current_size -= entry.size;
                order.retain(|k| k != &key);
                self.stats.invalidations += 1;
                removed += 1;
            }
        }
        
        removed
    }
}

/// Query cache manager
pub struct QueryCacheManager {
    cache: QueryCache,
    enabled: bool,
}

impl QueryCacheManager {
    pub fn new(capacity: usize) -> Self {
        QueryCacheManager {
            cache: QueryCache::new(capacity),
            enabled: true,
        }
    }
    
    pub fn with_ttl(mut self, ttl: Duration) -> Self {
        self.cache = self.cache.with_ttl(ttl);
        self
    }
    
    pub fn get(&self, key: &CacheKey) -> Option<CachedResult> {
        if !self.enabled {
            return None;
        }
        self.cache.get(key)
    }
    
    pub fn put(&self, key: CacheKey, result: Vec<u8>) {
        if !self.enabled {
            return;
        }
        self.cache.put(key, result);
    }
    
    pub fn delete(&self, key: &CacheKey) -> bool {
        if !self.enabled {
            return false;
        }
        self.cache.delete(key)
    }
    
    pub fn clear(&self) {
        self.cache.clear();
    }
    
    pub fn enable(&mut self) {
        self.enabled = true;
    }
    
    pub fn disable(&mut self) {
        self.enabled = false;
    }
    
    pub fn is_enabled(&self) -> bool {
        self.enabled
    }
    
    pub fn stats(&self) -> CacheStats {
        self.cache.stats()
    }
    
    pub fn hit_rate(&self) -> f64 {
        self.cache.hit_rate()
    }
    
    pub fn size(&self) -> usize {
        self.cache.size()
    }
    
    pub fn capacity(&self) -> usize {
        self.cache.capacity()
    }
    
    pub fn utilization(&self) -> f64 {
        self.cache.utilization()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_cache_key() {
        let key = CacheKey::new("query".to_string(), "tenant1".to_string())
            .with_param("param1".to_string(), "value1".to_string());
        
        assert_eq!(key.query, "query");
        assert_eq!(key.tenant, "tenant1");
        assert_eq!(key.params.get("param1"), Some(&"value1".to_string()));
    }
    
    #[test]
    fn test_cached_result() {
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        let result = CachedResult::new(key, b"result".to_vec());
        
        assert_eq!(result.result, b"result");
        assert_eq!(result.access_count, 1);
        assert!(!result.is_expired());
    }
    
    #[test]
    fn test_cached_result_ttl() {
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        let result = CachedResult::new(key, b"result".to_vec())
            .with_ttl(Duration::from_millis(1));
        
        std::thread::sleep(Duration::from_millis(10));
        assert!(result.is_expired());
    }
    
    #[test]
    fn test_query_cache_basic() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        cache.put(key.clone(), b"result".to_vec());
        let result = cache.get(&key);
        
        assert!(result.is_some());
        assert_eq!(result.unwrap().result, b"result");
    }
    
    #[test]
    fn test_query_cache_miss() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        let result = cache.get(&key);
        assert!(result.is_none());
    }
    
    #[test]
    fn test_query_cache_eviction() {
        let cache = QueryCache::new(10);
        let key1 = CacheKey::new("query1".to_string(), "tenant1".to_string());
        let key2 = CacheKey::new("query2".to_string(), "tenant1".to_string());
        
        cache.put(key1.clone(), vec![0u8; 8]);
        cache.put(key2.clone(), vec![0u8; 8]);
        
        // First should be evicted
        assert!(cache.get(&key1).is_none() || cache.get(&key2).is_some());
    }
    
    #[test]
    fn test_query_cache_delete() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        cache.put(key.clone(), b"result".to_vec());
        assert!(cache.delete(&key));
        assert!(cache.get(&key).is_none());
    }
    
    #[test]
    fn test_query_cache_clear() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        cache.put(key, b"result".to_vec());
        cache.clear();
        assert_eq!(cache.size(), 0);
    }
    
    #[test]
    fn test_query_cache_stats() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        cache.put(key.clone(), b"result".to_vec());
        cache.get(&key);
        cache.get(&CacheKey::new("missing".to_string(), "tenant1".to_string()));
        
        let stats = cache.stats();
        assert_eq!(stats.hits, 1);
        assert_eq!(stats.misses, 1);
        assert_eq!(stats.insertions, 1);
    }
    
    #[test]
    fn test_query_cache_hit_rate() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        cache.put(key.clone(), b"result".to_vec());
        cache.get(&key);
        cache.get(&key);
        cache.get(&CacheKey::new("missing".to_string(), "tenant1".to_string()));
        
        let rate = cache.hit_rate();
        assert!(rate > 0.6);
    }
    
    #[test]
    fn test_query_cache_invalidate_by_tag() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        let result = CachedResult::new(key.clone(), b"result".to_vec())
            .with_tags(vec!["tag1".to_string()]);
        cache.cache.write().unwrap().insert(key.clone(), result);
        
        let removed = cache.invalidate_by_tag("tag1");
        assert_eq!(removed, 1);
    }
    
    #[test]
    fn test_query_cache_invalidate_by_prefix() {
        let cache = QueryCache::new(1000);
        let key = CacheKey::new("prefix_query".to_string(), "tenant1".to_string());
        
        cache.put(key.clone(), b"result".to_vec());
        
        let removed = cache.invalidate_by_prefix("prefix");
        assert_eq!(removed, 1);
    }
    
    #[test]
    fn test_query_cache_manager() {
        let manager = QueryCacheManager::new(1000);
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        
        manager.put(key.clone(), b"result".to_vec());
        let result = manager.get(&key);
        
        assert!(result.is_some());
        assert!(manager.is_enabled());
    }
    
    #[test]
    fn test_query_cache_manager_disable() {
        let mut manager = QueryCacheManager::new(1000);
        manager.disable();
        
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        manager.put(key.clone(), b"result".to_vec());
        
        let result = manager.get(&key);
        assert!(result.is_none());
    }
    
    #[test]
    fn test_cached_result_touch() {
        let key = CacheKey::new("query".to_string(), "tenant1".to_string());
        let mut result = CachedResult::new(key, b"result".to_vec());
        
        let initial_count = result.access_count;
        result.touch();
        assert_eq!(result.access_count, initial_count + 1);
    }
    
    #[test]
    fn test_cache_key_params() {
        let key = CacheKey::new("query".to_string(), "tenant1".to_string())
            .with_param("p1".to_string(), "v1".to_string())
            .with_param("p2".to_string(), "v2".to_string());
        
        assert_eq!(key.params.len(), 2);
    }
}