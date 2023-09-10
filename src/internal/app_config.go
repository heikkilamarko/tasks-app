package internal

import "tasks-app/internal/shared"

func (a *App) loadConfig() error {
	config := &shared.Config{}

	if err := config.Load(); err != nil {
		return err
	}

	a.Config = config
	return nil
}
