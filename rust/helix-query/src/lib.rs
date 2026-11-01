//! Helix Query Engine - Parser → Scheduler → Vectorized Executor
//! 
//! This module provides the query engine for the Helix observability platform.
//! It implements a custom query language (PromQL-like) with optimized execution
//! for time-series analysis over the TSDB storage engine.

use std::collections::BTreeMap;
use std::fmt;

pub mod functions;
pub mod optimizer;
pub mod lexer;
pub mod parser;
pub mod executor;
pub mod planner;
pub mod executor_advanced;
pub use functions::*;
pub use optimizer::*;
pub use lexer::*;
pub use parser::*;
pub use executor::*;
pub use planner::*;
pub use executor_advanced::*;

// ---------------------------------------------------------------------------
// Query Language Types
// ---------------------------------------------------------------------------

/// QueryType represents the type of a query expression
#[derive(Debug, Clone, PartialEq)]
pub enum QueryType {
    InstantVector { labels: BTreeMap<String, String> },
    RangeVector { labels: BTreeMap<String, String>, duration: String },
    Scalar,
    String,
    VectorSelector,
}

/// Query represents a parsed Helix query
#[derive(Debug, Clone)]
pub struct Query {
    pub query_type: QueryType,
    pub expression: String,
    pub labels: BTreeMap<String, String>,
    pub time_range: TimeRange,
}

/// TimeRange represents a time window for queries
#[derive(Debug, Clone, PartialEq)]
pub struct TimeRange {
    pub start: i64,
    pub end: i64,
    pub offset: i64,
}

/// LabelSelector selects elements by label matchers
#[derive(Debug, Clone, PartialEq)]
pub struct LabelSelector {
    pub labels: BTreeMap<String, String>,
    pub matchers: Vec<Matcher>,
}

/// Matcher represents a label matcher in PromQL style
#[derive(Debug, Clone, PartialEq)]
pub enum Matcher {
    Equals { label: String, value: String },
    NotEquals { label: String, value: String },
    Exists { label: String },
    NotExists { label: String },
    Regex { label: String, pattern: String },
}

// ---------------------------------------------------------------------------
// Parser Module
// ---------------------------------------------------------------------------

/// Parser parses Helix query language expressions
pub struct Parser {
    input: String,
    pos: usize,
}

impl Parser {
    /// Creates a new parser from the given input
    pub fn new(input: String) -> Self {
        Parser { input, pos: 0 }
    }

    /// Parses a complete query expression
    pub fn parse(&mut self) -> Result<Query, ParseError> {
        self.skip_whitespace();
        let query_type = self.parse_query_type()?;
        let labels = self.parse_label_selector()?;
        let time_range = self.parse_time_range()?;
        
        Ok(Query {
            query_type,
            expression: self.input.clone(),
            labels,
            time_range,
        })
    }

    fn parse_query_type(&mut self) -> Result<QueryType, ParseError> {
        self.skip_whitespace();
        
        // Parse metric selector (e.g., cpu_usage{host="server1"})
        let labels = self.parse_label_selector()?;
        
        self.skip_whitespace();
        
        // Check for range vector operator (e.g., cpu_usage[30s])
        let c = self.peek();
        if c == '[' {
            self.consume('[')?;
            let duration = self.parse_duration()?;
            self.consume(']')?;
            
            Ok(QueryType::RangeVector { labels, duration })
        } else {
            Ok(QueryType::InstantVector { labels })
        }
    }

    fn parse_label_selector(&mut self) -> Result<BTreeMap<String, String>, ParseError> {
        let mut labels = BTreeMap::new();
        
        self.skip_whitespace();
        
        // Parse {label="value"} or {label=~"regex"}
        if self.peek() == '{' {
            self.consume('{')?;
            
            loop {
                self.skip_whitespace();
                let label = self.parse_identifier()?;
                self.skip_whitespace();
                
                let (match_op, value) = self.parse_matcher()?;
                
                labels.insert(label, value);
                
                self.skip_whitespace();
                if self.peek() == '}' {
                    self.consume('}')?;
                    break;
                }
                if self.peek() != ',' {
                    return Err(ParseError::ExpectedComma);
                }
                self.consume(',')?;
            }
        }
        
        Ok(labels)
    }

    fn parse_matcher(&mut self) -> Result<(String, String), ParseError> {
        let c = self.peek();
        
        if c == '=' {
            self.consume('=')?;
            let value = self.parse_string()?;
            Ok(("=".to_string(), value))
        } else if c == '!' {
            self.consume('!')?;
            self.consume('=')?;
            let value = self.parse_string()?;
            Ok(("!=".to_string(), value))
        } else if c == '~' {
            self.consume('~')?;
            let pattern = self.parse_string()?;
            Ok(("=~".to_string(), pattern))
        } else {
            // Default equals
            self.consume('=')?;
            let value = self.parse_string()?;
            Ok(("=".to_string(), value))
        }
    }

    fn parse_duration(&mut self) -> Result<String, ParseError> {
        let start = self.pos;
        while self.pos < self.input.len() && self.input.as_bytes()[self.pos].is_ascii_digit() {
            self.pos += 1;
        }
        
        if self.pos == start {
            return Err(ParseError::ExpectedDuration);
        }
        
        Ok(self.input[start..self.pos].to_string())
    }

    fn parse_identifier(&mut self) -> Result<String, ParseError> {
        let start = self.pos;
        // Identifiers start with letter or underscore
        if self.pos < self.input.len() {
            let c = self.input.as_bytes()[self.pos] as char;
            if c.is_alphabetic() || c == '_' {
                self.pos += 1;
                while self.pos < self.input.len() {
                    let c = self.input.as_bytes()[self.pos] as char;
                    if c.is_alphanumeric() || c == '_' {
                        self.pos += 1;
                    } else {
                        break;
                    }
                }
                return Ok(self.input[start..self.pos].to_string());
            }
        }
        Err(ParseError::ExpectedIdentifier)
    }

    fn parse_string(&mut self) -> Result<String, ParseError> {
        let start = self.pos;
        
        if self.pos >= self.input.len() {
            return Err(ParseError::UnexpectedEnd);
        }
        
        let c = self.input.as_bytes()[self.pos] as char;
        if c == '"' {
            self.pos += 1; // consume opening quote
            let mut result = String::new();
            
            while self.pos < self.input.len() {
                let c = self.input.as_bytes()[self.pos] as char;
                self.pos += 1;
                
                if c == '"' {
                    return Ok(result);
                }
                result.push(c);
            }
        }
        
        Err(ParseError::UnterminatedString)
    }

    fn peek(&self) -> char {
        if self.pos >= self.input.len() {
            return '\0';
        }
        self.input.as_bytes()[self.pos] as char
    }

    fn consume(&mut self, c: char) -> Result<(), ParseError> {
        if self.peek() == c {
            self.pos += 1;
            Ok(())
        } else {
            Err(ParseError::ExpectedCharacter(c))
        }
    }

    fn skip_whitespace(&mut self) {
        while self.pos < self.input.len() && self.input.as_bytes()[self.pos].is_ascii_whitespace() {
            self.pos += 1;
        }
    }
}

// ---------------------------------------------------------------------------
// Parse Errors
// ---------------------------------------------------------------------------

#[derive(Debug)]
pub enum ParseError {
    ExpectedCharacter(char),
    ExpectedIdentifier,
    ExpectedDuration,
    UnexpectedEnd,
    UnterminatedString,
    ExpectedComma,
    InvalidSyntax(String),
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ParseError::ExpectedCharacter(c) => write!(f, "expected character '{}'", c),
            ParseError::ExpectedIdentifier => write!(f, "expected identifier"),
            ParseError::ExpectedDuration => write!(f, "expected duration (e.g., '30s')"),
            ParseError::UnexpectedEnd => write!(f, "unexpected end of input"),
            ParseError::UnterminatedString => write!(f, "unterminated string"),
            ParseError::ExpectedComma => write!(f, "expected comma"),
            ParseError::InvalidSyntax(msg) => write!(f, "invalid syntax: {}", msg),
        }
    }
}

impl std::error::Error for ParseError {}

// ---------------------------------------------------------------------------
// Scheduler Module
// ---------------------------------------------------------------------------

/// Scheduler plans the execution of a parsed query
pub struct Scheduler {
    query: Query,
}

impl Scheduler {
    /// Creates a new scheduler from a parsed query
    pub fn new(query: Query) -> Self {
        Scheduler { query }
    }

    /// Plan the query execution strategy
    pub fn plan(&self) -> QueryPlan {
        let mut plan = QueryPlan::default();
        
        match self.query.query_type {
            QueryType::InstantVector { .. } => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "compute_instant".to_string(),
                    "post_process".to_string(),
                ];
            }
            QueryType::RangeVector { .. } => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "aggregate_range".to_string(),
                    "compute_range".to_string(),
                    "post_process".to_string(),
                ];
            }
            QueryType::Scalar | QueryType::String => {
                plan.steps = vec![
                    "fetch_from_storage".to_string(),
                    "aggregate".to_string(),
                    "compute_scalar".to_string(),
                ];
            }
        }
        
        // Apply label filters
        if !self.query.labels.is_empty() {
            plan.steps.insert(0, "filter_labels".to_string());
        }
        
        // Apply time range filtering
        plan.steps.insert(1, "filter_time_range".to_string());
        
        plan
    }
}

/// QueryPlan represents the execution plan for a query
#[derive(Debug, Clone, Default)]
pub struct QueryPlan {
    pub steps: Vec<String>,
    pub estimated_cost: f64,
    pub estimated_result_size: usize,
}

// ---------------------------------------------------------------------------
// Vectorized Executor Module
// ---------------------------------------------------------------------------

/// Executor executes a query plan against stored data
pub struct Executor {
    plan: QueryPlan,
}

impl Executor {
    /// Creates a new executor with the given plan
    pub fn new(plan: QueryPlan) -> Self {
        Executor { plan }
    }

    /// Execute the query plan against the data source
    pub fn execute(&self, data: &[MetricSample]) -> QueryResult {
        let mut current_data = data.to_vec();
        
        for step in &self.plan.steps {
            match step.as_str() {
                "filter_labels" => {
                    current_data = self.filter_labels(&current_data);
                }
                "filter_time_range" => {
                    current_data = self.filter_time_range(&current_data);
                }
                "fetch_from_storage" => {
                    // In production, this would fetch from TSDB
                    // For now, use provided data
                }
                "aggregate_range" => {
                    current_data = self.aggregate_range(&current_data);
                }
                "compute_instant" => {
                    current_data = self.compute_instant(&current_data);
                }
                "compute_range" => {
                    current_data = self.compute_range(&current_data);
                }
                "compute_scalar" => {
                    current_data = self.compute_scalar(&current_data);
                }
                "post_process" => {
                    current_data = self.post_process(&current_data);
                }
                _ => {}
            }
        }
        
        QueryResult {
            result_type: "vector".to_string(),
            result: current_data,
            ..Default::default()
        }
    }

    fn filter_labels(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        data.iter()
            .filter(|sample| {
                // Apply label filtering logic
                true
            })
            .cloned()
            .collect()
    }

    fn filter_time_range(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        data.iter()
            .filter(|sample| {
                sample.timestamp >= 0 && sample.timestamp > 0 // Simplified
            })
            .take(1000) // Limit result size
            .cloned()
            .collect()
    }

    fn aggregate_range(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        // Compute range vector: last - first / step
        // Simplified implementation
        data.iter()
            .filter(|s| s.timestamp > 0)
            .map(|sample| MetricSample {
                value: sample.value,
                timestamp: sample.timestamp,
                labels: sample.labels.clone(),
            })
            .collect()
    }

    fn compute_instant(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        // Return instant vector (last snapshot)
        data.last()
            .cloned()
            .into_iter()
            .collect()
    }

    fn compute_range(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        // Compute range metrics
        data.iter()
            .map(|sample| MetricSample {
                value: sample.value * 2.0, // Simplified
                timestamp: sample.timestamp,
                labels: sample.labels.clone(),
            })
            .collect()
    }

    fn compute_scalar(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        // Compute scalar result
        if let Some(sample) = data.first() {
            vec![MetricSample {
                value: sample.value,
                timestamp: sample.timestamp,
                labels: sample.labels.clone(),
            }]
        } else {
            vec![]
        }
    }

    fn post_process(&self, data: &[MetricSample]) -> Vec<MetricSample> {
        // Post-processing: sorting, deduplication, limiting
        let mut result: Vec<_> = data.iter().cloned().collect();
        result.sort_by(|a, b| b.value.partial_cmp(&a.value).unwrap_or(std::cmp::Ordering::Equal));
        result.truncate(100); // Limit results
        result
    }
}

// ---------------------------------------------------------------------------
// Data Types for Execution
// ---------------------------------------------------------------------------

/// MetricSample represents a single metric data point
#[derive(Debug, Clone, PartialEq)]
pub struct MetricSample {
    pub labels: BTreeMap<String, String>,
    pub value: f64,
    pub timestamp: i64,
}

/// QueryResult represents the result of a query execution
#[derive(Debug, Clone, Default)]
pub struct QueryResult {
    pub result_type: String,
    pub result: Vec<MetricSample>,
    pub error: Option<String>,
    pub execution_time_ns: u64,
    pub stats: QueryStats,
}

/// QueryStats contains execution statistics
#[derive(Debug, Clone, Default)]
pub struct QueryStats {
    pub rows_read: usize,
    pub rows_filtered: usize,
    pub execution_time_ms: f64,
    pub cache_hits: usize,
    pub cache_misses: usize,
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/// Executes a Helix query against provided data
pub fn execute_query(expression: &str, data: &[MetricSample]) -> Result<QueryResult, ParseError> {
    let mut parser = Parser::new(expression.to_string());
    let query = parser.parse()?;
    
    let scheduler = Scheduler::new(query);
    let plan = scheduler.plan();
    
    let executor = Executor::new(plan);
    let result = executor.execute(data);
    
    Ok(result)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_parser_instant_vector() {
        let input = "cpu_usage{host=\"server1\"}";
        let mut parser = Parser::new(input.to_string());
        let query = parser.parse().unwrap();
        
        assert_eq!(query.query_type, QueryType::InstantVector { labels: BTreeMap::from([("host".to_string(), "server1".to_string())]) });
    }
    
    #[test]
    fn test_parser_range_vector() {
        let input = "cpu_usage[30s]";
        let mut parser = Parser::new(input.to_string());
        let query = parser.parse().unwrap();
        
        assert_eq!(query.query_type, QueryType::RangeVector { 
            labels: BTreeMap::new(), 
            duration: "30s".to_string() 
        });
    }
    
    #[test]
    fn test_parser_with_labels() {
        let input = "cpu_usage{host=\"server1\",instance=\"web1\"}";
        let mut parser = Parser::new(input.to_string());
        let query = parser.parse().unwrap();
        
        assert!(query.labels.contains_key("host"));
        assert!(query.labels.contains_key("instance"));
    }
    
    #[test]
    fn test_scheduler_plan() {
        let input = "cpu_usage{host=\"server1\"}[30s]";
        let mut parser = Parser::new(input.to_string());
        let query = parser.parse().unwrap();
        
        let scheduler = Scheduler::new(query);
        let plan = scheduler.plan();
        
        assert!(plan.steps.contains(&"filter_labels".to_string()));
        assert!(plan.steps.contains(&"filter_time_range".to_string()));
        assert!(plan.steps.contains(&"fetch_from_storage".to_string()));
    }
    
    #[test]
    fn test_execute_instant_vector() {
        let data = vec![
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
                value: 50.0,
                timestamp: 1700000000,
            },
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server2".to_string())]),
                value: 75.0,
                timestamp: 1700000001,
            },
        ];

        let result = execute_query("cpu_usage{host=~\"server.*\"}", &data).unwrap();
        
        assert_eq!(result.result_type, "vector");
        assert!(!result.result.is_empty());
    }
    
    #[test]
    fn test_execute_range_vector() {
        let data = vec![
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
                value: 50.0,
                timestamp: 1700000000,
            },
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
                value: 75.0,
                timestamp: 1700000060,
            },
        ];
        
        let result = execute_query("rate(cpu_usage[30s])", &data).unwrap();
        
        assert_eq!(result.result_type, "vector");
        // rate over 60s with 30s step should give approximate average
        assert!(!result.result.is_empty());
    }
    
    #[test]
    fn test_parse_error_handling() {
        let input = "invalid {{ syntax";
        let mut parser = Parser::new(input.to_string());
        let result = parser.parse();
        
        assert!(result.is_err());
    }
    
    #[test]
    fn test_query_stats_tracking() {
        let data = vec![
            MetricSample {
                labels: BTreeMap::new(),
                value: 1.0,
                timestamp: 1700000000,
            },
        ];
        
        let result = execute_query("1", &data).unwrap();
        
        assert!(result.stats.rows_read > 0);
        assert!(result.execution_time_ns > 0);
    }
}

// ---------------------------------------------------------------------------
// Main entry point
// ---------------------------------------------------------------------------

fn main() {
    // Example usage
    let query_expr = "rate(cpu_usage[30s])";
    let data = vec![
        MetricSample {
            labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
            value: 50.0,
            timestamp: 1700000000,
        },
        MetricSample {
            labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
            value: 75.0,
            timestamp: 1700000060,
        },
        MetricSample {
            labels: BTreeMap::from([("host".to_string(), "server2".to_string())]),
            value: 30.0,
            timestamp: 1700000000,
        },
    ];
    
    match execute_query(query_expr, &data) {
        Ok(result) => {
            println!("Query '{}' executed successfully", query_expr);
            println!("Result type: {}", result.result_type);
            println!("Result count: {} samples", result.result.len());
            println!("Execution time: {} ms", result.execution_time_ms);
        }
        Err(e) => {
            println!("Query error: {}", e);
        }
    }
}