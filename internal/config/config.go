package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr              string
	MaxLimit          int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func New() (Config, error) {
	cfg := Config{
		Addr:              ":8090",
		MaxLimit:          100_000,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}

	var err error

	cfg.Addr = getenv("ADDR", cfg.Addr)
	if strings.TrimSpace(cfg.Addr) == "" {
		return Config{}, fmt.Errorf("ADDR must not be empty")
	}

	if cfg.MaxLimit, err = getenvInt("MAX_LIMIT", cfg.MaxLimit); err != nil {
		return Config{}, fmt.Errorf("MAX_LIMIT: %w", err)
	}
	if cfg.MaxLimit <= 0 {
		return Config{}, fmt.Errorf("MAX_LIMIT must be > 0")
	}

	if cfg.ReadHeaderTimeout, err = getenvDuration("READ_HEADER_TIMEOUT", cfg.ReadHeaderTimeout); err != nil {
		return Config{}, fmt.Errorf("READ_HEADER_TIMEOUT: %w", err)
	}
	if cfg.ReadHeaderTimeout <= 0 {
		return Config{}, fmt.Errorf("READ_HEADER_TIMEOUT must be > 0")
	}

	if cfg.ReadTimeout, err = getenvDuration("READ_TIMEOUT", cfg.ReadTimeout); err != nil {
		return Config{}, fmt.Errorf("READ_TIMEOUT: %w", err)
	}
	if cfg.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("READ_TIMEOUT must be > 0")
	}

	if cfg.WriteTimeout, err = getenvDuration("WRITE_TIMEOUT", cfg.WriteTimeout); err != nil {
		return Config{}, fmt.Errorf("WRITE_TIMEOUT: %w", err)
	}
	if cfg.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("WRITE_TIMEOUT must be > 0")
	}

	if cfg.IdleTimeout, err = getenvDuration("IDLE_TIMEOUT", cfg.IdleTimeout); err != nil {
		return Config{}, fmt.Errorf("IDLE_TIMEOUT: %w", err)
	}
	if cfg.IdleTimeout <= 0 {
		return Config{}, fmt.Errorf("IDLE_TIMEOUT must be > 0")
	}

	if cfg.ShutdownTimeout, err = getenvDuration("SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout); err != nil {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT: %w", err)
	}
	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be > 0")
	}

	return cfg, nil
}

func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

func getenvInt(k string, def int) (int, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return def, nil
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, fmt.Errorf("is set but empty")
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid int %q: %w", v, err)
	}
	return n, nil
}

func getenvDuration(k string, def time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return def, nil
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, fmt.Errorf("is set but empty")
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", v, err)
	}
	return d, nil
}
