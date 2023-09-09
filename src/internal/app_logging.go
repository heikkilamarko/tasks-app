package internal

import (
	"log/slog"
	"os"
)

func (a *App) initDefaultLogger() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		),
	)
}

func (a *App) initLogger() {
	level := slog.LevelInfo
	level.UnmarshalText([]byte(a.Config.LogLevel))

	a.Logger = slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)
}
