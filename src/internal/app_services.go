package internal

import (
	"context"
	"fmt"
	"slices"
	"tasks-app/internal/shared"
)

const (
	AppServiceDBPostgres    = "db:postgres"
	AppServiceMessagingNATS = "messaging:nats"
)

func (a *App) createServices(ctx context.Context) error {
	var err error

	if slices.Contains(a.Config.Shared.Services, AppServiceDBPostgres) {
		a.TaskRepository, err = shared.NewPostgresTaskRepository(ctx, a.Config)
		if err != nil {
			return fmt.Errorf("create service %s: %w", AppServiceDBPostgres, err)
		}
	}

	if slices.Contains(a.Config.Shared.Services, AppServiceMessagingNATS) {
		a.MessagingClient, err = shared.NewNATSMessagingClient(a.Config, a.Logger)
		if err != nil {
			return fmt.Errorf("create service %s: %w", AppServiceMessagingNATS, err)
		}
	}

	return nil
}

func (a *App) closeServices() []error {
	var errs []error

	if a.MessagingClient != nil {
		errs = append(errs, a.MessagingClient.Close())
	}

	if a.TaskRepository != nil {
		errs = append(errs, a.TaskRepository.Close())
	}

	return errs
}
