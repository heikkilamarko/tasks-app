package internal

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type GetUIHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *GetUIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepo.GetAll(r.Context())
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	data := struct{ Tasks []*Task }{Tasks: tasks}

	if err := UITemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

type GetUITasksHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *GetUITasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepo.GetAll(r.Context())
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	data := struct{ Tasks []*Task }{Tasks: tasks}

	if err := UITemplates.ExecuteTemplate(w, "tasks_table.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

type GetUITaskHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *GetUITaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	task, err := h.TaskRepo.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	if err := UITemplates.ExecuteTemplate(w, "tasks_table_task.html", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

type PutUITaskHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *PutUITaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if len(name) < 1 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	var expiresAt *time.Time
	expiresAtStr := r.FormValue("expires_at")
	if expiresAtStr != "" {
		expiresAtTemp, err := time.Parse("2006-01-02T15:04", expiresAtStr)
		if err != nil {
			http.Error(w, "invalid expires_at format", http.StatusBadRequest)
			return
		}
		expiresAt = &expiresAtTemp
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	task, err := h.TaskRepo.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	task.Name = name
	task.ExpiresAt = expiresAt

	err = h.TaskRepo.Update(r.Context(), task)
	if err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	if err := UITemplates.ExecuteTemplate(w, "tasks_table_task.html", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

type GetUITaskEditHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *GetUITaskEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	task, err := h.TaskRepo.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	if err := UITemplates.ExecuteTemplate(w, "tasks_table_task_edit.html", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

type DeleteUITaskHandler struct {
	TaskRepo TaskRepository
	Logger   *slog.Logger
}

func (h *DeleteUITaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	task, err := h.TaskRepo.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	err = h.TaskRepo.Delete(r.Context(), id)
	if err != nil {
		h.Logger.Error("delete task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
