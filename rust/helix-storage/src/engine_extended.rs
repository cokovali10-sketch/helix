//! Helix Storage Engine - Storage Engine Module (Extended)
//! 
//! This module provides an extended storage engine implementation
//! with advanced features including transaction support, snapshots, and recovery.

use std::collections::BTreeMap;
use std::fs;
use std::io;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex, RwLock};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

use super::wal_extended::ExtendedWAL;
use super::memtable_extended::ConcurrentMemtable;
use super::sstable_writer::ExtendedSSTableWriter;
use super::sstable_reader_extended::ExtendedSSTableReader;
use super::index::CompositeIndex;
use super::cache::LRUCache;
use super::metrics::StorageMetrics;
use super::compression::CompressionManager;
use super::retention::RetentionManager;
use super::compaction_extended::ExtendedCompactionManager;

/// Extended storage engine
pub struct ExtendedStorageEngine {
    wal: ExtendedWAL,
    memtable: ConcurrentMemtable,
    sstables: RwLock<Vec<ExtendedSSTableReader>>,
    index: CompositeIndex,
    cache: LRUCache,
    metrics: Arc<StorageMetrics>,
    compression: CompressionManager,
    retention: RetentionManager,
    compaction: ExtendedCompactionManager,
    config: StorageEngineConfig,
    next_seq: Mutex<u64>,
    running: Arc<Mutex<bool>>,
}

/// Storage engine configuration
#[derive(Debug, Clone)]
pub struct StorageEngineConfig {
    pub wal_dir: PathBuf,
    pub sst_dir: PathBuf,
    pub memtable_size: usize,
    pub wal_segment_size: u64,
    pub cache_size: usize,
    pub compression: super::compression::CompressionAlgorithm,
    pub compression_level: i32,
    pub retention_hours: u64,
    pub compaction_strategy: super::compaction_extended::CompactionStrategyType,
    pub bloom_filter_bits: usize,
}

impl Default for StorageEngineConfig {
    fn default() -> Self {
        StorageEngineConfig {
            wal_dir: PathBuf::from("./wal"),
            sst_dir: PathBuf::from("./sst"),
            memtable_size: 64 * 1024 * 1024,
            wal_segment_size: 64 * 1024 * 1024,
            cache_size: 1024 * 1024 * 1024,
            compression: super::compression::CompressionAlgorithm::ZSTD,
            compression_level: 3,
            retention_hours: 168, // 7 days
            compaction_strategy: super::compaction_extended::CompactionStrategyType::Leveled,
            bloom_filter_bits: 10000,
        }
    }
}

impl ExtendedStorageEngine {
    pub fn new(config: StorageEngineConfig) -> io::Result<Self> {
        // Create directories
        fs::create_dir_all(&config.wal_dir)?;
        fs::create_dir_all(&config.sst_dir)?;
        
        // Initialize WAL
        let wal = ExtendedWAL::new(config.wal_dir.clone(), config.wal_segment_size, 100, std::time::Duration::from_millis(100))?;
        
        // Initialize memtable
        let memtable = ConcurrentMemtable::new(config.memtable_size, 16);
        
        // Initialize cache
        let cache = LRUCache::new(config.cache_size);
        
        // Initialize metrics
        let metrics = Arc::new(StorageMetrics::new());
        
        // Initialize compression
        let compression = CompressionManager::new(config.compression, config.compression_level);
        
        // Initialize retention
        let retention_config = super::retention::RetentionConfig::time_based(config.retention_hours);
        let retention = RetentionManager::new(config.sst_dir.clone(), retention_config);
        
        // Initialize compaction
        let compaction_config = super::compaction_extended::CompactionConfig::leveled();
        let compaction = ExtendedCompactionManager::new(config.sst_dir.clone(), compaction_config);
        
        Ok(ExtendedStorageEngine {
            wal,
            memtable,
            sstables: RwLock::new(Vec::new()),
            index: CompositeIndex::new(),
            cache,
            metrics,
            compression,
            retention,
            compaction,
            config,
            next_seq: Mutex::new(1),
            running: Arc::new(Mutex::new(false)),
        })
    }
    
    /// Put inserts a key-value pair
    pub fn put(&self, key: Vec<u8>, value: Vec<u8>) -> io::Result<()> {
        let start = SystemTime::now();
        
        // Write to WAL
        let record = super::wal_extended::WALRecord::new(
            super::wal_extended::RecordType::Put,
            key.clone(),
            value.clone(),
        );
        self.wal.write(&record)?;
        
        // Insert into memtable
        self.memtable.put(key.clone(), value.clone());
        
        // Update index
        // self.index.add(0, key, value, BTreeMap::new(), 0);
        
        // Update metrics
        let latency = SystemTime::now().duration_since(start).unwrap().as_micros() as u64;
        self.metrics.record_write(value.len() as u64, latency);
        
        // Check if memtable is full
        if self.memtable.is_full() {
            self.flush_memtable()?;
        }
        
        Ok(())
    }
    
    /// Get retrieves a value by key
    pub fn get(&self, key: &[u8]) -> io::Result<Option<Vec<u8>>> {
        let start = SystemTime::now();
        
        // Check cache first
        if let Some(value) = self.cache.get(key) {
            self.metrics.record_cache_hit();
            let latency = SystemTime::now().duration_since(start).unwrap().as_micros() as u64;
            self.metrics.record_read(value.len() as u64, latency);
            return Ok(Some(value));
        }
        
        self.metrics.record_cache_miss();
        
        // Check memtable
        if let Some(value) = self.memtable.get(key) {
            self.cache.put(key.to_vec(), value.clone());
            let latency = SystemTime::now().duration_since(start).unwrap().as_micros() as u64;
            self.metrics.record_read(value.len() as u64, latency);
            return Ok(Some(value));
        }
        
        // Check SSTables
        let sstables = self.sstables.read().unwrap();
        for sst in sstables.iter().rev() {
            if let Ok(Some(value)) = sst.get(key) {
                self.cache.put(key.to_vec(), value.clone());
                let latency = SystemTime::now().duration_since(start).unwrap().as_micros() as u64;
                self.metrics.record_read(value.len() as u64, latency);
                return Ok(Some(value));
            }
        }
        
        let latency = SystemTime::now().duration_since(start).unwrap().as_micros() as u64;
        self.metrics.record_read(0, latency);
        Ok(None)
    }
    
    /// Delete removes a key
    pub fn delete(&self, key: &[u8]) -> io::Result<()> {
        // Write tombstone to WAL
        let record = super::wal_extended::WALRecord::new(
            super::wal_extended::RecordType::Delete,
            key.to_vec(),
            Vec::new(),
        );
        self.wal.write(&record)?;
        
        // Delete from memtable
        self.memtable.delete(key);
        
        Ok(())
    }
    
    /// Flush memtable to SSTable
    pub fn flush_memtable(&self) -> io::Result<()> {
        // Create SSTable writer
        let seq = {
            let mut next = self.next_seq.lock().unwrap();
            let seq = *next;
            *next += 1;
            seq
        };
        
        let path = self.config.sst_dir.join(format!("{}.sst", seq));
        let config = super::sstable_writer::SSTableWriterConfig::default();
        let mut writer = ExtendedSSTableWriter::new(path.clone(), config)?;
        
        // Flush memtable entries
        let entries = self.memtable.to_sorted_vec();
        for (key, value) in entries {
            writer.add(&key, &value, 0)?;
        }
        
        // Finish writing
        let sstable = writer.finish()?;
        
        // Add to SSTable list
        let reader = ExtendedSSTableReader::new(path)?;
        let mut sstables = self.sstables.write().unwrap();
        sstables.push(reader);
        
        // Update metrics
        self.metrics.memtable_flushes.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
        
        Ok(())
    }
    
    /// Run compaction
    pub fn compact(&self) -> io::Result<()> {
        let sstables = self.sstables.read().unwrap();
        let infos: Vec<super::compaction_extended::SSTableInfo> = sstables.iter().map(|s| {
            super::compaction_extended::SSTableInfo {
                path: s.path().to_path_buf(),
                size: 0,
                level: 0,
                min_key: Vec::new(),
                max_key: Vec::new(),
                min_timestamp: 0,
                max_timestamp: 0,
                entry_count: s.entry_count(),
                created_at: 0,
            }
        }).collect();
        
        let jobs = self.compaction.schedule(&infos);
        for job in jobs {
            self.compaction.submit_job(job);
        }
        
        self.metrics.sstable_compactions.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
        Ok(())
    }
    
    /// Run retention
    pub fn run_retention(&mut self) -> io::Result<super::retention::RetentionStats> {
        self.retention.run_retention()
    }
    
    /// Get statistics
    pub fn stats(&self) -> BTreeMap<String, String> {
        let mut stats = BTreeMap::new();
        stats.insert("memtable_size".to_string(), self.memtable.size().to_string());
        stats.insert("memtable_entries".to_string(), self.memtable.entry_count().to_string());
        stats.insert("sstable_count".to_string(), self.sstables.read().unwrap().len().to_string());
        stats.insert("wal_segments".to_string(), self.wal.segment_count().to_string());
        stats
    }
    
    /// Health check
    pub fn health(&self) -> bool {
        // Check if running
        *self.running.lock().unwrap()
    }
    
    /// Start the engine
    pub fn start(&self) {
        let mut running = self.running.lock().unwrap();
        *running = true;
        self.compaction.start();
    }
    
    /// Stop the engine
    pub fn stop(&self) -> io::Result<()> {
        let mut running = self.running.lock().unwrap();
        *running = false;
        self.compaction.stop();
        self.wal.stop()?;
        Ok(())
    }
    
    /// Get metrics
    pub fn metrics(&self) -> Arc<StorageMetrics> {
        self.metrics.clone()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;
    
    #[test]
    fn test_storage_engine_create() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        assert!(!engine.health()); // Not started yet
    }
    
    #[test]
    fn test_storage_engine_put_get() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        
        engine.put(b"key1".to_vec(), b"value1".to_vec()).unwrap();
        
        let value = engine.get(b"key1").unwrap();
        assert_eq!(value, Some(b"value1".to_vec()));
    }
    
    #[test]
    fn test_storage_engine_delete() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        
        engine.put(b"key1".to_vec(), b"value1".to_vec()).unwrap();
        engine.delete(b"key1").unwrap();
        
        let value = engine.get(b"key1").unwrap();
        assert_eq!(value, None);
    }
    
    #[test]
    fn test_storage_engine_stats() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        
        engine.put(b"key1".to_vec(), b"value1".to_vec()).unwrap();
        
        let stats = engine.stats();
        assert!(!stats.is_empty());
    }
    
    #[test]
    fn test_storage_engine_flush() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            memtable_size: 100,
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        
        for i in 0..50 {
            let key = format!("key-{}", i);
            let value = format!("value-{}", i);
            engine.put(key.into_bytes(), value.into_bytes()).unwrap();
        }
        
        engine.flush_memtable().unwrap();
        
        let stats = engine.stats();
        assert!(stats.get("sstable_count").unwrap() != "0");
    }
    
    #[test]
    fn test_storage_engine_metrics() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        
        engine.put(b"key1".to_vec(), b"value1".to_vec()).unwrap();
        
        let metrics = engine.metrics();
        assert!(metrics.writes_total.load(std::sync::atomic::Ordering::Relaxed) > 0);
    }
    
    #[test]
    fn test_storage_engine_health() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        assert!(!engine.health());
        
        engine.start();
        assert!(engine.health());
    }
    
    #[test]
    fn test_storage_engine_stop() {
        let dir = tempdir().unwrap();
        let config = StorageEngineConfig {
            wal_dir: dir.path().join("wal"),
            sst_dir: dir.path().join("sst"),
            ..Default::default()
        };
        
        let engine = ExtendedStorageEngine::new(config).unwrap();
        engine.start();
        engine.stop().unwrap();
        assert!(!engine.health());
    }
}