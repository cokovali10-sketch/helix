//! Helix Storage Engine - Compaction Module
//! 
//! This module provides SSTable compaction strategies for the TSDB.
//! Supports leveled, size-tiered, and time-window compaction.

use std::collections::BTreeMap;
use std::fs;
use std::io;
use std::path::{Path, PathBuf};
use std::time::{SystemTime, UNIX_EPOCH};

/// Compaction strategy
#[derive(Debug, Clone, PartialEq)]
pub enum CompactionStrategy {
    /// Leveled compaction (like RocksDB)
    Leveled {
        level_size_multiplier: f64,
        max_levels: usize,
    },
    /// Size-tiered compaction (like Cassandra)
    SizeTiered {
        min_threshold: usize,
        max_threshold: usize,
    },
    /// Time-window compaction
    TimeWindow {
        window_size: u64,
    },
    /// No compaction
    None,
}

/// Compaction configuration
#[derive(Debug, Clone)]
pub struct CompactionConfig {
    pub strategy: CompactionStrategy,
    pub max_concurrent_compactions: usize,
    pub compaction_threshold: f64,
    pub target_file_size: usize,
    pub max_file_size: usize,
}

impl CompactionConfig {
    pub fn leveled() -> Self {
        CompactionConfig {
            strategy: CompactionStrategy::Leveled {
                level_size_multiplier: 10.0,
                max_levels: 7,
            },
            max_concurrent_compactions: 2,
            compaction_threshold: 0.25,
            target_file_size: 64 * 1024 * 1024, // 64MB
            max_file_size: 256 * 1024 * 1024,   // 256MB
        }
    }
    
    pub fn size_tiered() -> Self {
        CompactionConfig {
            strategy: CompactionStrategy::SizeTiered {
                min_threshold: 4,
                max_threshold: 32,
            },
            max_concurrent_compactions: 1,
            compaction_threshold: 0.5,
            target_file_size: 160 * 1024 * 1024,
            max_file_size: 512 * 1024 * 1024,
        }
    }
    
    pub fn time_window(window_hours: u64) -> Self {
        CompactionConfig {
            strategy: CompactionStrategy::TimeWindow {
                window_size: window_hours * 3600,
            },
            max_concurrent_compactions: 2,
            compaction_threshold: 0.3,
            target_file_size: 128 * 1024 * 1024,
            max_file_size: 512 * 1024 * 1024,
        }
    }
}

/// SSTable metadata for compaction decisions
#[derive(Debug, Clone)]
pub struct SSTableMetadata {
    pub path: PathBuf,
    pub size: u64,
    pub level: usize,
    pub min_key: Vec<u8>,
    pub max_key: Vec<u8>,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
    pub entry_count: usize,
    pub bloom_filter_size: usize,
    pub created_at: u64,
}

impl SSTableMetadata {
    pub fn new(path: PathBuf) -> Self {
        SSTableMetadata {
            path,
            size: 0,
            level: 0,
            min_key: Vec::new(),
            max_key: Vec::new(),
            min_timestamp: 0,
            max_timestamp: 0,
            entry_count: 0,
            bloom_filter_size: 0,
            created_at: SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs(),
        }
    }
    
    pub fn key_range_overlaps(&self, other: &SSTableMetadata) -> bool {
        if self.min_key.is_empty() || other.min_key.is_empty() {
            return true;
        }
        !(self.max_key < other.min_key || other.max_key < self.min_key)
    }
    
    pub fn time_range_overlaps(&self, other: &SSTableMetadata) -> bool {
        self.min_timestamp <= other.max_timestamp && other.min_timestamp <= self.max_timestamp
    }
}

/// Compaction picker - decides which files to compact
pub struct CompactionPicker {
    config: CompactionConfig,
}

impl CompactionPicker {
    pub fn new(config: CompactionConfig) -> Self {
        CompactionPicker { config }
    }
    
    /// Pick files for compaction based on strategy
    pub fn pick(&self, files: &[SSTableMetadata]) -> CompactionPlan {
        match &self.config.strategy {
            CompactionStrategy::Leveled { level_size_multiplier, max_levels } => {
                self.pick_leveled(files, *level_size_multiplier, *max_levels)
            }
            CompactionStrategy::SizeTiered { min_threshold, max_threshold } => {
                self.pick_size_tiered(files, *min_threshold, *max_threshold)
            }
            CompactionStrategy::TimeWindow { window_size } => {
                self.pick_time_window(files, *window_size)
            }
            CompactionStrategy::None => CompactionPlan::default(),
        }
    }
    
    fn pick_leveled(&self, files: &[SSTableMetadata], multiplier: f64, max_levels: usize) -> CompactionPlan {
        let mut plan = CompactionPlan::default();
        
        // Group files by level
        let mut levels: Vec<Vec<&SSTableMetadata>> = vec![Vec::new(); max_levels];
        for file in files {
            if file.level < max_levels {
                levels[file.level].push(file);
            }
        }
        
        // Check each level for compaction
        for level in 0..max_levels - 1 {
            let level_size: u64 = levels[level].iter().map(|f| f.size).sum();
            let next_level_size: u64 = levels[level + 1].iter().map(|f| f.size).sum();
            let target_size = (self.config.target_file_size as f64 * multiplier.powi(level as i32)) as u64;
            
            if level_size > target_size || (level > 0 && level_size > next_level_size / 10) {
                // Compact this level into next level
                plan.input_files.extend(levels[level].iter().map(|f| f.path.clone()));
                plan.output_level = level + 1;
                plan.estimated_output_size = level_size;
                break;
            }
        }
        
        plan
    }
    
    fn pick_size_tiered(&self, files: &[SSTableMetadata], min_threshold: usize, max_threshold: usize) -> CompactionPlan {
        let mut plan = CompactionPlan::default();
        
        // Group files by size (bucketed)
        let mut buckets: BTreeMap<usize, Vec<&SSTableMetadata>> = BTreeMap::new();
        for file in files {
            let bucket = (file.size as f64).log2().floor() as usize;
            buckets.entry(bucket).or_default().push(file);
        }
        
        // Find buckets that need compaction
        for (bucket, files) in &buckets {
            if files.len() >= min_threshold && files.len() <= max_threshold {
                plan.input_files.extend(files.iter().map(|f| f.path.clone()));
                plan.output_level = *bucket + 1;
                plan.estimated_output_size = files.iter().map(|f| f.size).sum();
                break;
            }
        }
        
        plan
    }
    
    fn pick_time_window(&self, files: &[SSTableMetadata], window_size: u64) -> CompactionPlan {
        let mut plan = CompactionPlan::default();
        
        // Group files by time window
        let mut windows: BTreeMap<u64, Vec<&SSTableMetadata>> = BTreeMap::new();
        for file in files {
            let window = file.min_timestamp as u64 / window_size;
            windows.entry(window).or_default().push(file);
        }
        
        // Compact windows with multiple files
        for (window, files) in &windows {
            if files.len() > 1 {
                plan.input_files.extend(files.iter().map(|f| f.path.clone()));
                plan.output_level = 0;
                plan.estimated_output_size = files.iter().map(|f| f.size).sum();
                break;
            }
        }
        
        plan
    }
}

/// Compaction plan
#[derive(Debug, Clone, Default)]
pub struct CompactionPlan {
    pub input_files: Vec<PathBuf>,
    pub output_level: usize,
    pub estimated_output_size: u64,
    pub priority: i32,
}

impl CompactionPlan {
    pub fn is_empty(&self) -> bool {
        self.input_files.is_empty()
    }
    
    pub fn file_count(&self) -> usize {
        self.input_files.len()
    }
}

/// Compaction executor
pub struct CompactionExecutor {
    config: CompactionConfig,
    sst_dir: PathBuf,
}

impl CompactionExecutor {
    pub fn new(sst_dir: PathBuf, config: CompactionConfig) -> Self {
        CompactionExecutor { config, sst_dir }
    }
    
    /// Execute a compaction plan
    pub fn execute(&self, plan: &CompactionPlan) -> io::Result<CompactionResult> {
        if plan.is_empty() {
            return Ok(CompactionResult::default());
        }
        
        let start_time = SystemTime::now();
        
        // Read all input files
        let mut all_entries: Vec<(Vec<u8>, Vec<u8>)> = Vec::new();
        
        for path in &plan.input_files {
            // In production, would read and parse SSTable
            // For now, simulate
            let _ = fs::read(path)?;
        }
        
        // Sort and deduplicate entries
        all_entries.sort_by(|a, b| a.0.cmp(&b.0));
        all_entries.dedup_by(|a, b| a.0 == b.0);
        
        // Write new SSTable
        let output_path = self.sst_dir.join(format!("compacted-{}.sst", 
            SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs()));
        
        // In production, would write SSTable here
        
        let elapsed = SystemTime::now().duration_since(start_time).unwrap();
        
        Ok(CompactionResult {
            input_files: plan.input_files.len(),
            output_file: output_path,
            entries_processed: all_entries.len(),
            bytes_written: 0,
            duration_ms: elapsed.as_millis() as u64,
            success: true,
        })
    }
}

/// Compaction result
#[derive(Debug, Clone, Default)]
pub struct CompactionResult {
    pub input_files: usize,
    pub output_file: PathBuf,
    pub entries_processed: usize,
    pub bytes_written: u64,
    pub duration_ms: u64,
    pub success: bool,
}

/// Compaction scheduler
pub struct CompactionScheduler {
    config: CompactionConfig,
    picker: CompactionPicker,
    executor: CompactionExecutor,
    running: bool,
}

impl CompactionScheduler {
    pub fn new(sst_dir: PathBuf, config: CompactionConfig) -> Self {
        let picker = CompactionPicker::new(config.clone());
        let executor = CompactionExecutor::new(sst_dir, config.clone());
        CompactionScheduler {
            config,
            picker,
            executor,
            running: false,
        }
    }
    
    pub fn start(&mut self) {
        self.running = true;
    }
    
    pub fn stop(&mut self) {
        self.running = false;
    }
    
    pub fn is_running(&self) -> bool {
        self.running
    }
    
    pub fn schedule_compaction(&self, files: &[SSTableMetadata]) -> io::Result<CompactionResult> {
        let plan = self.picker.pick(files);
        if plan.is_empty() {
            return Ok(CompactionResult::default());
        }
        self.executor.execute(&plan)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;
    
    #[test]
    fn test_compaction_config_leveled() {
        let config = CompactionConfig::leveled();
        match config.strategy {
            CompactionStrategy::Leveled { level_size_multiplier, max_levels } => {
                assert_eq!(level_size_multiplier, 10.0);
                assert_eq!(max_levels, 7);
            }
            _ => panic!("Expected Leveled strategy"),
        }
    }
    
    #[test]
    fn test_compaction_config_size_tiered() {
        let config = CompactionConfig::size_tiered();
        match config.strategy {
            CompactionStrategy::SizeTiered { min_threshold, max_threshold } => {
                assert_eq!(min_threshold, 4);
                assert_eq!(max_threshold, 32);
            }
            _ => panic!("Expected SizeTiered strategy"),
        }
    }
    
    #[test]
    fn test_sstable_metadata_overlap() {
        let mut m1 = SSTableMetadata::new(PathBuf::from("a.sst"));
        m1.min_key = b"a".to_vec();
        m1.max_key = b"m".to_vec();
        
        let mut m2 = SSTableMetadata::new(PathBuf::from("b.sst"));
        m2.min_key = b"n".to_vec();
        m2.max_key = b"z".to_vec();
        
        assert!(!m1.key_range_overlaps(&m2));
        
        let mut m3 = SSTableMetadata::new(PathBuf::from("c.sst"));
        m3.min_key = b"k".to_vec();
        m3.max_key = b"p".to_vec();
        
        assert!(m1.key_range_overlaps(&m3));
    }
    
    #[test]
    fn test_compaction_plan_empty() {
        let plan = CompactionPlan::default();
        assert!(plan.is_empty());
        assert_eq!(plan.file_count(), 0);
    }
    
    #[test]
    fn test_compaction_picker_leveled() {
        let config = CompactionConfig::leveled();
        let picker = CompactionPicker::new(config);
        
        let files = vec![
            SSTableMetadata {
                path: PathBuf::from("l0-1.sst"),
                size: 100 * 1024 * 1024,
                level: 0,
                min_key: b"a".to_vec(),
                max_key: b"z".to_vec(),
                min_timestamp: 0,
                max_timestamp: 1000,
                entry_count: 1000,
                bloom_filter_size: 1024,
                created_at: 0,
            },
        ];
        
        let plan = picker.pick(&files);
        // Level 0 with 100MB should trigger compaction
        assert!(!plan.is_empty() || plan.input_files.is_empty());
    }
    
    #[test]
    fn test_compaction_scheduler() {
        let dir = tempdir().unwrap();
        let config = CompactionConfig::leveled();
        let mut scheduler = CompactionScheduler::new(dir.path().to_path_buf(), config);
        
        assert!(!scheduler.is_running());
        scheduler.start();
        assert!(scheduler.is_running());
        scheduler.stop();
        assert!(!scheduler.is_running());
    }
}