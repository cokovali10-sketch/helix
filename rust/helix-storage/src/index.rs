//! Helix Storage Engine - Index Module
//!
//! This module provides indexing capabilities for the TSDB.
//! Includes B-tree index, inverted index, and time-series index.

use std::collections::{BTreeMap, BTreeSet};
use std::sync::RwLock;

/// Index type
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum IndexType {
    BTree,
    Hash,
    Inverted,
    TimeSeries,
}

/// Index entry
#[derive(Debug, Clone)]
pub struct IndexEntry {
    pub key: Vec<u8>,
    pub value: Vec<u8>,
    pub timestamp: i64,
    pub series_id: u64,
}

/// B-tree index
pub struct BTreeIndex {
    data: RwLock<BTreeMap<Vec<u8>, IndexEntry>>,
    max_entries: usize,
}

impl BTreeIndex {
    pub fn new(max_entries: usize) -> Self {
        BTreeIndex {
            data: RwLock::new(BTreeMap::new()),
            max_entries,
        }
    }

    pub fn insert(&self, key: Vec<u8>, entry: IndexEntry) {
        let mut data = self.data.write().unwrap();
        data.insert(key, entry);
    }

    pub fn get(&self, key: &[u8]) -> Option<IndexEntry> {
        let data = self.data.read().unwrap();
        data.get(key).cloned()
    }

    pub fn range_scan(&self, start: &[u8], end: &[u8]) -> Vec<IndexEntry> {
        let data = self.data.read().unwrap();
        data.range(start.to_vec()..=end.to_vec())
            .map(|(_, v)| v.clone())
            .collect()
    }

    pub fn delete(&self, key: &[u8]) -> bool {
        let mut data = self.data.write().unwrap();
        data.remove(key).is_some()
    }

    pub fn len(&self) -> usize {
        let data = self.data.read().unwrap();
        data.len()
    }

    pub fn is_empty(&self) -> bool {
        self.len() == 0
    }

    pub fn clear(&self) {
        let mut data = self.data.write().unwrap();
        data.clear();
    }

    pub fn keys(&self) -> Vec<Vec<u8>> {
        let data = self.data.read().unwrap();
        data.keys().cloned().collect()
    }

    pub fn min_key(&self) -> Option<Vec<u8>> {
        let data = self.data.read().unwrap();
        data.keys().next().cloned()
    }

    pub fn max_key(&self) -> Option<Vec<u8>> {
        let data = self.data.read().unwrap();
        data.keys().last().cloned()
    }
}

/// Inverted index for label-based lookup
pub struct InvertedIndex {
    // label name -> label value -> set of series IDs
    index: RwLock<BTreeMap<String, BTreeMap<String, BTreeSet<u64>>>>,
    // series ID -> labels
    series_labels: RwLock<BTreeMap<u64, BTreeMap<String, String>>>,
}

impl InvertedIndex {
    pub fn new() -> Self {
        InvertedIndex {
            index: RwLock::new(BTreeMap::new()),
            series_labels: RwLock::new(BTreeMap::new()),
        }
    }

    pub fn add_series(&self, series_id: u64, labels: BTreeMap<String, String>) {
        let mut series = self.series_labels.write().unwrap();
        series.insert(series_id, labels.clone());

        let mut index = self.index.write().unwrap();
        for (key, value) in &labels {
            index
                .entry(key.clone())
                .or_default()
                .entry(value.clone())
                .or_default()
                .insert(series_id);
        }
    }

    pub fn remove_series(&self, series_id: u64) {
        let mut series = self.series_labels.write().unwrap();
        if let Some(labels) = series.remove(&series_id) {
            let mut index = self.index.write().unwrap();
            for (key, value) in &labels {
                if let Some(values) = index.get_mut(key) {
                    if let Some(ids) = values.get_mut(value) {
                        ids.remove(&series_id);
                    }
                }
            }
        }
    }

    pub fn lookup(&self, label: &str, value: &str) -> Vec<u64> {
        let index = self.index.read().unwrap();
        index
            .get(label)
            .and_then(|values| values.get(value))
            .map(|ids| ids.iter().cloned().collect())
            .unwrap_or_default()
    }

    pub fn lookup_eq(&self, matchers: &BTreeMap<String, String>) -> Vec<u64> {
        if matchers.is_empty() {
            let series = self.series_labels.read().unwrap();
            return series.keys().cloned().collect();
        }

        let mut result: Option<BTreeSet<u64>> = None;

        for (key, value) in matchers {
            let ids: BTreeSet<u64> = self.lookup(key, value).into_iter().collect();
            match &mut result {
                None => result = Some(ids),
                Some(existing) => {
                    existing.retain(|id| ids.contains(id));
                }
            }
        }

        result.map(|s| s.into_iter().collect()).unwrap_or_default()
    }

    pub fn label_values(&self, label: &str) -> Vec<String> {
        let index = self.index.read().unwrap();
        index
            .get(label)
            .map(|values| values.keys().cloned().collect())
            .unwrap_or_default()
    }

    pub fn label_names(&self) -> Vec<String> {
        let index = self.index.read().unwrap();
        index.keys().cloned().collect()
    }

    pub fn series_count(&self) -> usize {
        let series = self.series_labels.read().unwrap();
        series.len()
    }

    pub fn series_labels(&self, series_id: u64) -> Option<BTreeMap<String, String>> {
        let series = self.series_labels.read().unwrap();
        series.get(&series_id).cloned()
    }
}

impl Default for InvertedIndex {
    fn default() -> Self {
        Self::new()
    }
}

/// Time series index for time-based queries
pub struct TimeSeriesIndex {
    // timestamp -> set of series IDs
    by_time: RwLock<BTreeMap<i64, BTreeSet<u64>>>,
    // series ID -> timestamps
    series_times: RwLock<BTreeMap<u64, Vec<i64>>>,
}

impl TimeSeriesIndex {
    pub fn new() -> Self {
        TimeSeriesIndex {
            by_time: RwLock::new(BTreeMap::new()),
            series_times: RwLock::new(BTreeMap::new()),
        }
    }

    pub fn add(&self, series_id: u64, timestamp: i64) {
        let mut by_time = self.by_time.write().unwrap();
        by_time.entry(timestamp).or_default().insert(series_id);

        let mut series_times = self.series_times.write().unwrap();
        series_times.entry(series_id).or_default().push(timestamp);
    }

    pub fn remove(&self, series_id: u64, timestamp: i64) {
        let mut by_time = self.by_time.write().unwrap();
        if let Some(ids) = by_time.get_mut(&timestamp) {
            ids.remove(&series_id);
            if ids.is_empty() {
                by_time.remove(&timestamp);
            }
        }

        let mut series_times = self.series_times.write().unwrap();
        if let Some(times) = series_times.get_mut(&series_id) {
            times.retain(|&t| t != timestamp);
        }
    }

    pub fn time_range(&self, start: i64, end: i64) -> Vec<u64> {
        let by_time = self.by_time.read().unwrap();
        let mut result = BTreeSet::new();

        for (_ts, ids) in by_time.range(start..=end) {
            for id in ids {
                result.insert(*id);
            }
        }

        result.into_iter().collect()
    }

    pub fn series_timestamps(&self, series_id: u64) -> Vec<i64> {
        let series_times = self.series_times.read().unwrap();
        series_times.get(&series_id).cloned().unwrap_or_default()
    }

    pub fn min_timestamp(&self) -> Option<i64> {
        let by_time = self.by_time.read().unwrap();
        by_time.keys().next().cloned()
    }

    pub fn max_timestamp(&self) -> Option<i64> {
        let by_time = self.by_time.read().unwrap();
        by_time.keys().last().cloned()
    }

    pub fn count(&self) -> usize {
        let by_time = self.by_time.read().unwrap();
        by_time.values().map(|v| v.len()).sum()
    }
}

impl Default for TimeSeriesIndex {
    fn default() -> Self {
        Self::new()
    }
}

/// Composite index combining multiple index types
pub struct CompositeIndex {
    pub btree: BTreeIndex,
    pub inverted: InvertedIndex,
    pub timeseries: TimeSeriesIndex,
}

impl CompositeIndex {
    pub fn new() -> Self {
        CompositeIndex {
            btree: BTreeIndex::new(100000),
            inverted: InvertedIndex::new(),
            timeseries: TimeSeriesIndex::new(),
        }
    }

    pub fn add(
        &self,
        series_id: u64,
        key: Vec<u8>,
        value: Vec<u8>,
        labels: BTreeMap<String, String>,
        timestamp: i64,
    ) {
        self.btree.insert(
            key.clone(),
            IndexEntry {
                key,
                value,
                timestamp,
                series_id,
            },
        );
        self.inverted.add_series(series_id, labels);
        self.timeseries.add(series_id, timestamp);
    }

    pub fn remove(&self, series_id: u64, key: &[u8], timestamp: i64) {
        self.btree.delete(key);
        self.inverted.remove_series(series_id);
        self.timeseries.remove(series_id, timestamp);
    }

    pub fn search(
        &self,
        matchers: &BTreeMap<String, String>,
        time_range: Option<(i64, i64)>,
    ) -> Vec<u64> {
        let mut result = self.inverted.lookup_eq(matchers);

        if let Some((start, end)) = time_range {
            let time_ids: Vec<u64> = self.timeseries.time_range(start, end);
            let time_set: std::collections::BTreeSet<u64> = time_ids.into_iter().collect();
            result.retain(|id| time_set.contains(id));
        }

        result
    }
}

impl Default for CompositeIndex {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_btree_index() {
        let index = BTreeIndex::new(1000);

        let entry = IndexEntry {
            key: b"key1".to_vec(),
            value: b"value1".to_vec(),
            timestamp: 1000,
            series_id: 1,
        };
        index.insert(b"key1".to_vec(), entry);

        assert!(index.get(b"key1").is_some());
        assert_eq!(index.len(), 1);
    }

    #[test]
    fn test_btree_index_range() {
        let index = BTreeIndex::new(1000);

        for i in 0..10 {
            let key = format!("key-{:02}", i).into_bytes();
            let entry = IndexEntry {
                key: key.clone(),
                value: format!("value-{}", i).into_bytes(),
                timestamp: 1000 + i as i64,
                series_id: i as u64,
            };
            index.insert(key, entry);
        }

        let results = index.range_scan(b"key-03", b"key-07");
        assert_eq!(results.len(), 5);
    }

    #[test]
    fn test_inverted_index() {
        let index = InvertedIndex::new();

        let labels = BTreeMap::from([
            ("host".to_string(), "server1".to_string()),
            ("region".to_string(), "us-east".to_string()),
        ]);
        index.add_series(1, labels);

        let result = index.lookup("host", "server1");
        assert_eq!(result, vec![1]);

        let result = index.lookup("region", "us-east");
        assert_eq!(result, vec![1]);
    }

    #[test]
    fn test_inverted_index_lookup_eq() {
        let index = InvertedIndex::new();

        index.add_series(
            1,
            BTreeMap::from([
                ("host".to_string(), "server1".to_string()),
                ("region".to_string(), "us-east".to_string()),
            ]),
        );
        index.add_series(
            2,
            BTreeMap::from([
                ("host".to_string(), "server2".to_string()),
                ("region".to_string(), "us-east".to_string()),
            ]),
        );

        let matchers = BTreeMap::from([("region".to_string(), "us-east".to_string())]);
        let result = index.lookup_eq(&matchers);
        assert_eq!(result.len(), 2);
    }

    #[test]
    fn test_inverted_index_label_values() {
        let index = InvertedIndex::new();

        index.add_series(
            1,
            BTreeMap::from([("host".to_string(), "server1".to_string())]),
        );
        index.add_series(
            2,
            BTreeMap::from([("host".to_string(), "server2".to_string())]),
        );

        let values = index.label_values("host");
        assert_eq!(values.len(), 2);
    }

    #[test]
    fn test_timeseries_index() {
        let index = TimeSeriesIndex::new();

        index.add(1, 1000);
        index.add(1, 1010);
        index.add(2, 1005);

        let result = index.time_range(1000, 1005);
        assert!(result.contains(&1));
        assert!(result.contains(&2));
    }

    #[test]
    fn test_timeseries_index_min_max() {
        let index = TimeSeriesIndex::new();

        index.add(1, 1000);
        index.add(2, 2000);

        assert_eq!(index.min_timestamp(), Some(1000));
        assert_eq!(index.max_timestamp(), Some(2000));
    }

    #[test]
    fn test_composite_index() {
        let index = CompositeIndex::new();

        index.add(
            1,
            b"key1".to_vec(),
            b"value1".to_vec(),
            BTreeMap::from([("host".to_string(), "server1".to_string())]),
            1000,
        );

        let matchers = BTreeMap::from([("host".to_string(), "server1".to_string())]);
        let result = index.search(&matchers, Some((900, 1100)));
        assert_eq!(result, vec![1]);
    }

    #[test]
    fn test_btree_index_clear() {
        let index = BTreeIndex::new(1000);
        index.insert(
            b"key1".to_vec(),
            IndexEntry {
                key: b"key1".to_vec(),
                value: b"value1".to_vec(),
                timestamp: 1000,
                series_id: 1,
            },
        );

        assert!(!index.is_empty());
        index.clear();
        assert!(index.is_empty());
    }

    #[test]
    fn test_inverted_index_remove() {
        let index = InvertedIndex::new();
        index.add_series(
            1,
            BTreeMap::from([("host".to_string(), "server1".to_string())]),
        );
        assert_eq!(index.series_count(), 1);

        index.remove_series(1);
        assert_eq!(index.series_count(), 0);
    }
}
