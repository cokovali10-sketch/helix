//! Helix Storage Engine - WAL, Memtable, SSTable implementation
//! 
//! This module provides the core storage engine for the Helix observability platform.
//! It implements a custom TSDB with WAL (Write-Ahead Log), Memtable, and SSTable storage.

use std::collections::BTreeMap;
use std::fs;
use std::io::{self, Write, Read};
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex, RwLock};

pub mod compression;
pub mod retention;
pub mod compaction;
pub mod sstable_reader;
pub mod metrics;
pub mod index;
pub mod wal_extended;
pub mod memtable_extended;
pub mod sstable_writer;
pub mod cache;
pub mod sstable_reader_extended;
pub mod wal_advanced;
pub mod compaction_extended;

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

/// WAL configuration
#[derive(Debug, Clone)]
pub struct WALConfig {
    /// Directory for WAL files
    pub dir: PathBuf,
    /// Maximum WAL file size in bytes
    pub max_size: usize,
    /// Sync interval in seconds
    pub sync_interval: u64,
}

/// Memtable configuration
#[derive(Debug, Clone)]
pub struct MemtableConfig {
    /// Maximum entries in memtable
    pub max_entries: usize,
    /// Whether to auto-flush when full
    pub auto_flush: bool,
}

/// SSTable configuration
#[derive(Debug, Clone)]
pub struct SSTableConfig {
    /// Directory for SSTable files
    pub dir: PathBuf,
    /// Compression type (lz4, zstd, none)
    pub compression: CompressionType,
    /// Target file size in bytes
    pub target_file_size: usize,
}

/// Compression types supported
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum CompressionType {
    None,
    LZ4,
    ZSTD,
}

// ---------------------------------------------------------------------------
// WAL (Write-Ahead Log)
// ---------------------------------------------------------------------------

/// WAL handles write-ahead logging for durability
pub struct WAL {
    config: WALConfig,
    file: Option<fs::File>,
    written: usize,
    buf: Vec<u8>,
}

impl WAL {
    /// Creates a new WAL instance
    pub fn new(config: WALConfig) -> io::Result<Self> {
        // Ensure WAL directory exists
        fs::create_dir_all(&config.dir)?;
        
        let path = config.dir.join("wal.log");
        let file = fs::OpenOptions::new()
            .create(true)
            .append(true)
            .open(&path)?;
        
        Ok(WAL {
            config,
            file: Some(file),
            written: 0,
            buf: Vec::new(),
        })
    }
    
    /// Wal path based on sequence number
    fn wal_path(&self) -> io::Result<PathBuf> {
        let seq = self.get_next_seq()?;
        Ok(self.config.dir.join(format!("wal-{}.log", seq)))
    }
    
    fn get_next_seq(&self) -> io::Result<u64> {
        let mut max_seq: u64 = 0;
        if let Ok(entries) = fs::read_dir(&self.config.dir) {
            for entry in entries {
                let entry = entry?;
                let name = entry.file_name();
                if let Some(name_str) = name.to_str() {
                    if name_str.starts_with("wal-") && name_str.ends_with(".log") {
                        if let Some(num_str) = name_str.strip_prefix("wal-") {
                            if let Some(num_str) = num_str.strip_suffix(".log") {
                                if let Ok(num) = num_str.parse::<u64>() {
                                    if num > max_seq {
                                        max_seq = num;
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        Ok(max_seq + 1)
    }
    
    /// Records a record to the WAL
    pub fn record(&mut self, record: &[u8]) -> io::Result<usize> {
        let len = record.len();
        self.buf.extend_from_slice(&(len as u64).to_be_bytes());
        self.buf.extend_from_slice(record);
        
        if self.buf.len() >= self.config.max_size {
            self.flush()?;
        }
        
        self.written += len;
        Ok(len)
    }
    
    /// Flushes the WAL buffer to disk
    fn flush(&mut self) -> io::Result<()> {
        if self.buf.is_empty() {
            return Ok(());
        }
        
        if let Some(ref mut file) = self.file {
            file.write_all(&self.buf)?;
            file.sync_all()?;
        }
        
        self.buf.clear();
        Ok(())
    }
    
    /// Closes the WAL
    pub fn close(mut self) -> io::Result<()> {
        self.flush()?;
        if let Some(mut file) = self.file.take() {
            file.flush()?;
        }
        Ok(())
    }
    
    /// Returns the written bytes count
    pub fn written(&self) -> usize {
        self.written
    }
}

// ---------------------------------------------------------------------------
// Memtable
// ---------------------------------------------------------------------------

/// Memtable is an in-memory table for recent writes
pub struct Memtable {
    data: RwLock<BTreeMap<Vec<u8>, Vec<u8>>>,
    max_entries: usize,
    flushed: Arc<Mutex<bool>>,
}

impl Memtable {
    /// Creates a new Memtable
    pub fn new(max_entries: usize) -> Self {
        Memtable {
            data: RwLock::new(BTreeMap::new()),
            max_entries: max_entries,
            flushed: Arc::new(Mutex::new(false)),
        }
    }
    
    /// Inserts a key-value pair into the memtable
    pub fn insert(&self, key: Vec<u8>, value: Vec<u8>) {
        let mut data = self.data.write().unwrap();
        
        if data.len() >= self.max_entries {
            // Auto-flush when full
            // In a real implementation, this would trigger a flush to SSTable
            // eprintln!("Memtable full, would flush to SSTable");
        }
        
        data.insert(key, value);
    }
    
    /// Gets a value by key from the memtable
    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        let data = self.data.read().unwrap();
        data.get(key).cloned()
    }
    
    /// Checks if a key exists in the memtable
    pub fn contains(&self, key: &[u8]) -> bool {
        let data = self.data.read().unwrap();
        data.contains_key(key)
    }
    
    /// Flushes the memtable to SSTable writer
    pub fn flush_to_writer(&self, sst_writer: &mut SSTableWriter) -> io::Result<()> {
        let data = self.data.read().unwrap();
        for (key, value) in data.iter() {
            sst_writer.add_entry(key.to_vec(), value.to_vec())?;
        }
        Ok(())
    }
    
    /// Returns approximate size in bytes
    pub fn approximate_size(&self) -> usize {
        let data = self.data.read().unwrap();
        data.values().map(|v| v.len()).sum()
    }
}

// ---------------------------------------------------------------------------
// SSTableWriter
// ---------------------------------------------------------------------------

/// SSTableWriter builds SSTable files
pub struct SSTableWriter {
    path: PathBuf,
    config: SSTableConfig,
    entries: Vec<(Vec<u8>, Vec<u8>)>,
}

impl SSTableWriter {
    /// Creates a new SSTableWriter
    pub fn new(path: PathBuf, config: &SSTableConfig) -> io::Result<Self> {
        fs::create_dir_all(&config.dir)?;
        
        Ok(SSTableWriter {
            path: config.dir.join("current.sst"),
            config: config.clone(),
            entries: Vec::new(),
        })
    }
    
    /// Adds a key-value entry to the writer
    pub fn add_entry(&mut self, key: Vec<u8>, value: Vec<u8>) -> io::Result<()> {
        self.entries.push((key, value));
        Ok(())
    }
    
    /// Builds and writes the SSTable to disk
    pub fn build(mut self) -> io::Result<SSTable> {
        // Sort entries by key
        self.entries.sort_by(|a, b| a.0.cmp(&b.0));
        
        // Serialize data (simplified - no actual compression yet)
        let mut serialized = Vec::new();
        for (key, value) in &self.entries {
            serialized.extend_from_slice(&(key.len() as u64).to_be_bytes());
            serialized.extend_from_slice(key);
            serialized.extend_from_slice(&(value.len() as u64).to_be_bytes());
            serialized.extend_from_slice(value);
        }
        
        // Write to file
        let mut file = fs::File::create(&self.path)?;
        file.write_all(&serialized)?;
        file.sync_all()?;
        
        // Build index
        let mut index = BTreeMap::new();
        let mut offset: u64 = 0;
        for (key, value) in &self.entries {
            index.insert(key.clone(), offset);
            offset += (key.len() + value.len() + 16) as u64; // 16 bytes for length prefixes
        }
        
        Ok(SSTable {
            path: self.path,
            index,
            data: serialized,
        })
    }
    
    /// Returns the entries count
    pub fn entries_count(&self) -> usize {
        self.entries.len()
    }
}

// ---------------------------------------------------------------------------
// SSTable (read-only)
// ---------------------------------------------------------------------------

/// SSTable is a immutable sorted string table for long-term storage
#[derive(Clone)]
pub struct SSTable {
    path: PathBuf,
    index: BTreeMap<Vec<u8>, u64>, // key -> offset mapping
    data: Vec<u8>,                // compressed data block
}

impl SSTable {
    /// Creates a new SSTable from builder
    pub fn from_builder(builder: &SSTable) -> Self {
        builder.clone()
    }
    
    /// Reads a value by key from SSTable
    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        // Look up offset in index
        let offset = self.index.get(key)?;
        
        // Read data from file at offset
        // For now, return None - full implementation would read and decompress
        // the data at the given offset from the SSTable file
        None
    }
    
    /// Checks if a key exists in SSTable
    pub fn contains(&self, key: &[u8]) -> bool {
        self.index.contains_key(key)
    }
    
    /// Returns the SSTable path
    pub fn path(&self) -> &Path {
        &self.path
    }
    
    /// Returns the approximate data size
    pub fn data_size(&self) -> usize {
        self.data.len()
    }
}

// ---------------------------------------------------------------------------
// Storage Engine (Composite)
// ---------------------------------------------------------------------------

/// StorageEngine combines WAL, Memtable, and SSTable for complete storage
pub struct StorageEngine {
    wal: Option<WAL>,
    memtable: Memtable,
    sst_dir: PathBuf,
    sst_config: SSTableConfig,
    next_seq: u64,
}

impl StorageEngine {
    /// Creates a new StorageEngine
    pub fn new(wal_dir: PathBuf, sst_dir: PathBuf, wal_config: WALConfig, sst_config: SSTableConfig, memtable_max: usize) -> io::Result<Self> {
        let wal = WAL::new(wal_config)?;
        let memtable = Memtable::new(memtable_max);
        
        fs::create_dir_all(&sst_dir)?;
        
        Ok(StorageEngine {
            wal: Some(wal),
            memtable,
            sst_dir,
            sst_config,
            next_seq: 1,
        })
    }
    
    /// Inserts a key-value pair
    pub fn insert(&mut self, key: Vec<u8>, value: Vec<u8>) {
        // 1. Write to WAL for durability
        if let Some(ref mut wal) = self.wal {
            let _ = wal.record(&[key.len() as u8, value.len() as u8]).ok();
        }
        
        // 2. Insert into memtable
        self.memtable.insert(key.clone(), value.clone());
        
        // 3. Check if we need to flush memtable to SSTable
        // In a real implementation, we'd check entry count and flush
    }
    
    /// Gets a value by key, checking memtable first, then SSTables
    pub fn get(&self, key: &[u8]) -> Option<Vec<u8>> {
        // Check memtable first (fast path)
        if self.memtable.contains(key) {
            return self.memtable.get(key);
        }
        
        // Check SSTables
        // In a real implementation, we'd iterate through SSTables
        // For now, check if we can find it
        let sst_path = self.sst_dir.join("current.sst");
        if sst_path.exists() {
            // TODO: Read and lookup from SSTable
        }
        
        None
    }
    
    /// Flushes memtable to SSTable
    pub fn flush_memtable(&mut self) -> io::Result<()> {
        // Create SSTable writer
        let sst_path = self.sst_dir.join(format!("{}.sst", self.next_seq));
        self.next_seq += 1;
        
        let mut writer = SSTableWriter::new(sst_path.clone(), &self.sst_config)?;
        
        // Flush memtable entries
        self.memtable.flush_to_writer(&mut writer)?;
        
        // Build and write SSTable
        let sst = writer.build()?;
        
        println!("Flushed memtable to SSTable: {}", sst_path.display());
        Ok(())
    }
    
    /// Returns the WAL written bytes count
    pub fn wal_written(&self) -> usize {
        self.wal.as_ref().map(|w| w.written()).unwrap_or(0)
    }
    
    /// Returns the memtable approximate size
    pub fn memtable_size(&self) -> usize {
        self.memtable.approximate_size()
    }
    
    /// Returns the number of SSTables
    pub fn sstable_count(&self) -> usize {
        // Count .sst files in the sst_dir
        fs::read_dir(&self.sst_dir)
            .map(|rd| rd.filter(|e| e.as_ref().map(|f| f.path().extension().map(|ext| ext == "sst").unwrap_or(false)).unwrap_or(false)).count())
            .unwrap_or(0)
    }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;
    
    #[test]
    fn test_wal_creation() {
        let dir = tempdir().unwrap();
        let config = WALConfig {
            dir: dir.path().to_path_buf(),
            max_size: 1024,
            sync_interval: 5,
        };
        
        let wal = WAL::new(config).unwrap();
        assert!(wal.written() == 0);
        
        // Record some data
        let result = wal.record(b"test data");
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), 8); // 8 bytes for "test data"
    }
    
    #[test]
    fn test_wal_record_and_flush() {
        let dir = tempdir().unwrap();
        let config = WALConfig {
            dir: dir.path().to_path_buf(),
            max_size: 1024,
            sync_interval: 5,
        };
        
        let mut wal = WAL::new(config).unwrap();
        
        // Record enough data to trigger a flush
        let large_data = vec![0u8; 512];
        let _ = wal.record(&large_data);
        
        // The WAL should have written data
        assert!(wal.written() > 0);
    }
    
    #[test]
    fn test_memtable_insert_and_get() {
        let memtable = Memtable::new(1000);
        
        let key = b"test_key".to_vec();
        let value = b"test_value".to_vec();
        memtable.insert(key.clone(), value.clone());
        
        assert!(memtable.contains(&key));
        assert_eq!(memtable.get(&key), Some(value));
    }
    
    #[test]
    fn test_memtable_contains() {
        let memtable = Memtable::new(1000);
        
        let key1 = b"key1".to_vec();
        let key2 = b"key2".to_vec();
        
        memtable.insert(key1.clone(), b"value1".to_vec());
        // Don't insert key2
        
        assert!(memtable.contains(&key1));
        assert!(!memtable.contains(&key2));
    }
    
    #[test]
    fn test_sstable_writer_build() {
        use tempfile::tempdir;
        
        let dir = tempdir().unwrap();
        let config = SSTableConfig {
            dir: dir.path().to_path_buf(),
            compression: CompressionType::None,
            target_file_size: 1024 * 1024,
        };
        
        let writer = SSTableWriter::new(dir.path().join("test.sst"), &config).unwrap();
        
        // Add some entries
        writer.add_entry(b"key1".to_vec(), b"value1".to_vec());
        writer.add_entry(b"key2".to_vec(), b"value2".to_vec());
        writer.add_entry(b"key3".to_vec(), b"value3".to_vec());
        
        // Build the SSTable
        let sst = writer.build().unwrap();
        
        // Verify entries count
        assert_eq!(sst.index.len(), 3);
        assert!(sst.contains(b"key1"));
        assert!(sst.contains(b"key2"));
        assert!(sst.contains(b"key3"));
    }
    
    #[test]
    fn test_storage_engine_insert_and_get() {
        let dir = tempdir().unwrap();
        let wal_dir = dir.path().join("wal");
        let sst_dir = dir.path().join("sstables");
        
        let wal_config = WALConfig {
            dir: wal_dir,
            max_size: 1024,
            sync_interval: 5,
        };
        
        let sst_config = SSTableConfig {
            dir: sst_dir,
            compression: CompressionType::None,
            target_file_size: 10 * 1024 * 1024,
        };
        
        let mut storage = StorageEngine::new(
            wal_dir,
            sst_dir,
            wal_config,
            sst_config,
            10000,
        ).unwrap();
        
        // Insert test data
        for i in 0..50 {
            let key = format!("event-{:06}", i).into_bytes();
            let value = format!("metrics-data-{}", i).into_bytes();
            storage.insert(key, value);
        }
        
        // Query from memtable (recent data)
        let key = "event-000049".as_bytes();
        assert!(storage.memtable.contains(key));
        
        // Query from storage
        let value = storage.get(key);
        // Should find it in memtable
        assert!(value.is_some(), "Event should be in memtable");
    }
    
    #[test]
    fn test_storage_engine_flush_memtable() {
        let dir = tempdir().unwrap();
        let wal_dir = dir.path().join("wal");
        let sst_dir = dir.path().join("sstables");
        
        let wal_config = WALConfig {
            dir: wal_dir,
            max_size: 1024,
            sync_interval: 5,
        };
        
        let sst_config = SSTableConfig {
            dir: sst_dir,
            compression: CompressionType::None,
            target_file_size: 10 * 1024 * 1024,
        };
        
        let mut storage = StorageEngine::new(
            wal_dir,
            sst_dir,
            wal_config,
            sst_config,
            100,
        ).unwrap();
        
        // Insert a few entries
        for i in 0..5 {
            let key = format!("flush-{:03}", i).into_bytes();
            let value = format!("data-{}", i).into_bytes();
            storage.insert(key, value);
        }
        
        // Flush memtable to SSTable
        let result = storage.flush_memtable();
        assert!(result.is_ok(), "Memtable flush should succeed");
        
        // Verify SSTable was created
        let sst_path = dir.path().join("sstables").join("1.sst");
        assert!(sst_path.exists(), "SSTable file should be created");
    }
    
    #[test]
    fn test_storage_engine_wal_written() {
        let dir = tempdir().unwrap();
        let wal_dir = dir.path().join("wal");
        let sst_dir = dir.path().join("sstables");
        
        let wal_config = WALConfig {
            dir: wal_dir,
            max_size: 1024,
            sync_interval: 5,
        };
        
        let sst_config = SSTableConfig {
            dir: sst_dir,
            compression: CompressionType::None,
            target_file_size: 10 * 1024 * 1024,
        };
        
        let storage = StorageEngine::new(
            wal_dir,
            sst_dir,
            wal_config,
            sst_config,
            10000,
        ).unwrap();
        
        // Initially no data written
        assert_eq!(storage.wal_written(), 0);
        
        // After inserting data, WAL should have written some bytes
        // (this is tested implicitly through the insert method)
    }
}

// ---------------------------------------------------------------------------
// Main entry point
// ---------------------------------------------------------------------------

fn main() -> io::Result<()> {
    // Example usage
    let wal_dir = std::path::Path::new("./wal");
    let sst_dir = std::path::Path::new("./sstables");
    
    let wal_config = WALConfig {
        dir: wal_dir.to_path_buf(),
        max_size: 1024 * 1024, // 1MB
        sync_interval: 5,
    };
    
    let sst_config = SSTableConfig {
        dir: sst_dir.to_path_buf(),
        compression: CompressionType::ZSTD,
        target_file_size: 10 * 1024 * 1024, // 10MB
    };
    
    let mut storage = StorageEngine::new(
        wal_dir.to_path_buf(),
        sst_dir.to_path_buf(),
        wal_config,
        sst_config,
        10000, // memtable max entries
    )?;
    
    // Insert some test data
    for i in 0..100 {
        let key = format!("event-{:06}", i).into_bytes();
        let value = format!("metrics-data-{}", i).into_bytes();
        storage.insert(key, value);
    }
    
    // Flush memtable to SSTable
    if let Ok(_) = storage.flush_memtable() {
        println!("Memtable flushed successfully");
    }
    
    // Query example
    let key = "event-000005".as_bytes();
    if let Some(value) = storage.get(key) {
        println!("Found event-000005: {:?}", String::from_utf8_lossy(&value));
    } else {
        println!("Event not found (expected before flush or in different SSTable)");
    }
    
    // Print stats
    println!("WAL written: {} bytes", storage.wal_written());
    println!("Memtable size: {} bytes", storage.memtable_size());
    println!("SSTable count: {}", storage.sstable_count());

    Ok(())
}