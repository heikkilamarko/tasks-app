package internal

import (
	"log/slog"
	"slices"
	"tasks-app/internal/modules/emailnotifier"
	"tasks-app/internal/modules/taskchecker"
	"tasks-app/internal/modules/ui"
	"tasks-app/internal/shared"
)

const (
	AppModuleUI                = "ui"
	AppModuleTaskChecker       = "taskchecker"
	AppModuleEmailNotifierNull = "emailnotifier:null"
	AppModuleEmailNotifierSMTP = "emailnotifier:smtp"
)

func (a *App) createModules() error {
	modules := make(map[string]shared.AppModule)

	if slices.Contains(a.Config.Shared.Modules, AppModuleUI) {
		logger := a.Logger.With(slog.String("module", AppModuleUI))

		modules[AppModuleUI] = &ui.Module{
			Config:         a.Config,
			Logger:         logger,
			TaskRepository: a.TaskRepository,
			TaskAttachmentsRepository: &shared.FileTaskAttachmentsRepository{
				Config: a.Config,
			},
			FileExporter: &shared.ExcelFileExporter{
				Logger: logger,
			},
		}
	}

	if slices.Contains(a.Config.Shared.Modules, AppModuleTaskChecker) {
		logger := a.Logger.With(slog.String("module", AppModuleTaskChecker))

		modules[AppModuleTaskChecker] = &taskchecker.Module{
			Config:          a.Config,
			Logger:          logger,
			TaskRepository:  a.TaskRepository,
			MessagingClient: a.MessagingClient,
		}
	}

	if slices.Contains(a.Config.Shared.Modules, AppModuleEmailNotifierNull) {
		logger := a.Logger.With(slog.String("module", AppModuleEmailNotifierNull))

		modules[AppModuleEmailNotifierNull] = &emailnotifier.Module{
			Config:          a.Config,
			Logger:          logger,
			MessagingClient: a.MessagingClient,
			EmailClient: &emailnotifier.NullEmailClient{
				Logger: logger,
			},
		}
	}

	if slices.Contains(a.Config.Shared.Modules, AppModuleEmailNotifierSMTP) {
		logger := a.Logger.With(slog.String("module", AppModuleEmailNotifierSMTP))

		modules[AppModuleEmailNotifierSMTP] = &emailnotifier.Module{
			Config:          a.Config,
			Logger:          logger,
			MessagingClient: a.MessagingClient,
			EmailClient: &emailnotifier.SMTPEmailClient{
				Config: a.Config,
			},
		}
	}

	a.Modules = modules

	return nil
}
