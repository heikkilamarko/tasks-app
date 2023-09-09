package internal

import (
	"context"
	"fmt"
	"slices"
	"tasks-app/internal/shared"
)

func (a *App) createServices(ctx context.Context) error {
	var err error

	if slices.Contains(a.Config.Services, "db:postgres") {
		a.TaskRepository, err = shared.NewPostgresTaskRepository(ctx, shared.PostgresTaskRepositoryOptions{
			ConnectionString: a.Config.PostgresConnectionString,
			Logger:           a.Logger,
		})
		if err != nil {
			return fmt.Errorf("create service db:postgres: %w", err)
		}
	}

	if slices.Contains(a.Config.Services, "messaging:nats") {
		a.MessagingClient, err = shared.NewNATSMessagingClient(shared.NATSMessagingClientOptions{
			NATSURL:   a.Config.NATSURL,
			NATSToken: a.Config.NATSToken,
			Logger:    a.Logger,
		})
		if err != nil {
			return fmt.Errorf("create service messaging:nats: %w", err)
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
