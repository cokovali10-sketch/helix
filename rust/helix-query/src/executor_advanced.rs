//! Helix Query Engine - Query Executor Module
//!
//! This module provides the query execution engine with
//! vectorized operations and optimization.

use std::collections::BTreeMap;

use super::functions::FunctionRegistry;
use crate::parser::{Duration, Expr};

/// Execution context
#[derive(Debug, Clone)]
pub struct ExecutionContext {
    pub timestamp: i64,
    pub step: i64,
    pub range: i64,
    pub timeout: i64,
    pub max_rows: usize,
}

impl ExecutionContext {
    pub fn new(timestamp: i64) -> Self {
        ExecutionContext {
            timestamp,
            step: 1,
            range: 0,
            timeout: 30000,
            max_rows: 10000,
        }
    }

    pub fn with_step(mut self, step: i64) -> Self {
        self.step = step;
        self
    }

    pub fn with_range(mut self, range: i64) -> Self {
        self.range = range;
        self
    }

    pub fn with_timeout(mut self, timeout: i64) -> Self {
        self.timeout = timeout;
        self
    }

    pub fn with_max_rows(mut self, max_rows: usize) -> Self {
        self.max_rows = max_rows;
        self
    }
}

/// Execution result
#[derive(Debug, Clone)]
pub struct ExecutionResult {
    pub columns: Vec<String>,
    pub rows: Vec<ExecutionRow>,
    pub total_rows: usize,
    pub execution_time_ms: f64,
    pub stats: ExecutionStats,
}

/// Execution row
#[derive(Debug, Clone)]
pub struct ExecutionRow {
    pub values: Vec<ExecutionValue>,
}

/// Execution value
#[derive(Debug, Clone)]
pub enum ExecutionValue {
    Null,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Bytes(Vec<u8>),
    Timestamp(i64),
    Labels(BTreeMap<String, String>),
}

/// Execution statistics
#[derive(Debug, Clone, Default)]
pub struct ExecutionStats {
    pub rows_read: usize,
    pub rows_filtered: usize,
    pub rows_returned: usize,
    pub bytes_read: usize,
    pub execution_time_ms: f64,
    pub cache_hits: usize,
    pub cache_misses: usize,
}

/// Query executor
pub struct QueryExecutor {
    function_registry: FunctionRegistry,
    context: ExecutionContext,
}

impl QueryExecutor {
    pub fn new(context: ExecutionContext) -> Self {
        QueryExecutor {
            function_registry: FunctionRegistry::new(),
            context,
        }
    }

    /// Execute an expression
    pub fn execute(&self, expr: &Expr) -> Result<ExecutionResult, String> {
        let start = std::time::Instant::now();
        let rows = self.execute_expr(expr)?;
        let elapsed = start.elapsed();

        Ok(ExecutionResult {
            columns: vec![],
            rows: rows.clone(),
            total_rows: rows.len(),
            execution_time_ms: elapsed.as_millis() as f64,
            stats: ExecutionStats {
                rows_read: rows.len(),
                rows_filtered: 0,
                rows_returned: rows.len(),
                bytes_read: 0,
                execution_time_ms: elapsed.as_millis() as f64,
                cache_hits: 0,
                cache_misses: 0,
            },
        })
    }

    fn execute_expr(&self, expr: &Expr) -> Result<Vec<ExecutionRow>, String> {
        match expr {
            Expr::Number(n) => Ok(vec![ExecutionRow {
                values: vec![ExecutionValue::Float(*n)],
            }]),
            Expr::String(s) => Ok(vec![ExecutionRow {
                values: vec![ExecutionValue::String(s.clone())],
            }]),
            Expr::Identifier(name) => {
                // Look up identifier
                Ok(vec![ExecutionRow {
                    values: vec![ExecutionValue::String(name.clone())],
                }])
            }
            Expr::VectorSelector { name, matchers } => self.execute_vector_selector(name, matchers),
            Expr::MatrixSelector {
                name,
                matchers,
                range,
                offset,
            } => self.execute_matrix_selector(name, matchers, range, offset),
            Expr::FunctionCall { name, args } => self.execute_function(name, args),
            Expr::Aggregation {
                op,
                expr,
                grouping,
                without,
            } => self.execute_aggregation(op, expr, grouping, *without),
            Expr::BinaryOp {
                op,
                left,
                right,
                bool_modifier,
            } => self.execute_binary_op(op, left, right, *bool_modifier),
            Expr::UnaryOp { op, expr } => self.execute_unary_op(op, expr),
            Expr::Paren(inner) => self.execute_expr(inner),
            Expr::Subquery {
                expr,
                range,
                step,
                offset,
            } => self.execute_subquery(expr, range, step, offset),
        }
    }

    fn execute_vector_selector(
        &self,
        _name: &str,
        _matchers: &[crate::parser::Matcher],
    ) -> Result<Vec<ExecutionRow>, String> {
        // In production, would scan storage
        Ok(Vec::new())
    }

    fn execute_matrix_selector(
        &self,
        _name: &str,
        _matchers: &[crate::parser::Matcher],
        _range: &Duration,
        _offset: &Option<Duration>,
    ) -> Result<Vec<ExecutionRow>, String> {
        // In production, would scan storage with time range
        Ok(Vec::new())
    }

    fn execute_function(&self, name: &str, args: &[Expr]) -> Result<Vec<ExecutionRow>, String> {
        // Evaluate arguments
        let mut arg_results = Vec::new();
        for arg in args {
            arg_results.push(self.execute_expr(arg)?);
        }

        // Look up function
        let func = self
            .function_registry
            .get(name)
            .ok_or_else(|| format!("Unknown function: {}", name))?;

        // Convert to samples and call
        let samples: Vec<super::MetricSample> = arg_results
            .iter()
            .flat_map(|rows| {
                rows.iter().map(|r| super::MetricSample {
                    labels: BTreeMap::new(),
                    value: match &r.values[0] {
                        ExecutionValue::Float(v) => *v,
                        ExecutionValue::Int(v) => *v as f64,
                        _ => 0.0,
                    },
                    timestamp: self.context.timestamp,
                })
            })
            .collect();

        let ctx = super::functions::FunctionContext::new(
            self.context.timestamp,
            self.context.step,
            self.context.range,
        );
        let result = func.call(&samples, &ctx)?;

        // Convert result back
        Ok(result
            .iter()
            .map(|s| ExecutionRow {
                values: vec![ExecutionValue::Float(s.value)],
            })
            .collect())
    }

    fn execute_aggregation(
        &self,
        _op: &crate::parser::AggregationOp,
        expr: &Expr,
        _grouping: &[String],
        _without: bool,
    ) -> Result<Vec<ExecutionRow>, String> {
        let rows = self.execute_expr(expr)?;

        // Group and aggregate
        // Simplified: just return rows
        Ok(rows)
    }

    fn execute_binary_op(
        &self,
        op: &crate::parser::BinaryOperator,
        left: &Expr,
        right: &Expr,
        _bool_modifier: bool,
    ) -> Result<Vec<ExecutionRow>, String> {
        let left_rows = self.execute_expr(left)?;
        let right_rows = self.execute_expr(right)?;

        // Apply binary operation
        let mut result = Vec::new();
        for (l, r) in left_rows.iter().zip(right_rows.iter()) {
            let lval = match &l.values[0] {
                ExecutionValue::Float(v) => *v,
                ExecutionValue::Int(v) => *v as f64,
                _ => 0.0,
            };
            let rval = match &r.values[0] {
                ExecutionValue::Float(v) => *v,
                ExecutionValue::Int(v) => *v as f64,
                _ => 0.0,
            };

            let val = match op {
                crate::parser::BinaryOperator::Add => lval + rval,
                crate::parser::BinaryOperator::Sub => lval - rval,
                crate::parser::BinaryOperator::Mul => lval * rval,
                crate::parser::BinaryOperator::Div => {
                    if rval != 0.0 {
                        lval / rval
                    } else {
                        f64::NAN
                    }
                }
                crate::parser::BinaryOperator::Mod => lval % rval,
                crate::parser::BinaryOperator::Pow => lval.powf(rval),
                crate::parser::BinaryOperator::Equal => {
                    if (lval - rval).abs() < f64::EPSILON {
                        1.0
                    } else {
                        0.0
                    }
                }
                crate::parser::BinaryOperator::NotEqual => {
                    if (lval - rval).abs() >= f64::EPSILON {
                        1.0
                    } else {
                        0.0
                    }
                }
                crate::parser::BinaryOperator::GreaterThan => {
                    if lval > rval {
                        1.0
                    } else {
                        0.0
                    }
                }
                crate::parser::BinaryOperator::LessThan => {
                    if lval < rval {
                        1.0
                    } else {
                        0.0
                    }
                }
                crate::parser::BinaryOperator::GreaterThanOrEqual => {
                    if lval >= rval {
                        1.0
                    } else {
                        0.0
                    }
                }
                crate::parser::BinaryOperator::LessThanOrEqual => {
                    if lval <= rval {
                        1.0
                    } else {
                        0.0
                    }
                }
                _ => f64::NAN,
            };

            result.push(ExecutionRow {
                values: vec![ExecutionValue::Float(val)],
            });
        }

        Ok(result)
    }

    fn execute_unary_op(
        &self,
        op: &crate::parser::UnaryOperator,
        expr: &Expr,
    ) -> Result<Vec<ExecutionRow>, String> {
        let rows = self.execute_expr(expr)?;

        Ok(rows
            .into_iter()
            .map(|row| {
                let val = match &row.values[0] {
                    ExecutionValue::Float(v) => *v,
                    ExecutionValue::Int(v) => *v as f64,
                    _ => 0.0,
                };
                let result = match op {
                    crate::parser::UnaryOperator::Plus => val,
                    crate::parser::UnaryOperator::Minus => -val,
                };
                ExecutionRow {
                    values: vec![ExecutionValue::Float(result)],
                }
            })
            .collect())
    }

    fn execute_subquery(
        &self,
        expr: &Expr,
        _range: &Duration,
        _step: &Option<Duration>,
        _offset: &Option<Duration>,
    ) -> Result<Vec<ExecutionRow>, String> {
        // In production, would execute subquery with range
        self.execute_expr(expr)
    }
}

/// Query optimizer for execution
pub struct ExecutionOptimizer {
    rules: Vec<Box<dyn ExecutionRule>>,
}

/// Execution optimization rule
pub trait ExecutionRule: Send + Sync {
    fn name(&self) -> &'static str;
    fn apply(&self, plan: &mut PlanNode) -> bool;
}

/// Plan node for execution
#[derive(Debug, Clone)]
pub struct PlanNode {
    pub op: ExecutionOp,
    pub children: Vec<PlanNode>,
    pub params: BTreeMap<String, String>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ExecutionOp {
    Scan,
    Filter,
    Project,
    Aggregate,
    Sort,
    Limit,
    Offset,
    Join,
    Union,
    Distinct,
    Cache,
    Prefetch,
}

impl PlanNode {
    pub fn new(op: ExecutionOp) -> Self {
        PlanNode {
            op,
            children: Vec::new(),
            params: BTreeMap::new(),
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
}

impl ExecutionOptimizer {
    pub fn new() -> Self {
        let mut optimizer = ExecutionOptimizer { rules: Vec::new() };

        optimizer.add_rule(Box::new(PrefetchRule));
        optimizer.add_rule(Box::new(CacheRule));
        optimizer.add_rule(Box::new(VectorizationRule));
        optimizer.add_rule(Box::new(ParallelismRule));

        optimizer
    }

    pub fn add_rule(&mut self, rule: Box<dyn ExecutionRule>) {
        self.rules.push(rule);
    }

    pub fn optimize(&self, plan: &mut PlanNode) -> OptimizationResult {
        let mut applied = Vec::new();
        let mut iterations = 0;

        while iterations < 10 {
            let mut changed = false;
            for rule in &self.rules {
                if rule.apply(plan) {
                    applied.push(rule.name().to_string());
                    changed = true;
                }
            }
            if !changed {
                break;
            }
            iterations += 1;
        }

        OptimizationResult {
            applied_rules: applied,
            iterations,
        }
    }
}

impl Default for ExecutionOptimizer {
    fn default() -> Self {
        Self::new()
    }
}

/// Optimization result
#[derive(Debug, Clone)]
pub struct OptimizationResult {
    pub applied_rules: Vec<String>,
    pub iterations: usize,
}

/// Prefetch rule
pub struct PrefetchRule;

impl ExecutionRule for PrefetchRule {
    fn name(&self) -> &'static str {
        "prefetch"
    }
    fn apply(&self, plan: &mut PlanNode) -> bool {
        if plan.op == ExecutionOp::Scan && !plan.params.contains_key("prefetch") {
            plan.params
                .insert("prefetch".to_string(), "true".to_string());
            return true;
        }
        false
    }
}

/// Cache rule
pub struct CacheRule;

impl ExecutionRule for CacheRule {
    fn name(&self) -> &'static str {
        "cache"
    }
    fn apply(&self, plan: &mut PlanNode) -> bool {
        if plan.op == ExecutionOp::Scan && !plan.params.contains_key("cache") {
            plan.params.insert("cache".to_string(), "true".to_string());
            return true;
        }
        false
    }
}

/// Vectorization rule
pub struct VectorizationRule;

impl ExecutionRule for VectorizationRule {
    fn name(&self) -> &'static str {
        "vectorization"
    }
    fn apply(&self, plan: &mut PlanNode) -> bool {
        if !plan.params.contains_key("vectorized") {
            plan.params
                .insert("vectorized".to_string(), "true".to_string());
            return true;
        }
        false
    }
}

/// Parallelism rule
pub struct ParallelismRule;

impl ExecutionRule for ParallelismRule {
    fn name(&self) -> &'static str {
        "parallelism"
    }
    fn apply(&self, plan: &mut PlanNode) -> bool {
        if plan.op == ExecutionOp::Scan && !plan.params.contains_key("parallel") {
            plan.params
                .insert("parallel".to_string(), "true".to_string());
            plan.params.insert("threads".to_string(), "4".to_string());
            return true;
        }
        false
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execution_context() {
        let ctx = ExecutionContext::new(1000)
            .with_step(10)
            .with_range(300)
            .with_timeout(60000)
            .with_max_rows(5000);

        assert_eq!(ctx.timestamp, 1000);
        assert_eq!(ctx.step, 10);
        assert_eq!(ctx.range, 300);
        assert_eq!(ctx.timeout, 60000);
        assert_eq!(ctx.max_rows, 5000);
    }

    #[test]
    fn test_execute_number() {
        let ctx = ExecutionContext::new(1000);
        let executor = QueryExecutor::new(ctx);
        let result = executor.execute(&Expr::Number(42.0)).unwrap();
        assert_eq!(result.rows.len(), 1);
    }

    #[test]
    fn test_execute_string() {
        let ctx = ExecutionContext::new(1000);
        let executor = QueryExecutor::new(ctx);
        let result = executor
            .execute(&Expr::String("hello".to_string()))
            .unwrap();
        assert_eq!(result.rows.len(), 1);
    }

    #[test]
    fn test_execute_binary_add() {
        let ctx = ExecutionContext::new(1000);
        let executor = QueryExecutor::new(ctx);
        let expr = Expr::BinaryOp {
            op: crate::parser::BinaryOperator::Add,
            left: Box::new(Expr::Number(10.0)),
            right: Box::new(Expr::Number(5.0)),
            bool_modifier: false,
        };
        let result = executor.execute(&expr).unwrap();
        assert_eq!(result.rows.len(), 1);
    }

    #[test]
    fn test_execute_unary_minus() {
        let ctx = ExecutionContext::new(1000);
        let executor = QueryExecutor::new(ctx);
        let expr = Expr::UnaryOp {
            op: crate::parser::UnaryOperator::Minus,
            expr: Box::new(Expr::Number(5.0)),
        };
        let result = executor.execute(&expr).unwrap();
        assert_eq!(result.rows.len(), 1);
    }

    #[test]
    fn test_plan_node() {
        let node =
            PlanNode::new(ExecutionOp::Scan).with_param("table".to_string(), "metrics".to_string());

        assert_eq!(node.op, ExecutionOp::Scan);
        assert_eq!(node.params.get("table"), Some(&"metrics".to_string()));
    }

    #[test]
    fn test_plan_node_with_child() {
        let child = PlanNode::new(ExecutionOp::Scan);
        let parent = PlanNode::new(ExecutionOp::Filter).with_child(child);

        assert_eq!(parent.children.len(), 1);
    }

    #[test]
    fn test_query_builder() {
        let plan = crate::executor_advanced2::QueryBuilder::new()
            .scan("metrics")
            .filter("host = 'server1'")
            .project(vec!["value".to_string()])
            .limit(100)
            .build();

        assert!(plan.steps.last().unwrap().starts_with("limit:"))
    }

    #[test]
    fn test_execution_optimizer() {
        let optimizer = ExecutionOptimizer::new();
        let mut plan = PlanNode::new(ExecutionOp::Scan);
        let result = optimizer.optimize(&mut plan);
        assert!(!result.applied_rules.is_empty());
    }

    #[test]
    fn test_prefetch_rule() {
        let rule = PrefetchRule;
        let mut plan = PlanNode::new(ExecutionOp::Scan);
        assert!(rule.apply(&mut plan));
        assert_eq!(plan.params.get("prefetch"), Some(&"true".to_string()));
    }

    #[test]
    fn test_cache_rule() {
        let rule = CacheRule;
        let mut plan = PlanNode::new(ExecutionOp::Scan);
        assert!(rule.apply(&mut plan));
        assert_eq!(plan.params.get("cache"), Some(&"true".to_string()));
    }

    #[test]
    fn test_vectorization_rule() {
        let rule = VectorizationRule;
        let mut plan = PlanNode::new(ExecutionOp::Scan);
        assert!(rule.apply(&mut plan));
        assert_eq!(plan.params.get("vectorized"), Some(&"true".to_string()));
    }

    #[test]
    fn test_parallelism_rule() {
        let rule = ParallelismRule;
        let mut plan = PlanNode::new(ExecutionOp::Scan);
        assert!(rule.apply(&mut plan));
        assert_eq!(plan.params.get("parallel"), Some(&"true".to_string()));
        assert_eq!(plan.params.get("threads"), Some(&"4".to_string()));
    }

    #[test]
    fn test_execution_result() {
        let result = ExecutionResult {
            columns: vec!["host".to_string(), "value".to_string()],
            rows: vec![],
            total_rows: 0,
            execution_time_ms: 10.5,
            stats: ExecutionStats::default(),
        };
        assert_eq!(result.columns.len(), 2);
    }

    #[test]
    fn test_execution_value() {
        let val = ExecutionValue::Float(42.0);
        match val {
            ExecutionValue::Float(v) => assert_eq!(v, 42.0),
            _ => panic!("Expected Float"),
        }
    }
}
