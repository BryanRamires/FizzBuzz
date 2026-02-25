package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/fizzbuzz"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	goredis "github.com/redis/go-redis/v9"
)

type errorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	cfg   config.Config
	rdb   *goredis.Client
	stats *stats.Service
}

func NewHandler(cfg config.Config, rdb *goredis.Client, statsService *stats.Service) Handler {
	return Handler{
		cfg:   cfg,
		rdb:   rdb,
		stats: statsService,
	}
}

func (h Handler) FizzBuzz(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	int1, err := mustPositiveInt(q.Get("int1"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "int1 must be a positive integer")
		return
	}
	int2, err := mustPositiveInt(q.Get("int2"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "int2 must be a positive integer")
		return
	}
	limit, err := mustPositiveInt(q.Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "limit must be a positive integer")
		return
	}
	if limit > h.cfg.MaxLimit {
		writeError(w, http.StatusBadRequest, "limit is too large")
		return
	}

	str1 := q.Get("str1")
	str2 := q.Get("str2")
	if str1 == "" || str2 == "" {
		writeError(w, http.StatusBadRequest, "str1 and str2 must be non-empty strings")
		return
	}

	if hasControlChars(str1) || hasControlChars(str2) {
		writeError(w, http.StatusBadRequest, "str1 and str2 must not contain control characters")
		return
	}

	if len(str1) > h.cfg.MaxStrLen || len(str2) > h.cfg.MaxStrLen {
		writeError(w, http.StatusBadRequest, "str1 and str2 are too long")
		return
	}

	if h.stats != nil {
		h.stats.Record(stats.Key{
			Int1:  int1,
			Int2:  int2,
			Limit: limit,
			Str1:  str1,
			Str2:  str2,
		})
	}

	out := fizzbuzz.Generate(int1, int2, limit, str1, str2)
	writeJSON(w, http.StatusOK, out)
}

func (h Handler) StatsTop(w http.ResponseWriter, r *http.Request) {
	if h.stats == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	top, ok := h.stats.MostFrequent()
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, top)
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RedisEnabled {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		if err := h.rdb.Ping(ctx).Err(); err != nil {
			writeError(w, http.StatusServiceUnavailable, "redis not ready")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}
func hasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func mustPositiveInt(s string) (int, error) {
	if s == "" {
		return 0, errors.New("missing")
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, errors.New("invalid")
	}
	return n, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_, _ = w.Write(buf.Bytes())
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
