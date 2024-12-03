package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUI struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err := txc.TaskRepository.GetActive(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)
		vm.UI.Title = "Active"

		return h.Renderer.Render(w, "active_tasks.html", vm)
	})
}
