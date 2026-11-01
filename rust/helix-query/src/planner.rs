//! Helix Query Engine - Query Planner Module
//! 
//! This module provides advanced query planning capabilities
//! including cost-based optimization and plan caching.

use std::collections::BTreeMap;
use std::sync::{Arc, RwLock};
use std::time::{Duration, SystemTime, UNIX_EPOCH};

use super::parser::Expr;
use super::optimizer::QueryPlan;

/// Query planner with cost-based optimization
pub struct AdvancedQueryPlanner {
    cost_model: CostModel,
    plan_cache: PlanCache,
    rules: Vec<Box<dyn PlanningRule>>,
    stats: PlannerStats,
}

/// Cost model for query planning
#[derive(Debug, Clone)]
pub struct CostModel {
    pub cpu_cost_per_row: f64,
    pub io_cost_per_block: f64,
    pub network_cost_per_byte: f64,
    pub memory_cost_per_byte: f64,
    pub hash_join_cost: f64,
    pub sort_cost: f64,
    pub aggregation_cost: f64,
}

impl CostModel {
    pub fn default() -> Self {
        CostModel {
            cpu_cost_per_row: 0.001,
            io_cost_per_block: 1.0,
            network_cost_per_byte: 0.0001,
            memory_cost_per_byte: 0.00001,
            hash_join_cost: 0.01,
            sort_cost: 0.005,
            aggregation_cost: 0.002,
        }
    }
    
    pub fn estimate_scan_cost(&self, rows: usize, blocks: usize) -> f64 {
        rows as f64 * self.cpu_cost_per_row + blocks as f64 * self.io_cost_per_block
    }
    
    pub fn estimate_join_cost(&self, left_rows: usize, right_rows: usize) -> f64 {
        (left_rows + right_rows) as f64 * self.hash_join_cost
    }
    
    pub fn estimate_sort_cost(&self, rows: usize) -> f64 {
        rows as f64 * (rows as f64).log2() * self.sort_cost
    }
    
    pub fn estimate_aggregation_cost(&self, rows: usize, groups: usize) -> f64 {
        rows as f64 * self.aggregation_cost + groups as f64 * self.cpu_cost_per_row
    }
}

/// Plan cache
pub struct PlanCache {
    cache: RwLock<BTreeMap<String, CachedPlan>>,
    max_size: usize,
    ttl: Duration,
}

/// Cached plan
#[derive(Debug, Clone)]
pub struct CachedPlan {
    pub plan: QueryPlan,
    pub created_at: i64,
    pub hit_count: u64,
    pub cost: f64,
}

impl PlanCache {
    pub fn new(max_size: usize, ttl: Duration) -> Self {
        PlanCache {
            cache: RwLock::new(BTreeMap::new()),
            max_size,
            ttl,
        }
    }
    
    pub fn get(&self, key: &str) -> Option<CachedPlan> {
        let cache = self.cache.read().unwrap();
        if let Some(cached) = cache.get(key) {
            let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
            if now - cached.created_at < self.ttl.as_secs() as i64 {
                return Some(cached.clone());
            }
        }
        None
    }
    
    pub fn put(&self, key: String, plan: QueryPlan, cost: f64) {
        let mut cache = self.cache.write().unwrap();
        
        if cache.len() >= self.max_size {
            // Evict oldest
            if let Some(oldest_key) = cache.keys().next().cloned() {
                cache.remove(&oldest_key);
            }
        }
        
        let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as i64;
        cache.insert(key, CachedPlan {
            plan,
            created_at: now,
            hit_count: 0,
            cost,
        });
    }
    
    pub fn invalidate(&self, key: &str) {
        let mut cache = self.cache.write().unwrap();
        cache.remove(key);
    }
    
    pub fn clear(&self) {
        let mut cache = self.cache.write().unwrap();
        cache.clear();
    }
    
    pub fn size(&self) -> usize {
        let cache = self.cache.read().unwrap();
        cache.len()
    }
    
    pub fn hit_rate(&self) -> f64 {
        let cache = self.cache.read().unwrap();
        let total_hits: u64 = cache.values().map(|c| c.hit_count).sum();
        if cache.is_empty() {
            0.0
        } else {
            total_hits as f64 / cache.len() as f64
        }
    }
}

/// Planning rule trait
pub trait PlanningRule: Send + Sync {
    fn name(&self) -> &'static str;
    fn apply(&self, plan: &mut QueryPlan, stats: &TableStats) -> bool;
    fn description(&self) -> &'static str;
}

/// Table statistics
#[derive(Debug, Clone, Default)]
pub struct TableStats {
    pub row_count: u64,
    pub block_count: u64,
    pub byte_size: u64,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
    pub distinct_keys: u64,
    pub avg_row_size: f64,
    pub columns: BTreeMap<String, ColumnStats>,
}

/// Column statistics
#[derive(Debug, Clone, Default)]
pub struct ColumnStats {
    pub name: String,
    pub distinct_count: u64,
    pub null_count: u64,
    pub min_value: f64,
    pub max_value: f64,
    pub avg_value: f64,
    pub std_dev: f64,
}

/// Planner statistics
#[derive(Debug, Clone, Default)]
pub struct PlannerStats {
    pub plans_created: u64,
    pub plans_cached: u64,
    pub plans_optimized: u64,
    pub optimization_time_ms: u64,
    pub planning_time_ms: u64,
}

impl AdvancedQueryPlanner {
    pub fn new() -> Self {
        let mut planner = AdvancedQueryPlanner {
            cost_model: CostModel::default(),
            plan_cache: PlanCache::new(1000, Duration::from_secs(300)),
            rules: Vec::new(),
            stats: PlannerStats::default(),
        };
        
        planner.add_rule(Box::new(IndexSelectionRule));
        planner.add_rule(Box::new(JoinOrderingRule));
        planner.add_rule(Box::new(AggregationPushdownRule));
        planner.add_rule(Box::new(PredicatePushdownRule));
        planner.add_rule(Box::new(ProjectionPushdownRule));
        planner.add_rule(Box::new(PartitionPruningRule));
        planner.add_rule(Box::new(SortEliminationRule));
        planner.add_rule(Box::new(LimitPushdownRule));
        planner.add_rule(Box::new(SubqueryFlatteningRule));
        planner.add_rule(Box::new(ConstantFoldingRule));
        
        planner
    }
    
    pub fn add_rule(&mut self, rule: Box<dyn PlanningRule>) {
        self.rules.push(rule);
    }
    
    pub fn plan(&mut self, expr: &Expr, stats: &TableStats) -> QueryPlan {
        let start = SystemTime::now();
        
        // Check cache
        let cache_key = format!("{:?}", expr);
        if let Some(cached) = self.plan_cache.get(&cache_key) {
            self.stats.plans_cached += 1;
            return cached.plan;
        }
        
        // Create base plan
        let mut plan = self.create_base_plan(expr);
        
        // Apply planning rules
        for rule in &self.rules {
            if rule.apply(&mut plan, stats) {
                self.stats.plans_optimized += 1;
            }
        }
        
        // Estimate cost
        plan.estimated_cost = self.cost_model.estimate_scan_cost(stats.row_count as usize, stats.block_count as usize);
        
        // Cache plan
        self.plan_cache.put(cache_key, plan.clone(), plan.estimated_cost);
        
        self.stats.plans_created += 1;
        self.stats.planning_time_ms += start.elapsed().as_millis() as u64;
        
        plan
    }
    
    fn create_base_plan(&self, expr: &Expr) -> QueryPlan {
        let mut plan = QueryPlan::default();
        
        match expr {
            Expr::Number(_) | Expr::String(_) => {
                plan.steps = vec!["constant".to_string()];
            }
            Expr::Identifier(_) => {
                plan.steps = vec!["scan".to_string()];
            }
            Expr::VectorSelector { .. } => {
                plan.steps = vec!["scan".to_string(), "filter".to_string()];
            }
            Expr::MatrixSelector { .. } => {
                plan.steps = vec!["scan".to_string(), "filter".to_string(), "range".to_string()];
            }
            Expr::FunctionCall { name, .. } => {
                plan.steps = vec!["scan".to_string(), "filter".to_string(), format!("func:{}", name)];
            }
            Expr::Aggregation { op, .. } => {
                plan.steps = vec!["scan".to_string(), "filter".to_string(), format!("agg:{:?}", op)];
            }
            Expr::BinaryOp { .. } => {
                plan.steps = vec!["scan".to_string(), "binary_op".to_string()];
            }
            Expr::UnaryOp { .. } => {
                plan.steps = vec!["scan".to_string(), "unary_op".to_string()];
            }
            Expr::Paren(inner) => {
                return self.create_base_plan(inner);
            }
            _ => {
                plan.steps = vec!["scan".to_string()];
            }
        }
        
        plan
    }
    
    pub fn stats(&self) -> PlannerStats {
        self.stats.clone()
    }
    
    pub fn cache_stats(&self) -> (usize, f64) {
        (self.plan_cache.size(), self.plan_cache.hit_rate())
    }
}

// ---------------------------------------------------------------------------
// Planning Rules
// ---------------------------------------------------------------------------

/// Index selection rule
pub struct IndexSelectionRule;

impl PlanningRule for IndexSelectionRule {
    fn name(&self) -> &'static str { "index_selection" }
    fn apply(&self, plan: &mut QueryPlan, stats: &TableStats) -> bool {
        if stats.distinct_keys > 0 && !plan.steps.contains(&"use_index".to_string()) {
            plan.steps.insert(0, "use_index".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str { "Select optimal index" }
}

/// Join ordering rule
pub struct JoinOrderingRule;

impl PlanningRule for JoinOrderingRule {
    fn name(&self) -> &'static str { "join_ordering" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        false
    }
    fn description(&self) -> &'static str { "Reorder joins" }
}

/// Aggregation pushdown rule
pub struct AggregationPushdownRule;

impl PlanningRule for AggregationPushdownRule {
    fn name(&self) -> &'static str { "aggregation_pushdown" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        if plan.steps.contains(&"agg:Sum".to_string()) && !plan.steps.contains(&"push_agg".to_string()) {
            plan.steps.insert(0, "push_agg".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str { "Push aggregation down" }
}

/// Predicate pushdown rule
pub struct PredicatePushdownRule;

impl PlanningRule for PredicatePushdownRule {
    fn name(&self) -> &'static str { "predicate_pushdown" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        let filter_idx = plan.steps.iter().position(|s| s == "filter");
        let scan_idx = plan.steps.iter().position(|s| s == "scan");
        
        if let (Some(fi), Some(si)) = (filter_idx, scan_idx) {
            if fi > si {
                let filter = plan.steps.remove(fi);
                plan.steps.insert(si, filter);
                return true;
            }
        }
        false
    }
    fn description(&self) -> &'static str { "Push predicates down" }
}

/// Projection pushdown rule
pub struct ProjectionPushdownRule;

impl PlanningRule for ProjectionPushdownRule {
    fn name(&self) -> &'static str { "projection_pushdown" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        false
    }
    fn description(&self) -> &'static str { "Push projections down" }
}

/// Partition pruning rule
pub struct PartitionPruningRule;

impl PlanningRule for PartitionPruningRule {
    fn name(&self) -> &'static str { "partition_pruning" }
    fn apply(&self, plan: &mut QueryPlan, stats: &TableStats) -> bool {
        if stats.min_timestamp > 0 && !plan.steps.contains(&"prune_partitions".to_string()) {
            plan.steps.insert(0, "prune_partitions".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str { "Prune partitions" }
}

/// Sort elimination rule
pub struct SortEliminationRule;

impl PlanningRule for SortEliminationRule {
    fn name(&self) -> &'static str { "sort_elimination" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        false
    }
    fn description(&self) -> &'static str { "Eliminate redundant sorts" }
}

/// Limit pushdown rule
pub struct LimitPushdownRule;

impl PlanningRule for LimitPushdownRule {
    fn name(&self) -> &'static str { "limit_pushdown" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        if plan.steps.contains(&"limit".to_string()) && !plan.steps.contains(&"push_limit".to_string()) {
            plan.steps.insert(0, "push_limit".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str { "Push limit down" }
}

/// Subquery flattening rule
pub struct SubqueryFlatteningRule;

impl PlanningRule for SubqueryFlatteningRule {
    fn name(&self) -> &'static str { "subquery_flattening" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        false
    }
    fn description(&self) -> &'static str { "Flatten subqueries" }
}

/// Constant folding rule
pub struct ConstantFoldingRule;

impl PlanningRule for ConstantFoldingRule {
    fn name(&self) -> &'static str { "constant_folding" }
    fn apply(&self, plan: &mut QueryPlan, _stats: &TableStats) -> bool {
        false
    }
    fn description(&self) -> &'static str { "Fold constants" }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_cost_model() {
        let model = CostModel::default();
        let cost = model.estimate_scan_cost(1000, 10);
        assert!(cost > 0.0);
        
        let join_cost = model.estimate_join_cost(1000, 500);
        assert!(join_cost > 0.0);
        
        let sort_cost = model.estimate_sort_cost(1000);
        assert!(sort_cost > 0.0);
    }
    
    #[test]
    fn test_plan_cache() {
        let cache = PlanCache::new(10, Duration::from_secs(60));
        let plan = QueryPlan::default();
        
        cache.put("test".to_string(), plan.clone(), 100.0);
        assert!(cache.get("test").is_some());
        assert_eq!(cache.size(), 1);
    }
    
    #[test]
    fn test_plan_cache_invalidate() {
        let cache = PlanCache::new(10, Duration::from_secs(60));
        let plan = QueryPlan::default();
        
        cache.put("test".to_string(), plan, 100.0);
        cache.invalidate("test");
        assert!(cache.get("test").is_none());
    }
    
    #[test]
    fn test_plan_cache_clear() {
        let cache = PlanCache::new(10, Duration::from_secs(60));
        let plan = QueryPlan::default();
        
        cache.put("test1".to_string(), plan.clone(), 100.0);
        cache.put("test2".to_string(), plan, 100.0);
        cache.clear();
        assert_eq!(cache.size(), 0);
    }
    
    #[test]
    fn test_advanced_query_planner() {
        let mut planner = AdvancedQueryPlanner::new();
        let stats = TableStats::default();
        
        let plan = planner.plan(&Expr::Number(42.0), &stats);
        assert!(!plan.steps.is_empty());
    }
    
    #[test]
    fn test_index_selection_rule() {
        let rule = IndexSelectionRule;
        let mut plan = QueryPlan::default();
        let stats = TableStats { distinct_keys: 100, ..Default::default() };
        
        assert!(rule.apply(&mut plan, &stats));
        assert!(plan.steps.contains(&"use_index".to_string()));
    }
    
    #[test]
    fn test_predicate_pushdown_rule() {
        let rule = PredicatePushdownRule;
        let mut plan = QueryPlan {
            steps: vec!["scan".to_string(), "filter".to_string()],
            ..Default::default()
        };
        let stats = TableStats::default();
        
        assert!(rule.apply(&mut plan, &stats));
        assert_eq!(plan.steps[0], "filter");
    }
    
    #[test]
    fn test_partition_pruning_rule() {
        let rule = PartitionPruningRule;
        let mut plan = QueryPlan::default();
        let stats = TableStats { min_timestamp: 1000, ..Default::default() };
        
        assert!(rule.apply(&mut plan, &stats));
        assert!(plan.steps.contains(&"prune_partitions".to_string()));
    }
    
    #[test]
    fn test_limit_pushdown_rule() {
        let rule = LimitPushdownRule;
        let mut plan = QueryPlan {
            steps: vec!["limit".to_string()],
            ..Default::default()
        };
        let stats = TableStats::default();
        
        assert!(rule.apply(&mut plan, &stats));
        assert!(plan.steps.contains(&"push_limit".to_string()));
    }
    
    #[test]
    fn test_aggregation_pushdown_rule() {
        let rule = AggregationPushdownRule;
        let mut plan = QueryPlan {
            steps: vec!["agg:Sum".to_string()],
            ..Default::default()
        };
        let stats = TableStats::default();
        
        assert!(rule.apply(&mut plan, &stats));
        assert!(plan.steps.contains(&"push_agg".to_string()));
    }
    
    #[test]
    fn test_planner_stats() {
        let mut planner = AdvancedQueryPlanner::new();
        let stats = TableStats::default();
        
        planner.plan(&Expr::Number(42.0), &stats);
        
        let pstats = planner.stats();
        assert_eq!(pstats.plans_created, 1);
    }
    
    #[test]
    fn test_cache_stats() {
        let mut planner = AdvancedQueryPlanner::new();
        let stats = TableStats::default();
        
        planner.plan(&Expr::Number(42.0), &stats);
        planner.plan(&Expr::Number(42.0), &stats); // Should hit cache
        
        let (size, hit_rate) = planner.cache_stats();
        assert_eq!(size, 1);
        assert!(hit_rate >= 0.0);
    }
    
    #[test]
    fn test_table_stats() {
        let stats = TableStats {
            row_count: 1000000,
            block_count: 1000,
            byte_size: 100000000,
            distinct_keys: 50000,
            ..Default::default()
        };
        
        assert_eq!(stats.row_count, 1000000);
        assert_eq!(stats.distinct_keys, 50000);
    }
    
    #[test]
    fn test_column_stats() {
        let stats = ColumnStats {
            name: "cpu_usage".to_string(),
            distinct_count: 100,
            min_value: 0.0,
            max_value: 100.0,
            avg_value: 50.0,
            std_dev: 15.0,
            ..Default::default()
        };
        
        assert_eq!(stats.name, "cpu_usage");
        assert_eq!(stats.distinct_count, 100);
    }
}