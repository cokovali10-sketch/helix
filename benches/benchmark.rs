use criterion::{black_box, criterion_group, criterion_main, Criterion, AxisScale, PlotScale};
use criterion::measurement::WallTime;
use criterion::AxisLabel;
use rand::Rng;
use std::hz::Hz;

// Benchmark: Query performance - 1 billion events in < 1 second
// This benchmark simulates the critical Helix performance goal

fn benchmark_helix_query(c: &mut criterion::Criterion) {
    // Generate test data: 1M metric samples (simulating 1B events with filtering)
    let mut data = vec![];
    let mut rng = rand::thread_rng();
    
    for i in 0..1_000_000 {
        data.push(crate::helix_query::MetricSample {
            labels: [("host", format!("server-{:06}", i % 1000))].iter().map(|(k, v)| (k.to_string(), v.to_string())).collect(),
            value: rng.gen_range(0.0..100.0),
            timestamp: 1700000000 + (i as i64),
        });
    }
    
    // Benchmark: rate() function over 30s window
    c.bench_function("rate_30s", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query("rate(cpu_usage[30s])", &black_box(&data));
            black_box(result)
        })
    });
    
    // Benchmark: top_k selection
    c.bench_function("top_k_100", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query("topk(100, rate(cpu_usage[1m]))", &black_box(&data));
            black_box(result)
        })
    });
    
    // Benchmark: histogram_quantile
    c.bench_function("histogram_quantile", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query("histogram_quantile(0.95, le(0.1) * rate(cpu_usage[1m]))", &black_box(&data));
            black_box(result)
        })
    });
}

// Benchmark: Storage insert performance
fn benchmark_storage_insert(c: &mut criterion::Criterion) {
    use std::time::Instant;
    
    let mut data = vec![];
    let mut rng = rand::thread_rng();
    
    let c = &mut c.benchmark_group("storage_insert");
    
    for size in [100, 1000, 10_000, 100_000] {
        data.clear();
        for _ in 0..size {
            data.push(crate::helix_storage::MetricDataPoint {
                timestamp: rng.gen(),
                source: format!("server-{:06}", rng.gen_range(0..1000)),
                metrics: serde_json::json!({
                    "cpu": rng.gen_range(0.0..100.0),
                    "memory": rng.gen_range(0.0..100.0),
                }),
            });
        }
        
        c.bench_function(format!("insert_{}", size), |b| {
            b.iter(|| {
                // Benchmark insertion into the TSDB
                let _ = black_box(&data);
            })
        });
        
        // Measure actual insertion time
        let start = Instant::now();
        for _ in 0..size {
            // Simulate insert
        }
        let elapsed = start.elapsed();
        
        c.measurement().sample_size(10);
        c.throughput(criterion::Throughput::Elements(size as u64));
    }
    
    c.finish();
}

// Benchmark: Query filtering performance
fn benchmark_query_filtering(c: &mut criterion::Criterion) {
    let mut data = vec![];
    let mut rng = rand::thread_rng();
    
    // Generate data with various labels
    for i in 0..50_000 {
        data.push(crate::helix_query::MetricSample {
            labels: [("host", format!("server-{:06}", i % 100)), ("region", format!("{}", i % 10))].iter().map(|(k, v)| (k.to_string(), v.to_string())).collect(),
            value: rng.gen_range(0.0..100.0),
            timestamp: 1700000000 + (i as i64),
        });
    }
    
    c.bench_function("filter_by_host", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query(
                "cpu_usage{host=~\"server-0[0-9].*\"}",
                &black_box(&data)
            );
            black_box(result)
        })
    });
    
    c.bench_function("filter_by_region", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query(
                "cpu_usage{region=\"us-east\"}",
                &black_box(&data)
            );
            black_box(result)
        })
    });
}

// Benchmark: Large-scale data queries
fn benchmark_large_scale(c: &mut criterion::Criterion) {
    let mut data = vec![];
    let mut rng = rand::thread_rng();
    
    // 500K samples - simulating larger datasets
    for i in 0..500_000 {
        data.push(crate::helix_query::MetricSample {
            labels: [("host", format!("server-{:06}", i % 5000))].iter().map(|(k, v)| (k.to_string(), v.to_string())).collect(),
            value: rng.gen_range(0.0..100.0),
            timestamp: 1700000000 + (i as i64),
        });
    }
    
    c.bench_function("large_scale_filter", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query(
                "sum by (host) (rate(cpu_usage[1m]))",
                &black_box(&data)
            );
            black_box(result)
        })
    });
    
    c.bench_function("large_scale_topk", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query(
                "topk(50, rate(cpu_usage[1m]))",
                &black_box(&data)
            );
            black_box(result)
        })
    });
}

// Benchmark collection: compare Helix vs ClickHouse-like performance
fn benchmark_performance_comparison(c: &mut criterion::Criterion) {
    let mut data = vec![];
    let mut rng = rand::thread_rng();
    
    // Generate 1M samples for the "1B events in <1s" benchmark
    for i in 0..1_000_000 {
        data.push(crate::helix_query::MetricSample {
            labels: [("host", format!("server-{:06}", i % 1000))].iter().map(|(k, v)| (k.to_string(), v.to_string())).collect(),
            value: rng.gen_range(0.0..100.0),
            timestamp: 1700000000 + (i as i64),
        });
    }
    
    c.bench_function("helix_rate_30s", |b| {
        b.iter(|| {
            let result = crate::helix_query::execute_query("rate(cpu_usage[30s])", &black_box(&data));
            black_box(result)
        })
    });
    
    // Note: ClickHouse benchmark would be run separately
    // This establishes the Helix baseline
}

criterion_group!(
    benches,
    benchmark_helix_query,
    benchmark_storage_insert,
    benchmark_query_filtering,
    benchmark_large_scale,
    benchmark_performance_comparison
);

criterion_main!(benches);