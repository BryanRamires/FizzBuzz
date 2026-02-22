package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/fizzbuzz", h.FizzBuzz)

	r.Get("/stats", h.StatsTop)

	return r
}
