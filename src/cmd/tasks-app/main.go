package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks-app/internal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := &internal.App{}

	if err := app.Run(ctx); err != nil {
		slog.Error("run app", "error", err)
		os.Exit(1)
	}

	slog.Info("exit app")
}
