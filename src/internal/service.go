package internal

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	// PostgreSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service struct {
	Config          *Config
	Logger          *slog.Logger
	DB              *sql.DB
	TaskRepo        TaskRepository
	FileExporter    FileExporter
	MessagingClient MessagingClient
	EmailClient     EmailClient
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

	if err := s.initDB(ctx); err != nil {
		s.Logger.Error("init db", "error", err)
		os.Exit(1)
	}

	if err := s.initFileExporters(ctx); err != nil {
		s.Logger.Error("init file exporters", "error", err)
		os.Exit(1)
	}

	if err := s.initMessagingClient(ctx); err != nil {
		s.Logger.Error("init messaging client", "error", err)
		os.Exit(1)
	}

	if err := s.initEmailClient(ctx); err != nil {
		s.Logger.Error("init email client", "error", err)
		os.Exit(1)
	}

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

func (s *Service) initDB(ctx context.Context) error {
	db, err := sql.Open("pgx", s.Config.PostgresConnectionString)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	s.DB = db
	s.TaskRepo = &PostgresTaskRepository{s.DB}

	return nil
}

func (s *Service) initFileExporters(ctx context.Context) error {
	s.FileExporter = &ExcelFileExporter{s.Logger}
	return nil
}

func (s *Service) initMessagingClient(ctx context.Context) error {
	client, err := NewNATSMessagingClient(NATSMessagingClientOptions{
		NATSURL:   s.Config.NATSURL,
		NATSToken: s.Config.NATSToken,
		Logger:    s.Logger,
	})
	if err != nil {
		return err
	}

	s.MessagingClient = client

	return nil
}

func (s *Service) initEmailClient(ctx context.Context) error {
	s.EmailClient = &NullEmailClient{s.Logger}

	// s.EmailClient = &SMTPEmailClient{SMTPEmailClientOptions{
	// 	Host:        s.Config.SMTPHost,
	// 	Port:        s.Config.SMTPPort,
	// 	FromName:    s.Config.SMTPFromName,
	// 	FromAddress: s.Config.SMTPFromAddress,
	// 	Password:    s.Config.SMTPPassword,
	// }}

	return nil
}

func (s *Service) serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

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
		_ = s.DB.Close()

		return nil
	})

	return g.Wait()
}
