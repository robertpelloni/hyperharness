# Deployment Instructions - Borg / HyperHarness v0.5.0-alpha.5

## Build Requirements
- Go >= 1.24
- Node >= 20, pnpm >= 10 (for web dashboard & TS packages)

## Standard Build
To build both `hypercode` and `hyperharness` binaries with version string injection and non-destructive caching:
```bash
./build.sh
```

Or manually with Go:
```bash
go build -ldflags="-s -w -X github.com/robertpelloni/hyperharness/internal/buildinfo.Version=0.5.0-alpha.5" -o hypercode .
cp -f hypercode hyperharness
```

## Running the Daemon Backend (Fiber)
To start the Borg control plane server on port 8080:
```bash
./hypercode serve
```
Health endpoint is accessible at:
```bash
curl http://localhost:8080/health
```

## Running the Interactive Terminal UI (BubbleTea)
To launch the native TUI orchestrator and dashboard:
```bash
./hypercode tui
```

## Running the Web Dashboard (Next.js)
```bash
cd apps/web
pnpm install
pnpm dev
```

## Docker Deployment
```bash
docker build -t hyperharness:latest .
docker run -p 8080:8080 hyperharness:latest
```
