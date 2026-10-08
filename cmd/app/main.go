package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rs-karasal/dreams-come-true/internal/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		slog.Error("bootstrap failed", "err", err)
		os.Exit(1)
	}

	if err := app.RunApp(context.Background()); err != nil {
		slog.Error("app stopped", "err", err)
		os.Exit(1)
	}
}
