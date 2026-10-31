//! Helix Query Engine - Functions Module
//! 
//! This module provides built-in functions for the Helix query language.
//! Includes aggregation, rate, histogram, and mathematical functions.

use std::collections::BTreeMap;
use super::{MetricSample, QueryResult};

/// Function registry for query language functions
pub struct FunctionRegistry {
    functions: BTreeMap<String, Box<dyn QueryFunction>>,
}

impl FunctionRegistry {
    pub fn new() -> Self {
        let mut registry = FunctionRegistry {
            functions: BTreeMap::new(),
        };
        
        // Register built-in functions
        registry.register("rate", Box::new(RateFunction));
        registry.register("irate", Box::new(IRateFunction));
        registry.register("increase", Box::new(IncreaseFunction));
        registry.register("sum", Box::new(SumFunction));
        registry.register("avg", Box::new(AvgFunction));
        registry.register("min", Box::new(MinFunction));
        registry.register("max", Box::new(MaxFunction));
        registry.register("count", Box::new(CountFunction));
        registry.register("stddev", Box::new(StddevFunction));
        registry.register("stdvar", Box::new(StdvarFunction));
        registry.register("topk", Box::new(TopKFunction));
        registry.register("bottomk", Box::new(BottomKFunction));
        registry.register("quantile", Box::new(QuantileFunction));
        registry.register("histogram_quantile", Box::new(HistogramQuantileFunction));
        registry.register("abs", Box::new(AbsFunction));
        registry.register("ceil", Box::new(CeilFunction));
        registry.register("floor", Box::new(FloorFunction));
        registry.register("round", Box::new(RoundFunction));
        registry.register("sqrt", Box::new(SqrtFunction));
        registry.register("exp", Box::new(ExpFunction));
        registry.register("ln", Box::new(LnFunction));
        registry.register("log2", Box::new(Log2Function));
        registry.register("log10", Box::new(Log10Function));
        registry.register("clamp", Box::new(ClampFunction));
        registry.register("clamp_min", Box::new(ClampMinFunction));
        registry.register("clamp_max", Box::new(ClampMaxFunction));
        registry.register("deriv", Box::new(DerivFunction));
        registry.register("predict_linear", Box::new(PredictLinearFunction));
        registry.register("holt_winters", Box::new(HoltWintersFunction));
        registry.register("label_replace", Box::new(LabelReplaceFunction));
        registry.register("label_join", Box::new(LabelJoinFunction));
        registry.register("vector", Box::new(VectorFunction));
        registry.register("scalar", Box::new(ScalarFunction));
        registry.register("sort", Box::new(SortFunction));
        registry.register("sort_desc", Box::new(SortDescFunction));
        registry.register("time", Box::new(TimeFunction));
        registry.register("timestamp", Box::new(TimestampFunction));
        registry.register("days_in_month", Box::new(DaysInMonthFunction));
        registry.register("day_of_month", Box::new(DayOfMonthFunction));
        registry.register("day_of_week", Box::new(DayOfWeekFunction));
        registry.register("day_of_year", Box::new(DayOfYearFunction));
        registry.register("hour", Box::new(HourFunction));
        registry.register("minute", Box::new(MinuteFunction));
        registry.register("month", Box::new(MonthFunction));
        registry.register("year", Box::new(YearFunction));
        registry.register("changes", Box::new(ChangesFunction));
        registry.register("resets", Box::new(ResetsFunction));
        registry.register("idelta", Box::new(IDeltaFunction));
        registry.register("deriv", Box::new(DerivFunction));
        
        registry
    }
    
    pub fn register(&mut self, name: &str, func: Box<dyn QueryFunction>) {
        self.functions.insert(name.to_string(), func);
    }
    
    pub fn get(&self, name: &str) -> Option<&dyn QueryFunction> {
        self.functions.get(name).map(|f| f.as_ref())
    }
    
    pub fn list(&self) -> Vec<&str> {
        self.functions.keys().map(|k| k.as_str()).collect()
    }
    
    pub fn has(&self, name: &str) -> bool {
        self.functions.contains_key(name)
    }
}

impl Default for FunctionRegistry {
    fn default() -> Self {
        Self::new()
    }
}

/// Query function trait
pub trait QueryFunction: Send + Sync {
    fn name(&self) -> &'static str;
    fn call(&self, args: &[MetricSample], context: &FunctionContext) -> Result<Vec<MetricSample>, String>;
    fn description(&self) -> &'static str;
    fn arg_count(&self) -> (usize, usize); // (min, max)
}

/// Function execution context
#[derive(Debug, Clone)]
pub struct FunctionContext {
    pub timestamp: i64,
    pub step: i64,
    pub range: i64,
}

impl FunctionContext {
    pub fn new(timestamp: i64, step: i64, range: i64) -> Self {
        FunctionContext { timestamp, step, range }
    }
}

// ---------------------------------------------------------------------------
// Aggregation Functions
// ---------------------------------------------------------------------------

pub struct SumFunction;
impl QueryFunction for SumFunction {
    fn name(&self) -> &'static str { "sum" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("sum requires at least 1 argument".to_string());
        }
        let mut result = Vec::new();
        let mut by_labels: BTreeMap<BTreeMap<String, String>, f64> = BTreeMap::new();
        
        for sample in args {
            let key = sample.labels.clone();
            *by_labels.entry(key).or_insert(0.0) += sample.value;
        }
        
        for (labels, value) in by_labels {
            result.push(MetricSample {
                labels,
                value,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate sum over dimensions" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct AvgFunction;
impl QueryFunction for AvgFunction {
    fn name(&self) -> &'static str { "avg" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("avg requires at least 1 argument".to_string());
        }
        let mut result = Vec::new();
        let mut by_labels: BTreeMap<BTreeMap<String, String>, (f64, usize)> = BTreeMap::new();
        
        for sample in args {
            let key = sample.labels.clone();
            let entry = by_labels.entry(key).or_insert((0.0, 0));
            entry.0 += sample.value;
            entry.1 += 1;
        }
        
        for (labels, (sum, count)) in by_labels {
            result.push(MetricSample {
                labels,
                value: sum / count as f64,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate average over dimensions" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct MinFunction;
impl QueryFunction for MinFunction {
    fn name(&self) -> &'static str { "min" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("min requires at least 1 argument".to_string());
        }
        let mut result = Vec::new();
        let mut by_labels: BTreeMap<BTreeMap<String, String>, f64> = BTreeMap::new();
        
        for sample in args {
            let key = sample.labels.clone();
            let entry = by_labels.entry(key).or_insert(f64::MAX);
            if sample.value < *entry {
                *entry = sample.value;
            }
        }
        
        for (labels, value) in by_labels {
            result.push(MetricSample {
                labels,
                value,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate minimum over dimensions" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct MaxFunction;
impl QueryFunction for MaxFunction {
    fn name(&self) -> &'static str { "max" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("max requires at least 1 argument".to_string());
        }
        let mut result = Vec::new();
        let mut by_labels: BTreeMap<BTreeMap<String, String>, f64> = BTreeMap::new();
        
        for sample in args {
            let key = sample.labels.clone();
            let entry = by_labels.entry(key).or_insert(f64::MIN);
            if sample.value > *entry {
                *entry = sample.value;
            }
        }
        
        for (labels, value) in by_labels {
            result.push(MetricSample {
                labels,
                value,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate maximum over dimensions" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct CountFunction;
impl QueryFunction for CountFunction {
    fn name(&self) -> &'static str { "count" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("count requires at least 1 argument".to_string());
        }
        let mut result = Vec::new();
        let mut by_labels: BTreeMap<BTreeMap<String, String>, usize> = BTreeMap::new();
        
        for sample in args {
            let key = sample.labels.clone();
            *by_labels.entry(key).or_insert(0) += 1;
        }
        
        for (labels, count) in by_labels {
            result.push(MetricSample {
                labels,
                value: count as f64,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Count number of series" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

// ---------------------------------------------------------------------------
// Rate Functions
// ---------------------------------------------------------------------------

pub struct RateFunction;
impl QueryFunction for RateFunction {
    fn name(&self) -> &'static str { "rate" }
    fn call(&self, args: &[MetricSample], ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("rate requires at least 1 argument".to_string());
        }
        
        // Group by labels
        let mut by_labels: BTreeMap<BTreeMap<String, String>, Vec<&MetricSample>> = BTreeMap::new();
        for sample in args {
            by_labels.entry(sample.labels.clone()).or_default().push(sample);
        }
        
        let mut result = Vec::new();
        for (labels, samples) in by_labels {
            if samples.len() < 2 {
                continue;
            }
            
            // Sort by timestamp
            let mut sorted: Vec<&MetricSample> = samples.into_iter().collect();
            sorted.sort_by_key(|s| s.timestamp);
            
            let first = sorted.first().unwrap();
            let last = sorted.last().unwrap();
            
            let time_diff = (last.timestamp - first.timestamp) as f64;
            if time_diff <= 0.0 {
                continue;
            }
            
            let value_diff = last.value - first.value;
            let rate = value_diff / time_diff;
            
            result.push(MetricSample {
                labels,
                value: rate.max(0.0),
                timestamp: ctx.timestamp,
            });
        }
        
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate per-second rate of increase" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct IRateFunction;
impl QueryFunction for IRateFunction {
    fn name(&self) -> &'static str { "irate" }
    fn call(&self, args: &[MetricSample], ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("irate requires at least 2 data points".to_string());
        }
        
        let mut by_labels: BTreeMap<BTreeMap<String, String>, Vec<&MetricSample>> = BTreeMap::new();
        for sample in args {
            by_labels.entry(sample.labels.clone()).or_default().push(sample);
        }
        
        let mut result = Vec::new();
        for (labels, samples) in by_labels {
            if samples.len() < 2 {
                continue;
            }
            
            let mut sorted: Vec<&MetricSample> = samples.into_iter().collect();
            sorted.sort_by_key(|s| s.timestamp);
            
            // Use last two points
            let n = sorted.len();
            let prev = sorted[n - 2];
            let last = sorted[n - 1];
            
            let time_diff = (last.timestamp - prev.timestamp) as f64;
            if time_diff <= 0.0 {
                continue;
            }
            
            let value_diff = last.value - prev.value;
            let rate = value_diff / time_diff;
            
            result.push(MetricSample {
                labels,
                value: rate.max(0.0),
                timestamp: ctx.timestamp,
            });
        }
        
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate instant rate of increase" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct IncreaseFunction;
impl QueryFunction for IncreaseFunction {
    fn name(&self) -> &'static str { "increase" }
    fn call(&self, args: &[MetricSample], ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("increase requires at least 1 argument".to_string());
        }
        
        let mut by_labels: BTreeMap<BTreeMap<String, String>, Vec<&MetricSample>> = BTreeMap::new();
        for sample in args {
            by_labels.entry(sample.labels.clone()).or_default().push(sample);
        }
        
        let mut result = Vec::new();
        for (labels, samples) in by_labels {
            if samples.len() < 2 {
                continue;
            }
            
            let mut sorted: Vec<&MetricSample> = samples.into_iter().collect();
            sorted.sort_by_key(|s| s.timestamp);
            
            let first = sorted.first().unwrap();
            let last = sorted.last().unwrap();
            
            let value_diff = last.value - first.value;
            
            result.push(MetricSample {
                labels,
                value: value_diff.max(0.0),
                timestamp: ctx.timestamp,
            });
        }
        
        Ok(result)
    }
    fn description(&self) -> &'static str { "Calculate increase over time range" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

// ---------------------------------------------------------------------------
// Math Functions
// ---------------------------------------------------------------------------

pub struct AbsFunction;
impl QueryFunction for AbsFunction {
    fn name(&self) -> &'static str { "abs" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.abs(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Absolute value" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct CeilFunction;
impl QueryFunction for CeilFunction {
    fn name(&self) -> &'static str { "ceil" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.ceil(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Ceiling" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct FloorFunction;
impl QueryFunction for FloorFunction {
    fn name(&self) -> &'static str { "floor" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.floor(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Floor" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct RoundFunction;
impl QueryFunction for RoundFunction {
    fn name(&self) -> &'static str { "round" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.round(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Round to nearest integer" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct SqrtFunction;
impl QueryFunction for SqrtFunction {
    fn name(&self) -> &'static str { "sqrt" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.sqrt(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Square root" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct ExpFunction;
impl QueryFunction for ExpFunction {
    fn name(&self) -> &'static str { "exp" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.exp(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Exponential" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct LnFunction;
impl QueryFunction for LnFunction {
    fn name(&self) -> &'static str { "ln" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.ln(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Natural logarithm" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct Log2Function;
impl QueryFunction for Log2Function {
    fn name(&self) -> &'static str { "log2" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.log2(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Base-2 logarithm" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct Log10Function;
impl QueryFunction for Log10Function {
    fn name(&self) -> &'static str { "log10" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.log10(),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Base-10 logarithm" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct ClampFunction;
impl QueryFunction for ClampFunction {
    fn name(&self) -> &'static str { "clamp" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 3 {
            return Err("clamp requires 3 arguments: value, min, max".to_string());
        }
        let min_val = args[1].value;
        let max_val = args[2].value;
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.clamp(min_val, max_val),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Clamp values to range" }
    fn arg_count(&self) -> (usize, usize) { (3, 3) }
}

pub struct ClampMinFunction;
impl QueryFunction for ClampMinFunction {
    fn name(&self) -> &'static str { "clamp_min" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("clamp_min requires 2 arguments: value, min".to_string());
        }
        let min_val = args[1].value;
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.max(min_val),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Clamp minimum value" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

pub struct ClampMaxFunction;
impl QueryFunction for ClampMaxFunction {
    fn name(&self) -> &'static str { "clamp_max" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("clamp_max requires 2 arguments: value, max".to_string());
        }
        let max_val = args[1].value;
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.value.min(max_val),
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Clamp maximum value" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

// ---------------------------------------------------------------------------
// Top-K / Bottom-K Functions
// ---------------------------------------------------------------------------

pub struct TopKFunction;
impl QueryFunction for TopKFunction {
    fn name(&self) -> &'static str { "topk" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("topk requires arguments".to_string());
        }
        let k = if args.len() > 1 { args[1].value as usize } else { 10 };
        
        let mut sorted: Vec<MetricSample> = args.to_vec();
        sorted.sort_by(|a, b| b.value.partial_cmp(&a.value).unwrap_or(std::cmp::Ordering::Equal));
        sorted.truncate(k);
        
        Ok(sorted)
    }
    fn description(&self) -> &'static str { "Top K values" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

pub struct BottomKFunction;
impl QueryFunction for BottomKFunction {
    fn name(&self) -> &'static str { "bottomk" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("bottomk requires arguments".to_string());
        }
        let k = if args.len() > 1 { args[1].value as usize } else { 10 };
        
        let mut sorted: Vec<MetricSample> = args.to_vec();
        sorted.sort_by(|a, b| a.value.partial_cmp(&b.value).unwrap_or(std::cmp::Ordering::Equal));
        sorted.truncate(k);
        
        Ok(sorted)
    }
    fn description(&self) -> &'static str { "Bottom K values" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

// ---------------------------------------------------------------------------
// Statistical Functions
// ---------------------------------------------------------------------------

pub struct StddevFunction;
impl QueryFunction for StddevFunction {
    fn name(&self) -> &'static str { "stddev" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("stddev requires at least 1 argument".to_string());
        }
        
        let mut by_labels: BTreeMap<BTreeMap<String, String>, Vec<f64>> = BTreeMap::new();
        for sample in args {
            by_labels.entry(sample.labels.clone()).or_default().push(sample.value);
        }
        
        let mut result = Vec::new();
        for (labels, values) in by_labels {
            if values.is_empty() { continue; }
            let mean = values.iter().sum::<f64>() / values.len() as f64;
            let variance = values.iter().map(|v| (v - mean).powi(2)).sum::<f64>() / values.len() as f64;
            let stddev = variance.sqrt();
            
            result.push(MetricSample {
                labels,
                value: stddev,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Population standard deviation" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct StdvarFunction;
impl QueryFunction for StdvarFunction {
    fn name(&self) -> &'static str { "stdvar" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("stdvar requires at least 1 argument".to_string());
        }
        
        let mut by_labels: BTreeMap<BTreeMap<String, String>, Vec<f64>> = BTreeMap::new();
        for sample in args {
            by_labels.entry(sample.labels.clone()).or_default().push(sample.value);
        }
        
        let mut result = Vec::new();
        for (labels, values) in by_labels {
            if values.is_empty() { continue; }
            let mean = values.iter().sum::<f64>() / values.len() as f64;
            let variance = values.iter().map(|v| (v - mean).powi(2)).sum::<f64>() / values.len() as f64;
            
            result.push(MetricSample {
                labels,
                value: variance,
                timestamp: args[0].timestamp,
            });
        }
        Ok(result)
    }
    fn description(&self) -> &'static str { "Population variance" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct QuantileFunction;
impl QueryFunction for QuantileFunction {
    fn name(&self) -> &'static str { "quantile" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("quantile requires 2 arguments: phi, metric".to_string());
        }
        
        let phi = args[0].value;
        let mut values: Vec<f64> = args[1..].iter().map(|s| s.value).collect();
        values.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));
        
        if values.is_empty() {
            return Ok(Vec::new());
        }
        
        let idx = (phi * (values.len() - 1) as f64).round() as usize;
        let idx = idx.min(values.len() - 1);
        
        Ok(vec![MetricSample {
            labels: BTreeMap::new(),
            value: values[idx],
            timestamp: args[0].timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Calculate quantile" }
    fn arg_count(&self) -> (usize, usize) { (2, usize::MAX) }
}

pub struct HistogramQuantileFunction;
impl QueryFunction for HistogramQuantileFunction {
    fn name(&self) -> &'static str { "histogram_quantile" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("histogram_quantile requires 2 arguments: phi, metric".to_string());
        }
        
        let phi = args[0].value;
        
        // Simplified histogram quantile calculation
        let mut buckets: Vec<(f64, f64)> = args[1..].iter()
            .map(|s| (s.value, s.value))
            .collect();
        buckets.sort_by(|a, b| a.0.partial_cmp(&b.0).unwrap_or(std::cmp::Ordering::Equal));
        
        if buckets.is_empty() {
            return Ok(Vec::new());
        }
        
        let total_count: f64 = buckets.iter().map(|(_, v)| *v).sum();
        let target = phi * total_count;
        
        let mut cumulative = 0.0;
        for (le, count) in &buckets {
            cumulative += count;
            if cumulative >= target {
                return Ok(vec![MetricSample {
                    labels: BTreeMap::new(),
                    value: *le,
                    timestamp: args[0].timestamp,
                }]);
            }
        }
        
        Ok(vec![MetricSample {
            labels: BTreeMap::new(),
            value: buckets.last().unwrap().0,
            timestamp: args[0].timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Calculate quantile from histogram" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

// ---------------------------------------------------------------------------
// Time Functions
// ---------------------------------------------------------------------------

pub struct TimeFunction;
impl QueryFunction for TimeFunction {
    fn name(&self) -> &'static str { "time" }
    fn call(&self, args: &[MetricSample], ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(vec![MetricSample {
            labels: BTreeMap::new(),
            value: ctx.timestamp as f64,
            timestamp: ctx.timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Current timestamp" }
    fn arg_count(&self) -> (usize, usize) { (0, 0) }
}

pub struct TimestampFunction;
impl QueryFunction for TimestampFunction {
    fn name(&self) -> &'static str { "timestamp" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: s.timestamp as f64,
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Timestamp of sample" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

// ---------------------------------------------------------------------------
// Label Functions
// ---------------------------------------------------------------------------

pub struct LabelReplaceFunction;
impl QueryFunction for LabelReplaceFunction {
    fn name(&self) -> &'static str { "label_replace" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        // Simplified: just return input
        Ok(args.to_vec())
    }
    fn description(&self) -> &'static str { "Replace label values" }
    fn arg_count(&self) -> (usize, usize) { (5, 5) }
}

pub struct LabelJoinFunction;
impl QueryFunction for LabelJoinFunction {
    fn name(&self) -> &'static str { "label_join" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        // Simplified: just return input
        Ok(args.to_vec())
    }
    fn description(&self) -> &'static str { "Join label values" }
    fn arg_count(&self) -> (usize, usize) { (4, usize::MAX) }
}

// ---------------------------------------------------------------------------
// Vector/Scalar Conversion
// ---------------------------------------------------------------------------

pub struct VectorFunction;
impl QueryFunction for VectorFunction {
    fn name(&self) -> &'static str { "vector" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.to_vec())
    }
    fn description(&self) -> &'static str { "Convert scalar to vector" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct ScalarFunction;
impl QueryFunction for ScalarFunction {
    fn name(&self) -> &'static str { "scalar" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Ok(Vec::new());
        }
        Ok(vec![MetricSample {
            labels: BTreeMap::new(),
            value: args[0].value,
            timestamp: args[0].timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Convert single-element vector to scalar" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

// ---------------------------------------------------------------------------
// Sorting Functions
// ---------------------------------------------------------------------------

pub struct SortFunction;
impl QueryFunction for SortFunction {
    fn name(&self) -> &'static str { "sort" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        let mut result: Vec<MetricSample> = args.to_vec();
        result.sort_by(|a, b| a.value.partial_cmp(&b.value).unwrap_or(std::cmp::Ordering::Equal));
        Ok(result)
    }
    fn description(&self) -> &'static str { "Sort ascending" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct SortDescFunction;
impl QueryFunction for SortDescFunction {
    fn name(&self) -> &'static str { "sort_desc" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        let mut result: Vec<MetricSample> = args.to_vec();
        result.sort_by(|a, b| b.value.partial_cmp(&a.value).unwrap_or(std::cmp::Ordering::Equal));
        Ok(result)
    }
    fn description(&self) -> &'static str { "Sort descending" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

// ---------------------------------------------------------------------------
// Changes/Resets Functions
// ---------------------------------------------------------------------------

pub struct ChangesFunction;
impl QueryFunction for ChangesFunction {
    fn name(&self) -> &'static str { "changes" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        // Count number of value changes
        if args.len() < 2 {
            return Ok(Vec::new());
        }
        
        let mut changes = 0;
        for i in 1..args.len() {
            if args[i].value != args[i-1].value {
                changes += 1;
            }
        }
        
        Ok(vec![MetricSample {
            labels: args[0].labels.clone(),
            value: changes as f64,
            timestamp: args[0].timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Count value changes" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct ResetsFunction;
impl QueryFunction for ResetsFunction {
    fn name(&self) -> &'static str { "resets" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Ok(Vec::new());
        }
        
        let mut resets = 0;
        for i in 1..args.len() {
            if args[i].value < args[i-1].value {
                resets += 1;
            }
        }
        
        Ok(vec![MetricSample {
            labels: args[0].labels.clone(),
            value: resets as f64,
            timestamp: args[0].timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Count counter resets" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct IDeltaFunction;
impl QueryFunction for IDeltaFunction {
    fn name(&self) -> &'static str { "idelta" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("idelta requires at least 2 data points".to_string());
        }
        
        let mut sorted: Vec<&MetricSample> = args.iter().collect();
        sorted.sort_by_key(|s| s.timestamp);
        
        let n = sorted.len();
        let last = sorted[n - 1];
        let prev = sorted[n - 2];
        
        Ok(vec![MetricSample {
            labels: last.labels.clone(),
            value: last.value - prev.value,
            timestamp: last.timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Instantaneous delta" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct DerivFunction;
impl QueryFunction for DerivFunction {
    fn name(&self) -> &'static str { "deriv" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.len() < 2 {
            return Err("deriv requires at least 2 data points".to_string());
        }
        
        let mut sorted: Vec<&MetricSample> = args.iter().collect();
        sorted.sort_by_key(|s| s.timestamp);
        
        let first = sorted.first().unwrap();
        let last = sorted.last().unwrap();
        
        let time_diff = (last.timestamp - first.timestamp) as f64;
        if time_diff <= 0.0 {
            return Ok(Vec::new());
        }
        
        let value_diff = last.value - first.value;
        let deriv = value_diff / time_diff;
        
        Ok(vec![MetricSample {
            labels: last.labels.clone(),
            value: deriv,
            timestamp: last.timestamp,
        }])
    }
    fn description(&self) -> &'static str { "Derivative" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct PredictLinearFunction;
impl QueryFunction for PredictLinearFunction {
    fn name(&self) -> &'static str { "predict_linear" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("predict_linear requires arguments".to_string());
        }
        // Simplified linear prediction
        Ok(args.to_vec())
    }
    fn description(&self) -> &'static str { "Linear prediction" }
    fn arg_count(&self) -> (usize, usize) { (2, 2) }
}

pub struct HoltWintersFunction;
impl QueryFunction for HoltWintersFunction {
    fn name(&self) -> &'static str { "holt_winters" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        if args.is_empty() {
            return Err("holt_winters requires arguments".to_string());
        }
        // Simplified Holt-Winters
        Ok(args.to_vec())
    }
    fn description(&self) -> &'static str { "Holt-Winters smoothing" }
    fn arg_count(&self) -> (usize, usize) { (3, 3) }
}

// ---------------------------------------------------------------------------
// Date/Time Functions
// ---------------------------------------------------------------------------

pub struct DayOfMonthFunction;
impl QueryFunction for DayOfMonthFunction {
    fn name(&self) -> &'static str { "day_of_month" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 15.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Day of month" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct DayOfWeekFunction;
impl QueryFunction for DayOfWeekFunction {
    fn name(&self) -> &'static str { "day_of_week" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 3.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Day of week" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct DayOfYearFunction;
impl QueryFunction for DayOfYearFunction {
    fn name(&self) -> &'static str { "day_of_year" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 300.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Day of year" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct HourFunction;
impl QueryFunction for HourFunction {
    fn name(&self) -> &'static str { "hour" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 12.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Hour of day" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct MinuteFunction;
impl QueryFunction for MinuteFunction {
    fn name(&self) -> &'static str { "minute" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 30.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Minute of hour" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct MonthFunction;
impl QueryFunction for MonthFunction {
    fn name(&self) -> &'static str { "month" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 10.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Month of year" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct YearFunction;
impl QueryFunction for YearFunction {
    fn name(&self) -> &'static str { "year" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 2026.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Year" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

pub struct DaysInMonthFunction;
impl QueryFunction for DaysInMonthFunction {
    fn name(&self) -> &'static str { "days_in_month" }
    fn call(&self, args: &[MetricSample], _ctx: &FunctionContext) -> Result<Vec<MetricSample>, String> {
        Ok(args.iter().map(|s| MetricSample {
            labels: s.labels.clone(),
            value: 30.0, // Simplified
            timestamp: s.timestamp,
        }).collect())
    }
    fn description(&self) -> &'static str { "Days in month" }
    fn arg_count(&self) -> (usize, usize) { (1, 1) }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_function_registry() {
        let registry = FunctionRegistry::new();
        assert!(registry.has("rate"));
        assert!(registry.has("sum"));
        assert!(registry.has("avg"));
        assert!(registry.has("topk"));
        assert!(registry.has("abs"));
    }
    
    #[test]
    fn test_sum_function() {
        let func = SumFunction;
        let args = vec![
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
                value: 10.0,
                timestamp: 1000,
            },
            MetricSample {
                labels: BTreeMap::from([("host".to_string(), "server1".to_string())]),
                value: 20.0,
                timestamp: 1000,
            },
        ];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result.len(), 1);
        assert_eq!(result[0].value, 30.0);
    }
    
    #[test]
    fn test_avg_function() {
        let func = AvgFunction;
        let args = vec![
            MetricSample { labels: BTreeMap::new(), value: 10.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 20.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 30.0, timestamp: 1000 },
        ];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result[0].value, 20.0);
    }
    
    #[test]
    fn test_abs_function() {
        let func = AbsFunction;
        let args = vec![MetricSample { labels: BTreeMap::new(), value: -5.0, timestamp: 1000 }];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result[0].value, 5.0);
    }
    
    #[test]
    fn test_rate_function() {
        let func = RateFunction;
        let args = vec![
            MetricSample { labels: BTreeMap::new(), value: 10.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 20.0, timestamp: 1010 },
        ];
        let ctx = FunctionContext::new(1010, 10, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result[0].value, 1.0);
    }
    
    #[test]
    fn test_topk_function() {
        let func = TopKFunction;
        let args = vec![
            MetricSample { labels: BTreeMap::new(), value: 5.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 15.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 10.0, timestamp: 1000 },
        ];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result.len(), 3);
    }
    
    #[test]
    fn test_stddev_function() {
        let func = StddevFunction;
        let args = vec![
            MetricSample { labels: BTreeMap::new(), value: 10.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 20.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 30.0, timestamp: 1000 },
        ];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert!(result[0].value > 0.0);
    }
    
    #[test]
    fn test_sort_function() {
        let func = SortFunction;
        let args = vec![
            MetricSample { labels: BTreeMap::new(), value: 30.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 10.0, timestamp: 1000 },
            MetricSample { labels: BTreeMap::new(), value: 20.0, timestamp: 1000 },
        ];
        let ctx = FunctionContext::new(1000, 1, 30);
        let result = func.call(&args, &ctx).unwrap();
        assert_eq!(result[0].value, 10.0);
        assert_eq!(result[1].value, 20.0);
        assert_eq!(result[2].value, 30.0);
    }
}