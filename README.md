# FizzBuzz REST API (Go)

Production-ready FizzBuzz REST API implemented in Go, featuring a clean and minimalist architecture, configurable rules, structured logging, and request statistics tracking.

This project demonstrates idiomatic Go practices, production-grade HTTP handling, and maintainable software design.

---

## Quick Start (Docker Compose)

The easiest way to run the full stack (API + Redis) is using Docker Compose.

### Start the stack

```bash
docker compose up --build
```

The API will be available at:

http://localhost:8090

By default:

- The API runs on port **8090**
- Redis is started automatically
- The statistics backend can be switched via environment variables

### Stop the stack

```bash
docker compose down
```

---

## Local Development (optional)

Run the server locally without Docker (requires Go installed):

```bash
make run
```

---

## Development Workflow

### Format code

```bash
make fmt
```

### Run tests

```bash
make test
```

### Run race detector

```bash
make race
```

### Run linter (golangci-lint v2, pinned locally)

```bash
make lint
```

### Run vulnerability scan (govulncheck)

```bash
make vuln
```

### Run full CI pipeline locally

Runs exactly what CI executes:

```bash
make ci
```

Includes:

- go mod tidy
- gofmt
- tests
- race detector
- linting
- vulnerability scan

---

## API Endpoints

Quick examples:

### Liveness probe

Indicates that the application process is running.
Useful for potential use of kubernetes probes

```bash
curl -i http://localhost:8090/healthz
```

### Readiness probe

Indicates that the service is ready to receive traffic.

If Redis is enabled, this endpoint verifies the Redis connection.  
If Redis is disabled, it always returns success.
Useful for potential use of kubernetes probes

```bash
curl -i http://localhost:8090/readyz
```

### FizzBuzz

```bash
curl "http://localhost:8090/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz"
```

### Statistics

```bash
curl http://localhost:8090/stats
```

For full API details, see the OpenAPI specification in `docs/openapi.yaml`.

---

## API Documentation (OpenAPI)

The API contract is formally described using an OpenAPI specification.

Location:

```
docs/openapi.yaml
```

You can visualize it using Swagger Editor:

https://editor.swagger.io/

Simply copy the file contents into the editor to explore the interactive documentation and test endpoints.

---

## Configuration

Configuration is handled via environment variables.

A sample configuration file is provided:

```bash
cp config.example.env .env
```

You can customize behavior such as:

- HTTP port
- Maximum allowed `limit`
- Maximum string length
- Rate limiting
- Redis usage for statistics

---

### CORS (optional)

If you want to use Swagger Editor "Try it out" functionality from the browser, enable CORS:

```
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://editor.swagger.io
```

---

## Architecture

The project follows a simple and idiomatic Go structure focused on separation of concerns.

```
cmd/api                 → application entrypoint
internal/httpapi        → HTTP routing & handlers
internal/fizzbuzz       → core business logic
internal/stats          → statistics service
internal/stats/memory   → in-memory statistics repository
internal/stats/redis    → Redis statistics repository
internal/config         → configuration loading
```

### Design Principles

- Clear separation between transport and business logic
- Minimal abstractions (interfaces only at boundaries, see `internal/stats/repo.go`)
- Thread-safe statistics storage
- Production-ready HTTP server configuration
- Structured logging using slog
- Simple and maintainable codebase

---

## Statistics Storage

By default, statistics are stored in memory.

To prevent unbounded memory growth, the in-memory repository limits the number of distinct parameter combinations.

For multi-instance or production deployments, a Redis-backed repository is available and can be enabled via environment configuration.

---

## Production-Ready Features

### HTTP Hardening

- Graceful shutdown
- Read/write/idle timeouts
- Handler-level timeouts
- Panic recovery middleware
- Rate limiting per IP

### Observability

- Structured JSON logging (slog)
- Request ID propagation
- Request duration tracking
- Status code logging

### Concurrency Safety

- Thread-safe statistics repository
- No global mutable state
- Race detector validated

### Code Quality

- golangci-lint v2
- govulncheck
- Unit tests
- Race tests
- CI parity via `make ci`

---

## Potential Improvements

Possible future enhancements include:

- Prometheus metrics endpoint
- Distributed rate limiting
- Persistent statistics with TTL management
- Extended API versioning

---

## Author

Bryan Ramires