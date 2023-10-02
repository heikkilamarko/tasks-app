package ui

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"tasks-app/internal/shared"
	"time"

	"github.com/go-chi/chi/v5"
)

type TaskRequest struct {
	ID int
}

type NewTaskRequest struct {
	Name      string
	ExpiresAt *time.Time
}

type UpdateTaskRequest struct {
	ID        int
	Name      string
	ExpiresAt *time.Time
}

type TasksResponse struct {
	Tasks         []*shared.Task
	IsCreatingNew bool
}

func ParseTaskRequest(r *http.Request) (*TaskRequest, error) {
	var errs []error

	id, err := ParseTaskID(chi.URLParam(r, "id"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &TaskRequest{id}, nil
}

func ParseNewTaskRequest(r *http.Request) (*NewTaskRequest, error) {
	var errs []error

	name, err := ParseTaskName(r.FormValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	expiresAt, err := ParseTaskExpiresAt(r.FormValue("expires_at"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &NewTaskRequest{name, expiresAt}, nil
}

func ParseUpdateTaskRequest(r *http.Request) (*UpdateTaskRequest, error) {
	var errs []error

	id, err := ParseTaskID(chi.URLParam(r, "id"))
	if err != nil {
		errs = append(errs, err)
	}

	name, err := ParseTaskName(r.FormValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	expiresAt, err := ParseTaskExpiresAt(r.FormValue("expires_at"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &UpdateTaskRequest{id, name, expiresAt}, nil
}

func ParseTaskAttachments(r *http.Request, taskID int, taskAttachmentsRepository shared.TaskAttachmentsRepository) ([]string, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return nil, err
	}

	existingAttachments := r.Form["attachments"]
	newAttachments := r.MultipartForm.File["attachments"]

	if err = taskAttachmentsRepository.SaveAttachments(r.Context(), taskID, newAttachments); err != nil {
		return nil, err
	}

	var names []string

	for _, name := range existingAttachments {
		names = append(names, name)
	}

	for _, header := range newAttachments {
		names = append(names, header.Filename)
	}

	slices.Sort(names)

	return slices.Compact(names), nil
}

func ParseTaskID(value string) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil || v < 1 {
		return 0, errors.New("id: required, must be an integer greater than 0")
	}

	return v, nil
}

func ParseTaskName(value string) (string, error) {
	l := len(value)
	if l < 1 || 200 < l {
		return "", errors.New("name: required, must be between 1 and 200 characters")
	}
	return value, nil
}

func ParseTaskExpiresAt(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	v, err := ParseTime(value)
	if err != nil {
		return nil, errors.New("expires_at: must be in ISO 8601 format")
	}

	return v, nil
}

func ParseTime(t string) (*time.Time, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	v, err := time.ParseInLocation(timeFormatISO, t, l)
	if err != nil {
		return nil, err
	}

	v = v.UTC()

	return &v, nil
}
