# Helix Benchmark Suite

This directory contains performance benchmarks for the Helix observability platform, 
critical for verifying the "1 billion events in < 1 second" goal.

## Benchmarks Included

### 1. helix_query
- `rate_30s` - rate() function over 30s window
- `top_k_100` - top-k selection from rate results
- `histogram_quantile` - histogram percentile computation
- `filter_by_host` - label-based filtering
- `filter_by_region` - regional filtering

### 2. storage_insert
- Insert performance at various scales (100, 1K, 10K, 100K)
- Throughput measurement in elements/second

### 3. query_filtering
- Host pattern matching
- Region-based filtering

### 4. large_scale
- 500K sample queries
- Top-k on large datasets

### 4. performance_comparison
- Helix baseline establishment
- Structured for ClickHouse comparison

## Running Benchmarks

```bash
# Run all benchmarks
criterion

# Run specific benchmark
criterion --bench helix_rate_30s

# Run with HTML output
criterion --html report

# Compare before/after
criterion --incremental 0.1 --compare
```

## Performance Goals

| Benchmark | Target |
|-----------|--------|
| rate(cpu_usage[30s]) over 1M samples | < 1 second |
| top_k(100, rate(...)) | < 2 seconds |
| 1B events scanned/filtered | < 1 second (with proper indexing) |
| Insert 100K data points | < 5 seconds |

## Comparison with ClickHouse

This suite establishes Helix performance baselines. ClickHouse benchmarks 
are run separately and results are compared using:

```bash
criterion --incremental 0.1 --compare baseline.json
```

The critical metric is maintaining or exceeding ClickHouse's 
"1B events in < 1s" performance while providing custom TSDB flexibility.

## Adding New Benchmarks

1. Add a function starting with `benchmark_` 
2. The function should take `&mut Criterion` as argument
3. Use `criterion.group()` or `criterion.section()` for organization
4. Use `black_box()` to prevent dead code elimination
5. Document the benchmark purpose and expected targets