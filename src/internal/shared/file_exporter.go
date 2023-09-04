package shared

type FileExporter interface {
	Export(tasks []*Task) ([]byte, error)
}
