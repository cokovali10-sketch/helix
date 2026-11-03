//! Helix Query Engine - Optimizer Module
//!
//! This module provides query optimization for the Helix query language.
//! Includes rule-based and cost-based optimization.

use super::{Query, QueryPlan};
use std::collections::BTreeMap;

/// Query optimizer
pub struct QueryOptimizer {
    rules: Vec<Box<dyn OptimizationRule>>,
    cost_model: CostModel,
}

/// Optimization rule trait
pub trait OptimizationRule: Send + Sync {
    fn name(&self) -> &'static str;
    fn apply(&self, plan: &mut QueryPlan) -> bool;
    fn description(&self) -> &'static str;
}

/// Cost model for query optimization
pub struct CostModel {
    pub cpu_cost_per_row: f64,
    pub io_cost_per_block: f64,
    pub network_cost_per_byte: f64,
    pub memory_cost_per_byte: f64,
}

impl CostModel {
    #[allow(clippy::should_implement_trait)]
    pub fn default() -> Self {
        CostModel {
            cpu_cost_per_row: 0.001,
            io_cost_per_block: 1.0,
            network_cost_per_byte: 0.0001,
            memory_cost_per_byte: 0.00001,
        }
    }

    pub fn estimate_cost(&self, rows: usize, blocks: usize, bytes: usize) -> f64 {
        rows as f64 * self.cpu_cost_per_row
            + blocks as f64 * self.io_cost_per_block
            + bytes as f64 * self.network_cost_per_byte
    }
}

impl QueryOptimizer {
    pub fn new() -> Self {
        let mut optimizer = QueryOptimizer {
            rules: Vec::new(),
            cost_model: CostModel::default(),
        };

        // Register default optimization rules
        optimizer.add_rule(Box::new(PredicatePushdownRule));
        optimizer.add_rule(Box::new(ProjectionPushdownRule));
        optimizer.add_rule(Box::new(ConstantFoldingRule));
        optimizer.add_rule(Box::new(DeadCodeEliminationRule));
        optimizer.add_rule(Box::new(JoinReorderingRule));
        optimizer.add_rule(Box::new(AggregationPushdownRule));
        optimizer.add_rule(Box::new(PartitionPruningRule));
        optimizer.add_rule(Box::new(IndexSelectionRule));
        optimizer.add_rule(Box::new(CacheAwarenessRule));
        optimizer.add_rule(Box::new(VectorizationRule));

        optimizer
    }

    pub fn add_rule(&mut self, rule: Box<dyn OptimizationRule>) {
        self.rules.push(rule);
    }

    pub fn optimize(&self, plan: &mut QueryPlan) -> OptimizationResult {
        let mut applied_rules = Vec::new();
        let mut iterations = 0;
        let max_iterations = 10;

        let original_cost = plan.estimated_cost;

        while iterations < max_iterations {
            let mut changed = false;

            for rule in &self.rules {
                if rule.apply(plan) {
                    applied_rules.push(rule.name().to_string());
                    changed = true;
                }
            }

            if !changed {
                break;
            }
            iterations += 1;
        }

        OptimizationResult {
            original_cost,
            optimized_cost: plan.estimated_cost,
            applied_rules,
            iterations,
            improvement_percent: if original_cost > 0.0 {
                ((original_cost - plan.estimated_cost) / original_cost) * 100.0
            } else {
                0.0
            },
        }
    }
}

impl Default for QueryOptimizer {
    fn default() -> Self {
        Self::new()
    }
}

/// Optimization result
#[derive(Debug, Clone)]
pub struct OptimizationResult {
    pub original_cost: f64,
    pub optimized_cost: f64,
    pub applied_rules: Vec<String>,
    pub iterations: usize,
    pub improvement_percent: f64,
}

// ---------------------------------------------------------------------------
// Optimization Rules
// ---------------------------------------------------------------------------

/// Predicate pushdown rule
pub struct PredicatePushdownRule;

impl OptimizationRule for PredicatePushdownRule {
    fn name(&self) -> &'static str {
        "predicate_pushdown"
    }
    fn apply(&self, plan: &mut QueryPlan) -> bool {
        // Move filter operations earlier in the plan
        let filter_idx = plan.steps.iter().position(|s| s == "filter_labels");
        let fetch_idx = plan.steps.iter().position(|s| s == "fetch_from_storage");

        if let (Some(fi), Some(fi2)) = (filter_idx, fetch_idx) {
            if fi > fi2 {
                // Move filter before fetch
                let filter = plan.steps.remove(fi);
                plan.steps.insert(fi2, filter);
                return true;
            }
        }
        false
    }
    fn description(&self) -> &'static str {
        "Push filter predicates down to storage layer"
    }
}

/// Projection pushdown rule
pub struct ProjectionPushdownRule;

impl OptimizationRule for ProjectionPushdownRule {
    fn name(&self) -> &'static str {
        "projection_pushdown"
    }
    fn apply(&self, _plan: &mut QueryPlan) -> bool {
        // Push projection (column selection) down
        // Simplified: no-op for now
        false
    }
    fn description(&self) -> &'static str {
        "Push column projections down to storage layer"
    }
}

/// Constant folding rule
pub struct ConstantFoldingRule;

impl OptimizationRule for ConstantFoldingRule {
    fn name(&self) -> &'static str {
        "constant_folding"
    }
    fn apply(&self, _plan: &mut QueryPlan) -> bool {
        // Fold constant expressions at compile time
        false
    }
    fn description(&self) -> &'static str {
        "Evaluate constant expressions at compile time"
    }
}

/// Dead code elimination rule
pub struct DeadCodeEliminationRule;

impl OptimizationRule for DeadCodeEliminationRule {
    fn name(&self) -> &'static str {
        "dead_code_elimination"
    }
    fn apply(&self, plan: &mut QueryPlan) -> bool {
        // Remove unused computation steps
        let original_len = plan.steps.len();
        plan.steps.retain(|s| s != "unused");
        plan.steps.len() < original_len
    }
    fn description(&self) -> &'static str {
        "Remove unused computation steps"
    }
}

/// Join reordering rule
pub struct JoinReorderingRule;

impl OptimizationRule for JoinReorderingRule {
    fn name(&self) -> &'static str {
        "join_reordering"
    }
    fn apply(&self, _plan: &mut QueryPlan) -> bool {
        // Reorder joins for optimal performance
        false
    }
    fn description(&self) -> &'static str {
        "Reorder joins based on cardinality estimates"
    }
}

/// Aggregation pushdown rule
pub struct AggregationPushdownRule;

impl OptimizationRule for AggregationPushdownRule {
    fn name(&self) -> &'static str {
        "aggregation_pushdown"
    }
    fn apply(&self, _plan: &mut QueryPlan) -> bool {
        // Push aggregation down to storage
        false
    }
    fn description(&self) -> &'static str {
        "Push aggregation operations down to storage layer"
    }
}

/// Partition pruning rule
pub struct PartitionPruningRule;

impl OptimizationRule for PartitionPruningRule {
    fn name(&self) -> &'static str {
        "partition_pruning"
    }
    fn apply(&self, plan: &mut QueryPlan) -> bool {
        // Prune partitions based on time range
        let has_filter = plan.steps.iter().any(|s| s == "filter_time_range");
        if has_filter {
            // Add partition pruning step
            if !plan.steps.contains(&"prune_partitions".to_string()) {
                plan.steps.insert(0, "prune_partitions".to_string());
                return true;
            }
        }
        false
    }
    fn description(&self) -> &'static str {
        "Prune partitions based on query time range"
    }
}

/// Index selection rule
pub struct IndexSelectionRule;

impl OptimizationRule for IndexSelectionRule {
    fn name(&self) -> &'static str {
        "index_selection"
    }
    fn apply(&self, _plan: &mut QueryPlan) -> bool {
        // Select optimal index for query
        false
    }
    fn description(&self) -> &'static str {
        "Select optimal index based on query predicates"
    }
}

/// Cache awareness rule
pub struct CacheAwarenessRule;

impl OptimizationRule for CacheAwarenessRule {
    fn name(&self) -> &'static str {
        "cache_awareness"
    }
    fn apply(&self, plan: &mut QueryPlan) -> bool {
        // Add cache lookup step
        if !plan.steps.contains(&"check_cache".to_string()) {
            plan.steps.insert(0, "check_cache".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str {
        "Add cache lookup for frequently accessed data"
    }
}

/// Vectorization rule
pub struct VectorizationRule;

impl OptimizationRule for VectorizationRule {
    fn name(&self) -> &'static str {
        "vectorization"
    }
    fn apply(&self, plan: &mut QueryPlan) -> bool {
        // Enable vectorized execution
        if !plan.steps.contains(&"vectorize".to_string()) {
            plan.steps.push("vectorize".to_string());
            return true;
        }
        false
    }
    fn description(&self) -> &'static str {
        "Enable vectorized execution for SIMD optimization"
    }
}

// ---------------------------------------------------------------------------
// Query Planner
// ---------------------------------------------------------------------------

/// Query planner
pub struct QueryPlanner {
    optimizer: QueryOptimizer,
}

impl QueryPlanner {
    pub fn new() -> Self {
        QueryPlanner {
            optimizer: QueryOptimizer::new(),
        }
    }

    pub fn plan(&self, query: &Query) -> QueryPlan {
        let mut plan = QueryPlan::default();

        // Basic plan from query type
        match query.query_type {
            crate::QueryType::InstantVector { .. } => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "compute_instant".to_string(),
                    "post_process".to_string(),
                ];
            }
            crate::QueryType::RangeVector { .. } => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "aggregate_range".to_string(),
                    "compute_range".to_string(),
                    "post_process".to_string(),
                ];
            }
            crate::QueryType::VectorSelector => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "select_vector".to_string(),
                    "post_process".to_string(),
                ];
            }
            crate::QueryType::Scalar | crate::QueryType::String => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "aggregate".to_string(),
                    "compute_scalar".to_string(),
                ];
            }
        }

        // Add filters
        if !query.labels.is_empty() {
            plan.steps.insert(0, "filter_labels".to_string());
        }
        plan.steps.insert(1, "filter_time_range".to_string());

        // Optimize
        let _ = self.optimizer.optimize(&mut plan);

        plan
    }
}

impl Default for QueryPlanner {
    fn default() -> Self {
        Self::new()
    }
}

// ---------------------------------------------------------------------------
// Statistics for Cost Estimation
// ---------------------------------------------------------------------------

/// Table statistics for cost-based optimization
#[derive(Debug, Clone, Default)]
pub struct TableStats {
    pub row_count: u64,
    pub block_count: u64,
    pub byte_size: u64,
    pub min_timestamp: i64,
    pub max_timestamp: i64,
    pub distinct_keys: u64,
    pub avg_row_size: f64,
}

impl TableStats {
    pub fn new() -> Self {
        TableStats::default()
    }

    pub fn density(&self) -> f64 {
        if self.byte_size > 0 {
            self.row_count as f64 / self.byte_size as f64
        } else {
            0.0
        }
    }

    pub fn selectivity(&self, _predicate: &str) -> f64 {
        // Simplified selectivity estimation
        0.1
    }
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

impl ColumnStats {
    pub fn new(name: String) -> Self {
        ColumnStats {
            name,
            ..Default::default()
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_query_optimizer_creation() {
        let optimizer = QueryOptimizer::new();
        assert_eq!(optimizer.rules.len(), 10);
    }

    #[test]
    fn test_predicate_pushdown() {
        let rule = PredicatePushdownRule;
        let mut plan = QueryPlan {
            steps: vec![
                "fetch_from_storage".to_string(),
                "filter_labels".to_string(),
            ],
            ..Default::default()
        };
        assert!(rule.apply(&mut plan));
        assert_eq!(plan.steps[0], "filter_labels");
    }

    #[test]
    fn test_partition_pruning() {
        let rule = PartitionPruningRule;
        let mut plan = QueryPlan {
            steps: vec![
                "filter_time_range".to_string(),
                "fetch_from_storage".to_string(),
            ],
            ..Default::default()
        };
        assert!(rule.apply(&mut plan));
        assert!(plan.steps.contains(&"prune_partitions".to_string()));
    }

    #[test]
    fn test_cache_awareness() {
        let rule = CacheAwarenessRule;
        let mut plan = QueryPlan {
            steps: vec!["fetch_from_storage".to_string()],
            ..Default::default()
        };
        assert!(rule.apply(&mut plan));
        assert!(plan.steps.contains(&"check_cache".to_string()));
    }

    #[test]
    fn test_vectorization() {
        let rule = VectorizationRule;
        let mut plan = QueryPlan {
            steps: vec!["fetch_from_storage".to_string()],
            ..Default::default()
        };
        assert!(rule.apply(&mut plan));
        assert!(plan.steps.contains(&"vectorize".to_string()));
    }

    #[test]
    fn test_cost_model() {
        let model = CostModel::default();
        let cost = model.estimate_cost(1000, 10, 10000);
        assert!(cost > 0.0);
    }

    #[test]
    fn test_table_stats() {
        let stats = TableStats {
            row_count: 1000000,
            byte_size: 100000000,
            ..Default::default()
        };
        assert!(stats.density() > 0.0);
    }

    #[test]
    fn test_optimization_result() {
        let mut plan = QueryPlan {
            steps: vec!["fetch_from_storage".to_string()],
            estimated_cost: 100.0,
            ..Default::default()
        };

        let optimizer = QueryOptimizer::new();
        let result = optimizer.optimize(&mut plan);

        assert!(result.improvement_percent >= 0.0);
    }

    #[test]
    fn test_query_planner() {
        let planner = QueryPlanner::new();
        let query = Query {
            query_type: crate::QueryType::InstantVector {
                labels: BTreeMap::new(),
            },
            expression: "cpu_usage".to_string(),
            labels: BTreeMap::new(),
            time_range: crate::TimeRange {
                start: 0,
                end: 100,
                offset: 0,
            },
        };

        let plan = planner.plan(&query);
        assert!(!plan.steps.is_empty());
    }
}
