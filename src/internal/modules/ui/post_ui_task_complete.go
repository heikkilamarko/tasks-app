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

	var task *shared.Task

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		task, err = txc.TaskRepository.GetByID(r.Context(), req.ID)
		return err
	})

	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	task.SetCompleted()

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		err = txc.TaskRepository.Update(r.Context(), task)
		return err
	})

	if err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
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
