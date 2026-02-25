package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/httpapi"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	"github.com/BryanRamires/FizzBuzz/internal/stats/memory"
	redisrepo "github.com/BryanRamires/FizzBuzz/internal/stats/redis"
	goredis "github.com/redis/go-redis/v9"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	srv, err := newServer(cfg, logger)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", cfg.Addr)
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

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("stopped")
	return nil
}

func newServer(cfg config.Config, logger *slog.Logger) (*http.Server, error) {
	var (
		repo stats.Repository
		rdb  *goredis.Client
	)

	if cfg.RedisEnabled {
		rdb = goredis.NewClient(&goredis.Options{
			Addr:        cfg.RedisAddr,
			Password:    cfg.RedisPassword,
			DB:          cfg.RedisDB,
			DialTimeout: cfg.RedisDialTimeout,
		})

		pingCtx, cancel := context.WithTimeout(context.Background(), cfg.RedisDialTimeout)
		defer cancel()
		if err := rdb.Ping(pingCtx).Err(); err != nil {
			return nil, err
		}

		repo = redisrepo.NewRepo(cfg, rdb)
		logger.Info("stats backend", "type", "redis", "addr", cfg.RedisAddr, "db", cfg.RedisDB)
	} else {
		repo = memory.New()
		logger.Info("stats backend", "type", "memory")
	}

	svc, err := stats.NewService(repo)
	if err != nil {
		return nil, err
	}

	h := httpapi.NewHandler(cfg, rdb, svc)
	router := httpapi.NewRouter(cfg, logger, h)

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}, nil
}
