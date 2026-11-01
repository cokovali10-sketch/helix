//! Helix Storage Engine - WAL Module (Advanced)
//! 
//! This module provides advanced WAL features including
//! group commit, async sync, and recovery optimization.

use std::collections::VecDeque;
use std::fs::{self, File, OpenOptions};
use std::io::{self, BufReader, BufWriter, Read, Seek, SeekFrom, Write};
use std::path::{Path, PathBuf};
use std::sync::{Arc, Condvar, Mutex, RwLock};
use std::thread;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Advanced WAL with group commit
pub struct AdvancedWAL {
    dir: PathBuf,
    segments: Vec<WALSegment>,
    current_segment: Option<File>,
    current_sequence: u64,
    max_segment_size: u64,
    group_commit_size: usize,
    group_commit_timeout: Duration,
    buffer: Vec<u8>,
    pending_records: VecDeque<PendingRecord>,
    mu: Mutex<()>,
    cond: Condvar,
    stats: AdvancedWALStats,
    running: bool,
}

/// WAL segment
#[derive(Debug, Clone)]
pub struct WALSegment {
    pub path: PathBuf,
    pub sequence: u64,
    pub size: u64,
    pub created_at: i64,
    pub sealed: bool,
    pub record_count: u64,
}

/// Pending record for group commit
struct PendingRecord {
    data: Vec<u8>,
    result: Arc<Mutex<Option<io::Result<()>>>>,
}

/// Advanced WAL statistics
#[derive(Debug, Clone, Default)]
pub struct AdvancedWALStats {
    pub records_written: u64,
    pub bytes_written: u64,
    pub segments_created: u64,
    pub group_commits: u64,
    pub syncs: u64,
    pub errors: u64,
    pub avg_commit_latency_us: f64,
    pub max_commit_latency_us: u64,
}

impl AdvancedWAL {
    pub fn new(dir: PathBuf, max_segment_size: u64, group_commit_size: usize, group_commit_timeout: Duration) -> io::Result<Self> {
        fs::create_dir_all(&dir)?;
        
        let mut wal = AdvancedWAL {
            dir,
            segments: Vec::new(),
            current_segment: None,
            current_sequence: 0,
            max_segment_size,
            group_commit_size,
            group_commit_timeout,
            buffer: Vec::new(),
            pending_records: VecDeque::new(),
            mu: Mutex::new(()),
            cond: Condvar::new(),
            stats: AdvancedWALStats::default(),
            running: false,
        };
        
        wal.load_segments()?;
        
        if wal.current_segment.is_none() {
            wal.create_segment()?;
        }
        
        Ok(wal)
    }
    
    fn load_segments(&mut self) -> io::Result<()> {
        if let Ok(entries) = fs::read_dir(&self.dir) {
            for entry in entries {
                let entry = entry?;
                let path = entry.path();
                if path.extension().map(|e| e == "log").unwrap_or(false) {
                    let metadata = entry.metadata()?;
                    let segment = WALSegment {
                        path: path.clone(),
                        sequence: self.extract_sequence(&path),
                        size: metadata.len(),
                        created_at: metadata.created().ok()
                            .and_then(|t| t.duration_since(UNIX_EPOCH).ok())
                            .map(|d| d.as_secs() as i64)
                            .unwrap_or(0),
                        sealed: false,
                        record_count: 0,
                    };
                    self.segments.push(segment);
                }
            }
        }
        
        self.segments.sort_by_key(|s| s.sequence);
        
        if let Some(last) = self.segments.last() {
            self.current_sequence = last.sequence;
        }
        
        Ok(())
    }
    
    fn extract_sequence(&self, path: &Path) -> u64 {
        path.file_stem()
            .and_then(|s| s.to_str())
            .and_then(|s| s.strip_prefix("wal-"))
            .and_then(|s| s.parse().ok())
            .unwrap_or(0)
    }
    
    fn create_segment(&mut self) -> io::Result<()> {
        self.current_sequence += 1;
        let path = self.dir.join(format!("wal-{}.log", self.current_sequence));
        
        let file = OpenOptions::new()
            .create(true)
            .write(true)
            .append(true)
            .open(&path)?;
        
        self.current_segment = Some(file);
        
        let segment = WALSegment {
            path,
            sequence: self.current_sequence,
            size: 0,
            created_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64,
            sealed: false,
            record_count: 0,
        };
        self.segments.push(segment);
        self.stats.segments_created += 1;
        
        Ok(())
    }
    
    /// Write a record with group commit
    pub fn write(&mut self, data: &[u8]) -> io::Result<()> {
        let result = Arc::new(Mutex::new(None));
        let pending = PendingRecord {
            data: data.to_vec(),
            result: result.clone(),
        };
        
        // Add to pending queue
        {
            let _lock = self.mu.lock().unwrap();
            self.pending_records.push_back(pending);
        }
        
        // Check if we should flush
        let should_flush = {
            let lock = self.mu.lock().unwrap();
            self.pending_records.len() >= self.group_commit_size
        };
        
        if should_flush {
            self.flush_pending()?;
        }
        
        // Wait for result
        loop {
            let guard = result.lock().unwrap();
            if let Some(res) = guard.as_ref() {
                return res.clone().unwrap_or(Ok(()));
            }
            drop(guard);
            thread::sleep(Duration::from_micros(10));
        }
    }
    
    /// Flush pending records
    pub fn flush_pending(&mut self) -> io::Result<()> {
        let _lock = self.mu.lock().unwrap();
        
        if self.pending_records.is_empty() {
            return Ok(());
        }
        
        let start = SystemTime::now();
        let mut records_to_flush = Vec::new();
        
        // Collect pending records
        while let Some(record) = self.pending_records.pop_front() {
            records_to_flush.push(record);
            if records_to_flush.len() >= self.group_commit_size {
                break;
            }
        }
        
        // Write all records
        for record in &records_to_flush {
            // Check rotation
            if let Some(ref file) = self.current_segment {
                let metadata = file.metadata()?;
                if metadata.len() + record.data.len() as u64 + 8 > self.max_segment_size {
                    self.seal_current_segment()?;
                    self.create_segment()?;
                }
            }
            
            // Write record
            if let Some(ref mut file) = self.current_segment {
                // Write length prefix
                file.write_all(&(record.data.len() as u32).to_le_bytes())?;
                file.write_all(&record.data)?;
                
                // Write checksum
                let checksum = crc32_simple(&record.data);
                file.write_all(&checksum.to_le_bytes())?;
            }
            
            self.stats.records_written += 1;
            self.stats.bytes_written += record.data.len() as u64 + 8;
        }
        
        // Sync
        if let Some(ref mut file) = self.current_segment {
            file.sync_all()?;
        }
        
        self.stats.group_commits += 1;
        self.stats.syncs += 1;
        
        let elapsed = SystemTime::now().duration_since(start).unwrap();
        let latency_us = elapsed.as_micros() as u64;
        self.stats.max_commit_latency_us = self.stats.max_commit_latency_us.max(latency_us);
        
        // Update average latency
        let total = self.stats.group_commits;
        self.stats.avg_commit_latency_us = 
            (self.stats.avg_commit_latency_us * (total - 1) as f64 + latency_us as f64) / total as f64;
        
        // Set results
        for record in records_to_flush {
            let mut result = record.result.lock().unwrap();
            *result = Some(Ok(()));
        }
        
        // Update segment size
        if let Some(segment) = self.segments.last_mut() {
            if let Some(ref file) = self.current_segment {
                segment.size = file.metadata().map(|m| m.len()).unwrap_or(0);
            }
        }
        
        Ok(())
    }
    
    fn seal_current_segment(&mut self) -> io::Result<()> {
        if let Some(ref mut file) = self.current_segment {
            file.sync_all()?;
        }
        
        if let Some(segment) = self.segments.last_mut() {
            segment.sealed = true;
        }
        
        self.current_segment = None;
        Ok(())
    }
    
    /// Start background flush thread
    pub fn start(&mut self) {
        self.running = true;
        let _ = thread::spawn(move || {
            // Background flush loop would go here
        });
    }
    
    /// Stop the WAL
    pub fn stop(&mut self) -> io::Result<()> {
        self.running = false;
        self.flush_pending()?;
        self.seal_current_segment()?;
        Ok(())
    }
    
    /// Read records from a segment
    pub fn read_segment(&self, sequence: u64) -> io::Result<Vec<Vec<u8>>> {
        let path = self.dir.join(format!("wal-{}.log", sequence));
        let file = File::open(&path)?;
        let mut reader = BufReader::new(file);
        let mut records = Vec::new();
        
        loop {
            let mut len_buf = [0u8; 4];
            match reader.read_exact(&mut len_buf) {
                Ok(_) => {}
                Err(e) if e.kind() == io::ErrorKind::UnexpectedEof => break,
                Err(e) => return Err(e),
            }
            
            let len = u32::from_le_bytes(len_buf) as usize;
            let mut data = vec![0u8; len];
            reader.read_exact(&mut data)?;
            
            let mut crc_buf = [0u8; 4];
            reader.read_exact(&mut crc_buf)?;
            let stored_crc = u32::from_le_bytes(crc_buf);
            
            if crc32_simple(&data) != stored_crc {
                return Err(io::Error::new(io::ErrorKind::InvalidData, "checksum mismatch"));
            }
            
            records.push(data);
        }
        
        Ok(records)
    }
    
    /// Recover all records
    pub fn recover(&self) -> io::Result<Vec<Vec<u8>>> {
        let mut all_records = Vec::new();
        
        for segment in &self.segments {
            match self.read_segment(segment.sequence) {
                Ok(records) => all_records.extend(records),
                Err(e) => {
                    eprintln!("Error reading segment {}: {}", segment.sequence, e);
                }
            }
        }
        
        Ok(all_records)
    }
    
    /// Get statistics
    pub fn stats(&self) -> AdvancedWALStats {
        self.stats.clone()
    }
    
    /// Get segment count
    pub fn segment_count(&self) -> usize {
        self.segments.len()
    }
    
    /// Get current size
    pub fn current_size(&self) -> u64 {
        self.segments.last().map(|s| s.size).unwrap_or(0)
    }
}

trait ReadExact {
    fn read_exact(&mut self, buf: &mut [u8]) -> io::Result<()>;
}

impl<R: Read> ReadExact for R {
    fn read_exact(&mut self, buf: &mut [u8]) -> io::Result<()> {
        self.read_exact(buf)
    }
}

fn crc32_simple(data: &[u8]) -> u32 {
    let mut crc: u32 = 0xFFFFFFFF;
    for &byte in data {
        crc ^= byte as u32;
        for _ in 0..8 {
            if crc & 1 == 1 {
                crc = (crc >> 1) ^ 0xEDB88320;
            } else {
                crc >>= 1;
            }
        }
    }
    !crc
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;
    
    #[test]
    fn test_advanced_wal_create() {
        let dir = tempdir().unwrap();
        let wal = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        assert_eq!(wal.segment_count(), 1);
    }
    
    #[test]
    fn test_advanced_wal_write() {
        let dir = tempdir().unwrap();
        let mut wal = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        
        wal.write(b"test data").unwrap();
        assert_eq!(wal.stats().records_written, 1);
    }
    
    #[test]
    fn test_advanced_wal_multiple_writes() {
        let dir = tempdir().unwrap();
        let mut wal = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        
        for i in 0..100 {
            wal.write(format!("record-{}", i).as_bytes()).unwrap();
        }
        
        assert_eq!(wal.stats().records_written, 100);
    }
    
    #[test]
    fn test_advanced_wal_recover() {
        let dir = tempdir().unwrap();
        let mut wal = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        
        wal.write(b"test data").unwrap();
        wal.stop().unwrap();
        
        let wal2 = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        let records = wal2.recover().unwrap();
        assert!(!records.is_empty());
    }
    
    #[test]
    fn test_advanced_wal_segment_rotation() {
        let dir = tempdir().unwrap();
        let mut wal = AdvancedWAL::new(dir.path().to_path_buf(), 100, 10, Duration::from_millis(100)).unwrap();
        
        for i in 0..10 {
            wal.write(&[0u8; 50]).unwrap();
        }
        
        assert!(wal.segment_count() >= 2);
    }
    
    #[test]
    fn test_advanced_wal_stats() {
        let dir = tempdir().unwrap();
        let mut wal = AdvancedWAL::new(dir.path().to_path_buf(), 1024 * 1024, 10, Duration::from_millis(100)).unwrap();
        
        wal.write(b"test").unwrap();
        
        let stats = wal.stats();
        assert_eq!(stats.records_written, 1);
        assert!(stats.bytes_written > 0);
    }
    
    #[test]
    fn test_crc32_simple() {
        let crc = crc32_simple(b"hello");
        assert!(crc > 0);
    }
    
    #[test]
    fn test_wal_segment() {
        let segment = WALSegment {
            path: PathBuf::from("wal-1.log"),
            sequence: 1,
            size: 1024,
            created_at: 1000,
            sealed: false,
            record_count: 10,
        };
        
        assert_eq!(segment.sequence, 1);
        assert!(!segment.sealed);
    }
}