# Notes
The current Makefile commands are designed for Windows (PowerShell / curl.exe).
They may require adaptation for Linux or macOS environments.

# FizzBuzz
Production-ready FizzBuzz REST API in Go, featuring clean and minimalist architecture, configurable rules, and request statistics tracking.

# Run
make run

# Tests
make test

# Health check
curl.exe -i 'http://localhost:8090/healthz'

# FizzBuzz
curl 'http://localhost:8090/fizzbuzz?int1=0&int2=5&limit=16&str1=fizz&str2=buzz'

# Data
Stats are stored in memory for simplicity. To prevent unbounded growth, the in-memory repository caps the number of distinct parameter combinations. For real production deployments or untrusted traffic, use Redis with TTL / eviction and add rate limiting.