package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUICompletedTasks struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUICompletedTasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err := txc.TaskRepository.GetCompleted(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get completed tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)

		return h.Renderer.Render(w, "completed_tasks_table.html", vm)
	})
}
