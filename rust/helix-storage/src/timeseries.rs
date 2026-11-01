//! Helix Storage Engine - Time Series Module
//! 
//! This module provides time-series specific functionality
//! including time-based indexing, aggregation, and downsampling.

use std::collections::BTreeMap;
use std::sync::RwLock;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// Time series point
#[derive(Debug, Clone, PartialEq)]
pub struct TimeSeriesPoint {
    pub timestamp: i64,
    pub value: f64,
    pub labels: BTreeMap<String, String>,
}

impl TimeSeriesPoint {
    pub fn new(timestamp: i64, value: f64) -> Self {
        TimeSeriesPoint {
            timestamp,
            value,
            labels: BTreeMap::new(),
        }
    }
    
    pub fn with_labels(mut self, labels: BTreeMap<String, String>) -> Self {
        self.labels = labels;
        self
    }
}

/// Time series
#[derive(Debug, Clone)]
pub struct TimeSeries {
    pub name: String,
    pub labels: BTreeMap<String, String>,
    pub points: Vec<TimeSeriesPoint>,
}

impl TimeSeries {
    pub fn new(name: String, labels: BTreeMap<String, String>) -> Self {
        TimeSeries {
            name,
            labels,
            points: Vec::new(),
        }
    }
    
    pub fn add_point(&mut self, point: TimeSeriesPoint) {
        self.points.push(point);
        self.points.sort_by_key(|p| p.timestamp);
    }
    
    pub fn get_range(&self, start: i64, end: i64) -> Vec<&TimeSeriesPoint> {
        self.points.iter()
            .filter(|p| p.timestamp >= start && p.timestamp <= end)
            .collect()
    }
    
    pub fn last(&self) -> Option<&TimeSeriesPoint> {
        self.points.last()
    }
    
    pub fn first(&self) -> Option<&TimeSeriesPoint> {
        self.points.first()
    }
    
    pub fn count(&self) -> usize {
        self.points.len()
    }
    
    pub fn min(&self) -> Option<f64> {
        self.points.iter().map(|p| p.value).fold(None, |acc, v| {
            match acc {
                None => Some(v),
                Some(min) => Some(min.min(v)),
            }
        })
    }
    
    pub fn max(&self) -> Option<f64> {
        self.points.iter().map(|p| p.value).fold(None, |acc, v| {
            match acc {
                None => Some(v),
                Some(max) => Some(max.max(v)),
            }
        })
    }
    
    pub fn avg(&self) -> Option<f64> {
        if self.points.is_empty() {
            return None;
        }
        let sum: f64 = self.points.iter().map(|p| p.value).sum();
        Some(sum / self.points.len() as f64)
    }
    
    pub fn sum(&self) -> f64 {
        self.points.iter().map(|p| p.value).sum()
    }
    
    pub fn rate(&self, window: i64) -> Vec<TimeSeriesPoint> {
        let mut result = Vec::new();
        let mut prev: Option<&TimeSeriesPoint> = None;
        
        for point in &self.points {
            if let Some(p) = prev {
                let time_diff = (point.timestamp - p.timestamp) as f64;
                if time_diff > 0.0 {
                    let value_diff = point.value - p.value;
                    let rate = value_diff / time_diff;
                    result.push(TimeSeriesPoint::new(point.timestamp, rate).with_labels(point.labels.clone()));
                }
            }
            prev = Some(point);
        }
        
        result
    }
}

/// Time series database
pub struct TimeSeriesDB {
    series: RwLock<BTreeMap<String, TimeSeries>>,
    index: TimeSeriesIndex,
}

/// Time series index
pub struct TimeSeriesIndex {
    by_name: RwLock<BTreeMap<String, Vec<String>>>,
    by_label: RwLock<BTreeMap<String, BTreeMap<String, Vec<String>>>>,
}

impl TimeSeriesIndex {
    pub fn new() -> Self {
        TimeSeriesIndex {
            by_name: RwLock::new(BTreeMap::new()),
            by_label: RwLock::new(BTreeMap::new()),
        }
    }
    
    pub fn add(&self, name: &str, labels: &BTreeMap<String, String>) {
        let mut by_name = self.by_name.write().unwrap();
        by_name.entry(name.to_string()).or_insert_with(Vec::new).push(name.to_string());
        
        let mut by_label = self.by_label.write().unwrap();
        for (key, value) in labels {
            by_label
                .entry(key.clone())
                .or_insert_with(BTreeMap::new)
                .entry(value.clone())
                .or_insert_with(Vec::new)
                .push(name.to_string());
        }
    }
    
    pub fn find_by_name(&self, name: &str) -> Vec<String> {
        let by_name = self.by_name.read().unwrap();
        by_name.get(name).cloned().unwrap_or_default()
    }
    
    pub fn find_by_label(&self, key: &str, value: &str) -> Vec<String> {
        let by_label = self.by_label.read().unwrap();
        by_label.get(key)
            .and_then(|values| values.get(value))
            .cloned()
            .unwrap_or_default()
    }
}

impl TimeSeriesDB {
    pub fn new() -> Self {
        TimeSeriesDB {
            series: RwLock::new(BTreeMap::new()),
            index: TimeSeriesIndex::new(),
        }
    }
    
    pub fn create_series(&self, name: String, labels: BTreeMap<String, String>) -> &str {
        let mut series = self.series.write().unwrap();
        let key = format!("{}:{}", name, labels.iter().map(|(k, v)| format!("{}={}", k, v)).collect::<Vec<_>>().join(","));
        series.insert(key.clone(), TimeSeries::new(name, labels));
        &key
    }
    
    pub fn add_point(&self, name: &str, labels: BTreeMap<String, String>, point: TimeSeriesPoint) {
        let mut series = self.series.write().unwrap();
        let key = format!("{}:{}", name, labels.iter().map(|(k, v)| format!("{}={}", k, v)).collect::<Vec<_>>().join(","));
        if let Some(ts) = series.get_mut(&key) {
            ts.add_point(point);
        } else {
            let mut ts = TimeSeries::new(name.to_string(), labels.clone());
            ts.add_point(point);
            series.insert(key, ts);
            self.index.add(name, &labels);
        }
    }
    
    pub fn get_series(&self, name: &str, labels: &BTreeMap<String, String>) -> Option<TimeSeries> {
        let series = self.series.read().unwrap();
        let key = format!("{}:{}", name, labels.iter().map(|(k, v)| format!("{}={}", k, v)).collect::<Vec<_>>().join(","));
        series.get(&key).cloned()
    }
    
    pub fn query(&self, name: &str, start: i64, end: i64) -> Vec<TimeSeriesPoint> {
        let series = self.series.read().unwrap();
        let mut result = Vec::new();
        
        for (_, ts) in series.iter() {
            if ts.name == name {
                result.extend(ts.get_range(start, end).into_iter().cloned());
            }
        }
        
        result
    }
    
    pub fn list_series(&self) -> Vec<String> {
        let series = self.series.read().unwrap();
        series.keys().cloned().collect()
    }
    
    pub fn series_count(&self) -> usize {
        let series = self.series.read().unwrap();
        series.len()
    }
}

/// Downsampler for time series
pub struct Downsampler {
    strategy: DownsampleStrategy,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum DownsampleStrategy {
    Avg,
    Min,
    Max,
    Sum,
    Count,
    Last,
    First,
}

impl Downsampler {
    pub fn new(strategy: DownsampleStrategy) -> Self {
        Downsampler { strategy }
    }
    
    pub fn downsample(&self, series: &TimeSeries, interval: i64) -> TimeSeries {
        let mut result = TimeSeries::new(series.name.clone(), series.labels.clone());
        
        if series.points.is_empty() {
            return result;
        }
        
        let mut current_bucket: Vec<&TimeSeriesPoint> = Vec::new();
        let mut bucket_start = series.points[0].timestamp;
        
        for point in &series.points {
            if point.timestamp - bucket_start >= interval {
                // Process current bucket
                if !current_bucket.is_empty() {
                    let value = self.aggregate_bucket(&current_bucket);
                    result.add_point(TimeSeriesPoint::new(bucket_start, value).with_labels(point.labels.clone()));
                }
                current_bucket.clear();
                bucket_start = point.timestamp;
            }
            current_bucket.push(point);
        }
        
        // Process last bucket
        if !current_bucket.is_empty() {
            let value = self.aggregate_bucket(&current_bucket);
            result.add_point(TimeSeriesPoint::new(bucket_start, value));
        }
        
        result
    }
    
    fn aggregate_bucket(&self, points: &[&TimeSeriesPoint]) -> f64 {
        match self.strategy {
            DownsampleStrategy::Avg => {
                let sum: f64 = points.iter().map(|p| p.value).sum();
                sum / points.len() as f64
            }
            DownsampleStrategy::Min => {
                points.iter().map(|p| p.value).fold(f64::MAX, f64::min)
            }
            DownsampleStrategy::Max => {
                points.iter().map(|p| p.value).fold(f64::MIN, f64::max)
            }
            DownsampleStrategy::Sum => {
                points.iter().map(|p| p.value).sum()
            }
            DownsampleStrategy::Count => {
                points.len() as f64
            }
            DownsampleStrategy::Last => {
                points.last().map(|p| p.value).unwrap_or(0.0)
            }
            DownsampleStrategy::First => {
                points.first().map(|p| p.value).unwrap_or(0.0)
            }
        }
    }
}

/// Time series aggregator
pub struct TimeSeriesAggregator;

impl TimeSeriesAggregator {
    pub fn aggregate(series: &[TimeSeries], op: &str) -> TimeSeries {
        if series.is_empty() {
            return TimeSeries::new("".to_string(), BTreeMap::new());
        }
        
        let name = series[0].name.clone();
        let labels = series[0].labels.clone();
        let mut result = TimeSeries::new(name, labels);
        
        // Group by timestamp
        let mut by_timestamp: BTreeMap<i64, Vec<f64>> = BTreeMap::new();
        for ts in series {
            for point in &ts.points {
                by_timestamp.entry(point.timestamp).or_default().push(point.value);
            }
        }
        
        for (timestamp, values) in by_timestamp {
            let value = match op {
                "sum" => values.iter().sum(),
                "avg" => values.iter().sum::<f64>() / values.len() as f64,
                "min" => values.iter().fold(f64::MAX, |a, &b| a.min(b)),
                "max" => values.iter().fold(f64::MIN, |a, &b| a.max(b)),
                "count" => values.len() as f64,
                _ => values.iter().sum(),
            };
            result.add_point(TimeSeriesPoint::new(timestamp, value));
        }
        
        result
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_timeseries_point() {
        let point = TimeSeriesPoint::new(1000, 42.0);
        assert_eq!(point.timestamp, 1000);
        assert_eq!(point.value, 42.0);
    }
    
    #[test]
    fn test_timeseries_add_point() {
        let mut ts = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts.add_point(TimeSeriesPoint::new(1000, 50.0));
        ts.add_point(TimeSeriesPoint::new(1010, 55.0));
        ts.add_point(TimeSeriesPoint::new(1005, 52.0));
        
        assert_eq!(ts.count(), 3);
        assert_eq!(ts.points[0].timestamp, 1000);
        assert_eq!(ts.points[1].timestamp, 1005);
        assert_eq!(ts.points[2].timestamp, 1010);
    }
    
    #[test]
    fn test_timeseries_range() {
        let mut ts = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts.add_point(TimeSeriesPoint::new(1000, 50.0));
        ts.add_point(TimeSeriesPoint::new(1010, 55.0));
        ts.add_point(TimeSeriesPoint::new(1020, 60.0));
        
        let range = ts.get_range(1005, 1015);
        assert_eq!(range.len(), 1);
        assert_eq!(range[0].timestamp, 1010);
    }
    
    #[test]
    fn test_timeseries_min_max_avg() {
        let mut ts = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts.add_point(TimeSeriesPoint::new(1000, 10.0));
        ts.add_point(TimeSeriesPoint::new(1010, 20.0));
        ts.add_point(TimeSeriesPoint::new(1020, 30.0));
        
        assert_eq!(ts.min(), Some(10.0));
        assert_eq!(ts.max(), Some(30.0));
        assert_eq!(ts.avg(), Some(20.0));
        assert_eq!(ts.sum(), 60.0);
    }
    
    #[test]
    fn test_timeseries_rate() {
        let mut ts = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts.add_point(TimeSeriesPoint::new(1000, 10.0));
        ts.add_point(TimeSeriesPoint::new(1010, 20.0));
        
        let rate = ts.rate(10);
        assert_eq!(rate.len(), 1);
        assert_eq!(rate[0].value, 1.0);
    }
    
    #[test]
    fn test_timeseries_db() {
        let db = TimeSeriesDB::new();
        db.add_point("cpu", BTreeMap::new(), TimeSeriesPoint::new(1000, 50.0));
        db.add_point("cpu", BTreeMap::new(), TimeSeriesPoint::new(1010, 55.0));
        
        let series = db.get_series("cpu", &BTreeMap::new());
        assert!(series.is_some());
        assert_eq!(series.unwrap().count(), 2);
    }
    
    #[test]
    fn test_timeseries_db_query() {
        let db = TimeSeriesDB::new();
        db.add_point("cpu", BTreeMap::new(), TimeSeriesPoint::new(1000, 50.0));
        db.add_point("cpu", BTreeMap::new(), TimeSeriesPoint::new(1010, 55.0));
        db.add_point("cpu", BTreeMap::new(), TimeSeriesPoint::new(1020, 60.0));
        
        let result = db.query("cpu", 1005, 1015);
        assert_eq!(result.len(), 1);
    }
    
    #[test]
    fn test_downsampler() {
        let mut ts = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts.add_point(TimeSeriesPoint::new(1000, 10.0));
        ts.add_point(TimeSeriesPoint::new(1010, 20.0));
        ts.add_point(TimeSeriesPoint::new(1020, 30.0));
        ts.add_point(TimeSeriesPoint::new(1030, 40.0));
        
        let downsampler = Downsampler::new(DownsampleStrategy::Avg);
        let result = downsampler.downsample(&ts, 20);
        
        assert!(!result.points.is_empty());
    }
    
    #[test]
    fn test_timeseries_aggregator() {
        let mut ts1 = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts1.add_point(TimeSeriesPoint::new(1000, 10.0));
        ts1.add_point(TimeSeriesPoint::new(1010, 20.0));
        
        let mut ts2 = TimeSeries::new("cpu".to_string(), BTreeMap::new());
        ts2.add_point(TimeSeriesPoint::new(1000, 30.0));
        ts2.add_point(TimeSeriesPoint::new(1010, 40.0));
        
        let result = TimeSeriesAggregator::aggregate(&[ts1, ts2], "sum");
        assert_eq!(result.points.len(), 2);
        assert_eq!(result.points[0].value, 40.0);
        assert_eq!(result.points[1].value, 60.0);
    }
    
    #[test]
    fn test_timeseries_index() {
        let index = TimeSeriesIndex::new();
        let labels = BTreeMap::from([("host".to_string(), "server1".to_string())]);
        index.add("cpu", &labels);
        
        let found = index.find_by_name("cpu");
        assert_eq!(found.len(), 1);
        
        let found = index.find_by_label("host", "server1");
        assert_eq!(found.len(), 1);
    }
}