package internal

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// PostgreSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service struct {
	Config      *Config
	Logger      *slog.Logger
	DB          *sql.DB
	TaskRepo    TaskRepository
	EmailClient EmailClient
	Server      *http.Server
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

	if err := s.initEmailClient(ctx); err != nil {
		s.Logger.Error("init email client", "error", err)
		os.Exit(1)
	}

	s.initHTTPServer(ctx)

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

func (s *Service) initHTTPServer(ctx context.Context) {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Handle("/ui/static/*", http.FileServer(http.FS(UIStaticFS)))
	router.Method(http.MethodGet, "/ui", &GetUIHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks", &GetUITasksHandler{s.TaskRepo, s.Logger})
	router.NotFound(NotFound)

	s.Server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         s.Config.Addr,
		Handler:      router,
	}
}

func (s *Service) serve(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		<-ctx.Done()

		s.Logger.Info("application is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = s.Server.Shutdown(ctx)
		_ = s.DB.Close()

		errChan <- nil
	}()

	s.Logger.Info("application is running", "port", s.Server.Addr)

	if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return <-errChan
}
