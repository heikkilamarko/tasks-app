package internal

import (
	"slices"
	"tasks-app/internal/modules/emailnotifier"
	"tasks-app/internal/modules/taskchecker"
	"tasks-app/internal/modules/ui"
	"tasks-app/internal/modules/uinotifier"
	"tasks-app/internal/shared"
)

func (a *App) createModules() error {
	var modules []shared.AppModule

	if slices.Contains(a.Config.Modules, "taskchecker") {
		modules = append(modules, &taskchecker.TaskChecker{
			Config:          a.Config,
			Logger:          a.Logger,
			TaskRepository:  a.TaskRepository,
			MessagingClient: a.MessagingClient,
		})
	}

	if slices.Contains(a.Config.Modules, "emailnotifier") {
		var emailClient emailnotifier.EmailClient

		if slices.Contains(a.Config.Services, "emailnotifier:null") {
			emailClient = &emailnotifier.NullEmailClient{
				Logger: a.Logger,
			}
		}

		if slices.Contains(a.Config.Services, "emailnotifier:smtp") {
			emailClient = &emailnotifier.SMTPEmailClient{
				Options: emailnotifier.SMTPEmailClientOptions{
					Host:        a.Config.SMTPHost,
					Port:        a.Config.SMTPPort,
					FromName:    a.Config.SMTPFromName,
					FromAddress: a.Config.SMTPFromAddress,
					Password:    a.Config.SMTPPassword,
				}}
		}

		modules = append(modules, &emailnotifier.EmailNotifier{
			Config:          a.Config,
			Logger:          a.Logger,
			MessagingClient: a.MessagingClient,
			EmailClient:     emailClient,
		})
	}

	if slices.Contains(a.Config.Modules, "uinotifier") {
		modules = append(modules, &uinotifier.UINotifier{
			Config:          a.Config,
			Logger:          a.Logger,
			MessagingClient: a.MessagingClient,
		})
	}

	if slices.Contains(a.Config.Modules, "ui") {
		modules = append(modules, &ui.UI{
			Config:         a.Config,
			Logger:         a.Logger,
			TaskRepository: a.TaskRepository,
			FileExporter: &shared.ExcelFileExporter{
				Logger: a.Logger},
		})
	}

	a.Modules = modules

	return nil
}

func (a *App) closeModules() []error {
	var errs []error

	for _, m := range a.Modules {
		errs = append(errs, m.Close())
	}

	return errs
}
