package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITaskComplete struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *PostUITaskComplete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseTaskRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		task, err := txc.TaskRepository.GetByID(r.Context(), req.ID)
		if err != nil {
			return err
		}

		task.SetCompleted()

		return txc.TaskRepository.Update(r.Context(), task)
	})

	if err != nil {
		if err == shared.ErrNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			h.Logger.Error("complete task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	var tasks []*shared.Task

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

	h.Renderer.Render(w, "active_tasks_table.html", vm)
}
