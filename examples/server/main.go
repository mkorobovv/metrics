package main

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mkorobovv/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	registry := prometheus.NewRegistry()
	serverMetrics := metrics.NewServerMetrics()

	registry.MustRegister(serverMetrics)
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	r.HandleFunc("/health", serverMetrics.WrapHandler(healthHandler())).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		logger.Warn("Shutdown signal received")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Error during shutdown", "err", err)
			os.Exit(1)
		}
		logger.Info("Server gracefully stopped")
		return
	}
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		flag := rand.Intn(2)
		if flag == 1 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
