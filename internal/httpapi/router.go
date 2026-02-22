package httpapi

import (
	"log"
	"net/http"

	"github.com/BryanRamires/FizzBuzz/internal/stats"
	"github.com/BryanRamires/FizzBuzz/internal/stats/memory"
	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	repo := memory.New()
	svc, err := stats.NewService(repo)
	if err != nil {
		log.Fatal(err)
	}
	h := NewHandler(svc)

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/fizzbuzz", h.FizzBuzz)

	r.Get("/stats", h.StatsTop)

	return r
}
