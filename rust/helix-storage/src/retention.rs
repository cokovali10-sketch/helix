//! Helix Storage Engine - Retention Module
//!
//! This module provides data retention policies for the TSDB.
//! Supports time-based, size-based, and custom retention strategies.

use std::fs;
use std::io;
use std::path::{Path, PathBuf};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Retention policy type
#[derive(Debug, Clone, PartialEq)]
pub enum RetentionPolicy {
    /// Keep data for a specified duration
    TimeBased { duration: Duration },
    /// Keep data up to a specified size
    SizeBased { max_bytes: u64 },
    /// Keep data for duration and up to size
    TimeAndSize { duration: Duration, max_bytes: u64 },
    /// Keep data for a specified number of series
    SeriesBased { max_series: usize },
    /// Custom retention with predicate
    Custom { predicate: String },
    /// No retention (keep forever)
    None,
}

/// Retention configuration
#[derive(Debug, Clone)]
pub struct RetentionConfig {
    pub policy: RetentionPolicy,
    pub check_interval: Duration,
    pub dry_run: bool,
    pub verbose: bool,
}

impl RetentionConfig {
    pub fn new(policy: RetentionPolicy) -> Self {
        RetentionConfig {
            policy,
            check_interval: Duration::from_secs(3600), // 1 hour
            dry_run: false,
            verbose: false,
        }
    }

    pub fn time_based(hours: u64) -> Self {
        RetentionConfig::new(RetentionPolicy::TimeBased {
            duration: Duration::from_secs(hours * 3600),
        })
    }

    pub fn none() -> Self {
        RetentionConfig::new(RetentionPolicy::None)
    }

    pub fn size_based(gb: u64) -> Self {
        RetentionConfig::new(RetentionPolicy::SizeBased {
            max_bytes: gb * 1024 * 1024 * 1024,
        })
    }

    pub fn time_and_size(hours: u64, gb: u64) -> Self {
        RetentionConfig::new(RetentionPolicy::TimeAndSize {
            duration: Duration::from_secs(hours * 3600),
            max_bytes: gb * 1024 * 1024 * 1024,
        })
    }
}

/// Retention manager
pub struct RetentionManager {
    config: RetentionConfig,
    sst_dir: PathBuf,
    stats: RetentionStats,
}

/// Retention statistics
#[derive(Debug, Clone, Default)]
pub struct RetentionStats {
    pub files_scanned: usize,
    pub files_deleted: usize,
    pub bytes_freed: u64,
    pub series_removed: usize,
    pub oldest_retained_timestamp: i64,
    pub newest_retained_timestamp: i64,
}

impl RetentionManager {
    pub fn new(sst_dir: PathBuf, config: RetentionConfig) -> Self {
        RetentionManager {
            config,
            sst_dir,
            stats: RetentionStats::default(),
        }
    }

    /// Check if a file should be retained based on policy
    pub fn should_retain(
        &self,
        _file_path: &Path,
        _file_size: u64,
        _oldest_ts: i64,
        newest_ts: i64,
    ) -> bool {
        match &self.config.policy {
            RetentionPolicy::TimeBased { duration } => {
                let cutoff = SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64
                    - duration.as_secs() as i64;
                newest_ts >= cutoff
            }
            RetentionPolicy::SizeBased { max_bytes: _ } => {
                // Size-based retention handled at aggregate level
                true
            }
            RetentionPolicy::TimeAndSize {
                duration,
                max_bytes: _,
            } => {
                let cutoff = SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64
                    - duration.as_secs() as i64;
                newest_ts >= cutoff
            }
            RetentionPolicy::SeriesBased { max_series: _ } => true,
            RetentionPolicy::Custom { predicate: _ } => true,
            RetentionPolicy::None => true,
        }
    }

    /// Run retention check and delete expired files
    pub fn run_retention(&mut self) -> io::Result<RetentionStats> {
        let mut stats = RetentionStats::default();

        if !self.sst_dir.exists() {
            return Ok(stats);
        }

        // Collect SSTable files
        let mut sst_files: Vec<(PathBuf, u64, i64, i64)> = Vec::new();

        for entry in fs::read_dir(&self.sst_dir)? {
            let entry = entry?;
            let path = entry.path();

            if path.extension().map(|e| e == "sst").unwrap_or(false) {
                let metadata = entry.metadata()?;
                let size = metadata.len();

                // Read oldest/newest timestamps from file (simplified)
                let oldest_ts = 0i64;
                let newest_ts = SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64;

                sst_files.push((path, size, oldest_ts, newest_ts));
                stats.files_scanned += 1;
            }
        }

        // Sort by newest timestamp (oldest first)
        sst_files.sort_by_key(|f| f.3);

        // Calculate total size
        let total_size: u64 = sst_files.iter().map(|f| f.1).sum();

        // Determine which files to delete
        let mut files_to_delete: Vec<PathBuf> = Vec::new();

        match &self.config.policy {
            RetentionPolicy::TimeBased { duration } => {
                let cutoff = SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64
                    - duration.as_secs() as i64;

                for (path, _, _, newest_ts) in &sst_files {
                    if *newest_ts < cutoff {
                        files_to_delete.push(path.clone());
                    }
                }
            }
            RetentionPolicy::SizeBased { max_bytes } => {
                if total_size > *max_bytes {
                    let mut freed = 0u64;
                    let target_free = total_size - max_bytes;

                    for (path, size, _, _) in &sst_files {
                        if freed >= target_free {
                            break;
                        }
                        files_to_delete.push(path.clone());
                        freed += size;
                    }
                }
            }
            RetentionPolicy::TimeAndSize {
                duration,
                max_bytes,
            } => {
                let cutoff = SystemTime::now()
                    .duration_since(UNIX_EPOCH)
                    .unwrap()
                    .as_secs() as i64
                    - duration.as_secs() as i64;

                let mut freed = 0u64;
                let target_free = total_size.saturating_sub(*max_bytes);

                for (path, size, _, newest_ts) in &sst_files {
                    if *newest_ts < cutoff || freed < target_free {
                        files_to_delete.push(path.clone());
                        freed += size;
                    }
                }
            }
            _ => {}
        }

        // Delete files
        for path in &files_to_delete {
            if self.config.dry_run {
                if self.config.verbose {
                    println!("DRY RUN: Would delete {}", path.display());
                }
            } else {
                fs::remove_file(path)?;
                stats.files_deleted += 1;

                if let Ok(metadata) = fs::metadata(path) {
                    stats.bytes_freed += metadata.len();
                }
            }
        }

        self.stats = stats.clone();
        Ok(stats)
    }

    /// Get retention statistics
    pub fn stats(&self) -> &RetentionStats {
        &self.stats
    }

    /// Estimate data age
    pub fn estimate_age(&self, timestamp: i64) -> Duration {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        let age_secs = now - timestamp;
        Duration::from_secs(age_secs.max(0) as u64)
    }

    /// Check if data is expired
    pub fn is_expired(&self, timestamp: i64) -> bool {
        match &self.config.policy {
            RetentionPolicy::TimeBased { duration } => self.estimate_age(timestamp) > *duration,
            RetentionPolicy::TimeAndSize { duration, .. } => {
                self.estimate_age(timestamp) > *duration
            }
            _ => false,
        }
    }
}

/// Retention policy builder
pub struct RetentionPolicyBuilder {
    policy: RetentionPolicy,
}

impl RetentionPolicyBuilder {
    pub fn new() -> Self {
        RetentionPolicyBuilder {
            policy: RetentionPolicy::None,
        }
    }

    pub fn time_based(mut self, hours: u64) -> Self {
        self.policy = RetentionPolicy::TimeBased {
            duration: Duration::from_secs(hours * 3600),
        };
        self
    }

    pub fn size_based(mut self, gb: u64) -> Self {
        self.policy = RetentionPolicy::SizeBased {
            max_bytes: gb * 1024 * 1024 * 1024,
        };
        self
    }

    pub fn time_and_size(mut self, hours: u64, gb: u64) -> Self {
        self.policy = RetentionPolicy::TimeAndSize {
            duration: Duration::from_secs(hours * 3600),
            max_bytes: gb * 1024 * 1024 * 1024,
        };
        self
    }

    pub fn none(mut self) -> Self {
        self.policy = RetentionPolicy::None;
        self
    }

    pub fn build(self) -> RetentionPolicy {
        self.policy
    }
}

impl Default for RetentionPolicyBuilder {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_retention_policy_builder() {
        let policy = RetentionPolicyBuilder::new().time_based(24).build();

        match policy {
            RetentionPolicy::TimeBased { duration } => {
                assert_eq!(duration.as_secs(), 24 * 3600);
            }
            _ => panic!("Expected TimeBased policy"),
        }
    }

    #[test]
    fn test_retention_config() {
        let config = RetentionConfig::time_based(24);

        match config.policy {
            RetentionPolicy::TimeBased { duration } => {
                assert_eq!(duration.as_secs(), 24 * 3600);
            }
            _ => panic!("Expected TimeBased policy"),
        }
    }

    #[test]
    fn test_retention_manager_creation() {
        let dir = tempdir().unwrap();
        let config = RetentionConfig::time_based(24);
        let manager = RetentionManager::new(dir.path().to_path_buf(), config);

        assert_eq!(manager.stats().files_scanned, 0);
    }

    #[test]
    fn test_is_expired() {
        let config = RetentionConfig::time_based(1); // 1 hour
        let manager = RetentionManager::new(std::path::PathBuf::from("./test"), config);

        // Data from 2 hours ago should be expired
        let two_hours_ago = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64
            - 7200;

        assert!(manager.is_expired(two_hours_ago));
    }

    #[test]
    fn test_estimate_age() {
        let config = RetentionConfig::none();
        let manager = RetentionManager::new(std::path::PathBuf::from("./test"), config);

        let one_hour_ago = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64
            - 3600;

        let age = manager.estimate_age(one_hour_ago);
        assert!(age.as_secs() >= 3590 && age.as_secs() <= 3610);
    }
}
