package internal

import (
	"log/slog"
	"os"
)

func (a *App) createLogger() (err error) {
	level := slog.LevelWarn

	levelEnv := os.Getenv("APP_SHARED_LOG_LEVEL")
	if levelEnv != "" {
		err = level.UnmarshalText([]byte(levelEnv))
	}

	a.Logger = slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)

	slog.SetDefault(a.Logger)
	slog.SetLogLoggerLevel(level)

	return err
}
