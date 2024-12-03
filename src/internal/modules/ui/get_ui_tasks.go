package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasks struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUITasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err := txc.TaskRepository.GetActive(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)

		return h.Renderer.Render(w, "active_tasks_table.html", vm)
	})
}
