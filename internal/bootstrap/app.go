package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appconfig "github.com/rs-karasal/dreams-come-true/internal/config"
)

type App struct {
	config appconfig.Config
	logger *slog.Logger
}

func New() (*App, error) {
	cfg, err := appconfig.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return &App{
		config: cfg,
		logger: logger,
	}, nil
}

func (a *App) RunApp(ctx context.Context) error {
	// context that can be cancelled in next situations:
	//
	// os.Interrupt => Ctrl + C
	// SIGTERM => Docker and Kubernetes send to stop container
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// router: chi + middleware
	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})

	// 4. http.Server with timeouts
	srv := &http.Server{
		Addr:              a.config.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("server starting", "addr", a.config.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	// wait for a signal OR a server error
	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
	}

	// graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
