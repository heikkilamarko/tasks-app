package ui

import "tasks-app/internal/shared"

type TasksViewModel struct {
	Tasks         []*shared.Task
	IsCreatingNew bool
}
