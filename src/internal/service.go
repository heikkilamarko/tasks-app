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

type Service struct {
	Config          *shared.Config
	Logger          *slog.Logger
	TaskRepo        shared.TaskRepository
	MessagingClient shared.MessagingClient
}

func (s *Service) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s.initDefaultLogger()

	if err := s.loadConfig(); err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	s.initLogger()

	s.Logger.Info("application is starting up...")

	if err := s.serve(ctx); err != nil {
		s.Logger.Error("serve", "error", err)
		os.Exit(1)
	}

	s.Logger.Info("application is shut down")
}

func (s *Service) loadConfig() error {
	c := &shared.Config{}
	if err := c.Load(); err != nil {
		return err
	}

	slog.Debug("app config", slog.Any("config", c))

	s.Config = c

	return nil
}

func (s *Service) initDefaultLogger() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func (s *Service) initLogger() {
	level := slog.LevelInfo
	level.UnmarshalText([]byte(s.Config.LogLevel))
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	s.Logger = logger
}

func (s *Service) serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	var err error

	if s.TaskRepo, err = shared.NewPostgresTaskRepository(ctx, shared.PostgresTaskRepositoryOptions{
		ConnectionString: s.Config.PostgresConnectionString,
		Logger:           s.Logger,
	}); err != nil {
		return err
	}

	if s.MessagingClient, err = shared.NewNATSMessagingClient(shared.NATSMessagingClientOptions{
		NATSURL:   s.Config.NATSURL,
		NATSToken: s.Config.NATSToken,
		Logger:    s.Logger,
	}); err != nil {
		return err
	}

	var modules []shared.AppModule

	isAllModules := len(s.Config.Modules) == 0

	if isAllModules || slices.Contains(s.Config.Modules, "taskchecker") {
		modules = append(modules, &taskchecker.TaskChecker{
			Config:          s.Config,
			Logger:          s.Logger,
			TaskRepository:  s.TaskRepo,
			MessagingClient: s.MessagingClient,
		})
	}

	if isAllModules || slices.Contains(s.Config.Modules, "emailnotifier") {
		modules = append(modules, &emailnotifier.EmailNotifier{
			Config:          s.Config,
			Logger:          s.Logger,
			MessagingClient: s.MessagingClient,
			EmailClient: &emailnotifier.NullEmailClient{
				Logger: s.Logger,
			},
			// EmailClient: &emailnotifier.SMTPEmailClient{
			// 	Options: emailnotifier.SMTPEmailClientOptions{
			// 		Host:        s.Config.SMTPHost,
			// 		Port:        s.Config.SMTPPort,
			// 		FromName:    s.Config.SMTPFromName,
			// 		FromAddress: s.Config.SMTPFromAddress,
			// 		Password:    s.Config.SMTPPassword,
			// 	}},
		})
	}

	if isAllModules || slices.Contains(s.Config.Modules, "uinotifier") {
		modules = append(modules, &uinotifier.UINotifier{
			Config:          s.Config,
			Logger:          s.Logger,
			MessagingClient: s.MessagingClient,
		})
	}

	if isAllModules || slices.Contains(s.Config.Modules, "ui") {
		modules = append(modules, &ui.HTTPService{
			Config:   s.Config,
			Logger:   s.Logger,
			TaskRepo: s.TaskRepo,
			FileExporter: &shared.ExcelFileExporter{
				Logger: s.Logger},
		})
	}

	for _, m := range modules {
		m := m
		g.Go(func() error {
			s.Logger.Info("run app module", slog.String("module", m.Name()))
			return m.Run(ctx)
		})
	}

	g.Go(func() error {
		<-ctx.Done()

		s.Logger.Info("application is shutting down...")

		for _, m := range modules {
			m.Close()
		}

		_ = s.MessagingClient.Close()
		_ = s.TaskRepo.Close()

		return nil
	})

	return g.Wait()
}
