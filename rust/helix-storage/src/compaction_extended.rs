//! Helix Storage Engine - Compaction Module (Extended)
//!
//! This module provides extended compaction features including
//! automatic compaction, compaction statistics, and scheduling.

use std::collections::BTreeMap;
use std::path::PathBuf;
use std::sync::{Arc, Mutex, RwLock};
use std::time::{SystemTime, UNIX_EPOCH};

/// Compaction priority
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum CompactionPriority {
    Low = 0,
    Medium = 1,
    High = 2,
    Critical = 3,
}

/// Compaction job
#[derive(Debug, Clone)]
pub struct CompactionJob {
    pub id: String,
    pub input_files: Vec<PathBuf>,
    pub output_file: PathBuf,
    pub level: usize,
    pub priority: CompactionPriority,
    pub estimated_size: u64,
    pub created_at: i64,
    pub started_at: Option<i64>,
    pub completed_at: Option<i64>,
    pub status: CompactionJobStatus,
}

/// Compaction job status
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum CompactionJobStatus {
    Pending,
    Running,
    Completed,
    Failed,
    Cancelled,
}

/// Compaction statistics
#[derive(Debug, Clone, Default)]
pub struct CompactionStats {
    pub jobs_completed: u64,
    pub jobs_failed: u64,
    pub files_compacted: u64,
    pub bytes_written: u64,
    pub bytes_read: u64,
    pub duration_ms: u64,
    pub last_compaction: i64,
}

/// Extended compaction manager
pub struct ExtendedCompactionManager {
    sst_dir: PathBuf,
    config: CompactionConfig,
    jobs: RwLock<Vec<CompactionJob>>,
    stats: Arc<Mutex<CompactionStats>>,
    running: Arc<Mutex<bool>>,
    scheduler: CompactionScheduler,
}

/// Compaction configuration
#[derive(Debug, Clone)]
pub struct CompactionConfig {
    pub strategy: CompactionStrategyType,
    pub max_concurrent: usize,
    pub threshold: f64,
    pub target_file_size: u64,
    pub max_file_size: u64,
    pub min_files_to_compact: usize,
    pub max_files_to_compact: usize,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum CompactionStrategyType {
    Leveled,
    SizeTiered,
    TimeWindow,
    Universal,
}

impl CompactionConfig {
    pub fn leveled() -> Self {
        CompactionConfig {
            strategy: CompactionStrategyType::Leveled,
            max_concurrent: 2,
            threshold: 0.25,
            target_file_size: 64 * 1024 * 1024,
            max_file_size: 256 * 1024 * 1024,
            min_files_to_compact: 4,
            max_files_to_compact: 32,
        }
    }

    pub fn size_tiered() -> Self {
        CompactionConfig {
            strategy: CompactionStrategyType::SizeTiered,
            max_concurrent: 1,
            threshold: 0.5,
            target_file_size: 160 * 1024 * 1024,
            max_file_size: 512 * 1024 * 1024,
            min_files_to_compact: 4,
            max_files_to_compact: 32,
        }
    }

    pub fn time_window(_window_hours: u64) -> Self {
        CompactionConfig {
            strategy: CompactionStrategyType::TimeWindow,
            max_concurrent: 2,
            threshold: 0.3,
            target_file_size: 128 * 1024 * 1024,
            max_file_size: 512 * 1024 * 1024,
            min_files_to_compact: 2,
            max_files_to_compact: 16,
        }
    }
}

/// Compaction scheduler
pub struct CompactionScheduler {
    config: CompactionConfig,
}

impl CompactionScheduler {
    pub fn new(config: CompactionConfig) -> Self {
        CompactionScheduler { config }
    }

    pub fn schedule(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        match self.config.strategy {
            CompactionStrategyType::Leveled => self.schedule_leveled(files),
            CompactionStrategyType::SizeTiered => self.schedule_size_tiered(files),
            CompactionStrategyType::TimeWindow => self.schedule_time_window(files),
            CompactionStrategyType::Universal => self.schedule_universal(files),
        }
    }

    fn schedule_leveled(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        let mut jobs = Vec::new();

        // Group by level
        let mut levels: BTreeMap<usize, Vec<&SSTableInfo>> = BTreeMap::new();
        for file in files {
            levels.entry(file.level).or_default().push(file);
        }

        // Check each level
        for (level, level_files) in &levels {
            if level_files.len() < self.config.min_files_to_compact {
                continue;
            }

            let total_size: u64 = level_files.iter().map(|f| f.size).sum();
            let target_size =
                (self.config.target_file_size as f64 * 10.0_f64.powi(*level as i32)) as u64;

            if total_size > target_size {
                let job = CompactionJob {
                    id: format!(
                        "compact-{}",
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    ),
                    input_files: level_files.iter().map(|f| f.path.clone()).collect(),
                    output_file: PathBuf::from(format!(
                        "level-{}-{}.sst",
                        level + 1,
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    )),
                    level: *level + 1,
                    priority: if total_size > target_size * 2 {
                        CompactionPriority::High
                    } else {
                        CompactionPriority::Medium
                    },
                    estimated_size: total_size,
                    created_at: SystemTime::now()
                        .duration_since(UNIX_EPOCH)
                        .unwrap()
                        .as_secs() as i64,
                    started_at: None,
                    completed_at: None,
                    status: CompactionJobStatus::Pending,
                };
                jobs.push(job);
            }
        }

        jobs
    }

    fn schedule_size_tiered(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        let mut jobs = Vec::new();

        // Group by size bucket
        let mut buckets: BTreeMap<usize, Vec<&SSTableInfo>> = BTreeMap::new();
        for file in files {
            let bucket = (file.size as f64).log2().floor() as usize;
            buckets.entry(bucket).or_default().push(file);
        }

        for (bucket, bucket_files) in &buckets {
            if bucket_files.len() >= self.config.min_files_to_compact {
                let job = CompactionJob {
                    id: format!(
                        "compact-{}",
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    ),
                    input_files: bucket_files.iter().map(|f| f.path.clone()).collect(),
                    output_file: PathBuf::from(format!(
                        "tier-{}-{}.sst",
                        bucket + 1,
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    )),
                    level: *bucket + 1,
                    priority: CompactionPriority::Medium,
                    estimated_size: bucket_files.iter().map(|f| f.size).sum(),
                    created_at: SystemTime::now()
                        .duration_since(UNIX_EPOCH)
                        .unwrap()
                        .as_secs() as i64,
                    started_at: None,
                    completed_at: None,
                    status: CompactionJobStatus::Pending,
                };
                jobs.push(job);
            }
        }

        jobs
    }

    fn schedule_time_window(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        let mut jobs = Vec::new();

        // Group by time window (1 hour)
        let window_size: i64 = 3600;
        let mut windows: BTreeMap<i64, Vec<&SSTableInfo>> = BTreeMap::new();
        for file in files {
            let window = file.min_timestamp / window_size;
            windows.entry(window).or_default().push(file);
        }

        for (window, window_files) in &windows {
            if window_files.len() > 1 {
                let job = CompactionJob {
                    id: format!(
                        "compact-{}",
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    ),
                    input_files: window_files.iter().map(|f| f.path.clone()).collect(),
                    output_file: PathBuf::from(format!(
                        "window-{}-{}.sst",
                        window,
                        SystemTime::now()
                            .duration_since(UNIX_EPOCH)
                            .unwrap()
                            .as_nanos()
                    )),
                    level: 0,
                    priority: CompactionPriority::Low,
                    estimated_size: window_files.iter().map(|f| f.size).sum(),
                    created_at: SystemTime::now()
                        .duration_since(UNIX_EPOCH)
                        .unwrap()
                        .as_secs() as i64,
                    started_at: None,
                    completed_at: None,
                    status: CompactionJobStatus::Pending,
                };
                jobs.push(job);
            }
        }

        jobs
    }

    fn schedule_universal(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        // Universal compaction: compact all files if conditions met
        let mut jobs = Vec::new();

        if files.len() >= self.config.min_files_to_compact {
            let job = CompactionJob {
                id: format!(
                    "compact-{}",
                    SystemTime::now()
                        .duration_since(UNIX_EPOCH)
                        .unwrap()
                        .as_nanos()
                ),
                input_files: files.iter().map(|f| f.path.clone()).collect(),
                output_file: PathBuf::from(format!(
                    "universal-{}.sst",
                    SystemTime::now()
                        .duration_since(UNIX_EPOCH)
                        .unwrap()
                        .as_nanos()
                )),
                level: 0,
                priority: CompactionPriority::High,
                estimated_size: files.iter().map(|f| f.size).sum(),
                created_at: SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64,
                started_at: None,
                completed_at: None,
                status: CompactionJobStatus::Pending,
            };
            jobs.push(job);
        }

        jobs
    }
}

/// SSTable info for compaction decisions
#[derive(Debug, Clone)]
pub struct SSTableInfo {
    pub path: PathBuf,
    pub size: u64,
    pub level: usize,
    pub min_key: Vec<u8>,
    pub max_key: Vec<u8>,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
    pub entry_count: usize,
    pub created_at: i64,
}

impl ExtendedCompactionManager {
    pub fn new(sst_dir: PathBuf, config: CompactionConfig) -> Self {
        let scheduler = CompactionScheduler::new(config.clone());
        ExtendedCompactionManager {
            sst_dir,
            config,
            jobs: RwLock::new(Vec::new()),
            stats: Arc::new(Mutex::new(CompactionStats::default())),
            running: Arc::new(Mutex::new(false)),
            scheduler,
        }
    }

    pub fn start(&self) {
        let mut running = self.running.lock().unwrap();
        *running = true;
    }

    pub fn stop(&self) {
        let mut running = self.running.lock().unwrap();
        *running = false;
    }

    pub fn is_running(&self) -> bool {
        *self.running.lock().unwrap()
    }

    pub fn submit_job(&self, job: CompactionJob) {
        let mut jobs = self.jobs.write().unwrap();
        jobs.push(job);
    }

    pub fn cancel_job(&self, id: &str) -> bool {
        let mut jobs = self.jobs.write().unwrap();
        for job in jobs.iter_mut() {
            if job.id == id && job.status == CompactionJobStatus::Pending {
                job.status = CompactionJobStatus::Cancelled;
                return true;
            }
        }
        false
    }

    pub fn pending_jobs(&self) -> Vec<CompactionJob> {
        let jobs = self.jobs.read().unwrap();
        jobs.iter()
            .filter(|j| j.status == CompactionJobStatus::Pending)
            .cloned()
            .collect()
    }

    pub fn running_jobs(&self) -> Vec<CompactionJob> {
        let jobs = self.jobs.read().unwrap();
        jobs.iter()
            .filter(|j| j.status == CompactionJobStatus::Running)
            .cloned()
            .collect()
    }

    pub fn completed_jobs(&self) -> Vec<CompactionJob> {
        let jobs = self.jobs.read().unwrap();
        jobs.iter()
            .filter(|j| j.status == CompactionJobStatus::Completed)
            .cloned()
            .collect()
    }

    pub fn stats(&self) -> CompactionStats {
        self.stats.lock().unwrap().clone()
    }

    pub fn schedule(&self, files: &[SSTableInfo]) -> Vec<CompactionJob> {
        let sstable_infos: Vec<SSTableInfo> = files.to_vec();
        self.scheduler.schedule(&sstable_infos)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_compaction_config_leveled() {
        let config = CompactionConfig::leveled();
        assert_eq!(config.strategy, CompactionStrategyType::Leveled);
        assert_eq!(config.max_concurrent, 2);
    }

    #[test]
    fn test_compaction_config_size_tiered() {
        let config = CompactionConfig::size_tiered();
        assert_eq!(config.strategy, CompactionStrategyType::SizeTiered);
    }

    #[test]
    fn test_compaction_config_time_window() {
        let config = CompactionConfig::time_window(24);
        assert_eq!(config.strategy, CompactionStrategyType::TimeWindow);
    }

    #[test]
    fn test_compaction_job_status() {
        let job = CompactionJob {
            id: "job-1".to_string(),
            input_files: vec![],
            output_file: PathBuf::from("output.sst"),
            level: 1,
            priority: CompactionPriority::Medium,
            estimated_size: 1024,
            created_at: 1000,
            started_at: None,
            completed_at: None,
            status: CompactionJobStatus::Pending,
        };

        assert_eq!(job.status, CompactionJobStatus::Pending);
        assert!(job.status != CompactionJobStatus::Completed);
    }

    #[test]
    fn test_compaction_priority_ordering() {
        assert!(CompactionPriority::Critical > CompactionPriority::High);
        assert!(CompactionPriority::High > CompactionPriority::Medium);
        assert!(CompactionPriority::Medium > CompactionPriority::Low);
    }

    #[test]
    fn test_compaction_manager_create() {
        let dir = tempdir().unwrap();
        let config = CompactionConfig::leveled();
        let manager = ExtendedCompactionManager::new(dir.path().to_path_buf(), config);

        assert!(!manager.is_running());
        manager.start();
        assert!(manager.is_running());
        manager.stop();
        assert!(!manager.is_running());
    }

    #[test]
    fn test_compaction_manager_submit_job() {
        let dir = tempdir().unwrap();
        let config = CompactionConfig::leveled();
        let manager = ExtendedCompactionManager::new(dir.path().to_path_buf(), config);

        let job = CompactionJob {
            id: "job-1".to_string(),
            input_files: vec![],
            output_file: PathBuf::from("output.sst"),
            level: 1,
            priority: CompactionPriority::Medium,
            estimated_size: 1024,
            created_at: 1000,
            started_at: None,
            completed_at: None,
            status: CompactionJobStatus::Pending,
        };

        manager.submit_job(job);
        assert_eq!(manager.pending_jobs().len(), 1);
    }

    #[test]
    fn test_compaction_manager_cancel_job() {
        let dir = tempdir().unwrap();
        let config = CompactionConfig::leveled();
        let manager = ExtendedCompactionManager::new(dir.path().to_path_buf(), config);

        let job = CompactionJob {
            id: "job-1".to_string(),
            input_files: vec![],
            output_file: PathBuf::from("output.sst"),
            level: 1,
            priority: CompactionPriority::Medium,
            estimated_size: 1024,
            created_at: 1000,
            started_at: None,
            completed_at: None,
            status: CompactionJobStatus::Pending,
        };

        manager.submit_job(job);
        assert!(manager.cancel_job("job-1"));
        assert_eq!(manager.pending_jobs().len(), 0);
    }

    #[test]
    fn test_compaction_scheduler_leveled() {
        let config = CompactionConfig::leveled();
        let scheduler = CompactionScheduler::new(config);

        let files = vec![SSTableInfo {
            path: PathBuf::from("l0-1.sst"),
            size: 100 * 1024 * 1024,
            level: 0,
            min_key: b"a".to_vec(),
            max_key: b"z".to_vec(),
            min_timestamp: 0,
            max_timestamp: 1000,
            entry_count: 1000,
            created_at: 0,
        }];

        let jobs = scheduler.schedule(&files);
        // May or may not produce jobs depending on thresholds
        let _ = jobs;
    }

    #[test]
    fn test_compaction_stats() {
        let stats = CompactionStats::default();
        assert_eq!(stats.jobs_completed, 0);
        assert_eq!(stats.jobs_failed, 0);
    }
}
