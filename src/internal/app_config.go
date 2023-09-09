package internal

import (
	"log/slog"
	"tasks-app/internal/shared"
)

func (a *App) loadConfig() error {
	c := &shared.Config{}
	if err := c.Load(); err != nil {
		return err
	}

	slog.Debug("app config", slog.Any("config", c))

	a.Config = c

	return nil
}
