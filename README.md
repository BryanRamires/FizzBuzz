# FizzBuzz REST API (Go)

Production-ready FizzBuzz REST API implemented in Go, featuring a clean and minimalist architecture, configurable rules, structured logging, and request statistics tracking.

This project demonstrates idiomatic Go practices, production-grade HTTP handling, and maintainable software design.

---

## Quick Start

### Run the server

```bash
make run
```

The server starts on:

http://localhost:8090

---

### Run tests

```bash
make test
```

---

## API Endpoints

### Health Check

Used for monitoring and readiness probes.

```bash
curl -i http://localhost:8090/healthz
```

---

### FizzBuzz

Generates a configurable FizzBuzz sequence.

```bash
curl "http://localhost:8090/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz"
```

Parameters:

| Parameter | Description |
|----------|-------------|
| int1 | First divisor |
| int2 | Second divisor |
| limit | Maximum number |
| str1 | Replacement for multiples of int1 |
| str2 | Replacement for multiples of int2 |

---

### Statistics

Returns the most frequent request parameters and hit count.

```bash
curl http://localhost:8090/stats
```

---

## Configuration

Configuration is handled via environment variables.

A sample file is provided:

```bash
cp config.example.env .env
```

### Linux / macOS

```bash
set -a
source .env
set +a
make run
```

---

### Windows (PowerShell)

```powershell
cp config.example.env .env
Get-Content .env | ForEach-Object {
    if ($_ -match "^\s*([^#=]+)\s*=\s*(.*)$") {
        [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    }
}
make run
```

---

## Architecture

The project follows a simple and idiomatic Go structure focused on separation of concerns.

```
cmd/api                 → application entrypoint
internal/httpapi        → HTTP routing & handlers
internal/fizzbuzz       → core business logic
internal/stats          → statistics service
internal/stats/memory   → statistics repository
internal/config         → configuration loading
```

### Design Principles

- Clear separation between transport and business logic
- Minimal abstractions (interfaces only at boundaries) -> see internal/stats/repo.go
- Thread-safe in-memory storage
- Production-ready HTTP server configuration
- Structured logging using slog
- Simple and maintainable codebase

---

## Statistics Storage

Statistics are stored in memory for simplicity.

To prevent unbounded memory growth, the in-memory repository limits the number of distinct parameter combinations.

For real production deployments:

- Use Redis with TTL/eviction
- Add rate limiting
- Consider persistence and monitoring

---

## Production-Ready Features

### HTTP Hardening

- Graceful shutdown
- Read/write/idle timeouts
- Panic recovery middleware

### Observability

- Structured JSON logging (slog)
- Request ID propagation
- Request duration tracking
- Status code logging

### Concurrency Safety

- Thread-safe statistics repository
- No global mutable state

---

## Potential Improvements

Possible future enhancements include:

- Docker containerization for easier deployment and reproducible environments.
- Redis-backed statistics repository (TTL/eviction).
- Configurable rate limiting middleware.
- OpenAPI specification for endpoint documentation.
- Basic metrics endpoint (Prometheus).

---

## Notes

The Makefile commands are primarily designed for Windows (PowerShell / curl.exe).  
Minor adjustments may be required for Linux or macOS environments.

---

## Author

Bryan Ramires