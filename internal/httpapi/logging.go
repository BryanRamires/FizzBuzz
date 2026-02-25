package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			path := r.URL.Path
			reqID := middleware.GetReqID(r.Context())
			durUs := time.Since(start).Microseconds()
			status := ww.Status()

			if (path == "/healthz" || path == "/readyz") && status == http.StatusOK {
				return
			}

			attrs := []any{
				"method", r.Method,
				"url", r.URL.String(),
				"status", status,
				"bytes", ww.BytesWritten(),
				"duration_us", durUs,
				"request_id", reqID,
			}

			if path == "/healthz" || path == "/readyz" {
				logger.Warn("healthcheck failed", attrs...)
				return
			}

			logger.Info("http request", attrs...)
		})
	}
}
