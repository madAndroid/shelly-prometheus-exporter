# Copilot Instructions for shelly-prometheus-exporter

## Project Overview
- This is a Go-based Prometheus exporter for monitoring Shelly smart devices (e.g., 2.5 roller shutter, 1PM plugs).
- The exporter scrapes device status via HTTP, parses the response, and exposes metrics on a `/metrics` endpoint for Prometheus.
- Configuration is provided via `config.yaml` (see example in repo).

## Key Components
- `main.go`: Entry point. Loads config, registers metrics, starts device polling and HTTP server.
- `config.go`: Loads YAML config using Viper. Defines `configuration` and `device` structs.
- `endpoint_status.go`: Defines the `StatusResponse` struct and related types for device status JSON.
- `http.go`: Handles HTTP requests to devices, parses responses, updates Prometheus metrics.
- `metric.go`: Defines and registers Prometheus metrics (gauges, counters) for device data.

## Build & Run
- Use the Makefile for common tasks:
  - `make run` — Build and run locally
  - `make watch` — Hot reload (requires `air`)
  - `make test` — Run all tests
  - `make docker-build` — Build Docker image
  - `make docker-run` — Run Docker container (bind-mounts `config.yaml`)
- Dockerfiles for x86 and ARM (`Dockerfile.arm32v7`) are provided for multi-arch builds.

## Patterns & Conventions
- All Go code is in the `main` package; no sub-packages.
- Device polling interval and HTTP timeouts are set in `config.yaml`.
- Device credentials and addresses are configured per device in YAML.
- Prometheus metrics use consistent labels: `name`, `address`, `type`.
- Device status is fetched in a background goroutine and metrics are updated atomically.
- Errors in device polling increment a Prometheus counter and are logged to stdout.

## External Dependencies
- Prometheus Go client (`github.com/prometheus/client_golang`)
- Viper for config (`github.com/spf13/viper`)

## Example config.yaml
```yaml
port: 9123
requestTimeout: 5s
scrapeInterval: 60s
devices:
  - IPAddress: "192.168.88.38"
    displayName: "OfficePlug"
    type: "1pm"
```

## Testing
- Run `make test` to execute all Go tests with coverage.

## Extending
- To add new metrics, update `metric.go` and extend the parsing logic in `http.go`.
- To support new device types, update the `StatusResponse` struct in `endpoint_status.go`.

## References
- See `README.md` for a high-level project summary.
- See `Makefile` for all supported developer commands.
