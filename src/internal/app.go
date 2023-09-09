package internal

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks-app/internal/shared"
	"time"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Config          *shared.Config
	Logger          *slog.Logger
	TaskRepository  shared.TaskRepository
	MessagingClient shared.MessagingClient
	Modules         []shared.AppModule
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a.initDefaultLogger()

	if err := a.loadConfig(); err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	a.initLogger()

	if err := a.createServices(ctx); err != nil {
		slog.Error("create services", "error", err)
		os.Exit(1)
	}

	if err := a.createModules(); err != nil {
		slog.Error("create modules", "error", err)
		os.Exit(1)
	}

	if err := a.serve(ctx); err != nil {
		a.Logger.Error("serve", "error", err)
		os.Exit(1)
	}

	a.Logger.Info("exit app")
}

func (a *App) serve(ctx context.Context) error {
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

		var errs []error

		errs = append(errs, a.closeModules()...)

		time.Sleep(5 * time.Second)

		errs = append(errs, a.closeServices()...)

		if err := errors.Join(errs...); err != nil {
			a.Logger.Error("graceful shutdown", "error", err)
		}

		return nil
	})

	return g.Wait()
}
