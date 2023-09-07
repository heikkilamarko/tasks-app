package shared

import "net/http"

type FileExporter interface {
	ExportTasks(w http.ResponseWriter, tasks []*Task, name string) error
}
