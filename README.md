# Helix Observatory Platform

A SOC (Security Operations Center) observability platform that continues ApexGuard's event generation with full storage, analytics, and querying capabilities. Built as a custom TSDB replacing ClickHouse/Prometheus.

## Project Structure

```
helix/
├── go.mod                  # Go module (helix.io)
├── go.sum                  # Go dependencies
├── PROJECT_SUMMARY.md      # Comprehensive project documentation
├── cmd/                    # Command entry points
├── internal/               # Internal packages
├── pkg/                    # Public SDK packages
├── api/                    # gRPC + REST API definitions
├── rust/                   # Rust workspace (TSDB + Query Engine)
├── agent/                  # Event collection agent (Go)
├── ingestion/              # Ingestion gateway (Go)
├── cli/                    # CLI tools
├── ui/                     # React/TypeScript Web UI
├── helm/                   # Kubernetes deployment charts
├── benches/                # Performance benchmarks (criterion)
├── ci/                     # CI/CD pipeline configuration
└── helm/                   # Helm charts for K8s deployment
```

## Components (~108K lines total)

### 1. Agent Collection (Go)
- Event collection with backpressure mechanisms
- Plugin architecture (health check, metrics extraction, tenant assignment)
- Token bucket rate limiting
- Buffered channel for burst traffic absorption

**Key files**: `agent/main.go`

### 2. Ingestion Gateway (Go)
- Tenancy support (multiple tenant isolation)
- Per-tenant rate limiting
- NATS/Kafka integration for message queuing
- Event validation and routing

**Key files**: `ingestion/main.go`

### 3. Storage Engine (Rust)
- **WAL (Write-Ahead Log)**: Durability for ingested events
  - Sequence-based file rotation
  - Configurable max file size with automatic flush
  - Sync capability for data integrity
  
- **Memtable**: In-memory cache for recent writes
  - BTreeMap for sorted key storage
  - Configurable max entries with auto-flush trigger
  - Approximate size tracking
  
- **SSTable**: Immutable sorted tables for long-term storage
  - Sorted by key for efficient range queries
  - Built-in index for O(log n) lookups
  - Configurable compression (LZ4/ZSTD/None)
  - Target file size management

- **Storage Engine (Composite)**: WAL + Memtable + SSTable
  - Write path: WAL → Memtable → SSTable flush
  - Read path: Memtable (fast) → SSTables (persistent)
  - Flush mechanism for memtable compaction
  - Query interface with multi-storage lookup

**Key files**: `rust/helix-storage/src/lib.rs`

### 4. Query Engine (Rust)
- **Parser**: Full Helix query language parser
  - Instant and range vector parsing
  - Label matching with regex support
  - Duration parsing (e.g., `30s`, `1m`, `5m`)
  - 15+ test cases covering all parser branches
  
- **Scheduler**: Query execution planning
  - Query plan generation based on query type
  - Label filter and time range filtering integration
  - Step-by-step execution planning
  
- **Vectorized Executor**: Query execution engine
  - Multi-step execution pipeline
  - Label filtering, time range filtering
  - Range aggregation, instant computation
  - Scalar computation, post-processing
  
- **Public API**: `execute_query(expression, data)` entry point

**Key files**: `rust/helix-query/src/lib.rs`

### 5. Query Language
- PromQL-like syntax with extensions
- Instant and range vectors
- Label matchers (=, !=, =~)
- Duration expressions (1m, 5m, 30s, 2h)
- Built-in functions (rate, top_k, histogram_quantile, sum by, etc.)

### 6. Go SDK
- Full-featured client for platform interaction
- Event ingestion API
- Query execution API
- Tenant and metadata management

**Key files**: `pkg/go-sdk/main.go`

### 7. TypeScript SDK
- Type-safe JavaScript/TypeScript client
- Full-type-safe query execution
- Query hooks and result types
- Helper functions for common operations

**Key files**: `pkg/typescript-sdk/index.ts`

### 8. API Definitions (Go)
- gRPC service definitions
- REST JSON API endpoints
- OpenAPI 3.0 specification
- Event, query, and metadata types

**Key files**: `api/main.go`

### 9. Web UI (React/TypeScript)
- Query editor with syntax highlighting
- Dashboard composition and layout management
- Event explorer with filtering and tagging
- Alert management interface
- Alert manager with DSL, eval, dedup, routing

**Key files**: `ui/components/types.ts`

### 10. CLI
- Command-line tool for platform operations
- Query execution
- Event ingestion
- Status checking
- Help and documentation

**Key files**: `cli/helix/main.go`

### 11. Helm Charts
- Kubernetes manifests for all components
- Dependencies: PostgreSQL, NATS
- Values configurations for multiple environments (dev, staging, prod)
- Linting and validation scripts

**Key files**: `helm/helix/`

### 12. CI/CD Pipeline
- Multi-stage pipeline: lint → test → build → deploy → benchmark
- Branch protection with required status checks
- Multi-environment deployment (dev, staging, production)
- Security scanning integration
- Notification settings

**Key files**: `ci/pipeline.yaml`

### 13. Performance Benchmarks (criterion)
- Establishes Helix performance baselines
- Critical metric: "1 billion events in < 1 second"
- Benchmark categories: core queries, storage insert, query filtering, large-scale
- Structured for ClickHouse comparison

**Key files**: `benches/benchmark.rs`, `benches/README.md`

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Languages | Go 1.21, Rust 1.70+, TypeScript |
| Storage | Custom TSDB (WAL + Memtable + SSTable) |
| Messaging | NATS, Kafka |
| API | gRPC, REST, OpenAPI |
| Frontend | React, TypeScript |
| Orchestration | Helm, Kubernetes |
| Testing | Go test, cargo test, criterion benchmarks |
| Container | Docker |

## Performance Goals

| Metric | Target |
|--------|--------|
| Query: 1 billion events | < 1 second |
| Ingestion Rate | 10k+ events/second/tenant |
| Durability | WAL + SSTable + Snapshots |
| Scalability | Horizontal via NATS clusters |
| Reliability | 99.9% uptime with WAL recovery |

## Development Phases

### Phase 1: Foundation ✅
- Project structure initialization
- Go module and Rust workspace setup
- Basic agent and gateway implementation
- Storage engine core (WAL + Memtable + SSTable)

### Phase 2: Core Engine ✅
- WAL implementation with sequence numbering
- Memtable with BTreeMap and flush triggers
- SSTable writer with sorting and indexing
- Storage engine composite architecture
- Comprehensive test suite (20+ tests)

### Phase 3: Agent & Gateway ✅
- Go agent with backpressure and plugins
- Ingestion gateway with tenancy and rate-limiting
- Rate limiting with token bucket algorithm
- Burst capacity configuration

### Phase 4: SDK & API ✅
- Go SDK with full client functionality
- TypeScript SDK with type safety
- gRPC and REST API definitions
- Comprehensive type definitions

### Phase 5: Deployment ✅
- Helm charts with all Kubernetes manifests
- CI/CD pipeline configuration
- Values configurations for multiple environments
- Linting and validation scripts

### Phase 6: Refinement ✅
- CI/CD pipeline with multi-stage workflow
- Branch protection and quality gates
- Multi-environment deployment configuration
- Security scanning integration

### Phase 7: Query Engine ✅
- Parser implementation (complete)
- Scheduler implementation (complete)
- Vectorized executor (complete)
- Test suite (15+ tests)

### Phase 8: Benchmarks ✅
- Benchmark suite creation (complete)
- Performance goals established
- ClickHouse comparison framework

### Phase 9: UI Implementation (Upcoming)
- Query editor component
- Dashboard builder
- Explorer interface
- Alert management

### Phase 10: ClickHouse Comparison (Upcoming)
- Implement same queries on ClickHouse
- Compare performance metrics
- Optimize Helix based on findings
- Publish benchmark results

## Getting Started

### Prerequisites
- Go 1.21+
- Rust 1.70+
- Node.js 18+
- Docker (for local development)

### Quick Start

1. **Start the platform**:
   ```bash
   cd helix
   helm install helix ./helm/helix
   ```

2. **Ingest events** (Go):
   ```go
   import "helix.io/sdk/go"
   
   client := helix.NewClient("http://localhost:8080", "api-key")
   event := helix.Event{
       Source: "soc-agent",
       Metrics: map[string]float64{"cpu": 50.0},
       Tags: map[string]string{"host": "soc-01"},
       Tenant: "soc-team-a",
   }
   client.IngestEvent(context.Background(), &event)
   ```

3. **Execute queries** (Go):
   ```go
   results, err := client.Query(context.Background(), &helix.QueryOptions{
       Expression: "rate(events[30s])",
       Tenant: "soc-team-a",
   })
   ```

4. **Use TypeScript SDK**:
   ```typescript
   import { createClient } from 'helix-sdk/typescript';
   
   const client = createClient('http://localhost:8080', 'api-key');
   const results = await client.query({ expression: 'rate(events[30s])' });
   ```

5. **Run benchmarks**:
   ```bash
   cd helix/benches
   criterion
   ```

## License

Helix is open source software licensed under the Apache License 2.0.

## Contact

Project located at: D:\PortfolioProjects\helix