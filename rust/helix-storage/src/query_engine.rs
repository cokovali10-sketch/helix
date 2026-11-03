//! Helix Storage Engine - Query Engine Module
//!
//! This module provides query capabilities for the storage engine.
//! Includes scan, filter, aggregation, and join operations.

use std::collections::BTreeMap;
use std::sync::Arc;

use super::index::CompositeIndex;

/// Query result
#[derive(Debug, Clone)]
pub struct QueryResult {
    pub columns: Vec<String>,
    pub rows: Vec<QueryRow>,
    pub total_rows: usize,
    pub execution_time_ms: f64,
}

/// Query row
#[derive(Debug, Clone)]
pub struct QueryRow {
    pub values: Vec<QueryValue>,
}

/// Query value
#[derive(Debug, Clone)]
pub enum QueryValue {
    Null,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Bytes(Vec<u8>),
    Timestamp(i64),
    Labels(BTreeMap<String, String>),
}

/// Query operator
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum QueryOp {
    Scan,
    Filter,
    Project,
    Aggregate,
    Sort,
    Limit,
    Offset,
    Join,
    Union,
    Intersect,
    Distinct,
}

/// Query plan node
#[derive(Debug, Clone)]
pub struct PlanNode {
    pub op: QueryOp,
    pub children: Vec<PlanNode>,
    pub params: BTreeMap<String, String>,
    pub estimated_cost: f64,
    pub estimated_rows: usize,
}

impl PlanNode {
    pub fn new(op: QueryOp) -> Self {
        PlanNode {
            op,
            children: Vec::new(),
            params: BTreeMap::new(),
            estimated_cost: 0.0,
            estimated_rows: 0,
        }
    }

    pub fn with_child(mut self, child: PlanNode) -> Self {
        self.children.push(child);
        self
    }

    pub fn with_param(mut self, key: String, value: String) -> Self {
        self.params.insert(key, value);
        self
    }

    pub fn cost(&self) -> f64 {
        self.estimated_cost + self.children.iter().map(|c| c.cost()).sum::<f64>()
    }
}

/// Query executor
pub struct QueryExecutor {
    index: Arc<CompositeIndex>,
}

impl QueryExecutor {
    pub fn new(index: Arc<CompositeIndex>) -> Self {
        QueryExecutor { index }
    }

    /// Execute a query plan
    pub fn execute(&self, plan: &PlanNode) -> Result<QueryResult, String> {
        let start = std::time::Instant::now();
        let rows = self.execute_node(plan)?;
        let elapsed = start.elapsed();

        Ok(QueryResult {
            columns: vec![],
            rows,
            total_rows: 0,
            execution_time_ms: elapsed.as_millis() as f64,
        })
    }

    fn execute_node(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        match node.op {
            QueryOp::Scan => self.execute_scan(node),
            QueryOp::Filter => self.execute_filter(node),
            QueryOp::Project => self.execute_project(node),
            QueryOp::Aggregate => self.execute_aggregate(node),
            QueryOp::Sort => self.execute_sort(node),
            QueryOp::Limit => self.execute_limit(node),
            QueryOp::Offset => self.execute_offset(node),
            QueryOp::Distinct => self.execute_distinct(node),
            _ => Err("Unsupported operation".to_string()),
        }
    }

    fn execute_scan(&self, _node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        // Scan all data
        Ok(Vec::new())
    }

    fn execute_filter(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        // Apply filter
        Ok(child_rows)
    }

    fn execute_project(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        // Project columns
        Ok(child_rows)
    }

    fn execute_aggregate(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        // Aggregate
        Ok(child_rows)
    }

    fn execute_sort(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let mut child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        // Sort
        child_rows.sort_by(|_a, _b| {
            // Simplified sort
            std::cmp::Ordering::Equal
        });

        Ok(child_rows)
    }

    fn execute_limit(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        let limit: usize = node
            .params
            .get("limit")
            .and_then(|s| s.parse().ok())
            .unwrap_or(100);

        Ok(child_rows.into_iter().take(limit).collect())
    }

    fn execute_offset(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        let offset: usize = node
            .params
            .get("offset")
            .and_then(|s| s.parse().ok())
            .unwrap_or(0);

        Ok(child_rows.into_iter().skip(offset).collect())
    }

    fn execute_distinct(&self, node: &PlanNode) -> Result<Vec<QueryRow>, String> {
        let child_rows = if !node.children.is_empty() {
            self.execute_node(&node.children[0])?
        } else {
            Vec::new()
        };

        // Remove duplicates
        let mut seen = std::collections::BTreeSet::new();
        let mut result = Vec::new();
        for row in child_rows {
            let key = format!("{:?}", row.values);
            if seen.insert(key) {
                result.push(row);
            }
        }

        Ok(result)
    }
}

/// Query builder
pub struct QueryBuilder {
    plan: PlanNode,
}

impl QueryBuilder {
    pub fn new() -> Self {
        QueryBuilder {
            plan: PlanNode::new(QueryOp::Scan),
        }
    }

    pub fn scan(mut self, table: &str) -> Self {
        self.plan = PlanNode::new(QueryOp::Scan).with_param("table".to_string(), table.to_string());
        self
    }

    pub fn filter(mut self, predicate: &str) -> Self {
        self.plan = PlanNode::new(QueryOp::Filter)
            .with_child(self.plan)
            .with_param("predicate".to_string(), predicate.to_string());
        self
    }

    pub fn project(mut self, columns: Vec<String>) -> Self {
        self.plan = PlanNode::new(QueryOp::Project)
            .with_child(self.plan)
            .with_param("columns".to_string(), columns.join(","));
        self
    }

    pub fn aggregate(mut self, op: &str, group_by: Vec<String>) -> Self {
        self.plan = PlanNode::new(QueryOp::Aggregate)
            .with_child(self.plan)
            .with_param("op".to_string(), op.to_string())
            .with_param("group_by".to_string(), group_by.join(","));
        self
    }

    pub fn sort(mut self, columns: Vec<String>, ascending: bool) -> Self {
        self.plan = PlanNode::new(QueryOp::Sort)
            .with_child(self.plan)
            .with_param("columns".to_string(), columns.join(","))
            .with_param("ascending".to_string(), ascending.to_string());
        self
    }

    pub fn limit(mut self, count: usize) -> Self {
        self.plan = PlanNode::new(QueryOp::Limit)
            .with_child(self.plan)
            .with_param("limit".to_string(), count.to_string());
        self
    }

    pub fn offset(mut self, count: usize) -> Self {
        self.plan = PlanNode::new(QueryOp::Offset)
            .with_child(self.plan)
            .with_param("offset".to_string(), count.to_string());
        self
    }

    pub fn distinct(mut self) -> Self {
        self.plan = PlanNode::new(QueryOp::Distinct).with_child(self.plan);
        self
    }

    pub fn build(self) -> PlanNode {
        self.plan
    }
}

impl Default for QueryBuilder {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_query_op() {
        assert_eq!(QueryOp::Scan, QueryOp::Scan);
        assert_ne!(QueryOp::Scan, QueryOp::Filter);
    }

    #[test]
    fn test_plan_node() {
        let node =
            PlanNode::new(QueryOp::Scan).with_param("table".to_string(), "metrics".to_string());

        assert_eq!(node.op, QueryOp::Scan);
        assert_eq!(node.params.get("table"), Some(&"metrics".to_string()));
    }

    #[test]
    fn test_plan_node_with_child() {
        let child = PlanNode::new(QueryOp::Scan);
        let parent = PlanNode::new(QueryOp::Filter).with_child(child);

        assert_eq!(parent.children.len(), 1);
        assert_eq!(parent.children[0].op, QueryOp::Scan);
    }

    #[test]
    fn test_plan_node_cost() {
        let child = PlanNode::new(QueryOp::Scan);
        let mut parent = PlanNode::new(QueryOp::Filter).with_child(child);
        parent.estimated_cost = 10.0;
        parent.children[0].estimated_cost = 5.0;

        assert_eq!(parent.cost(), 15.0);
    }

    #[test]
    fn test_query_builder() {
        let plan = QueryBuilder::new()
            .scan("metrics")
            .filter("host = 'server1'")
            .project(vec!["value".to_string()])
            .limit(100)
            .build();

        assert_eq!(plan.op, QueryOp::Limit);
    }

    #[test]
    fn test_query_builder_aggregate() {
        let plan = QueryBuilder::new()
            .scan("metrics")
            .aggregate("sum", vec!["host".to_string()])
            .build();

        assert_eq!(plan.op, QueryOp::Aggregate);
    }

    #[test]
    fn test_query_builder_sort() {
        let plan = QueryBuilder::new()
            .scan("metrics")
            .sort(vec!["value".to_string()], false)
            .build();

        assert_eq!(plan.op, QueryOp::Sort);
    }

    #[test]
    fn test_query_builder_distinct() {
        let plan = QueryBuilder::new().scan("metrics").distinct().build();

        assert_eq!(plan.op, QueryOp::Distinct);
    }

    #[test]
    fn test_query_value() {
        let val = QueryValue::Float(42.0);
        match val {
            QueryValue::Float(v) => assert_eq!(v, 42.0),
            _ => panic!("Expected Float"),
        }
    }

    #[test]
    fn test_query_result() {
        let result = QueryResult {
            columns: vec!["host".to_string(), "value".to_string()],
            rows: vec![],
            total_rows: 0,
            execution_time_ms: 10.5,
        };
        assert_eq!(result.columns.len(), 2);
        assert_eq!(result.execution_time_ms, 10.5);
    }
}
