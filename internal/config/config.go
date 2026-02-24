package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr               string
	MaxLimit           int
	ReadHeaderTimeout  time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	HTTPHandlerTimeout time.Duration

	RateLimitWindow   time.Duration
	RateLimitFizzBuzz int
	RateLimitStats    int

	CORSEnabled        bool
	CORSAllowedOrigins []string
}

func New() (Config, error) {
	cfg := Config{
		Addr:               ":8090",
		MaxLimit:           100_000,
		ReadHeaderTimeout:  5 * time.Second,
		ReadTimeout:        5 * time.Second,
		WriteTimeout:       10 * time.Second,
		IdleTimeout:        60 * time.Second,
		ShutdownTimeout:    10 * time.Second,
		HTTPHandlerTimeout: 30 * time.Second,

		RateLimitWindow:   1 * time.Minute,
		RateLimitFizzBuzz: 30,
		RateLimitStats:    10,

		CORSEnabled: false,
		// Do not default to "*" to avoid unintentionally opening the API.
		// CORS must always be explicitly configured.
		CORSAllowedOrigins: []string{},
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
	if cfg.HTTPHandlerTimeout, err = getenvDuration("HTTP_HANDLER_TIMEOUT", cfg.HTTPHandlerTimeout); err != nil {
		return Config{}, fmt.Errorf("HTTP_HANDLER_TIMEOUT: %w", err)
	}
	if cfg.HTTPHandlerTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_HANDLER_TIMEOUT must be > 0")
	}
	if cfg.RateLimitWindow, err = getenvDuration("RATE_LIMIT_WINDOW", cfg.RateLimitWindow); err != nil {
		return Config{}, fmt.Errorf("RATE_LIMIT_WINDOW: %w", err)
	}
	if cfg.RateLimitWindow <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_WINDOW must be > 0")
	}

	if cfg.RateLimitFizzBuzz, err = getenvInt("RATE_LIMIT_FIZZBUZZ", cfg.RateLimitFizzBuzz); err != nil {
		return Config{}, fmt.Errorf("RATE_LIMIT_FIZZBUZZ: %w", err)
	}
	if cfg.RateLimitFizzBuzz <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_FIZZBUZZ must be > 0")
	}

	if cfg.RateLimitStats, err = getenvInt("RATE_LIMIT_STATS", cfg.RateLimitStats); err != nil {
		return Config{}, fmt.Errorf("RATE_LIMIT_STATS: %w", err)
	}
	if cfg.RateLimitStats <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_STATS must be > 0")
	}
	if cfg.CORSEnabled, err = getenvBool("CORS_ENABLED", cfg.CORSEnabled); err != nil {
		return Config{}, fmt.Errorf("CORS_ENABLED: %w", err)
	}

	origins := getenv("CORS_ALLOWED_ORIGINS", "")
	if strings.TrimSpace(origins) != "" {
		cfg.CORSAllowedOrigins = splitCSV(origins)
	}
	if cfg.CORSEnabled && len(cfg.CORSAllowedOrigins) == 0 {
		return Config{}, fmt.Errorf("CORS_ALLOWED_ORIGINS must be set when CORS_ENABLED=true")
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

func getenvBool(k string, def bool) (bool, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return def, nil
	}
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return false, fmt.Errorf("is set but empty")
	}
	switch v {
	case "1", "true", "yes", "y", "on":
		return true, nil
	case "0", "false", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool %q (expected true/false)", v)
	}
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}
