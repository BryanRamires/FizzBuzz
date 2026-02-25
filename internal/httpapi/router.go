package httpapi

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func NewRouter(cfg config.Config, logger *slog.Logger, h Handler) http.Handler {
	r := chi.NewRouter()

	if cfg.CORSEnabled {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: cfg.CORSAllowedOrigins,
			AllowedMethods: []string{"GET", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type", "X-Request-Id"},
			ExposedHeaders: []string{"X-Request-Id"},
			MaxAge:         300,
		}))
	}

	r.Use(middleware.RequestID)
	// RealIP assumes requests come through a trusted proxy that sets X-Forwarded-For / X-Real-IP.
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.HTTPHandlerTimeout))
	r.Use(LoggingMiddleware(logger))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "ok"); err != nil {
			http.Error(w, "write error", http.StatusInternalServerError)
		}
	})

	r.Get("/readyz", h.Readyz)

	// In multi-instance production setups, rate limiting is typically enforced at the edge
	// (API gateway / ingress / WAF). This in-app limiter provides basic protection per instance.
	r.With(httprate.LimitByIP(cfg.RateLimitFizzBuzz, cfg.RateLimitWindow)).Get("/fizzbuzz", h.FizzBuzz)

	r.With(httprate.LimitByIP(cfg.RateLimitStats, cfg.RateLimitWindow)).Get("/stats", h.StatsTop)

	return r
}
