//! Helix Storage Engine - Metrics and Observability Module
//! 
//! This module provides metrics collection and observability for the storage engine.
//! Includes Prometheus-compatible metrics, tracing, and profiling.

use std::collections::BTreeMap;
use std::sync::atomic::{AtomicU64, AtomicI64, Ordering};
use std::sync::Arc;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Storage metrics
pub struct StorageMetrics {
    // Write metrics
    pub writes_total: AtomicU64,
    pub writes_bytes: AtomicU64,
    pub write_errors: AtomicU64,
    pub write_latency_sum: AtomicU64,
    pub write_latency_count: AtomicU64,
    
    // Read metrics
    pub reads_total: AtomicU64,
    pub reads_bytes: AtomicU64,
    pub read_errors: AtomicU64,
    pub read_latency_sum: AtomicU64,
    pub read_latency_count: AtomicU64,
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
    
    // WAL metrics
    pub wal_size_bytes: AtomicU64,
    pub wal_segments: AtomicU64,
    pub wal_syncs: AtomicU64,
    
    // Retention metrics
    pub retention_deleted_files: AtomicU64,
    pub retention_freed_bytes: AtomicU64,
    
    // General
    pub uptime_seconds: AtomicU64,
    pub start_time: AtomicI64,
}

impl StorageMetrics {
    pub fn new() -> Self {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        StorageMetrics {
            writes_total: AtomicU64::new(0),
            writes_bytes: AtomicU64::new(0),
            write_errors: AtomicU64::new(0),
            write_latency_sum: AtomicU64::new(0),
            write_latency_count: AtomicU64::new(0),
            reads_total: AtomicU64::new(0),
            reads_bytes: AtomicU64::new(0),
            read_errors: AtomicU64::new(0),
            read_latency_sum: AtomicU64::new(0),
            read_latency_count: AtomicU64::new(0),
            cache_hits: AtomicU64::new(0),
            cache_misses: AtomicU64::new(0),
            memtable_size_bytes: AtomicU64::new(0),
            memtable_entries: AtomicU64::new(0),
            memtable_flushes: AtomicU64::new(0),
            sstable_count: AtomicU64::new(0),
            sstable_size_bytes: AtomicU64::new(0),
            sstable_compactions: AtomicU64::new(0),
            wal_size_bytes: AtomicU64::new(0),
            wal_segments: AtomicU64::new(0),
            wal_syncs: AtomicU64::new(0),
            retention_deleted_files: AtomicU64::new(0),
            retention_freed_bytes: AtomicU64::new(0),
            uptime_seconds: AtomicU64::new(0),
            start_time: AtomicI64::new(now),
        }
    }
    
    pub fn record_write(&self, bytes: u64, latency_us: u64) {
        self.writes_total.fetch_add(1, Ordering::Relaxed);
        self.writes_bytes.fetch_add(bytes, Ordering::Relaxed);
        self.write_latency_sum.fetch_add(latency_us, Ordering::Relaxed);
        self.write_latency_count.fetch_add(1, Ordering::Relaxed);
    }
    
    pub fn record_write_error(&self) {
        self.write_errors.fetch_add(1, Ordering::Relaxed);
    }
    
    pub fn record_read(&self, bytes: u64, latency_us: u64) {
        self.reads_total.fetch_add(1, Ordering::Relaxed);
        self.reads_bytes.fetch_add(bytes, Ordering::Relaxed);
        self.read_latency_sum.fetch_add(latency_us, Ordering::Relaxed);
        self.read_latency_count.fetch_add(1, Ordering::Relaxed);
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
    
    pub fn record_compaction(&self) {
        self.sstable_compactions.fetch_add(1, Ordering::Relaxed);
    }
    
    pub fn record_retention(&self, files: u64, bytes: u64) {
        self.retention_deleted_files.fetch_add(files, Ordering::Relaxed);
        self.retention_freed_bytes.fetch_add(bytes, Ordering::Relaxed);
    }
    
    pub fn update_memtable(&self, size: u64, entries: u64) {
        self.memtable_size_bytes.store(size, Ordering::Relaxed);
        self.memtable_entries.store(entries, Ordering::Relaxed);
    }
    
    pub fn update_sstable(&self, count: u64, size: u64) {
        self.sstable_count.store(count, Ordering::Relaxed);
        self.sstable_size_bytes.store(size, Ordering::Relaxed);
    }
    
    pub fn update_wal(&self, size: u64, segments: u64) {
        self.wal_size_bytes.store(size, Ordering::Relaxed);
        self.wal_segments.store(segments, Ordering::Relaxed);
    }
    
    pub fn update_uptime(&self) {
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs();
        let start = self.start_time.load(Ordering::Relaxed) as u64;
        self.uptime_seconds.store(now - start as u64, Ordering::Relaxed);
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
    
    pub fn to_prometheus(&self) -> String {
        let mut output = String::new();
        
        output.push_str("# HELP helix_storage_writes_total Total number of writes\n");
        output.push_str("# TYPE helix_storage_writes_total counter\n");
        output.push_str(&format!("helix_storage_writes_total {}\n", self.writes_total.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_writes_bytes Total bytes written\n");
        output.push_str("# TYPE helix_storage_writes_bytes counter\n");
        output.push_str(&format!("helix_storage_writes_bytes {}\n", self.writes_bytes.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_write_errors Total write errors\n");
        output.push_str("# TYPE helix_storage_write_errors counter\n");
        output.push_str(&format!("helix_storage_write_errors {}\n", self.write_errors.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_read_latency_us Average read latency in microseconds\n");
        output.push_str("# TYPE helix_storage_read_latency_us gauge\n");
        output.push_str(&format!("helix_storage_read_latency_us {}\n", self.avg_read_latency_us()));
        
        output.push_str("# HELP helix_storage_write_latency_us Average write latency in microseconds\n");
        output.push_str("# TYPE helix_storage_write_latency_us gauge\n");
        output.push_str(&format!("helix_storage_write_latency_us {}\n", self.avg_write_latency_us()));
        
        output.push_str("# HELP helix_storage_cache_hit_rate Cache hit rate\n");
        output.push_str("# TYPE helix_storage_cache_hit_rate gauge\n");
        output.push_str(&format!("helix_storage_cache_hit_rate {}\n", self.cache_hit_rate()));
        
        output.push_str("# HELP helix_storage_memtable_size_bytes Memtable size in bytes\n");
        output.push_str("# TYPE helix_storage_memtable_size_bytes gauge\n");
        output.push_str(&format!("helix_storage_memtable_size_bytes {}\n", self.memtable_size_bytes.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_memtable_entries Number of entries in memtable\n");
        output.push_str("# TYPE helix_storage_memtable_entries gauge\n");
        output.push_str(&format!("helix_storage_memtable_entries {}\n", self.memtable_entries.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_sstable_count Number of SSTables\n");
        output.push_str("# TYPE helix_storage_sstable_count gauge\n");
        output.push_str(&format!("helix_storage_sstable_count {}\n", self.sstable_count.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_sstable_size_bytes Total SSTable size in bytes\n");
        output.push_str("# TYPE helix_storage_sstable_size_bytes gauge\n");
        output.push_str(&format!("helix_storage_sstable_size_bytes {}\n", self.sstable_size_bytes.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_wal_size_bytes WAL size in bytes\n");
        output.push_str("# TYPE helix_storage_wal_size_bytes gauge\n");
        output.push_str(&format!("helix_storage_wal_size_bytes {}\n", self.wal_size_bytes.load(Ordering::Relaxed)));
        
        output.push_str("# HELP helix_storage_uptime_seconds Uptime in seconds\n");
        output.push_str("# TYPE helix_storage_uptime_seconds gauge\n");
        output.push_str(&format!("helix_storage_uptime_seconds {}\n", self.uptime_seconds.load(Ordering::Relaxed)));
        
        output
    }
}

/// Metrics registry
pub struct MetricsRegistry {
    metrics: Arc<StorageMetrics>,
    custom: BTreeMap<String, Arc<AtomicU64>>,
}

impl MetricsRegistry {
    pub fn new() -> Self {
        MetricsRegistry {
            metrics: Arc::new(StorageMetrics::new()),
            custom: BTreeMap::new(),
        }
    }
    
    pub fn storage(&self) -> Arc<StorageMetrics> {
        self.metrics.clone()
    }
    
    pub fn register_custom(&mut self, name: String) {
        self.custom.insert(name, Arc::new(AtomicU64::new(0)));
    }
    
    pub fn increment_custom(&self, name: &str) {
        if let Some(counter) = self.custom.get(name) {
            counter.fetch_add(1, Ordering::Relaxed);
        }
    }
    
    pub fn get_custom(&self, name: &str) -> u64 {
        self.custom.get(name).map(|c| c.load(Ordering::Relaxed)).unwrap_or(0)
    }
}

/// Tracing span
#[derive(Debug, Clone)]
pub struct Span {
    pub trace_id: String,
    pub span_id: String,
    pub parent_id: Option<String>,
    pub operation: String,
    pub start_time: i64,
    pub end_time: Option<i64>,
    pub tags: BTreeMap<String, String>,
    pub logs: Vec<SpanLog>,
}

#[derive(Debug, Clone)]
pub struct SpanLog {
    pub timestamp: i64,
    pub message: String,
    pub fields: BTreeMap<String, String>,
}

impl Span {
    pub fn new(trace_id: String, operation: String) -> Self {
        Span {
            trace_id,
            span_id: format!("span-{}", SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos()),
            parent_id: None,
            operation,
            start_time: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis() as i64,
            end_time: None,
            tags: BTreeMap::new(),
            logs: Vec::new(),
        }
    }
    
    pub fn with_parent(mut self, parent_id: String) -> Self {
        self.parent_id = Some(parent_id);
        self
    }
    
    pub fn add_tag(&mut self, key: String, value: String) {
        self.tags.insert(key, value);
    }
    
    pub fn add_log(&mut self, message: String) {
        self.logs.push(SpanLog {
            timestamp: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis() as i64,
            message,
            fields: BTreeMap::new(),
        });
    }
    
    pub fn finish(&mut self) {
        self.end_time = Some(SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis() as i64);
    }
    
    pub fn duration_ms(&self) -> i64 {
        match self.end_time {
            Some(end) => end - self.start_time,
            None => 0,
        }
    }
}

/// Tracer
pub struct Tracer {
    spans: Arc<Mutex<Vec<Span>>>,
    service_name: String,
}

impl Tracer {
    pub fn new(service_name: String) -> Self {
        Tracer {
            spans: Arc::new(Mutex::new(Vec::new())),
            service_name,
        }
    }
    
    pub fn start_span(&self, operation: String) -> Span {
        let trace_id = format!("trace-{}", SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos());
        Span::new(trace_id, operation)
    }
    
    pub fn record_span(&self, span: Span) {
        let mut spans = self.spans.lock().unwrap();
        spans.push(span);
    }
    
    pub fn get_spans(&self) -> Vec<Span> {
        let spans = self.spans.lock().unwrap();
        spans.clone()
    }
    
    pub fn clear(&self) {
        let mut spans = self.spans.lock().unwrap();
        spans.clear();
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_storage_metrics() {
        let metrics = StorageMetrics::new();
        
        metrics.record_write(100, 50);
        metrics.record_write(200, 100);
        
        assert_eq!(metrics.writes_total.load(Ordering::Relaxed), 2);
        assert_eq!(metrics.writes_bytes.load(Ordering::Relaxed), 300);
        assert_eq!(metrics.avg_write_latency_us(), 75.0);
    }
    
    #[test]
    fn test_cache_hit_rate() {
        let metrics = StorageMetrics::new();
        
        metrics.record_cache_hit();
        metrics.record_cache_hit();
        metrics.record_cache_miss();
        
        assert!((metrics.cache_hit_rate() - 0.666).abs() < 0.01);
    }
    
    #[test]
    fn test_prometheus_output() {
        let metrics = StorageMetrics::new();
        metrics.record_write(100, 50);
        
        let output = metrics.to_prometheus();
        assert!(output.contains("helix_storage_writes_total"));
        assert!(output.contains("helix_storage_writes_bytes"));
    }
    
    #[test]
    fn test_span() {
        let mut span = Span::new("trace-1".to_string(), "test_op".to_string());
        span.add_tag("key".to_string(), "value".to_string());
        span.add_log("test log".to_string());
        span.finish();
        
        assert!(span.duration_ms() >= 0);
        assert_eq!(span.tags.get("key"), Some(&"value".to_string()));
        assert_eq!(span.logs.len(), 1);
    }
    
    #[test]
    fn test_tracer() {
        let tracer = Tracer::new("test-service".to_string());
        let mut span = tracer.start_span("operation".to_string());
        span.finish();
        tracer.record_span(span);
        
        assert_eq!(tracer.get_spans().len(), 1);
    }
    
    #[test]
    fn test_metrics_registry() {
        let mut registry = MetricsRegistry::new();
        registry.register_custom("custom_counter".to_string());
        registry.increment_custom("custom_counter");
        registry.increment_custom("custom_counter");
        
        assert_eq!(registry.get_custom("custom_counter"), 2);
    }
}