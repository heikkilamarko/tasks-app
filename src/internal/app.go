package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks-app/internal/shared"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Logger          *slog.Logger
	Config          *shared.Config
	TaskRepository  shared.TaskRepository
	MessagingClient shared.MessagingClient
	Modules         []shared.AppModule
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.init(ctx); err != nil {
		a.Logger.Error("app init", "error", err)
		os.Exit(1)
	}

	if err := a.run(ctx); err != nil {
		a.Logger.Error("app run", "error", err)
		os.Exit(1)
	}

	if err := a.close(ctx); err != nil {
		a.Logger.Error("app close", "error", err)
	}

	a.Logger.Info("app exit")
}

func (a *App) init(ctx context.Context) error {
	if err := a.createLogger(); err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	if err := a.loadConfig(); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := a.createServices(ctx); err != nil {
		return fmt.Errorf("create services: %w", err)
	}

	if err := a.createModules(); err != nil {
		return fmt.Errorf("create modules: %w", err)
	}

	return nil
}

func (a *App) close(ctx context.Context) error {
	if err := errors.Join(a.closeServices()...); err != nil {
		return fmt.Errorf("close services: %w", err)
	}

	return nil
}

func (a *App) run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, m := range a.Modules {
		m := m
		g.Go(func() error {
			a.Logger.Info("run app module", slog.String("module", m.Name()))
			return m.Run(ctx)
		})
	}

	g.Go(func() error {
		<-ctx.Done()

		a.Logger.Info("graceful shutdown")

		if err := errors.Join(a.closeModules()...); err != nil {
			a.Logger.Error("graceful shutdown", "error", err)
		}

		return nil
	})

	return g.Wait()
}
