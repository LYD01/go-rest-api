// Entry point: config, wiring, graceful shutdown
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYD01/my-go-app/internal/albums"
	"github.com/LYD01/my-go-app/internal/config"
	"github.com/LYD01/my-go-app/internal/middleware"
	"github.com/LYD01/my-go-app/internal/respond"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	store := albums.NewStore()
	handler := albums.NewHandler(store)

	albumMux := http.NewServeMux()
	albumMux.HandleFunc("GET /albums", handler.GetAll())
	albumMux.HandleFunc("GET /albums/{id}", handler.GetById())
	albumMux.HandleFunc("POST /albums", handler.Create())

	protected := middleware.RequireApiKey(cfg.APIKey)(albumMux)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("/", protected)

	stack := middleware.RequestId(
		middleware.Logger(
			middleware.Recoverer(mux),
		),
	)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      stack,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server started", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
