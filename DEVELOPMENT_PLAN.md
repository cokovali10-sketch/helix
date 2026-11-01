# INSTRUCTIONS FOR AI CONTINUING HELIX DEVELOPMENT

## Project Location
D:\PortfolioProjects\helix

## Goal
Continue developing the Helix SOC observability platform. Target: ~106,000 lines of code. Currently at ~56,300 lines. Need to add ~50,000 more lines of functional, useful code.

## How to Work
1. Write code in batches (large files, 300-600 lines each)
2. Commit after each batch: `git add -A` then `git commit -m "..."` then `git push origin master`
3. Use separate commands (no && chaining in this shell)
4. Write functional, useful code - not stubs
5. Add tests for each new module
6. Keep code compiling (Go builds, Rust needs cargo verification)

## Component Targets and Current State

### 1. Agent Collection (Go) - Target: 8K, Current: ~2.5K
**Need to add: ~5.5K lines**
- More plugins (Kubernetes, Docker, AWS, GCP, Azure collectors)
- More data sources (file, HTTP, syslog, journald, windows event log)
- More exporters (Kafka, NATS, gRPC, WebSocket)
- Advanced filtering and routing
- Configuration management
- Health checks and diagnostics
- Service discovery and registration
- Load balancing and failover
- Metrics collection and export
- Tracing and profiling

### 2. Ingestion Gateway (Go) - Target: 7K, Current: ~2K
**Need to add: ~5K lines**
- Advanced rate limiting (sliding window, leaky bucket, token bucket)
- Tenant management (CRUD, quotas, billing)
- NATS/Kafka integration (publish, subscribe, stream)
- Batch processing and buffering
- Event validation and transformation
- Event routing and filtering
- Monitoring and alerting
- Health checks and diagnostics
- Configuration management

### 3. Storage Engine (Rust) - Target: 15K, Current: ~8K
**Need to add: ~7K lines**
- WAL (write-ahead log) with recovery
- Memtable (skip list, arena, concurrent)
- SSTable (writer, reader, merger, splitter, validator)
- Compression (LZ4, ZSTD, Snappy, Gzip)
- Retention (time-based, size-based, series-based)
- Compaction (leveled, size-tiered, time-window, universal)
- Index (B-tree, inverted, time-series, composite)
- Cache (LRU, write-through, write-back)
- Bloom filter (standard, counting, scalable, cuckoo, quotient)
- Block manager (allocation, caching, prefetch)
- Time series (points, series, DB, downsampler, aggregator)
- Metrics (Prometheus, JSON, tracing)
- Query engine (scan, filter, project, aggregate, sort, limit)
- Transaction support
- Snapshot and backup
- Recovery and repair

### 4. Query Engine (Rust) - Target: 12K, Current: ~3K
**Need to add: ~9K lines**
- Lexer (tokenizer, keywords, operators)
- Parser (AST, expressions, aggregations, selectors)
- Planner (cost-based, plan cache, 10+ rules)
- Executor (vectorized, binary/unary ops, functions)
- Optimizer (predicate pushdown, projection pushdown, etc.)
- Cache (LRU, TTL, tag/prefix invalidation)
- Functions (40+ built-in functions)
- Subqueries and joins
- Time series operations

### 5. Query Language (Rust) - Target: 6K, Current: ~1.5K
**Need to add: ~4.5K lines**
- Lexer (tokens, keywords, operators, durations)
- Parser (expressions, aggregations, selectors)
- Types (InstantVector, RangeVector, Scalar, String)
- Functions (rate, irate, increase, sum, avg, min, max, etc.)
- Operators (arithmetic, comparison, logical, set)
- Aggregations (by, without, group_left, group_right)
- Time ranges and offsets
- Subqueries

### 6. Alerting (Rust + Go) - Target: 7K, Current: ~1.5K
**Need to add: ~5.5K lines**
- Alert DSL (parser, evaluator)
- Alert rules (CRUD, validation)
- Alert evaluation (periodic, on-demand)
- De-duplication (fingerprinting, grouping)
- Routing (matchers, receivers, routes)
- Silence rules (time-based, label-based)
- Notification channels (email, Slack, webhook, PagerDuty)
- Alert history and tracking
- Alert templates

### 7. Catalog/Metadata/RBAC/API - Target: 12K, Current: ~2K
**Need to add: ~10K lines**
- Metadata storage (metrics, dashboards, alerts, tenants)
- RBAC (roles, permissions, users, groups)
- API (gRPC + REST + OpenAPI)
- Search and autocomplete
- Audit logging
- Configuration management
- Feature flags
- Rate limiting
- Authentication (JWT, sessions)
- Authorization (middleware, checks)

### 8. Web UI (React/TS) - Target: 18K, Current: ~6K
**Need to add: ~12K lines**
- Query editor (Monaco, syntax highlighting, autocomplete)
- Dashboards (grid, panels, drag-and-drop)
- Explorer (metrics, events, logs)
- Alerts (rules, notifications, history)
- Settings (profile, notifications, API, advanced)
- Charts (line, bar, area, heatmap, table)
- Components (MetricCard, ChartPanel, DataTable, FilterBar, TimeRangePicker, etc.)
- Contexts (Theme, Auth, Query, Toast)
- Hooks (useQuery, useMetrics, useAlerts, etc.)
- Services (API client, WebSocket, SSE)

### 9. CLI + SDK (Go/TS/Python) - Target: 6K, Current: ~8K ✅
Already exceeded target. Focus on quality and tests.

### 10. Tests/Benches/Helm/Migrations - Target: 15K, Current: ~4K
**Need to add: ~11K lines**
- Go tests (unit, integration, benchmark)
- Rust tests (unit, integration, benchmark)
- React tests (component, integration)
- Helm charts (templates, values, tests)
- Database migrations
- CI/CD pipeline
- Docker configurations
- Documentation

## Development Workflow

### For Each New Module:
1. Create the module file
2. Add comprehensive implementation (300-600 lines)
3. Add tests (50-100 lines)
4. Commit: `git add -A` then `git commit -m "..."` then `git push origin master`
5. Move to next module

### Priority Order:
1. Rust storage engine (most critical - custom TSDB)
2. Rust query engine (query capabilities)
3. Go agent (event collection)
4. Go ingestion (event ingestion)
5. React UI (visualization)
6. Go API (service layer)
7. Alerting (monitoring)
8. Catalog/RBAC (metadata)
9. Tests and benchmarks
10. Documentation

### Code Quality Guidelines:
- Write functional, useful code - not stubs
- Add proper error handling
- Add comprehensive tests
- Use proper types and interfaces
- Follow Go/Rust/TypeScript conventions
- Add comments only when necessary
- Keep functions focused and small
- Use existing patterns from the codebase

### File Naming Convention:
- Go: `package/feature.go` and `package/feature_test.go`
- Rust: `rust/crate/src/feature.rs` and `rust/crate/src/feature_tests.rs`
- React: `ui/src/components/Category/Component.tsx`
- Tests: `ui/src/components/Category/Component.test.tsx`

### Commit Message Format:
"Add: [Component] [feature] ([details])"

Examples:
- "Add: Rust storage compression (LZ4, ZSTD, Snappy, Gzip) + 7 tests"
- "Add: Go agent health checker, metrics collector, profiler + 10 tests"
- "Add: React TimeSeriesChart with Chart.js, formatting, threshold"

### Important Notes:
- Cannot use && in shell commands (use separate commands)
- Windows Defender blocks test binary execution
- No cargo/rustc on system (Rust code can't be verified locally)
- Go builds and vets successfully
- Use forge toolkit at D:\Forge-Factority for auditing
- Git remote: https://github.com/cokovali10-sketch/helix.git
- Branch: master

### Current File Structure:
```
helix/
├── agent/          # Go agent (main.go, agent.go, plugins.go, health.go, metrics.go, service.go, config.go, controller.go, logcollector.go, servicemesh.go, monitor.go, scheduler.go, cluster.go, cluster_advanced.go, dashboard.go, dashboard2.go, dashboard3.go, exporter.go, orchestrator.go, supervisor.go)
├── ingestion/      # Go ingestion (main.go, gateway.go, advanced.go, advanced2.go, monitor.go)
├── storage/        # Go storage (engine.go, compression.go, transaction.go, lsm.go, lsm_advanced.go, batch.go, integrity.go, query.go, query_advanced.go, buffer.go, index_advanced.go, cache_advanced.go, snapshot.go, metrics_advanced.go)
├── api/            # Go API (main.go, server.go, grpc.go, auth.go, handlers.go, handlers_extra.go, search.go, streaming.go, debug.go, dashboard.go, config_handler.go, config_handler2.go, webhook.go, auth_handler.go)
├── catalog/        # Go catalog (metadata.go)
├── catalogapi/     # Go catalog API (handler.go, handler_test.go)
├── cli/            # Go CLI (main.go, commands.go)
├── pkg/            # SDKs (go-sdk/, typescript-sdk/, python-sdk/)
├── rust/           # Rust workspace
│   ├── helix-storage/src/  # 20+ modules
│   ├── helix-query/src/    # 8+ modules
│   └── helix-alerting/src/ # 1 module
├── ui/             # React/TS UI
│   ├── src/components/     # 15+ components
│   ├── src/pages/          # 8 pages
│   ├── src/context/        # 3 contexts
│   └── src/styles/         # CSS
├── helm/           # Kubernetes charts
├── benches/        # Performance benchmarks
├── ci/             # CI/CD pipeline
├── .github/workflows/ # GitHub Actions
└── docs/           # Documentation
```

### How to Continue:
1. Pick a component from the priority list
2. Write a substantial file (300-600 lines)
3. Add tests (50-100 lines)
4. Commit and push
5. Move to next component
6. Repeat until 100K lines reached

### Example Code Patterns:
- Use existing types and interfaces from the codebase
- Follow Go/Rust/TypeScript conventions
- Add comprehensive error handling
- Include unit tests in the same file
- Use sync.Mutex for thread safety
- Use context.Context for cancellation

### Verification Commands:
- Go build: `go build ./...`
- Go vet: `go vet ./...`
- Go test: `go test ./...` (blocked by Windows Defender)
- Rust build: `cargo build` (not available on system)
- Rust test: `cargo test` (not available on system)

### Current README Status:
README has been updated with honest statistics. Should be updated again when significant progress is made.

### GitHub Status:
- Repository: https://github.com/cokovali10-sketch/helix
- Branch: master
- All changes committed and pushed
- CI/CD configured via GitHub Actions

## END OF INSTRUCTIONS