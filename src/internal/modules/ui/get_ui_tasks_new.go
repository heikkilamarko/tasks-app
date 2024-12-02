package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksNew struct {
	TxProvider shared.TxProvider
	Renderer   Renderer
	Logger     *slog.Logger
}

func (h *GetUITasksNew) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxProvider.Transact(func(adapters shared.TxAdapters) error {
		tasks, err := adapters.TaskRepository.GetActive(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)
		vm.IsCreatingNew = true

		return h.Renderer.Render(w, "active_tasks_table.html", vm)
	})
}
