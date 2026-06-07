# Deployment Instructions - HyperHarness v0.4.5

## Build Requirements
- Go >= 1.24

## Standard Build
```bash
go build -o hypercode main.go
```

## Running the Daemon
To start the control plane server:
```bash
./hypercode serve
```

## Running the TUI
To launch the interactive dashboard:
```bash
./hypercode tui
```

## Docker Deployment
```bash
docker build -t hyperharness:latest .
docker run -p 8080:8080 hyperharness:latest
```
