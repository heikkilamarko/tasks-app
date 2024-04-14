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
			return fmt.Errorf("create postgres connection: %w", err)
		}
	}

	if a.Config.IsServiceEnabled(AppServiceAttachmentsNATS) || a.Config.IsServiceEnabled(AppServiceMessagingNATS) {
		a.NATSConn, err = shared.NewNATSConn(a.Config, a.Logger)
		if err != nil {
			return fmt.Errorf("create nats connection: %w", err)
		}
	}

	if a.Config.IsServiceEnabled(AppServiceDBPostgres) {
		a.TxManager = shared.NewSQLTxManager(a.DB)
		a.TaskRepository = shared.NewPostgresTaskRepository(a.DB)
	}

	if a.Config.IsServiceEnabled(AppServiceAttachmentsNATS) {
		a.TaskAttachmentsRepository, err = shared.NewNATSTaskAttachmentsRepository(a.NATSConn, a.Logger)
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
		a.MessagingClient, err = shared.NewNATSMessagingClient(a.NATSConn, a.Logger)
		if err != nil {
			return fmt.Errorf("create service %s: %w", AppServiceMessagingNATS, err)
		}
	}

	return nil
}

func (a *App) closeServices() []error {
	var errs []error

	if a.NATSConn != nil {
		errs = append(errs, a.NATSConn.Drain())
	}

	if a.DB != nil {
		errs = append(errs, a.DB.Close())
	}

	return errs
}
