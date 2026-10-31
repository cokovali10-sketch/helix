//! Helix Query Engine - Vectorized Executor
//! 
//! This module provides a vectorized execution engine for the Helix query language.
//! Uses columnar data processing for SIMD-like performance.

use std::collections::BTreeMap;
use std::sync::Arc;

use super::parser::{Expr, Duration};
use super::functions::FunctionRegistry;

/// Columnar data batch
#[derive(Debug, Clone)]
pub struct DataBatch {
    pub columns: Vec<Column>,
    pub num_rows: usize,
}

/// Column in a data batch
#[derive(Debug, Clone)]
pub enum Column {
    Float64(Vec<f64>),
    Int64(Vec<i64>),
    String(Vec<String>),
    Timestamp(Vec<i64>),
    Labels(Vec<BTreeMap<String, String>>),
}

impl Column {
    pub fn len(&self) -> usize {
        match self {
            Column::Float64(v) => v.len(),
            Column::Int64(v) => v.len(),
            Column::String(v) => v.len(),
            Column::Timestamp(v) => v.len(),
            Column::Labels(v) => v.len(),
        }
    }
    
    pub fn is_empty(&self) -> bool {
        self.len() == 0
    }
}

/// Vectorized executor
pub struct VectorizedExecutor {
    function_registry: FunctionRegistry,
    batch_size: usize,
}

impl VectorizedExecutor {
    pub fn new(batch_size: usize) -> Self {
        VectorizedExecutor {
            function_registry: FunctionRegistry::new(),
            batch_size,
        }
    }
    
    /// Execute an expression against a data batch
    pub fn execute(&self, expr: &Expr, batch: &DataBatch) -> Result<DataBatch, String> {
        match expr {
            Expr::Number(n) => {
                let mut values = vec![*n; batch.num_rows];
                Ok(DataBatch {
                    columns: vec![Column::Float64(values)],
                    num_rows: batch.num_rows,
                })
            }
            Expr::Identifier(name) => {
                // Look up column by name
                self.get_column(batch, name)
            }
            Expr::VectorSelector { name, matchers } => {
                self.execute_vector_selector(name, matchers, batch)
            }
            Expr::FunctionCall { name, args } => {
                self.execute_function(name, args, batch)
            }
            Expr::BinaryOp { op, left, right, bool_modifier } => {
                self.execute_binary_op(*op, left, right, *bool_modifier, batch)
            }
            Expr::UnaryOp { op, expr } => {
                self.execute_unary_op(*op, expr, batch)
            }
            Expr::Aggregation { op, expr, grouping, without } => {
                self.execute_aggregation(*op, expr, grouping, *without, batch)
            }
            Expr::Paren(inner) => {
                self.execute(inner, batch)
            }
            _ => Err("Unsupported expression type".to_string()),
        }
    }
    
    fn get_column(&self, batch: &DataBatch, name: &str) -> Result<DataBatch, String> {
        // Find column by name (simplified)
        if batch.columns.is_empty() {
            return Err("No columns in batch".to_string());
        }
        Ok(batch.clone())
    }
    
    fn execute_vector_selector(&self, name: &str, matchers: &[super::parser::Matcher], batch: &DataBatch) -> Result<DataBatch, String> {
        // Filter rows based on matchers
        let mut mask = vec![true; batch.num_rows];
        
        for matcher in matchers {
            match matcher.op {
                super::parser::MatchOp::Equal => {
                    // Filter rows where label equals value
                }
                super::parser::MatchOp::NotEqual => {
                    // Filter rows where label doesn't equal value
                }
                super::parser::MatchOp::RegexMatch => {
                    // Filter rows matching regex
                }
                super::parser::MatchOp::RegexNotMatch => {
                    // Filter rows not matching regex
                }
            }
        }
        
        // Apply mask to batch
        Ok(batch.clone())
    }
    
    fn execute_function(&self, name: &str, args: &[Expr], batch: &DataBatch) -> Result<DataBatch, String> {
        // Evaluate arguments
        let mut arg_results = Vec::new();
        for arg in args {
            arg_results.push(self.execute(arg, batch)?);
        }
        
        // Look up function
        let func = self.function_registry.get(name)
            .ok_or_else(|| format!("Unknown function: {}", name))?;
        
        // Convert batch to samples
        let samples = self.batch_to_samples(&arg_results[0]);
        
        // Call function
        let ctx = super::functions::FunctionContext::new(0, 1, 0);
        let result = func.call(&samples, &ctx)?;
        
        // Convert result back to batch
        self.samples_to_batch(&result)
    }
    
    fn execute_binary_op(&self, op: super::parser::BinaryOperator, left: &Expr, right: &Expr, bool_modifier: bool, batch: &DataBatch) -> Result<DataBatch, String> {
        let left_batch = self.execute(left, batch)?;
        let right_batch = self.execute(right, batch)?;
        
        // Apply binary operation element-wise
        let left_vals = match &left_batch.columns[0] {
            Column::Float64(v) => v.clone(),
            _ => return Err("Expected float64 column".to_string()),
        };
        let right_vals = match &right_batch.columns[0] {
            Column::Float64(v) => v.clone(),
            _ => return Err("Expected float64 column".to_string()),
        };
        
        let mut result = Vec::with_capacity(left_vals.len());
        for i in 0..left_vals.len() {
            let val = match op {
                super::parser::BinaryOperator::Add => left_vals[i] + right_vals[i],
                super::parser::BinaryOperator::Sub => left_vals[i] - right_vals[i],
                super::parser::BinaryOperator::Mul => left_vals[i] * right_vals[i],
                super::parser::BinaryOperator::Div => {
                    if right_vals[i] != 0.0 {
                        left_vals[i] / right_vals[i]
                    } else {
                        f64::NAN
                    }
                }
                super::parser::BinaryOperator::Mod => left_vals[i] % right_vals[i],
                super::parser::BinaryOperator::Pow => left_vals[i].powf(right_vals[i]),
                super::parser::BinaryOperator::Equal => {
                    if (left_vals[i] - right_vals[i]).abs() < f64::EPSILON { 1.0 } else { 0.0 }
                }
                super::parser::BinaryOperator::NotEqual => {
                    if (left_vals[i] - right_vals[i]).abs() >= f64::EPSILON { 1.0 } else { 0.0 }
                }
                super::parser::BinaryOperator::GreaterThan => {
                    if left_vals[i] > right_vals[i] { 1.0 } else { 0.0 }
                }
                super::parser::BinaryOperator::LessThan => {
                    if left_vals[i] < right_vals[i] { 1.0 } else { 0.0 }
                }
                super::parser::BinaryOperator::GreaterThanOrEqual => {
                    if left_vals[i] >= right_vals[i] { 1.0 } else { 0.0 }
                }
                super::parser::BinaryOperator::LessThanOrEqual => {
                    if left_vals[i] <= right_vals[i] { 1.0 } else { 0.0 }
                }
                _ => f64::NAN,
            };
            result.push(val);
        }
        
        Ok(DataBatch {
            columns: vec![Column::Float64(result)],
            num_rows: result.len(),
        })
    }
    
    fn execute_unary_op(&self, op: super::parser::UnaryOperator, expr: &Expr, batch: &DataBatch) -> Result<DataBatch, String> {
        let inner = self.execute(expr, batch)?;
        let values = match &inner.columns[0] {
            Column::Float64(v) => v.clone(),
            _ => return Err("Expected float64 column".to_string()),
        };
        
        let result: Vec<f64> = values.iter().map(|v| match op {
            super::parser::UnaryOperator::Plus => *v,
            super::parser::UnaryOperator::Minus => -*v,
        }).collect();
        
        Ok(DataBatch {
            columns: vec![Column::Float64(result)],
            num_rows: result.len(),
        })
    }
    
    fn execute_aggregation(&self, op: super::parser::AggregationOp, expr: &Expr, grouping: &[String], without: bool, batch: &DataBatch) -> Result<DataBatch, String> {
        let inner = self.execute(expr, batch)?;
        
        // Group by labels and aggregate
        let values = match &inner.columns[0] {
            Column::Float64(v) => v.clone(),
            _ => return Err("Expected float64 column".to_string()),
        };
        
        // Simplified aggregation - just return sum
        let sum: f64 = values.iter().sum();
        
        Ok(DataBatch {
            columns: vec![Column::Float64(vec![sum])],
            num_rows: 1,
        })
    }
    
    fn batch_to_samples(&self, batch: &DataBatch) -> Vec<super::MetricSample> {
        let values = match &batch.columns[0] {
            Column::Float64(v) => v.clone(),
            _ => vec![],
        };
        
        values.iter().enumerate().map(|(i, v)| {
            super::MetricSample {
                labels: BTreeMap::new(),
                value: *v,
                timestamp: 0,
            }
        }).collect()
    }
    
    fn samples_to_batch(&self, samples: &[super::MetricSample]) -> Result<DataBatch, String> {
        let values: Vec<f64> = samples.iter().map(|s| s.value).collect();
        Ok(DataBatch {
            columns: vec![Column::Float64(values.clone())],
            num_rows: values.len(),
        })
    }
}

/// Query execution plan
#[derive(Debug, Clone)]
pub struct ExecutionPlan {
    pub steps: Vec<ExecutionStep>,
    pub estimated_rows: usize,
    pub estimated_cost: f64,
}

/// Execution step
#[derive(Debug, Clone)]
pub enum ExecutionStep {
    Scan { table: String, columns: Vec<String> },
    Filter { predicate: String },
    Project { columns: Vec<String> },
    Aggregate { op: String, group_by: Vec<String> },
    Sort { columns: Vec<String>, ascending: bool },
    Limit { count: usize },
    Offset { count: usize },
    Join { left: String, right: String, on: String },
    Union { inputs: Vec<String> },
    Intersect { inputs: Vec<String> },
}

/// Plan builder
pub struct PlanBuilder {
    plan: ExecutionPlan,
}

impl PlanBuilder {
    pub fn new() -> Self {
        PlanBuilder {
            plan: ExecutionPlan {
                steps: Vec::new(),
                estimated_rows: 0,
                estimated_cost: 0.0,
            },
        }
    }
    
    pub fn scan(mut self, table: &str, columns: Vec<String>) -> Self {
        self.plan.steps.push(ExecutionStep::Scan {
            table: table.to_string(),
            columns,
        });
        self
    }
    
    pub fn filter(mut self, predicate: &str) -> Self {
        self.plan.steps.push(ExecutionStep::Filter {
            predicate: predicate.to_string(),
        });
        self
    }
    
    pub fn project(mut self, columns: Vec<String>) -> Self {
        self.plan.steps.push(ExecutionStep::Project { columns });
        self
    }
    
    pub fn aggregate(mut self, op: &str, group_by: Vec<String>) -> Self {
        self.plan.steps.push(ExecutionStep::Aggregate {
            op: op.to_string(),
            group_by,
        });
        self
    }
    
    pub fn sort(mut self, columns: Vec<String>, ascending: bool) -> Self {
        self.plan.steps.push(ExecutionStep::Sort { columns, ascending });
        self
    }
    
    pub fn limit(mut self, count: usize) -> Self {
        self.plan.steps.push(ExecutionStep::Limit { count });
        self
    }
    
    pub fn offset(mut self, count: usize) -> Self {
        self.plan.steps.push(ExecutionStep::Offset { count });
        self
    }
    
    pub fn build(self) -> ExecutionPlan {
        self.plan
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_vectorized_executor_creation() {
        let executor = VectorizedExecutor::new(1024);
        assert_eq!(executor.batch_size, 1024);
    }
    
    #[test]
    fn test_execute_number() {
        let executor = VectorizedExecutor::new(1024);
        let batch = DataBatch {
            columns: vec![Column::Float64(vec![1.0, 2.0, 3.0])],
            num_rows: 3,
        };
        let result = executor.execute(&Expr::Number(42.0), &batch).unwrap();
        match &result.columns[0] {
            Column::Float64(v) => {
                assert_eq!(v, &vec![42.0, 42.0, 42.0]);
            }
            _ => panic!("Expected Float64 column"),
        }
    }
    
    #[test]
    fn test_execute_binary_add() {
        let executor = VectorizedExecutor::new(1024);
        let batch = DataBatch {
            columns: vec![Column::Float64(vec![1.0, 2.0, 3.0])],
            num_rows: 3,
        };
        let expr = Expr::BinaryOp {
            op: super::parser::BinaryOperator::Add,
            left: Box::new(Expr::Number(10.0)),
            right: Box::new(Expr::Number(5.0)),
            bool_modifier: false,
        };
        let result = executor.execute(&expr, &batch).unwrap();
        match &result.columns[0] {
            Column::Float64(v) => {
                assert_eq!(v, &vec![15.0, 15.0, 15.0]);
            }
            _ => panic!("Expected Float64 column"),
        }
    }
    
    #[test]
    fn test_execute_unary_minus() {
        let executor = VectorizedExecutor::new(1024);
        let batch = DataBatch {
            columns: vec![Column::Float64(vec![1.0, 2.0, 3.0])],
            num_rows: 3,
        };
        let expr = Expr::UnaryOp {
            op: super::parser::UnaryOperator::Minus,
            expr: Box::new(Expr::Number(5.0)),
        };
        let result = executor.execute(&expr, &batch).unwrap();
        match &result.columns[0] {
            Column::Float64(v) => {
                assert_eq!(v, &vec![-5.0, -5.0, -5.0]);
            }
            _ => panic!("Expected Float64 column"),
        }
    }
    
    #[test]
    fn test_plan_builder() {
        let plan = PlanBuilder::new()
            .scan("metrics", vec!["host".to_string(), "value".to_string()])
            .filter("host = 'server1'")
            .project(vec!["value".to_string()])
            .aggregate("sum", vec!["host".to_string()])
            .sort(vec!["value".to_string()], false)
            .limit(100)
            .build();
        
        assert_eq!(plan.steps.len(), 6);
    }
    
    #[test]
    fn test_column_len() {
        let col = Column::Float64(vec![1.0, 2.0, 3.0]);
        assert_eq!(col.len(), 3);
        assert!(!col.is_empty());
        
        let empty = Column::Float64(vec![]);
        assert!(empty.is_empty());
    }
    
    #[test]
    fn test_data_batch() {
        let batch = DataBatch {
            columns: vec![Column::Float64(vec![1.0, 2.0])],
            num_rows: 2,
        };
        assert_eq!(batch.num_rows, 2);
    }
}