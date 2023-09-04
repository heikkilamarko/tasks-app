package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	Config          *Config
	Logger          *slog.Logger
	TaskRepo        TaskRepository
	MessagingClient MessagingClient
	EmailClient     EmailClient
	FileExporter    FileExporter
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
	c := &Config{}
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

	if s.TaskRepo, err = NewPostgresTaskRepository(ctx, PostgresTaskRepositoryOptions{
		ConnectionString: s.Config.PostgresConnectionString,
		Logger:           s.Logger,
	}); err != nil {
		return err
	}

	if s.MessagingClient, err = NewNATSMessagingClient(NATSMessagingClientOptions{
		NATSURL:   s.Config.NATSURL,
		NATSToken: s.Config.NATSToken,
		Logger:    s.Logger,
	}); err != nil {
		return err
	}

	s.EmailClient = &NullEmailClient{s.Logger}

	// s.EmailClient = &SMTPEmailClient{SMTPEmailClientOptions{
	// 	Host:        s.Config.SMTPHost,
	// 	Port:        s.Config.SMTPPort,
	// 	FromName:    s.Config.SMTPFromName,
	// 	FromAddress: s.Config.SMTPFromAddress,
	// 	Password:    s.Config.SMTPPassword,
	// }}

	s.FileExporter = &ExcelFileExporter{s.Logger}

	appModules := []AppModule{
		&TaskChecker{s.Config, s.Logger, s.TaskRepo, s.MessagingClient},
		&UINotifier{s.Config, s.Logger, s.MessagingClient},
		&EmailNotifier{s.Config, s.Logger, s.MessagingClient, s.EmailClient},
		&HTTPService{s.Config, s.Logger, s.TaskRepo, s.FileExporter, nil},
	}

	for _, am := range appModules {
		am := am
		g.Go(func() error { return am.Run(ctx) })
	}

	g.Go(func() error {
		<-ctx.Done()
		s.Logger.Info("application is shutting down...")

		for _, am := range appModules {
			am.Close()
		}

		_ = s.MessagingClient.Close()
		_ = s.TaskRepo.Close()

		return nil
	})

	return g.Wait()
}
