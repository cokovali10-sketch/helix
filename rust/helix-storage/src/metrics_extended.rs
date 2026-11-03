//! Helix Storage Engine - Metrics Module (Extended)
//!
//! This module provides extended metrics collection and observability.

use std::collections::BTreeMap;
use std::sync::atomic::{AtomicI64, AtomicU64, Ordering};
use std::sync::Arc;
use std::time::{SystemTime, UNIX_EPOCH};

/// Extended storage metrics
pub struct ExtendedStorageMetrics {
    // Write metrics
    pub writes_total: AtomicU64,
    pub writes_bytes: AtomicU64,
    pub write_errors: AtomicU64,
    pub write_latency_sum: AtomicU64,
    pub write_latency_count: AtomicU64,
    pub write_latency_max: AtomicU64,

    // Read metrics
    pub reads_total: AtomicU64,
    pub reads_bytes: AtomicU64,
    pub read_errors: AtomicU64,
    pub read_latency_sum: AtomicU64,
    pub read_latency_count: AtomicU64,
    pub read_latency_max: AtomicU64,
    pub cache_hits: AtomicU64,
    pub cache_misses: AtomicU64,

    // Memtable metrics
    pub memtable_size_bytes: AtomicU64,
    pub memtable_entries: AtomicU64,
    pub memtable_flushes: AtomicU64,

    // SSTable metrics
    pub sstable_count: AtomicU64,
    pub sstable_size_bytes: AtomicU64,
    pub sstable_compactions: AtomicU64,
    pub sstable_reads: AtomicU64,
    pub sstable_writes: AtomicU64,

    // WAL metrics
    pub wal_size_bytes: AtomicU64,
    pub wal_segments: AtomicU64,
    pub wal_syncs: AtomicU64,
    pub wal_records: AtomicU64,

    // Retention metrics
    pub retention_deleted_files: AtomicU64,
    pub retention_freed_bytes: AtomicU64,
    pub retention_runs: AtomicU64,

    // Compaction metrics
    pub compaction_jobs: AtomicU64,
    pub compaction_bytes_read: AtomicU64,
    pub compaction_bytes_written: AtomicU64,
    pub compaction_duration_ms: AtomicU64,

    // General
    pub uptime_seconds: AtomicU64,
    pub start_time: AtomicI64,
    pub active_connections: AtomicU64,
    pub total_connections: AtomicU64,
}

impl ExtendedStorageMetrics {
    pub fn new() -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;
        ExtendedStorageMetrics {
            writes_total: AtomicU64::new(0),
            writes_bytes: AtomicU64::new(0),
            write_errors: AtomicU64::new(0),
            write_latency_sum: AtomicU64::new(0),
            write_latency_count: AtomicU64::new(0),
            write_latency_max: AtomicU64::new(0),
            reads_total: AtomicU64::new(0),
            reads_bytes: AtomicU64::new(0),
            read_errors: AtomicU64::new(0),
            read_latency_sum: AtomicU64::new(0),
            read_latency_count: AtomicU64::new(0),
            read_latency_max: AtomicU64::new(0),
            cache_hits: AtomicU64::new(0),
            cache_misses: AtomicU64::new(0),
            memtable_size_bytes: AtomicU64::new(0),
            memtable_entries: AtomicU64::new(0),
            memtable_flushes: AtomicU64::new(0),
            sstable_count: AtomicU64::new(0),
            sstable_size_bytes: AtomicU64::new(0),
            sstable_compactions: AtomicU64::new(0),
            sstable_reads: AtomicU64::new(0),
            sstable_writes: AtomicU64::new(0),
            wal_size_bytes: AtomicU64::new(0),
            wal_segments: AtomicU64::new(0),
            wal_syncs: AtomicU64::new(0),
            wal_records: AtomicU64::new(0),
            retention_deleted_files: AtomicU64::new(0),
            retention_freed_bytes: AtomicU64::new(0),
            retention_runs: AtomicU64::new(0),
            compaction_jobs: AtomicU64::new(0),
            compaction_bytes_read: AtomicU64::new(0),
            compaction_bytes_written: AtomicU64::new(0),
            compaction_duration_ms: AtomicU64::new(0),
            uptime_seconds: AtomicU64::new(0),
            start_time: AtomicI64::new(now),
            active_connections: AtomicU64::new(0),
            total_connections: AtomicU64::new(0),
        }
    }

    pub fn record_write(&self, bytes: u64, latency_us: u64) {
        self.writes_total.fetch_add(1, Ordering::Relaxed);
        self.writes_bytes.fetch_add(bytes, Ordering::Relaxed);
        self.write_latency_sum
            .fetch_add(latency_us, Ordering::Relaxed);
        self.write_latency_count.fetch_add(1, Ordering::Relaxed);
        // Update max
        let current_max = self.write_latency_max.load(Ordering::Relaxed);
        if latency_us > current_max {
            self.write_latency_max.store(latency_us, Ordering::Relaxed);
        }
    }

    pub fn record_write_error(&self) {
        self.write_errors.fetch_add(1, Ordering::Relaxed);
    }

    pub fn record_read(&self, bytes: u64, latency_us: u64) {
        self.reads_total.fetch_add(1, Ordering::Relaxed);
        self.reads_bytes.fetch_add(bytes, Ordering::Relaxed);
        self.read_latency_sum
            .fetch_add(latency_us, Ordering::Relaxed);
        self.read_latency_count.fetch_add(1, Ordering::Relaxed);
        let current_max = self.read_latency_max.load(Ordering::Relaxed);
        if latency_us > current_max {
            self.read_latency_max.store(latency_us, Ordering::Relaxed);
        }
    }

    pub fn record_read_error(&self) {
        self.read_errors.fetch_add(1, Ordering::Relaxed);
    }

    pub fn record_cache_hit(&self) {
        self.cache_hits.fetch_add(1, Ordering::Relaxed);
    }

    pub fn record_cache_miss(&self) {
        self.cache_misses.fetch_add(1, Ordering::Relaxed);
    }

    pub fn record_memtable_flush(&self) {
        self.memtable_flushes.fetch_add(1, Ordering::Relaxed);
    }

    pub fn record_compaction(&self, bytes_read: u64, bytes_written: u64, duration_ms: u64) {
        self.sstable_compactions.fetch_add(1, Ordering::Relaxed);
        self.compaction_jobs.fetch_add(1, Ordering::Relaxed);
        self.compaction_bytes_read
            .fetch_add(bytes_read, Ordering::Relaxed);
        self.compaction_bytes_written
            .fetch_add(bytes_written, Ordering::Relaxed);
        self.compaction_duration_ms
            .fetch_add(duration_ms, Ordering::Relaxed);
    }

    pub fn record_retention(&self, files: u64, bytes: u64) {
        self.retention_deleted_files
            .fetch_add(files, Ordering::Relaxed);
        self.retention_freed_bytes
            .fetch_add(bytes, Ordering::Relaxed);
        self.retention_runs.fetch_add(1, Ordering::Relaxed);
    }

    pub fn update_memtable(&self, size: u64, entries: u64) {
        self.memtable_size_bytes.store(size, Ordering::Relaxed);
        self.memtable_entries.store(entries, Ordering::Relaxed);
    }

    pub fn update_sstable(&self, count: u64, size: u64) {
        self.sstable_count.store(count, Ordering::Relaxed);
        self.sstable_size_bytes.store(size, Ordering::Relaxed);
    }

    pub fn update_wal(&self, size: u64, segments: u64, records: u64) {
        self.wal_size_bytes.store(size, Ordering::Relaxed);
        self.wal_segments.store(segments, Ordering::Relaxed);
        self.wal_records.store(records, Ordering::Relaxed);
    }

    pub fn update_uptime(&self) {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        let start = self.start_time.load(Ordering::Relaxed) as u64;
        self.uptime_seconds.store(now - start, Ordering::Relaxed);
    }

    pub fn connection_opened(&self) {
        self.active_connections.fetch_add(1, Ordering::Relaxed);
        self.total_connections.fetch_add(1, Ordering::Relaxed);
    }

    pub fn connection_closed(&self) {
        self.active_connections.fetch_sub(1, Ordering::Relaxed);
    }

    pub fn avg_write_latency_us(&self) -> f64 {
        let sum = self.write_latency_sum.load(Ordering::Relaxed);
        let count = self.write_latency_count.load(Ordering::Relaxed);
        if count > 0 {
            sum as f64 / count as f64
        } else {
            0.0
        }
    }

    pub fn avg_read_latency_us(&self) -> f64 {
        let sum = self.read_latency_sum.load(Ordering::Relaxed);
        let count = self.read_latency_count.load(Ordering::Relaxed);
        if count > 0 {
            sum as f64 / count as f64
        } else {
            0.0
        }
    }

    pub fn cache_hit_rate(&self) -> f64 {
        let hits = self.cache_hits.load(Ordering::Relaxed);
        let misses = self.cache_misses.load(Ordering::Relaxed);
        let total = hits + misses;
        if total > 0 {
            hits as f64 / total as f64
        } else {
            0.0
        }
    }

    pub fn write_error_rate(&self) -> f64 {
        let errors = self.write_errors.load(Ordering::Relaxed);
        let total = self.writes_total.load(Ordering::Relaxed);
        let denom = errors + total;
        if denom > 0 {
            errors as f64 / denom as f64
        } else {
            0.0
        }
    }

    pub fn read_error_rate(&self) -> f64 {
        let errors = self.read_errors.load(Ordering::Relaxed);
        let total = self.reads_total.load(Ordering::Relaxed);
        if total > 0 {
            errors as f64 / total as f64
        } else {
            0.0
        }
    }

    pub fn to_prometheus(&self) -> String {
        let mut output = String::new();

        // Write metrics
        output.push_str("# HELP helix_storage_writes_total Total number of writes\n");
        output.push_str("# TYPE helix_storage_writes_total counter\n");
        output.push_str(&format!(
            "helix_storage_writes_total {}\n",
            self.writes_total.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_writes_bytes Total bytes written\n");
        output.push_str("# TYPE helix_storage_writes_bytes counter\n");
        output.push_str(&format!(
            "helix_storage_writes_bytes {}\n",
            self.writes_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_write_errors Total write errors\n");
        output.push_str("# TYPE helix_storage_write_errors counter\n");
        output.push_str(&format!(
            "helix_storage_write_errors {}\n",
            self.write_errors.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_write_latency_avg_us Average write latency in microseconds\n",
        );
        output.push_str("# TYPE helix_storage_write_latency_avg_us gauge\n");
        output.push_str(&format!(
            "helix_storage_write_latency_avg_us {}\n",
            self.avg_write_latency_us()
        ));

        output.push_str(
            "# HELP helix_storage_write_latency_max_us Max write latency in microseconds\n",
        );
        output.push_str("# TYPE helix_storage_write_latency_max_us gauge\n");
        output.push_str(&format!(
            "helix_storage_write_latency_max_us {}\n",
            self.write_latency_max.load(Ordering::Relaxed)
        ));

        // Read metrics
        output.push_str("# HELP helix_storage_reads_total Total number of reads\n");
        output.push_str("# TYPE helix_storage_reads_total counter\n");
        output.push_str(&format!(
            "helix_storage_reads_total {}\n",
            self.reads_total.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_reads_bytes Total bytes read\n");
        output.push_str("# TYPE helix_storage_reads_bytes counter\n");
        output.push_str(&format!(
            "helix_storage_reads_bytes {}\n",
            self.reads_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_read_errors Total read errors\n");
        output.push_str("# TYPE helix_storage_read_errors counter\n");
        output.push_str(&format!(
            "helix_storage_read_errors {}\n",
            self.read_errors.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_read_latency_avg_us Average read latency in microseconds\n",
        );
        output.push_str("# TYPE helix_storage_read_latency_avg_us gauge\n");
        output.push_str(&format!(
            "helix_storage_read_latency_avg_us {}\n",
            self.avg_read_latency_us()
        ));

        output.push_str(
            "# HELP helix_storage_read_latency_max_us Max read latency in microseconds\n",
        );
        output.push_str("# TYPE helix_storage_read_latency_max_us gauge\n");
        output.push_str(&format!(
            "helix_storage_read_latency_max_us {}\n",
            self.read_latency_max.load(Ordering::Relaxed)
        ));

        // Cache metrics
        output.push_str("# HELP helix_storage_cache_hits Total cache hits\n");
        output.push_str("# TYPE helix_storage_cache_hits counter\n");
        output.push_str(&format!(
            "helix_storage_cache_hits {}\n",
            self.cache_hits.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_cache_misses Total cache misses\n");
        output.push_str("# TYPE helix_storage_cache_misses counter\n");
        output.push_str(&format!(
            "helix_storage_cache_misses {}\n",
            self.cache_misses.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_cache_hit_rate Cache hit rate\n");
        output.push_str("# TYPE helix_storage_cache_hit_rate gauge\n");
        output.push_str(&format!(
            "helix_storage_cache_hit_rate {}\n",
            self.cache_hit_rate()
        ));

        // Memtable metrics
        output.push_str("# HELP helix_storage_memtable_size_bytes Memtable size in bytes\n");
        output.push_str("# TYPE helix_storage_memtable_size_bytes gauge\n");
        output.push_str(&format!(
            "helix_storage_memtable_size_bytes {}\n",
            self.memtable_size_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_memtable_entries Number of entries in memtable\n");
        output.push_str("# TYPE helix_storage_memtable_entries gauge\n");
        output.push_str(&format!(
            "helix_storage_memtable_entries {}\n",
            self.memtable_entries.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_memtable_flushes Total memtable flushes\n");
        output.push_str("# TYPE helix_storage_memtable_flushes counter\n");
        output.push_str(&format!(
            "helix_storage_memtable_flushes {}\n",
            self.memtable_flushes.load(Ordering::Relaxed)
        ));

        // SSTable metrics
        output.push_str("# HELP helix_storage_sstable_count Number of SSTables\n");
        output.push_str("# TYPE helix_storage_sstable_count gauge\n");
        output.push_str(&format!(
            "helix_storage_sstable_count {}\n",
            self.sstable_count.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_sstable_size_bytes Total SSTable size in bytes\n");
        output.push_str("# TYPE helix_storage_sstable_size_bytes gauge\n");
        output.push_str(&format!(
            "helix_storage_sstable_size_bytes {}\n",
            self.sstable_size_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_sstable_compactions Total SSTable compactions\n");
        output.push_str("# TYPE helix_storage_sstable_compactions counter\n");
        output.push_str(&format!(
            "helix_storage_sstable_compactions {}\n",
            self.sstable_compactions.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_sstable_reads Total SSTable reads\n");
        output.push_str("# TYPE helix_storage_sstable_reads counter\n");
        output.push_str(&format!(
            "helix_storage_sstable_reads {}\n",
            self.sstable_reads.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_sstable_writes Total SSTable writes\n");
        output.push_str("# TYPE helix_storage_sstable_writes counter\n");
        output.push_str(&format!(
            "helix_storage_sstable_writes {}\n",
            self.sstable_writes.load(Ordering::Relaxed)
        ));

        // WAL metrics
        output.push_str("# HELP helix_storage_wal_size_bytes WAL size in bytes\n");
        output.push_str("# TYPE helix_storage_wal_size_bytes gauge\n");
        output.push_str(&format!(
            "helix_storage_wal_size_bytes {}\n",
            self.wal_size_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_wal_segments Number of WAL segments\n");
        output.push_str("# TYPE helix_storage_wal_segments gauge\n");
        output.push_str(&format!(
            "helix_storage_wal_segments {}\n",
            self.wal_segments.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_wal_syncs Total WAL syncs\n");
        output.push_str("# TYPE helix_storage_wal_syncs counter\n");
        output.push_str(&format!(
            "helix_storage_wal_syncs {}\n",
            self.wal_syncs.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_wal_records Total WAL records\n");
        output.push_str("# TYPE helix_storage_wal_records counter\n");
        output.push_str(&format!(
            "helix_storage_wal_records {}\n",
            self.wal_records.load(Ordering::Relaxed)
        ));

        // Retention metrics
        output.push_str(
            "# HELP helix_storage_retention_deleted_files Total files deleted by retention\n",
        );
        output.push_str("# TYPE helix_storage_retention_deleted_files counter\n");
        output.push_str(&format!(
            "helix_storage_retention_deleted_files {}\n",
            self.retention_deleted_files.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_retention_freed_bytes Total bytes freed by retention\n",
        );
        output.push_str("# TYPE helix_storage_retention_freed_bytes counter\n");
        output.push_str(&format!(
            "helix_storage_retention_freed_bytes {}\n",
            self.retention_freed_bytes.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_retention_runs Total retention runs\n");
        output.push_str("# TYPE helix_storage_retention_runs counter\n");
        output.push_str(&format!(
            "helix_storage_retention_runs {}\n",
            self.retention_runs.load(Ordering::Relaxed)
        ));

        // Compaction metrics
        output.push_str("# HELP helix_storage_compaction_jobs Total compaction jobs\n");
        output.push_str("# TYPE helix_storage_compaction_jobs counter\n");
        output.push_str(&format!(
            "helix_storage_compaction_jobs {}\n",
            self.compaction_jobs.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_compaction_bytes_read Total bytes read during compaction\n",
        );
        output.push_str("# TYPE helix_storage_compaction_bytes_read counter\n");
        output.push_str(&format!(
            "helix_storage_compaction_bytes_read {}\n",
            self.compaction_bytes_read.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_compaction_bytes_written Total bytes written during compaction\n",
        );
        output.push_str("# TYPE helix_storage_compaction_bytes_written counter\n");
        output.push_str(&format!(
            "helix_storage_compaction_bytes_written {}\n",
            self.compaction_bytes_written.load(Ordering::Relaxed)
        ));

        output.push_str(
            "# HELP helix_storage_compaction_duration_ms Total compaction duration in ms\n",
        );
        output.push_str("# TYPE helix_storage_compaction_duration_ms counter\n");
        output.push_str(&format!(
            "helix_storage_compaction_duration_ms {}\n",
            self.compaction_duration_ms.load(Ordering::Relaxed)
        ));

        // General metrics
        output.push_str("# HELP helix_storage_uptime_seconds Uptime in seconds\n");
        output.push_str("# TYPE helix_storage_uptime_seconds gauge\n");
        output.push_str(&format!(
            "helix_storage_uptime_seconds {}\n",
            self.uptime_seconds.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_active_connections Active connections\n");
        output.push_str("# TYPE helix_storage_active_connections gauge\n");
        output.push_str(&format!(
            "helix_storage_active_connections {}\n",
            self.active_connections.load(Ordering::Relaxed)
        ));

        output.push_str("# HELP helix_storage_total_connections Total connections\n");
        output.push_str("# TYPE helix_storage_total_connections counter\n");
        output.push_str(&format!(
            "helix_storage_total_connections {}\n",
            self.total_connections.load(Ordering::Relaxed)
        ));

        output
    }

    pub fn to_json(&self) -> String {
        let mut map = BTreeMap::new();
        map.insert("writes_total", self.writes_total.load(Ordering::Relaxed));
        map.insert("writes_bytes", self.writes_bytes.load(Ordering::Relaxed));
        map.insert("write_errors", self.write_errors.load(Ordering::Relaxed));
        map.insert("reads_total", self.reads_total.load(Ordering::Relaxed));
        map.insert("reads_bytes", self.reads_bytes.load(Ordering::Relaxed));
        map.insert("read_errors", self.read_errors.load(Ordering::Relaxed));
        map.insert("cache_hits", self.cache_hits.load(Ordering::Relaxed));
        map.insert("cache_misses", self.cache_misses.load(Ordering::Relaxed));
        map.insert(
            "memtable_size_bytes",
            self.memtable_size_bytes.load(Ordering::Relaxed),
        );
        map.insert(
            "memtable_entries",
            self.memtable_entries.load(Ordering::Relaxed),
        );
        map.insert("sstable_count", self.sstable_count.load(Ordering::Relaxed));
        map.insert(
            "sstable_size_bytes",
            self.sstable_size_bytes.load(Ordering::Relaxed),
        );
        map.insert(
            "wal_size_bytes",
            self.wal_size_bytes.load(Ordering::Relaxed),
        );
        map.insert(
            "uptime_seconds",
            self.uptime_seconds.load(Ordering::Relaxed),
        );

        // Simple JSON serialization
        let mut json = String::from("{");
        let mut first = true;
        for (k, v) in map {
            if !first {
                json.push(',');
            }
            json.push_str(&format!("\"{}\":{}", k, v));
            first = false;
        }
        json.push('}');
        json
    }
}

impl Default for ExtendedStorageMetrics {
    fn default() -> Self {
        Self::new()
    }
}

/// Metrics collector
pub struct MetricsCollector {
    metrics: Arc<ExtendedStorageMetrics>,
    collectors: Vec<Box<dyn MetricCollector>>,
}

/// Metric collector trait
pub trait MetricCollector: Send + Sync {
    fn name(&self) -> &'static str;
    fn collect(&self) -> BTreeMap<String, u64>;
}

impl MetricsCollector {
    pub fn new(metrics: Arc<ExtendedStorageMetrics>) -> Self {
        MetricsCollector {
            metrics,
            collectors: Vec::new(),
        }
    }

    pub fn add_collector(&mut self, collector: Box<dyn MetricCollector>) {
        self.collectors.push(collector);
    }

    pub fn collect_all(&self) -> BTreeMap<String, u64> {
        let mut result = BTreeMap::new();
        for collector in &self.collectors {
            let collected = collector.collect();
            for (k, v) in collected {
                result.insert(k, v);
            }
        }
        result
    }

    pub fn metrics(&self) -> Arc<ExtendedStorageMetrics> {
        self.metrics.clone()
    }
}

/// CPU metric collector
pub struct CPUMetricCollector;

impl MetricCollector for CPUMetricCollector {
    fn name(&self) -> &'static str {
        "cpu"
    }
    fn collect(&self) -> BTreeMap<String, u64> {
        let mut result = BTreeMap::new();
        result.insert("cpu_usage".to_string(), 50);
        result
    }
}

/// Memory metric collector
pub struct MemoryMetricCollector;

impl MetricCollector for MemoryMetricCollector {
    fn name(&self) -> &'static str {
        "memory"
    }
    fn collect(&self) -> BTreeMap<String, u64> {
        let mut result = BTreeMap::new();
        result.insert("memory_usage".to_string(), 1024 * 1024 * 1024);
        result
    }
}

/// Disk metric collector
pub struct DiskMetricCollector;

impl MetricCollector for DiskMetricCollector {
    fn name(&self) -> &'static str {
        "disk"
    }
    fn collect(&self) -> BTreeMap<String, u64> {
        let mut result = BTreeMap::new();
        result.insert("disk_usage".to_string(), 50 * 1024 * 1024 * 1024);
        result.insert("disk_available".to_string(), 100 * 1024 * 1024 * 1024);
        result
    }
}

/// Network metric collector
pub struct NetworkMetricCollector;

impl MetricCollector for NetworkMetricCollector {
    fn name(&self) -> &'static str {
        "network"
    }
    fn collect(&self) -> BTreeMap<String, u64> {
        let mut result = BTreeMap::new();
        result.insert("network_rx_bytes".to_string(), 1000000);
        result.insert("network_tx_bytes".to_string(), 500000);
        result
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extended_metrics() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.record_write(100, 50);
        metrics.record_write(200, 100);

        assert_eq!(metrics.writes_total.load(Ordering::Relaxed), 2);
        assert_eq!(metrics.writes_bytes.load(Ordering::Relaxed), 300);
        assert_eq!(metrics.avg_write_latency_us(), 75.0);
    }

    #[test]
    fn test_cache_hit_rate() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.record_cache_hit();
        metrics.record_cache_hit();
        metrics.record_cache_miss();

        assert!((metrics.cache_hit_rate() - 0.666).abs() < 0.01);
    }

    #[test]
    fn test_prometheus_output() {
        let metrics = ExtendedStorageMetrics::new();
        metrics.record_write(100, 50);

        let output = metrics.to_prometheus();
        assert!(output.contains("helix_storage_writes_total"));
        assert!(output.contains("helix_storage_writes_bytes"));
    }

    #[test]
    fn test_json_output() {
        let metrics = ExtendedStorageMetrics::new();
        metrics.record_write(100, 50);

        let json = metrics.to_json();
        assert!(json.contains("writes_total"));
    }

    #[test]
    fn test_metrics_collector() {
        let metrics = Arc::new(ExtendedStorageMetrics::new());
        let mut collector = MetricsCollector::new(metrics);

        collector.add_collector(Box::new(CPUMetricCollector));
        collector.add_collector(Box::new(MemoryMetricCollector));

        let collected = collector.collect_all();
        assert!(collected.contains_key("cpu_usage"));
        assert!(collected.contains_key("memory_usage"));
    }

    #[test]
    fn test_connection_tracking() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.connection_opened();
        metrics.connection_opened();
        assert_eq!(metrics.active_connections.load(Ordering::Relaxed), 2);
        assert_eq!(metrics.total_connections.load(Ordering::Relaxed), 2);

        metrics.connection_closed();
        assert_eq!(metrics.active_connections.load(Ordering::Relaxed), 1);
    }

    #[test]
    fn test_error_rates() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.record_write(100, 50);
        metrics.record_write_error();

        let rate = metrics.write_error_rate();
        assert!((rate - 0.5).abs() < 0.01);
    }

    #[test]
    fn test_retention_metrics() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.record_retention(5, 1024 * 1024);
        assert_eq!(metrics.retention_deleted_files.load(Ordering::Relaxed), 5);
        assert_eq!(
            metrics.retention_freed_bytes.load(Ordering::Relaxed),
            1024 * 1024
        );
        assert_eq!(metrics.retention_runs.load(Ordering::Relaxed), 1);
    }

    #[test]
    fn test_compaction_metrics() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.record_compaction(1024, 512, 100);
        assert_eq!(metrics.compaction_jobs.load(Ordering::Relaxed), 1);
        assert_eq!(metrics.compaction_bytes_read.load(Ordering::Relaxed), 1024);
        assert_eq!(
            metrics.compaction_bytes_written.load(Ordering::Relaxed),
            512
        );
        assert_eq!(metrics.compaction_duration_ms.load(Ordering::Relaxed), 100);
    }

    #[test]
    fn test_wal_metrics() {
        let metrics = ExtendedStorageMetrics::new();

        metrics.update_wal(1024 * 1024, 5, 1000);
        assert_eq!(metrics.wal_size_bytes.load(Ordering::Relaxed), 1024 * 1024);
        assert_eq!(metrics.wal_segments.load(Ordering::Relaxed), 5);
        assert_eq!(metrics.wal_records.load(Ordering::Relaxed), 1000);
    }
}
