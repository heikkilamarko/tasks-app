package internal

import (
	"log/slog"
	"net/http"
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

	if err := UITemplates.ExecuteTemplate(w, "tasks.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
