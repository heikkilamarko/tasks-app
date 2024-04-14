package internal

import (
	"context"
	"fmt"
	"tasks-app/internal/shared"
)

const (
	AppServiceDBPostgres      = "db:postgres"
	AppServiceAttachmentsFile = "attachments:file"
	AppServiceAttachmentsNATS = "attachments:nats"
	AppServiceMessagingNATS   = "messaging:nats"
)

func (a *App) createServices(ctx context.Context) error {
	var err error

	if a.Config.IsServiceEnabled(AppServiceDBPostgres) {
		a.DB, err = shared.NewPostgresDB(ctx, a.Config)
		if err != nil {
			return fmt.Errorf("create service %s: %w", AppServiceDBPostgres, err)
		}
		a.TxManager = shared.NewSQLTxManager(a.DB)
		a.TaskRepository = shared.NewPostgresTaskRepository(a.DB)
	}

	if a.Config.IsServiceEnabled(AppServiceAttachmentsNATS) {
		a.TaskAttachmentsRepository, err = shared.NewNATSTaskAttachmentsRepository(a.Config, a.Logger)
		if err != nil {
			return fmt.Errorf("create service %s: %w", AppServiceAttachmentsNATS, err)
		}
	}

	if a.Config.IsServiceEnabled(AppServiceAttachmentsFile) {
		a.TaskAttachmentsRepository = &shared.FileTaskAttachmentsRepository{
			Config: a.Config,
		}
	}

	if a.Config.IsServiceEnabled(AppServiceMessagingNATS) {
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

	if a.TaskAttachmentsRepository != nil {
		errs = append(errs, a.TaskAttachmentsRepository.Close())
	}

	if a.DB != nil {
		errs = append(errs, a.DB.Close())
	}

	return errs
}
