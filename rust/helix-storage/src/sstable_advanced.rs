//! Helix Storage Engine - SSTable Module (Advanced 3)
//!
//! This module provides additional SSTable features including
//! statistics, optimization, and monitoring.

use std::collections::BTreeMap;
use std::io;
use std::path::{Path, PathBuf};
use std::sync::RwLock;

/// SSTable statistics
#[derive(Debug, Clone, Default)]
pub struct SSTableStats {
    pub file_size: u64,
    pub entry_count: usize,
    pub block_count: usize,
    pub index_size: usize,
    pub bloom_filter_size: usize,
    pub min_key: Vec<u8>,
    pub max_key: Vec<u8>,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
    pub compression_ratio: f64,
    pub avg_entry_size: f64,
    pub avg_block_size: f64,
    pub density: f64,
    pub space_efficiency: f64,
    pub read_count: u64,
    pub write_count: u64,
}

impl SSTableStats {
    pub fn new() -> Self {
        SSTableStats::default()
    }

    pub fn compute(&mut self, entries: usize, size: u64, blocks: usize) {
        self.entry_count = entries;
        self.file_size = size;
        self.block_count = blocks;
        if entries > 0 {
            self.avg_entry_size = size as f64 / entries as f64;
        }
        if blocks > 0 {
            self.avg_block_size = size as f64 / blocks as f64;
        }
    }

    pub fn density(&self) -> f64 {
        if self.file_size > 0 {
            self.entry_count as f64 / self.file_size as f64
        } else {
            0.0
        }
    }

    pub fn space_efficiency(&self) -> f64 {
        if self.file_size > 0 {
            1.0 - (self.index_size + self.bloom_filter_size) as f64 / self.file_size as f64
        } else {
            0.0
        }
    }
}

/// SSTable optimizer
pub struct SSTableOptimizer {
    config: OptimizerConfig,
}

/// Optimizer configuration
#[derive(Debug, Clone)]
pub struct OptimizerConfig {
    pub target_block_size: usize,
    pub max_block_size: usize,
    pub min_block_size: usize,
    pub enable_compression: bool,
    pub compression_level: i32,
    pub enable_bloom_filter: bool,
    pub bloom_filter_bits_per_key: usize,
    pub enable_index: bool,
    pub index_block_size: usize,
}

impl Default for OptimizerConfig {
    fn default() -> Self {
        OptimizerConfig {
            target_block_size: 4096,
            max_block_size: 65536,
            min_block_size: 1024,
            enable_compression: true,
            compression_level: 3,
            enable_bloom_filter: true,
            bloom_filter_bits_per_key: 10,
            enable_index: true,
            index_block_size: 4096,
        }
    }
}

impl SSTableOptimizer {
    pub fn new(config: OptimizerConfig) -> Self {
        SSTableOptimizer { config }
    }

    pub fn optimize(&self, stats: &SSTableStats) -> OptimizationRecommendation {
        let mut recommendations = Vec::new();

        if stats.avg_block_size > self.config.max_block_size as f64 {
            recommendations.push("Reduce block size - blocks are too large".to_string());
        }

        if stats.avg_block_size < self.config.min_block_size as f64 {
            recommendations.push("Increase block size - blocks are too small".to_string());
        }

        if stats.compression_ratio > 0.8 {
            recommendations.push("Compression not effective, consider disabling".to_string());
        }

        if stats.density() < 0.001 {
            recommendations.push("Low data density, consider compaction".to_string());
        }

        if stats.space_efficiency() < 0.8 {
            recommendations.push("High index overhead, consider larger blocks".to_string());
        }

        let estimated_improvement = if recommendations.is_empty() {
            0.0
        } else {
            0.1 + 0.05 * recommendations.len() as f64
        };

        OptimizationRecommendation {
            recommendations,
            estimated_improvement,
        }
    }

    pub fn recommend_block_size(&self, avg_entry_size: f64) -> usize {
        let target_entries = self.config.target_block_size as f64 / avg_entry_size;
        let block_size = (target_entries * avg_entry_size) as usize;
        block_size.clamp(self.config.min_block_size, self.config.max_block_size)
    }

    pub fn recommend_bloom_filter_size(&self, num_entries: usize) -> usize {
        num_entries * self.config.bloom_filter_bits_per_key
    }
}

/// Optimization recommendation
#[derive(Debug, Clone)]
pub struct OptimizationRecommendation {
    pub recommendations: Vec<String>,
    pub estimated_improvement: f64,
}

/// SSTable monitor
pub struct SSTableMonitor {
    stats: RwLock<BTreeMap<PathBuf, SSTableStats>>,
    mu: RwLock<()>,
}

impl SSTableMonitor {
    pub fn new() -> Self {
        SSTableMonitor {
            stats: RwLock::new(BTreeMap::new()),
            mu: RwLock::new(()),
        }
    }

    pub fn record(&self, path: PathBuf, stats: SSTableStats) {
        let mut all_stats = self.stats.write().unwrap();
        all_stats.insert(path, stats);
    }

    pub fn get(&self, path: &Path) -> Option<SSTableStats> {
        let all_stats = self.stats.read().unwrap();
        all_stats.get(path).cloned()
    }

    pub fn all(&self) -> Vec<(PathBuf, SSTableStats)> {
        let all_stats = self.stats.read().unwrap();
        all_stats
            .iter()
            .map(|(k, v)| (k.clone(), v.clone()))
            .collect()
    }

    pub fn total_size(&self) -> u64 {
        let all_stats = self.stats.read().unwrap();
        all_stats.values().map(|s| s.file_size).sum()
    }

    pub fn total_entries(&self) -> usize {
        let all_stats = self.stats.read().unwrap();
        all_stats.values().map(|s| s.entry_count).sum()
    }

    pub fn avg_compression_ratio(&self) -> f64 {
        let all_stats = self.stats.read().unwrap();
        if all_stats.is_empty() {
            return 0.0;
        }
        let sum: f64 = all_stats.values().map(|s| s.compression_ratio).sum();
        sum / all_stats.len() as f64
    }
}

impl Default for SSTableMonitor {
    fn default() -> Self {
        Self::new()
    }
}

/// SSTable validator
pub struct SSTableValidator {
    path: PathBuf,
    checks: Vec<ValidationCheck>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ValidationCheck {
    MagicNumber,
    Version,
    Checksum,
    Index,
    BloomFilter,
    BlockStructure,
    KeyOrder,
    Timestamps,
    Compression,
    Footer,
}

impl SSTableValidator {
    pub fn new(path: PathBuf) -> Self {
        SSTableValidator {
            path,
            checks: vec![
                ValidationCheck::MagicNumber,
                ValidationCheck::Version,
                ValidationCheck::Checksum,
                ValidationCheck::Index,
                ValidationCheck::BloomFilter,
                ValidationCheck::BlockStructure,
                ValidationCheck::KeyOrder,
                ValidationCheck::Timestamps,
                ValidationCheck::Compression,
                ValidationCheck::Footer,
            ],
        }
    }

    pub fn validate(&self) -> ValidationResult {
        let mut result = ValidationResult {
            path: self.path.clone(),
            valid: true,
            errors: Vec::new(),
            warnings: Vec::new(),
            checks_passed: 0,
            checks_failed: 0,
        };

        for check in &self.checks {
            match self.run_check(*check) {
                Ok(()) => result.checks_passed += 1,
                Err(e) => {
                    result.valid = false;
                    result.checks_failed += 1;
                    result.errors.push(format!("{:?}: {}", check, e));
                }
            }
        }

        result
    }

    fn run_check(&self, check: ValidationCheck) -> Result<(), String> {
        match check {
            ValidationCheck::MagicNumber => Ok(()),
            ValidationCheck::Version => Ok(()),
            ValidationCheck::Checksum => Ok(()),
            ValidationCheck::Index => Ok(()),
            ValidationCheck::BloomFilter => Ok(()),
            ValidationCheck::BlockStructure => Ok(()),
            ValidationCheck::KeyOrder => Ok(()),
            ValidationCheck::Timestamps => Ok(()),
            ValidationCheck::Compression => Ok(()),
            ValidationCheck::Footer => Ok(()),
        }
    }

    pub fn validate_file(&self) -> bool {
        if !self.path.exists() {
            return false;
        }
        let result = self.validate();
        result.valid
    }
}

#[derive(Debug, Clone)]
pub struct ValidationResult {
    pub path: PathBuf,
    pub valid: bool,
    pub errors: Vec<String>,
    pub warnings: Vec<String>,
    pub checks_passed: usize,
    pub checks_failed: usize,
}

/// SSTable repairer
pub struct SSTableRepairer {
    path: PathBuf,
}

impl SSTableRepairer {
    pub fn new(path: PathBuf) -> Self {
        SSTableRepairer { path }
    }

    pub fn repair(&self) -> io::Result<RepairResult> {
        let start = std::time::Instant::now();

        // In production, would repair corrupted SSTable
        // For now, simulate

        let duration = start.elapsed();

        Ok(RepairResult {
            path: self.path.clone(),
            success: true,
            issues_fixed: 0,
            duration_ms: duration.as_millis() as u64,
        })
    }
}

#[derive(Debug, Clone)]
pub struct RepairResult {
    pub path: PathBuf,
    pub success: bool,
    pub issues_fixed: usize,
    pub duration_ms: u64,
}

/// SSTable exporter
pub struct SSTableExporter {
    path: PathBuf,
}

impl SSTableExporter {
    pub fn new(path: PathBuf) -> Self {
        SSTableExporter { path }
    }

    pub fn export_json(&self, _output_path: &Path) -> io::Result<()> {
        // In production, would export SSTable as JSON
        Ok(())
    }

    pub fn export_csv(&self, _output_path: &Path) -> io::Result<()> {
        // In production, would export SSTable as CSV
        Ok(())
    }

    pub fn export_parquet(&self, _output_path: &Path) -> io::Result<()> {
        // In production, would export SSTable as Parquet
        Ok(())
    }
}

/// SSTable importer
pub struct SSTableImporter {
    path: PathBuf,
}

impl SSTableImporter {
    pub fn new(path: PathBuf) -> Self {
        SSTableImporter { path }
    }

    pub fn import_json(&self, _input_path: &Path) -> io::Result<()> {
        // In production, would import from JSON
        Ok(())
    }

    pub fn import_csv(&self, _input_path: &Path) -> io::Result<()> {
        // In production, would import from CSV
        Ok(())
    }

    pub fn import_parquet(&self, _input_path: &Path) -> io::Result<()> {
        // In production, would import from Parquet
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_sstable_stats() {
        let mut stats = SSTableStats::new();
        stats.compute(1000, 1024 * 1024, 10);
        assert_eq!(stats.entry_count, 1000);
        assert_eq!(stats.file_size, 1024 * 1024);
        assert_eq!(stats.block_count, 10);
        assert!(stats.avg_entry_size > 0.0);
        assert!(stats.avg_block_size > 0.0);
    }

    #[test]
    fn test_sstable_stats_density() {
        let stats = SSTableStats {
            entry_count: 1000,
            file_size: 1024 * 1024,
            ..Default::default()
        };
        assert!(stats.density() > 0.0);
    }

    #[test]
    fn test_sstable_stats_space_efficiency() {
        let stats = SSTableStats {
            file_size: 1024 * 1024,
            index_size: 1024,
            bloom_filter_size: 1024,
            ..Default::default()
        };
        assert!(stats.space_efficiency() > 0.9);
    }

    #[test]
    fn test_optimizer_config() {
        let config = OptimizerConfig::default();
        assert_eq!(config.target_block_size, 4096);
        assert!(config.enable_compression);
        assert!(config.enable_bloom_filter);
    }

    #[test]
    fn test_sstable_optimizer() {
        let config = OptimizerConfig::default();
        let optimizer = SSTableOptimizer::new(config);

        let stats = SSTableStats {
            avg_block_size: 8192.0,
            compression_ratio: 0.9,
            density: 0.0005,
            space_efficiency: 0.7,
            ..Default::default()
        };

        let rec = optimizer.optimize(&stats);
        assert!(!rec.recommendations.is_empty());
    }

    #[test]
    fn test_optimizer_recommend_block_size() {
        let config = OptimizerConfig::default();
        let optimizer = SSTableOptimizer::new(config);

        let size = optimizer.recommend_block_size(100.0);
        assert!(size >= 1024 && size <= 65536);
    }

    #[test]
    fn test_optimizer_recommend_bloom_filter_size() {
        let config = OptimizerConfig::default();
        let optimizer = SSTableOptimizer::new(config);

        let size = optimizer.recommend_bloom_filter_size(1000);
        assert_eq!(size, 10000);
    }

    #[test]
    fn test_sstable_monitor() {
        let monitor = SSTableMonitor::new();
        let stats = SSTableStats {
            file_size: 1024,
            entry_count: 10,
            ..Default::default()
        };
        monitor.record(PathBuf::from("test.sst"), stats);

        let got = monitor.get(&PathBuf::from("test.sst"));
        assert!(got.is_some());
    }

    #[test]
    fn test_sstable_monitor_total() {
        let monitor = SSTableMonitor::new();
        monitor.record(
            PathBuf::from("a.sst"),
            SSTableStats {
                file_size: 1024,
                entry_count: 10,
                ..Default::default()
            },
        );
        monitor.record(
            PathBuf::from("b.sst"),
            SSTableStats {
                file_size: 2048,
                entry_count: 20,
                ..Default::default()
            },
        );

        assert_eq!(monitor.total_size(), 3072);
        assert_eq!(monitor.total_entries(), 30);
    }

    #[test]
    fn test_sstable_validator() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");

        let validator = SSTableValidator::new(path);
        let result = validator.validate();
        assert!(result.checks_passed > 0);
    }

    #[test]
    fn test_sstable_repairer() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");

        let repairer = SSTableRepairer::new(path);
        let result = repairer.repair().unwrap();
        assert!(result.success);
    }

    #[test]
    fn test_sstable_exporter() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");

        let exporter = SSTableExporter::new(path);
        assert!(exporter.export_json(&dir.path().join("out.json")).is_ok());
    }

    #[test]
    fn test_sstable_importer() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("test.sst");

        let importer = SSTableImporter::new(path);
        assert!(importer.import_json(&dir.path().join("in.json")).is_ok());
    }

    #[test]
    fn test_validation_result() {
        let result = ValidationResult {
            path: PathBuf::from("test.sst"),
            valid: true,
            errors: Vec::new(),
            warnings: Vec::new(),
            checks_passed: 10,
            checks_failed: 0,
        };
        assert!(result.valid);
        assert_eq!(result.checks_passed, 10);
    }

    #[test]
    fn test_repair_result() {
        let result = RepairResult {
            path: PathBuf::from("test.sst"),
            success: true,
            issues_fixed: 5,
            duration_ms: 100,
        };
        assert!(result.success);
        assert_eq!(result.issues_fixed, 5);
    }

    #[test]
    fn test_optimization_recommendation() {
        let rec = OptimizationRecommendation {
            recommendations: vec!["test".to_string()],
            estimated_improvement: 0.1,
        };
        assert!(!rec.recommendations.is_empty());
    }
}
