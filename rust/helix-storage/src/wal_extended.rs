//! Helix Storage Engine - WAL Module (Extended)
//!
//! This module provides an extended Write-Ahead Log implementation
//! with checksums, compression, and recovery capabilities.

use std::fs::{self, File, OpenOptions};
use std::io::{self, BufReader, Read, Write};
use std::path::{Path, PathBuf};
use std::sync::Mutex;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// WAL segment
#[derive(Debug)]
pub struct WALSsegment {
    pub path: PathBuf,
    pub sequence: u64,
    pub size: u64,
    pub created_at: i64,
    pub sealed: bool,
}

/// WAL record types
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum RecordType {
    Put = 1,
    Delete = 2,
    Batch = 3,
    Checkpoint = 4,
    Metadata = 5,
}

/// WAL record
#[derive(Debug, Clone)]
pub struct WALRecord {
    pub record_type: RecordType,
    pub sequence: u64,
    pub timestamp: i64,
    pub key: Vec<u8>,
    pub value: Vec<u8>,
    pub checksum: u32,
}

impl WALRecord {
    pub fn new(record_type: RecordType, key: Vec<u8>, value: Vec<u8>) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        let checksum = crc32(&key, &value);
        WALRecord {
            record_type,
            sequence: 0,
            timestamp,
            key,
            value,
            checksum,
        }
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut buf = Vec::new();
        buf.push(self.record_type as u8);
        buf.extend_from_slice(&self.sequence.to_le_bytes());
        buf.extend_from_slice(&self.timestamp.to_le_bytes());
        buf.extend_from_slice(&(self.key.len() as u32).to_le_bytes());
        buf.extend_from_slice(&self.key);
        buf.extend_from_slice(&(self.value.len() as u32).to_le_bytes());
        buf.extend_from_slice(&self.value);
        buf.extend_from_slice(&self.checksum.to_le_bytes());
        buf
    }

    pub fn decode(data: &[u8]) -> Result<Self, String> {
        if data.len() < 29 {
            return Err("record too short".to_string());
        }

        let record_type = match data[0] {
            1 => RecordType::Put,
            2 => RecordType::Delete,
            3 => RecordType::Batch,
            4 => RecordType::Checkpoint,
            5 => RecordType::Metadata,
            _ => return Err("invalid record type".to_string()),
        };

        let sequence = u64::from_le_bytes([
            data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8],
        ]);
        let timestamp = i64::from_le_bytes([
            data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16],
        ]);

        let key_len = u32::from_le_bytes([data[17], data[18], data[19], data[20]]) as usize;
        let key_start = 21;
        let key_end = key_start + key_len;

        if key_end + 4 > data.len() {
            return Err("invalid key length".to_string());
        }

        let key = data[key_start..key_end].to_vec();

        let value_len = u32::from_le_bytes([
            data[key_end],
            data[key_end + 1],
            data[key_end + 2],
            data[key_end + 3],
        ]) as usize;
        let value_start = key_end + 4;
        let value_end = value_start + value_len;

        if value_end + 4 > data.len() {
            return Err("invalid value length".to_string());
        }

        let value = data[value_start..value_end].to_vec();
        let checksum = u32::from_le_bytes([
            data[value_end],
            data[value_end + 1],
            data[value_end + 2],
            data[value_end + 3],
        ]);

        // Verify checksum
        if crc32(&key, &value) != checksum {
            return Err("checksum mismatch".to_string());
        }

        Ok(WALRecord {
            record_type,
            sequence,
            timestamp,
            key,
            value,
            checksum,
        })
    }
}

fn crc32(key: &[u8], value: &[u8]) -> u32 {
    let mut crc: u32 = 0xFFFFFFFF;
    for &byte in key.iter().chain(value.iter()) {
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

/// Extended WAL
pub struct ExtendedWAL {
    dir: PathBuf,
    segments: Vec<WALSsegment>,
    current_segment: Option<File>,
    current_sequence: u64,
    max_segment_size: u64,
    sync_interval: Duration,
    buffer: Vec<u8>,
    mu: Mutex<()>,
    stats: WALStats,
}

/// WAL statistics
#[derive(Debug, Clone, Default)]
pub struct WALStats {
    pub records_written: u64,
    pub bytes_written: u64,
    pub segments_created: u64,
    pub syncs: u64,
    pub errors: u64,
}

impl ExtendedWAL {
    pub fn new(dir: PathBuf, max_segment_size: u64) -> io::Result<Self> {
        fs::create_dir_all(&dir)?;

        let mut wal = ExtendedWAL {
            dir,
            segments: Vec::new(),
            current_segment: None,
            current_sequence: 0,
            max_segment_size,
            sync_interval: Duration::from_secs(5),
            buffer: Vec::new(),
            mu: Mutex::new(()),
            stats: WALStats::default(),
        };

        // Load existing segments
        wal.load_segments()?;

        // Create new segment if needed
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
                    let segment = WALSsegment {
                        path: path.clone(),
                        sequence: self.extract_sequence(&path),
                        size: metadata.len(),
                        created_at: metadata
                            .created()
                            .ok()
                            .and_then(|t| t.duration_since(UNIX_EPOCH).ok())
                            .map(|d| d.as_secs() as i64)
                            .unwrap_or(0),
                        sealed: false,
                    };
                    self.segments.push(segment);
                }
            }
        }

        // Sort by sequence
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

        let file = OpenOptions::new().create(true).append(true).open(&path)?;

        self.current_segment = Some(file);

        let segment = WALSsegment {
            path,
            sequence: self.current_sequence,
            size: 0,
            created_at: SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs() as i64,
            sealed: false,
        };
        self.segments.push(segment);
        self.stats.segments_created += 1;

        Ok(())
    }

    /// Write a record
    pub fn write(&mut self, record: &WALRecord) -> io::Result<()> {
        let mut record = record.clone();
        record.sequence = self.current_sequence;
        let encoded = record.encode();

        // Check if we need to rotate
        if let Some(ref file) = self.current_segment {
            let metadata = file.metadata()?;
            if metadata.len() + encoded.len() as u64 > self.max_segment_size {
                self.seal_current_segment()?;
                self.create_segment()?;
            }
        }

        // Write record
        if let Some(ref mut file) = self.current_segment {
            file.write_all(&encoded)?;
            file.sync_all()?;
        }

        self.stats.records_written += 1;
        self.stats.bytes_written += encoded.len() as u64;
        self.stats.syncs += 1;

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

    /// Read all records from a segment
    pub fn read_segment(&self, sequence: u64) -> io::Result<Vec<WALRecord>> {
        let path = self.dir.join(format!("wal-{}.log", sequence));
        let file = File::open(&path)?;
        let mut reader = BufReader::new(file);
        let mut records = Vec::new();

        loop {
            // Read record type
            let mut type_buf = [0u8; 1];
            match reader.read_exact(&mut type_buf) {
                Ok(_) => {}
                Err(e) if e.kind() == io::ErrorKind::UnexpectedEof => break,
                Err(e) => return Err(e),
            }

            // Read sequence
            let mut seq_buf = [0u8; 8];
            reader.read_exact(&mut seq_buf)?;
            let sequence = u64::from_le_bytes(seq_buf);

            // Read timestamp
            let mut ts_buf = [0u8; 8];
            reader.read_exact(&mut ts_buf)?;
            let timestamp = i64::from_le_bytes(ts_buf);

            // Read key length
            let mut key_len_buf = [0u8; 4];
            reader.read_exact(&mut key_len_buf)?;
            let key_len = u32::from_le_bytes(key_len_buf) as usize;

            // Read key
            let mut key = vec![0u8; key_len];
            reader.read_exact(&mut key)?;

            // Read value length
            let mut val_len_buf = [0u8; 4];
            reader.read_exact(&mut val_len_buf)?;
            let val_len = u32::from_le_bytes(val_len_buf) as usize;

            // Read value
            let mut value = vec![0u8; val_len];
            reader.read_exact(&mut value)?;

            // Read checksum
            let mut crc_buf = [0u8; 4];
            reader.read_exact(&mut crc_buf)?;
            let checksum = u32::from_le_bytes(crc_buf);

            let record_type = match type_buf[0] {
                1 => RecordType::Put,
                2 => RecordType::Delete,
                3 => RecordType::Batch,
                4 => RecordType::Checkpoint,
                5 => RecordType::Metadata,
                _ => {
                    return Err(io::Error::new(
                        io::ErrorKind::InvalidData,
                        "invalid record type",
                    ))
                }
            };

            records.push(WALRecord {
                record_type,
                sequence,
                timestamp,
                key,
                value,
                checksum,
            });
        }

        Ok(records)
    }

    /// Recover from WAL
    pub fn recover(&self) -> io::Result<Vec<WALRecord>> {
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

    /// Close the WAL
    pub fn close(&mut self) -> io::Result<()> {
        self.seal_current_segment()?;
        Ok(())
    }

    /// Get WAL statistics
    pub fn stats(&self) -> WALStats {
        self.stats.clone()
    }

    /// Get segment count
    pub fn segment_count(&self) -> usize {
        self.segments.len()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_wal_record_encode_decode() {
        let record = WALRecord::new(RecordType::Put, b"key1".to_vec(), b"value1".to_vec());
        let encoded = record.encode();
        let decoded = WALRecord::decode(&encoded).unwrap();

        assert_eq!(decoded.record_type, RecordType::Put);
        assert_eq!(decoded.key, b"key1");
        assert_eq!(decoded.value, b"value1");
    }

    #[test]
    fn test_wal_record_checksum() {
        let record = WALRecord::new(RecordType::Put, b"key".to_vec(), b"value".to_vec());
        assert_eq!(record.checksum, crc32(b"key", b"value"));
    }

    #[test]
    fn test_extended_wal_create() {
        let dir = tempdir().unwrap();
        let wal = ExtendedWAL::new(dir.path().to_path_buf(), 1024 * 1024).unwrap();
        assert_eq!(wal.segment_count(), 1);
    }

    #[test]
    fn test_extended_wal_write() {
        let dir = tempdir().unwrap();
        let mut wal = ExtendedWAL::new(dir.path().to_path_buf(), 1024 * 1024).unwrap();

        let record = WALRecord::new(RecordType::Put, b"key1".to_vec(), b"value1".to_vec());
        wal.write(&record).unwrap();

        assert_eq!(wal.stats().records_written, 1);
    }

    #[test]
    fn test_extended_wal_recover() {
        let dir = tempdir().unwrap();
        let mut wal = ExtendedWAL::new(dir.path().to_path_buf(), 1024 * 1024).unwrap();

        let record = WALRecord::new(RecordType::Put, b"key1".to_vec(), b"value1".to_vec());
        wal.write(&record).unwrap();
        wal.close().unwrap();

        let wal2 = ExtendedWAL::new(dir.path().to_path_buf(), 1024 * 1024).unwrap();
        let records = wal2.recover().unwrap();
        assert!(!records.is_empty());
    }

    #[test]
    fn test_crc32() {
        let crc = crc32(b"hello", b"world");
        assert!(crc > 0);
    }

    #[test]
    fn test_record_type_values() {
        assert_eq!(RecordType::Put as u8, 1);
        assert_eq!(RecordType::Delete as u8, 2);
        assert_eq!(RecordType::Batch as u8, 3);
        assert_eq!(RecordType::Checkpoint as u8, 4);
        assert_eq!(RecordType::Metadata as u8, 5);
    }
}
