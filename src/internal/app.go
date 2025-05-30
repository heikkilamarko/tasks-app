package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"

	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Logger                    *slog.Logger
	Config                    *shared.Config
	DB                        *sql.DB
	NATSConn                  *nats.Conn
	TxManager                 shared.TxManager
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	MessagingClient           shared.MessagingClient
	Modules                   map[string]shared.AppModule
}

func (a *App) Run(ctx context.Context) error {
	if err := a.init(ctx); err != nil {
		return fmt.Errorf("init: %w", err)
	}

	if err := a.run(ctx); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if err := a.close(); err != nil {
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

	for k, m := range a.Modules {
		g.Go(func() error {
			a.Logger.Info("run app module", slog.String("module", k))
			defer a.Logger.Info("exit app module", slog.String("module", k))
			return m.Run(ctx)
		})
	}

	return g.Wait()
}

func (a *App) close() error {
	if err := errors.Join(a.closeServices()...); err != nil {
		return fmt.Errorf("close services: %w", err)
	}

	return nil
}
