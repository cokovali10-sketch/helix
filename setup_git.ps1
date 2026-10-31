git init
git config user.name "HelixUser"
git config user.email "helix@users.noreply.github.com"
git branch -M main

# Add all files
git add -A

# Create initial commit
git commit -m "Initial: Helix SOC observability platform structure

- Go module (helix.io)
- Agent collection with backpressure
- Ingestion gateway with tenancy/rate-limit/NATS
- Rust TSDB engine (WAL + Memtable + SSTable)
- Query engine (parser + scheduler + executor)
- Go and TypeScript SDKs
- gRPC + REST API definitions
- React/TypeScript UI types
- CLI tools
- Helm charts for K8s deployment
- CI/CD pipeline configuration
- Performance benchmarks
- Full documentation"

echo "Initial commit created successfully"