# FizzBuzz REST API (Go)

Production-ready FizzBuzz REST API implemented in Go, featuring a clean and minimalist architecture, configurable rules, structured logging, and request statistics tracking.

This project demonstrates idiomatic Go practices, production-grade HTTP handling, and maintainable software design.

---

## Quick Start (Docker)

### Build the image

```bash
docker build -t fizzbuzz-api .
```

### Run the container

```bash
cp config.example.env .env
docker run --env-file .env -p 8090:8090 fizzbuzz-api
```

The API will be available at:

http://localhost:8090

---

## Local Development (optional)

Run the server locally without Docker (requires Go installed):

```bash
make run
```

### Run tests

```bash
make test
```

---

## API Endpoints

Quick examples:

### Health check
curl -i http://localhost:8090/healthz

### FizzBuzz
curl "http://localhost:8090/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz"

### Statistics
curl http://localhost:8090/stats

For full API details, see the OpenAPI specification in docs/openapi.yaml.

## API Documentation (OpenAPI)

The API contract is formally described using an OpenAPI specification.

Location:

docs/openapi.yaml

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

---

### CORS (optional)

If you want to use Swagger Editor "Try it out" function from the browser, enable CORS:

CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://editor.swagger.io

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
- Minimal abstractions (interfaces only at boundaries, see internal/stats/repo.go)
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

- Redis-backed statistics repository (TTL/eviction)
- Basic metrics endpoint (Prometheus)

---

## Author

Bryan Ramires