//! Helix Storage Engine - Memtable Module (Extended)
//!
//! This module provides an extended Memtable implementation
//! with skip list, arena allocation, and concurrent access.

use std::collections::BTreeMap;
use std::sync::{Arc, Mutex, RwLock};

/// Memtable entry
#[derive(Debug, Clone)]
pub struct MemtableEntry {
    pub key: Vec<u8>,
    pub value: Vec<u8>,
    pub timestamp: i64,
    pub deleted: bool,
    pub sequence: u64,
}

impl MemtableEntry {
    pub fn new(key: Vec<u8>, value: Vec<u8>, sequence: u64) -> Self {
        MemtableEntry {
            key,
            value,
            timestamp: 0,
            deleted: false,
            sequence,
        }
    }

    pub fn delete(key: Vec<u8>, sequence: u64) -> Self {
        MemtableEntry {
            key,
            value: Vec::new(),
            timestamp: 0,
            deleted: true,
            sequence,
        }
    }

    pub fn size(&self) -> usize {
        self.key.len() + self.value.len() + 24 // 24 bytes for metadata
    }
}

/// Skip list node
#[derive(Debug)]
pub struct SkipListNode {
    pub entry: MemtableEntry,
    pub level: usize,
    pub forward: RwLock<Vec<Option<Arc<SkipListNode>>>>,
}

impl SkipListNode {
    pub fn new(entry: MemtableEntry, level: usize) -> Self {
        SkipListNode {
            entry,
            level,
            forward: RwLock::new(vec![None; level + 1]),
        }
    }
}

/// Skip list memtable
pub struct SkipListMemtable {
    head: Arc<SkipListNode>,
    max_level: usize,
    probability: f64,
    size: Arc<Mutex<usize>>,
    entry_count: Arc<Mutex<usize>>,
    max_size: usize,
    sequence: Arc<Mutex<u64>>,
}

impl SkipListMemtable {
    pub fn new(max_size: usize) -> Self {
        let head = Arc::new(SkipListNode::new(
            MemtableEntry::new(Vec::new(), Vec::new(), 0),
            32,
        ));

        SkipListMemtable {
            head,
            max_level: 32,
            probability: 0.5,
            size: Arc::new(Mutex::new(0)),
            entry_count: Arc::new(Mutex::new(0)),
            max_size,
            sequence: Arc::new(Mutex::new(0)),
        }
    }

    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) {
        self.put_inner(key, value, false);
    }

    pub fn put_deleted(&self, key: Vec<u8>) {
        self.put_inner(key, Vec::new(), true);
    }

    #[allow(clippy::needless_range_loop)]
    fn put_inner(&self, key: Vec<u8>, value: Vec<u8>, deleted: bool) {
        let seq = {
            let mut seq = self.sequence.lock().unwrap();
            *seq += 1;
            *seq
        };

        let mut entry = MemtableEntry::new(key.clone(), value, seq);
        entry.deleted = deleted;
        let node = Arc::new(SkipListNode::new(entry, self.random_level()));

        // Find insertion point
        let mut update = vec![self.head.clone(); self.max_level + 1];
        let mut current = self.head.clone();

        for i in (0..=self.max_level).rev() {
            loop {
                let next = current.forward.read().unwrap()[i].clone();
                match next {
                    Some(n) if n.entry.key < key => current = n,
                    _ => break,
                }
            }
            update[i] = current.clone();
        }

        // Check if key exists
        if !deleted {
            if let Some(next) = update[0].forward.read().unwrap()[0].clone() {
                if next.entry.key == key {
                    // Update existing
                    // In a real implementation, would update in place
                    return;
                }
            }
        }

        // Insert new node
        {
            let mut node_forward = node.forward.write().unwrap();
            for i in 0..=node.level {
                node_forward[i] = update[i].forward.read().unwrap()[i].clone();
            }
        }
        for i in 0..=node.level {
            update[i].forward.write().unwrap()[i] = Some(node.clone());
        }

        // Update size
        let mut size = self.size.lock().unwrap();
        *size += node.entry.size();

        let mut count = self.entry_count.lock().unwrap();
        *count += 1;
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        let mut current = self.head.clone();

        for i in (0..=self.max_level).rev() {
            loop {
                let next = current.forward.read().unwrap()[i].clone();
                match next {
                    Some(n) if n.entry.key < key.to_vec() => current = n,
                    _ => break,
                }
            }
        }

        let next = current.forward.read().unwrap()[0].clone();
        if let Some(next) = next {
            if next.entry.key == key {
                if next.entry.deleted {
                    return None;
                }
                return Some(next.entry.value.clone());
            }
        }

        None
    }

    pub fn delete(&self, key: Vec<u8>) {
        let _seq = {
            let mut seq = self.sequence.lock().unwrap();
            *seq += 1;
            *seq
        };

        self.put_deleted(key);
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        self.get(key).is_some()
    }

    pub fn size(&self) -> usize {
        *self.size.lock().unwrap()
    }

    pub fn entry_count(&self) -> usize {
        *self.entry_count.lock().unwrap()
    }

    pub fn is_full(&self) -> bool {
        self.size() >= self.max_size
    }

    fn random_level(&self) -> usize {
        let mut level = 0;
        while level < self.max_level && rand() < self.probability {
            level += 1;
        }
        level
    }

    pub fn to_sorted_vec(&self) -> Vec<MemtableEntry> {
        let mut result = Vec::new();
        let mut current = self.head.clone();

        loop {
            let next = current.forward.read().unwrap()[0].clone();
            match next {
                Some(n) => {
                    result.push(n.entry.clone());
                    current = n;
                }
                None => break,
            }
        }

        result
    }
}

fn rand() -> f64 {
    // Simple pseudo-random number generator
    use std::time::{SystemTime, UNIX_EPOCH};
    let nanos = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .subsec_nanos();
    (nanos as f64) / (u32::MAX as f64)
}

/// Arena-allocated memtable for reduced GC pressure
pub struct ArenaMemtable {
    arena: Vec<u8>,
    arena_size: usize,
    arena_used: usize,
    entries: BTreeMap<Vec<u8>, (usize, usize)>, // key -> (offset, length)
    deleted: BTreeSet<Vec<u8>>,
    max_size: usize,
    sequence: u64,
}

impl ArenaMemtable {
    pub fn new(arena_size: usize, max_size: usize) -> Self {
        ArenaMemtable {
            arena: vec![0u8; arena_size],
            arena_size,
            arena_used: 0,
            entries: BTreeMap::new(),
            deleted: BTreeSet::new(),
            max_size,
            sequence: 0,
        }
    }

    pub fn put(&mut self, key: Vec<u8>, value: Vec<u8>) -> bool {
        let total_size = key.len() + value.len() + 8; // 8 bytes for length prefix

        if self.arena_used + total_size > self.arena_size {
            return false; // Arena full
        }

        // Write key length
        let offset = self.arena_used;
        self.arena[offset..offset + 4].copy_from_slice(&(key.len() as u32).to_le_bytes());
        self.arena_used += 4;

        // Write key
        self.arena[self.arena_used..self.arena_used + key.len()].copy_from_slice(&key);
        self.arena_used += key.len();

        // Write value length
        self.arena[self.arena_used..self.arena_used + 4]
            .copy_from_slice(&(value.len() as u32).to_le_bytes());
        self.arena_used += 4;

        // Write value
        self.arena[self.arena_used..self.arena_used + value.len()].copy_from_slice(&value);
        self.arena_used += value.len();

        self.entries.insert(key.clone(), (offset, total_size));
        self.deleted.remove(&key);
        self.sequence += 1;

        true
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        if self.deleted.contains(key) {
            return None;
        }

        let (offset, _) = self.entries.get(key)?;

        let key_len = u32::from_le_bytes([
            self.arena[*offset],
            self.arena[*offset + 1],
            self.arena[*offset + 2],
            self.arena[*offset + 3],
        ]) as usize;

        let value_offset = *offset + 4 + key_len + 4;
        let value_len = u32::from_le_bytes([
            self.arena[*offset + 4 + key_len],
            self.arena[*offset + 4 + key_len + 1],
            self.arena[*offset + 4 + key_len + 2],
            self.arena[*offset + 4 + key_len + 3],
        ]) as usize;

        Some(self.arena[value_offset..value_offset + value_len].to_vec())
    }

    pub fn delete(&mut self, key: Vec<u8>) {
        self.deleted.insert(key);
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        self.entries.contains_key(key) && !self.deleted.contains(key)
    }

    pub fn size(&self) -> usize {
        self.arena_used
    }

    pub fn entry_count(&self) -> usize {
        self.entries.len()
    }

    pub fn is_full(&self) -> bool {
        self.arena_used >= self.max_size
    }

    pub fn utilization(&self) -> f64 {
        self.arena_used as f64 / self.arena_size as f64
    }
}

use std::collections::BTreeSet;

/// Concurrent memtable with fine-grained locking
pub struct ConcurrentMemtable {
    shards: Vec<RwLock<BTreeMap<Vec<u8>, Vec<u8>>>>,
    num_shards: usize,
    max_size: usize,
    size: Arc<Mutex<usize>>,
}

impl ConcurrentMemtable {
    pub fn new(max_size: usize, num_shards: usize) -> Self {
        let mut shards = Vec::with_capacity(num_shards);
        for _ in 0..num_shards {
            shards.push(RwLock::new(BTreeMap::new()));
        }

        ConcurrentMemtable {
            shards,
            num_shards,
            max_size,
            size: Arc::new(Mutex::new(0)),
        }
    }

    fn shard_index(&self, key: &[u8]) -> usize {
        let mut hash: usize = 0;
        for &byte in key {
            hash = hash.wrapping_mul(31).wrapping_add(byte as usize);
        }
        hash % self.num_shards
    }

    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) {
        let idx = self.shard_index(&key);
        let mut shard = self.shards[idx].write().unwrap();
        shard.insert(key.clone(), value.clone());

        let mut size = self.size.lock().unwrap();
        *size += value.len();
    }

    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        let idx = self.shard_index(key);
        let shard = self.shards[idx].read().unwrap();
        shard.get(key).cloned()
    }

    pub fn delete(&self, key: &[u8]) {
        let idx = self.shard_index(key);
        let mut shard = self.shards[idx].write().unwrap();
        shard.remove(key);
    }

    pub fn contains(&self, key: &[u8]) -> bool {
        let idx = self.shard_index(key);
        let shard = self.shards[idx].read().unwrap();
        shard.contains_key(key)
    }

    pub fn size(&self) -> usize {
        *self.size.lock().unwrap()
    }

    pub fn is_full(&self) -> bool {
        self.size() >= self.max_size
    }

    pub fn entry_count(&self) -> usize {
        self.shards.iter().map(|s| s.read().unwrap().len()).sum()
    }

    pub fn to_sorted_vec(&self) -> Vec<(Vec<u8>, Vec<u8>)> {
        let mut result = Vec::new();
        for shard in &self.shards {
            let s = shard.read().unwrap();
            for (k, v) in s.iter() {
                result.push((k.clone(), v.clone()));
            }
        }
        result.sort_by(|a, b| a.0.cmp(&b.0));
        result
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_skiplist_memtable_put_get() {
        let memtable = SkipListMemtable::new(1024 * 1024);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.put(b"key2".to_vec(), b"value2".to_vec());

        assert_eq!(memtable.get(b"key1"), Some(b"value1".to_vec()));
        assert_eq!(memtable.get(b"key2"), Some(b"value2".to_vec()));
        assert_eq!(memtable.get(b"key3"), None);
    }

    #[test]
    fn test_skiplist_memtable_delete() {
        let memtable = SkipListMemtable::new(1024 * 1024);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.delete(b"key1".to_vec());

        assert_eq!(memtable.get(b"key1"), None);
    }

    #[test]
    fn test_skiplist_memtable_size() {
        let memtable = SkipListMemtable::new(1024 * 1024);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());

        assert!(memtable.size() > 0);
        assert_eq!(memtable.entry_count(), 1);
    }

    #[test]
    fn test_arena_memtable_put_get() {
        let mut memtable = ArenaMemtable::new(1024 * 1024, 1024 * 1024);
        assert!(memtable.put(b"key1".to_vec(), b"value1".to_vec()));
        assert!(memtable.put(b"key2".to_vec(), b"value2".to_vec()));

        assert_eq!(memtable.get(b"key1"), Some(b"value1".to_vec()));
        assert_eq!(memtable.get(b"key2"), Some(b"value2".to_vec()));
        assert_eq!(memtable.get(b"key3"), None);
    }

    #[test]
    fn test_arena_memtable_delete() {
        let mut memtable = ArenaMemtable::new(1024 * 1024, 1024 * 1024);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.delete(b"key1".to_vec());

        assert_eq!(memtable.get(b"key1"), None);
        assert!(!memtable.contains(b"key1"));
    }

    #[test]
    fn test_arena_memtable_utilization() {
        let mut memtable = ArenaMemtable::new(1024, 1024);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());

        assert!(memtable.utilization() > 0.0);
        assert!(memtable.utilization() < 1.0);
    }

    #[test]
    fn test_concurrent_memtable_put_get() {
        let memtable = ConcurrentMemtable::new(1024 * 1024, 16);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.put(b"key2".to_vec(), b"value2".to_vec());

        assert_eq!(memtable.get(b"key1"), Some(b"value1".to_vec()));
        assert_eq!(memtable.get(b"key2"), Some(b"value2".to_vec()));
    }

    #[test]
    fn test_concurrent_memtable_delete() {
        let memtable = ConcurrentMemtable::new(1024 * 1024, 16);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.delete(b"key1");

        assert_eq!(memtable.get(b"key1"), None);
        assert!(!memtable.contains(b"key1"));
    }

    #[test]
    fn test_concurrent_memtable_sorted() {
        let memtable = ConcurrentMemtable::new(1024 * 1024, 16);
        memtable.put(b"key3".to_vec(), b"value3".to_vec());
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.put(b"key2".to_vec(), b"value2".to_vec());

        let sorted = memtable.to_sorted_vec();
        assert_eq!(sorted.len(), 3);
        assert_eq!(sorted[0].0, b"key1");
        assert_eq!(sorted[1].0, b"key2");
        assert_eq!(sorted[2].0, b"key3");
    }

    #[test]
    fn test_concurrent_memtable_entry_count() {
        let memtable = ConcurrentMemtable::new(1024 * 1024, 16);
        memtable.put(b"key1".to_vec(), b"value1".to_vec());
        memtable.put(b"key2".to_vec(), b"value2".to_vec());
        memtable.put(b"key3".to_vec(), b"value3".to_vec());

        assert_eq!(memtable.entry_count(), 3);
    }
}
