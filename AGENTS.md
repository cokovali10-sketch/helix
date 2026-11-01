# Helix Project Context

## Project Overview
Helix is a SOC (Security Operations Center) observability platform with custom TSDB.
- Location: D:\PortfolioProjects\helix
- GitHub: https://github.com/cokovali10-sketch/helix
- Language: Go, Rust, TypeScript/React

## Tech Stack
- **Agent**: Go (event collection, backpressure, plugins)
- **Ingestion**: Go (tenancy, rate-limiting, NATS/Kafka)
- **Storage**: Rust (WAL, Memtable, SSTable, compression, retention)
- **Query Engine**: Rust (parser, scheduler, vectorized executor)
- **Query Language**: PromQL-like (lexer, types, functions)
- **Alerting**: Go + Rust (DSL, eval, dedup, routing)
- **API**: Go (gRPC + REST + OpenAPI)
- **Web UI**: React/TypeScript (query editor, dashboards, explorer)
- **CLI + SDK**: Go, TypeScript, Python
- **Infrastructure**: Helm, Kubernetes, CI/CD

## Build Commands
- Go build: `go build ./...`
- Go vet: `go vet ./...`
- Go test: `go test ./...`
- Rust build: `cargo build` (in rust/ directory)
- Rust test: `cargo test` (in rust/ directory)
- UI: `cd ui && npm install && npm run build`

## Code Conventions
- Go: standard Go conventions, no comments unless necessary
- Rust: standard Rust conventions
- TypeScript: strict mode, functional components
- No emojis in code
- Commit after each change

## Important Notes
- Windows Defender may block test binaries
- No cargo/rustc available on system for verification
- Use forge toolkit at D:\Forge-Factority for auditing
- Git remote: https://github.com/cokovali10-sketch/helix.git
- Branch: master
- Do NOT create duplicate files
- Do NOT inflate line counts with padding
- Keep code compiling (go build ./... must pass)
