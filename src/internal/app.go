package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Logger          *slog.Logger
	Config          *shared.Config
	TaskRepository  shared.TaskRepository
	MessagingClient shared.MessagingClient
	Modules         map[string]shared.AppModule
}

func (a *App) Run(ctx context.Context) error {
	if err := a.init(ctx); err != nil {
		return fmt.Errorf("init: %w", err)
	}

	if err := a.run(ctx); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if err := a.close(ctx); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
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

func (a *App) run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for key, module := range a.Modules {
		key, module := key, module
		g.Go(func() error {
			a.Logger.Info("run app module", slog.String("module", key))
			defer a.Logger.Info("exit app module", slog.String("module", key))
			return module.Run(ctx)
		})
	}

	return g.Wait()
}

func (a *App) close(ctx context.Context) error {
	if err := errors.Join(a.closeServices()...); err != nil {
		return fmt.Errorf("close services: %w", err)
	}

	return nil
}
