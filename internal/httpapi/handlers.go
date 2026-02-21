package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BryanRamires/FizzBuzz/internal/fizzbuzz"
)

type errorResponse struct {
	Error string `json:"error"`
}

const maxLimitDefault = 100_000

type Handler struct {
	MaxLimit int
}

func NewHandler() Handler {
	return Handler{MaxLimit: maxLimitDefault}
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
	if limit > h.MaxLimit {
		writeError(w, http.StatusBadRequest, "limit is too large")
		return
	}

	str1 := q.Get("str1")
	str2 := q.Get("str2")
	if str1 == "" || str2 == "" {
		writeError(w, http.StatusBadRequest, "str1 and str2 must be non-empty strings")
		return
	}

	out := fizzbuzz.Generate(int1, int2, limit, str1, str2)
	writeJSON(w, http.StatusOK, out)
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
