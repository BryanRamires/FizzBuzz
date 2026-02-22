package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/httpapi"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	"github.com/BryanRamires/FizzBuzz/internal/stats/memory"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	srv, err := newServer(cfg)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}

	log.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("stopped")
	return nil
}

func newServer(cfg config.Config) (*http.Server, error) {
	repo := memory.New()
	svc, err := stats.NewService(repo)
	if err != nil {
		return nil, err
	}

	h := httpapi.NewHandler(cfg.MaxLimit, svc)
	router := httpapi.NewRouter(h)

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}, nil
}
