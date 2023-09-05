package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"tasks-app/internal/modules/emailnotifier"
	"tasks-app/internal/modules/taskchecker"
	"tasks-app/internal/modules/ui"
	"tasks-app/internal/modules/uinotifier"
	"tasks-app/internal/shared"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Config          *shared.Config
	Logger          *slog.Logger
	TaskRepo        shared.TaskRepository
	MessagingClient shared.MessagingClient
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

	a.Logger.Info("app is starting up...")

	if err := a.serve(ctx); err != nil {
		a.Logger.Error("serve", "error", err)
		os.Exit(1)
	}

	a.Logger.Info("app is shut down")
}

func (a *App) loadConfig() error {
	c := &shared.Config{}
	if err := c.Load(); err != nil {
		return err
	}

	slog.Debug("app config", slog.Any("config", c))

	a.Config = c

	return nil
}

func (a *App) initDefaultLogger() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func (a *App) initLogger() {
	level := slog.LevelInfo
	level.UnmarshalText([]byte(a.Config.LogLevel))
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	a.Logger = logger
}

func (a *App) serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	var err error

	if a.TaskRepo, err = shared.NewPostgresTaskRepository(ctx, shared.PostgresTaskRepositoryOptions{
		ConnectionString: a.Config.PostgresConnectionString,
		Logger:           a.Logger,
	}); err != nil {
		return err
	}

	if a.MessagingClient, err = shared.NewNATSMessagingClient(shared.NATSMessagingClientOptions{
		NATSURL:   a.Config.NATSURL,
		NATSToken: a.Config.NATSToken,
		Logger:    a.Logger,
	}); err != nil {
		return err
	}

	var modules []shared.AppModule

	isAllModules := len(a.Config.Modules) == 0

	if isAllModules || slices.Contains(a.Config.Modules, "taskchecker") {
		modules = append(modules, &taskchecker.TaskChecker{
			Config:          a.Config,
			Logger:          a.Logger,
			TaskRepository:  a.TaskRepo,
			MessagingClient: a.MessagingClient,
		})
	}

	if isAllModules || slices.Contains(a.Config.Modules, "emailnotifier") {
		modules = append(modules, &emailnotifier.EmailNotifier{
			Config:          a.Config,
			Logger:          a.Logger,
			MessagingClient: a.MessagingClient,
			EmailClient: &emailnotifier.NullEmailClient{
				Logger: a.Logger,
			},
			// EmailClient: &emailnotifier.SMTPEmailClient{
			// 	Options: emailnotifier.SMTPEmailClientOptions{
			// 		Host:        a.Config.SMTPHost,
			// 		Port:        a.Config.SMTPPort,
			// 		FromName:    a.Config.SMTPFromName,
			// 		FromAddress: a.Config.SMTPFromAddress,
			// 		Password:    a.Config.SMTPPassword,
			// 	}},
		})
	}

	if isAllModules || slices.Contains(a.Config.Modules, "uinotifier") {
		modules = append(modules, &uinotifier.UINotifier{
			Config:          a.Config,
			Logger:          a.Logger,
			MessagingClient: a.MessagingClient,
		})
	}

	if isAllModules || slices.Contains(a.Config.Modules, "ui") {
		modules = append(modules, &ui.UI{
			Config:   a.Config,
			Logger:   a.Logger,
			TaskRepo: a.TaskRepo,
			FileExporter: &shared.ExcelFileExporter{
				Logger: a.Logger},
		})
	}

	for _, m := range modules {
		m := m
		g.Go(func() error {
			a.Logger.Info("run app module", slog.String("module", m.Name()))
			return m.Run(ctx)
		})
	}

	g.Go(func() error {
		<-ctx.Done()

		a.Logger.Info("app is shutting down...")

		for _, m := range modules {
			m.Close()
		}

		_ = a.MessagingClient.Close()
		_ = a.TaskRepo.Close()

		return nil
	})

	return g.Wait()
}
