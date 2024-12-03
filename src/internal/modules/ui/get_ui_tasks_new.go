package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksNew struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUITasksNew) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tasks []*shared.Task
	var err error

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err = txc.TaskRepository.GetActive(r.Context(), 0, 50)
		return err
	})

	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTasksResponse(r, tasks)
	vm.IsCreatingNew = true

	h.Renderer.Render(w, "active_tasks_table.html", vm)
}
