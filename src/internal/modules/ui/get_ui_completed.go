package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUICompleted struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUICompleted) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err := txc.TaskRepository.GetCompleted(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)
		vm.UI.Title = "Completed"

		return h.Renderer.Render(w, "completed_tasks.html", vm)
	})
}
